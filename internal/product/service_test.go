package product

import (
	"context"
	"testing"
)

type fakeRepo struct{ items []Product }

func (f *fakeRepo) Upsert(ctx context.Context, p Product) error {
	f.items = append(f.items, p)
	return nil
}
func (f *fakeRepo) GetPaginated(ctx context.Context, limit, offset int) ([]Product, int, error) {
	return f.items, len(f.items), nil
}
func (f *fakeRepo) GetByID(ctx context.Context, id string) (Product, error) { return f.items[0], nil }
func (f *fakeRepo) Summary(ctx context.Context) (Summary, error) {
	return Summary{TotalProducts: len(f.items)}, nil
}

func TestServiceListDefaults(t *testing.T) {
	repo := &fakeRepo{items: []Product{{ID: "1"}}}
	svc := NewService(repo)
	items, total, err := svc.List(context.Background(), 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || total != 1 {
		t.Fatalf("unexpected: %v %d", items, total)
	}
}
