package product

import "time"

type Product struct {
	ID            string    `json:"id"`
	SKU           string    `json:"sku"`
	Name          string    `json:"product_name"`
	Category      string    `json:"product_category"`
	Description   string    `json:"product_description"`
	BrandName     string    `json:"brand_name"`
	StockQuantity int       `json:"stock_quantity"`
	Manufacturer  string    `json:"manufacturer"`
	WeightGrams   int       `json:"weight_grams"`
	Color         string    `json:"color"`
	PriceCents    int       `json:"price_cents"`
	Currency      string    `json:"currency"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Summary struct {
	TotalProducts   int     `json:"total_products"`
	AvgPriceCents   float64 `json:"avg_price_cents"`
	MinPriceCents   int     `json:"min_price_cents"`
	MaxPriceCents   int     `json:"max_price_cents"`
	CategoriesCount int     `json:"categories_count"`
}
