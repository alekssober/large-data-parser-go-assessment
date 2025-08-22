package csvimporter

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aleks/products-service/internal/product"
	"github.com/google/uuid"
)

type Importer struct{ Repo product.Repository }

func New(repo product.Repository) *Importer { return &Importer{Repo: repo} }

type row struct {
	Name         string
	Category     string
	Price        string
	Description  string
	BrandName    string
	StockQty     string
	Manufacturer string
	SKU          string
	Weight       string
	Color        string
	Currency     string
}

func (im *Importer) FromCSV(ctx context.Context, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))
	r.TrimLeadingSpace = true

	headers, err := r.Read()
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}
	idx := indexCols(headers)
	if idx["sku"] < 0 || idx["product_name"] < 0 || idx["product_category"] < 0 || idx["product_price"] < 0 {
		return errors.New("missing one of required headers: sku,product_name,product_category,product_price")
	}

	defaultCurrency := strings.ToUpper(strings.TrimSpace(os.Getenv("CSV_DEFAULT_CURRENCY")))
	if defaultCurrency == "" {
		defaultCurrency = "USD"
	}

	dedup := make(map[string]struct{})
	line := 1
	for {
		line++
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("skip line %d: %v", line, err)
			continue
		}

		row := row{
			SKU:          safe(rec, idx["sku"]),
			Name:         safe(rec, idx["product_name"]),
			Category:     safe(rec, idx["product_category"]),
			Price:        safe(rec, idx["product_price"]),
			Description:  safe(rec, idx["product_description"]),
			BrandName:    safe(rec, idx["brand_name"]),
			StockQty:     safe(rec, idx["stock_quantity"]),
			Manufacturer: safe(rec, idx["manufacturer"]),
			Weight:       safe(rec, idx["weight"]),
			Color:        safe(rec, idx["color"]),
			Currency:     defaultCurrency,
		}

		p, ok := normalize(row)
		if !ok {
			continue
		}

		if _, seen := dedup[p.SKU]; seen {
			continue
		}
		dedup[p.SKU] = struct{}{}

		if err := im.Repo.Upsert(ctx, p); err != nil {
			log.Printf("upsert sku=%s failed: %v", p.SKU, err)
		}
	}
	return nil
}

func indexCols(headers []string) map[string]int {
	m := map[string]int{
		"product_name": -1, "product_category": -1, "product_price": -1, "product_description": -1,
		"brand_name": -1, "stock_quantity": -1, "manufacturer": -1, "sku": -1, "weight": -1,
		"color": -1, "currency": -1,
	}
	for i, h := range headers {
		key := strings.ToLower(strings.TrimSpace(h))
		if _, ok := m[key]; ok {
			m[key] = i
		}
	}
	return m
}

func safe(rec []string, i int) string {
	if i >= 0 && i < len(rec) {
		return rec[i]
	}
	return ""
}

func normalize(r row) (product.Product, bool) {
	if strings.TrimSpace(r.SKU) == "" ||
		strings.TrimSpace(r.Name) == "" ||
		strings.TrimSpace(r.Category) == "" ||
		strings.TrimSpace(r.Price) == "" {
		return product.Product{}, false
	}

	cents, ok := parsePriceToCents(r.Price)
	if !ok || cents < 0 {
		return product.Product{}, false
	}

	qty, err := strconv.Atoi(strings.TrimSpace(strings.ReplaceAll(r.StockQty, ",", "")))
	if err != nil || qty < 0 {
		qty = 0
	}

	wg, ok := parseFloatToInt(r.Weight)
	if !ok || wg < 0 {
		wg = 0
	}

	currency := strings.ToUpper(strings.TrimSpace(r.Currency))
	if currency == "" {
		currency = "USD"
	}

	return product.Product{
		ID:            uuid.NewString(),
		SKU:           strings.TrimSpace(r.SKU),
		Name:          strings.TrimSpace(r.Name),
		Category:      strings.TrimSpace(r.Category),
		Description:   strings.TrimSpace(r.Description),
		BrandName:     strings.TrimSpace(r.BrandName),
		StockQuantity: qty,
		Manufacturer:  strings.TrimSpace(r.Manufacturer),
		WeightGrams:   wg,
		Color:         strings.TrimSpace(r.Color),
		PriceCents:    cents,
		Currency:      currency,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}, true
}

func parseFloatToInt(s string) (int, bool) {
	s = strings.TrimSpace(strings.ReplaceAll(s, ",", ""))
	if s == "" {
		return 0, false
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	return int(math.Round(f)), true
}

func parsePriceToCents(s string) (int, bool) {
	s = strings.TrimSpace(strings.ReplaceAll(s, ",", ""))
	if s == "" {
		return 0, false
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	cents := int(math.Round(f * 100))
	return cents, true
}
