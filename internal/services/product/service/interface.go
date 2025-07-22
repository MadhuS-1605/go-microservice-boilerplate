package service

import (
	"context"
	"go-microservice-boilerplate/internal/services/product/model"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *model.CreateProductRequest) (*model.Product, error)
	GetProduct(ctx context.Context, id string) (*model.Product, error)
	UpdateProduct(ctx context.Context, id string, req *model.UpdateProductRequest) (*model.Product, error)
	DeleteProduct(ctx context.Context, id string) error
	ListProducts(ctx context.Context, page, limit int, search, category string) ([]*model.Product, int64, error)
}
