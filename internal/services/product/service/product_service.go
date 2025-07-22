package service

import (
	"context"
	"fmt"

	"go-microservice-boilerplate/internal/services/product/model"
	"go-microservice-boilerplate/internal/services/product/repository"
)

type productService struct {
	repo  repository.ProductRepository
	cache repository.ProductCache
}

func NewProductService(repo repository.ProductRepository, cache repository.ProductCache) ProductService {
	return &productService{
		repo:  repo,
		cache: cache,
	}
}

func (s *productService) CreateProduct(ctx context.Context, req *model.CreateProductRequest) (*model.Product, error) {
	// Check if product with same SKU already exists
	existingProduct, err := s.repo.GetBySKU(ctx, req.SKU)
	if err == nil && existingProduct != nil {
		return nil, fmt.Errorf("product with SKU %s already exists", req.SKU)
	}

	product := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    req.Quantity,
		Category:    req.Category,
		SKU:         req.SKU,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Cache the product
	cacheKey := fmt.Sprintf("product:%s", product.ID.Hex())
	s.cache.Set(ctx, cacheKey, product, 3600) // 1 hour

	return product, nil
}

func (s *productService) GetProduct(ctx context.Context, id string) (*model.Product, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("product:%s", id)
	if product, err := s.cache.Get(ctx, cacheKey); err == nil && product != nil {
		return product, nil
	}

	// Get from database
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Cache the result
	s.cache.Set(ctx, cacheKey, product, 3600)

	return product, nil
}

func (s *productService) UpdateProduct(ctx context.Context, id string, req *model.UpdateProductRequest) (*model.Product, error) {
	// Get existing product
	product, err := s.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Quantity >= 0 {
		product.Quantity = req.Quantity
	}
	if req.Category != "" {
		product.Category = req.Category
	}

	if err := s.repo.Update(ctx, id, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	// Update cache
	cacheKey := fmt.Sprintf("product:%s", id)
	s.cache.Set(ctx, cacheKey, product, 3600)

	return product, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	// Remove from cache
	cacheKey := fmt.Sprintf("product:%s", id)
	s.cache.Delete(ctx, cacheKey)

	return nil
}

func (s *productService) ListProducts(ctx context.Context, page, limit int, search, category string) ([]*model.Product, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	products, total, err := s.repo.List(ctx, page, limit, search, category)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}

	return products, total, nil
}
