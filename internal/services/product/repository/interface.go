package repository

import (
	"context"
	"go-microservice-boilerplate/internal/services/product/model"
)

// ProductRepository defines the contract for product data operations
type ProductRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, product *model.Product) error
	GetByID(ctx context.Context, id string) (*model.Product, error)
	GetBySKU(ctx context.Context, sku string) (*model.Product, error)
	Update(ctx context.Context, id string, product *model.Product) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int, search, category string) ([]*model.Product, int64, error)
}

// ProductCache defines the contract for product caching operations
type ProductCache interface {
	Set(ctx context.Context, key string, product *model.Product, expiration int) error
	Get(ctx context.Context, key string) (*model.Product, error)
	Delete(ctx context.Context, key string) error
	SetList(ctx context.Context, key string, products []*model.Product, expiration int) error
	GetList(ctx context.Context, key string) ([]*model.Product, error)
	InvalidatePattern(ctx context.Context, pattern string) error
}
