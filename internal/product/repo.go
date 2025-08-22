package product

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Upsert(ctx context.Context, p Product) error
	GetPaginated(ctx context.Context, limit, offset int) ([]Product, int, error)
	GetByID(ctx context.Context, id string) (Product, error)
	Summary(ctx context.Context) (Summary, error)
}

type PGRepo struct{ pool *pgxpool.Pool }

func NewPGRepo(pool *pgxpool.Pool) *PGRepo { return &PGRepo{pool: pool} }

func (r *PGRepo) Upsert(ctx context.Context, p Product) error {
	const q = `
INSERT INTO products (id, sku, name, category, description, brand_name, stock_quantity, manufacturer, weight_grams, color, price_cents, currency)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
ON CONFLICT (sku) DO UPDATE SET
name = EXCLUDED.name,
category = EXCLUDED.category,
description = EXCLUDED.description,
brand_name = EXCLUDED.brand_name,
stock_quantity = EXCLUDED.stock_quantity,
manufacturer = EXCLUDED.manufacturer,
weight_grams = EXCLUDED.weight_grams,
color = EXCLUDED.color,
price_cents = EXCLUDED.price_cents,
currency = EXCLUDED.currency,
updated_at = NOW();`
	_, err := r.pool.Exec(ctx, q, p.ID, p.SKU, p.Name, p.Category, p.Description, p.BrandName, p.StockQuantity, p.Manufacturer, p.WeightGrams, p.Color, p.PriceCents, p.Currency)
	return err
}

func (r *PGRepo) GetPaginated(ctx context.Context, limit, offset int) ([]Product, int, error) {
	const listQ = `SELECT id, sku, name, category, description, brand_name, stock_quantity, manufacturer, weight_grams, color, price_cents, currency, created_at, updated_at
FROM products ORDER BY created_at DESC LIMIT $1 OFFSET $2;`
	const countQ = `SELECT COUNT(*) FROM products;`

	rows, err := r.pool.Query(ctx, listQ, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var res []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Category, &p.Description, &p.BrandName, &p.StockQuantity, &p.Manufacturer, &p.WeightGrams, &p.Color, &p.PriceCents, &p.Currency, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		res = append(res, p)
	}
	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}

	var total int
	if err := r.pool.QueryRow(ctx, countQ).Scan(&total); err != nil {
		return nil, 0, err
	}
	return res, total, nil
}

func (r *PGRepo) GetByID(ctx context.Context, id string) (Product, error) {
	const q = `SELECT id, sku, name, category, description, brand_name, stock_quantity, manufacturer, weight_grams, color, price_cents, currency, created_at, updated_at FROM products WHERE id = $1;`
	var p Product
	err := r.pool.QueryRow(ctx, q, id).Scan(&p.ID, &p.SKU, &p.Name, &p.Category, &p.Description, &p.BrandName, &p.StockQuantity, &p.Manufacturer, &p.WeightGrams, &p.Color, &p.PriceCents, &p.Currency, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Product{}, err
	}
	return p, err
}

func (r *PGRepo) Summary(ctx context.Context) (Summary, error) {
	const q = `
SELECT COUNT(*), AVG(price_cents), MIN(price_cents), MAX(price_cents), COUNT(DISTINCT category)
FROM products;`
	var s Summary
	err := r.pool.QueryRow(ctx, q).Scan(&s.TotalProducts, &s.AvgPriceCents, &s.MinPriceCents, &s.MaxPriceCents, &s.CategoriesCount)
	return s, err
}
