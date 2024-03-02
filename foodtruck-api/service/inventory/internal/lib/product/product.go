package product

import (
	"context"
	"foodtruck/service/inventory/internal/lib/db/code_gen/ent"
)

type ProductService struct {
	client *ent.Client
	ctx    context.Context
}

func New(client *ent.Client) *ProductService {
	return &ProductService{
		client: client,
		ctx:    context.Background(),
	}
}

func (p *ProductService) Create(name string, description string) (*ent.Product, error) {
	return p.client.Product.Create().
		SetName(name).
		SetDescription(description).
		Save(p.ctx)
}
