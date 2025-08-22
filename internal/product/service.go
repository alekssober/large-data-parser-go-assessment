package product

import (
	"context"
)

type Service interface {
	List(ctx context.Context, page, pageSize int) ([]Product, int, error)
	Get(ctx context.Context, id string) (Product, error)
	Stats(ctx context.Context) (Summary, error)
}

type service struct{ repo Repository }

func NewService(r Repository) Service { return &service{repo: r} }

func (s *service) List(ctx context.Context, page, pageSize int) ([]Product, int, error) {
	if pageSize <= 0 {
		pageSize = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pageSize
	return s.repo.GetPaginated(ctx, pageSize, offset)
}

func (s *service) Get(ctx context.Context, id string) (Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) Stats(ctx context.Context) (Summary, error) {
	return s.repo.Summary(ctx)
}
