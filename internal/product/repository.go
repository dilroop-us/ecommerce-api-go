package product

import (
	"context"

	"github.com/dilroop-us/ecommerce-go/internal/db"
)

type Repository struct{ q *db.Queries }

func NewRepository(q *db.Queries) *Repository { return &Repository{q: q} }

func (r *Repository) List(ctx context.Context) ([]db.Product, error) {
	return r.q.ListProducts(ctx)
}

func (r *Repository) Create(ctx context.Context, name string, price float64) (db.Product, error) {
	return r.q.CreateProduct(ctx, db.CreateProductParams{Name: name, Price: price})
}
