package product

import (
	"context"

	"github.com/dilroop-us/ecommerce-go/internal/db"
)

// No Product struct here — we use the one from model.go

type Service struct{ repo *Repository }

func NewService(repo *Repository) *Service { return &Service{repo: repo} }

func toDomain(p db.Product) Product {
	return Product{
		ID:    p.ID.String(),
		Name:  p.Name,
		Price: p.Price, // will be float64 after step 2
	}
}

func (s *Service) List(ctx context.Context) ([]Product, error) {
	rows, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]Product, len(rows))
	for i := range rows {
		out[i] = toDomain(rows[i])
	}
	return out, nil
}

func (s *Service) Create(ctx context.Context, name string, price float64) (Product, error) {
	row, err := s.repo.Create(ctx, name, price)
	if err != nil {
		return Product{}, err
	}
	return toDomain(row), nil
}
