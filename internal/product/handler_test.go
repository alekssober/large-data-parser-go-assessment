package product

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

type fakeSvc struct{}

func (fakeSvc) List(ctx context.Context, page, pageSize int) ([]Product, int, error) {
	return []Product{{ID: "x"}}, 1, nil
}
func (fakeSvc) Get(ctx context.Context, id string) (Product, error) { return Product{ID: id}, nil }
func (fakeSvc) Stats(ctx context.Context) (Summary, error)          { return Summary{TotalProducts: 1}, nil }

func TestRoutes(t *testing.T) {
	h := NewHandler(fakeSvc{})
	r := chi.NewRouter()
	h.Routes(r)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
}
