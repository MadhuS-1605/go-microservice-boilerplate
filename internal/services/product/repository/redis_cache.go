package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-microservice-boilerplate/internal/database"
	"go-microservice-boilerplate/internal/services/product/model"
)

type productCache struct {
	client *database.Redis
}

// NewRedisProductCache creates a new product cache instance
func NewRedisProductCache(redis *database.Redis) ProductCache {
	return &productCache{
		client: redis,
	}
}

// Set stores a product in cache
func (c *productCache) Set(ctx context.Context, key string, product *model.Product, expiration int) error {
	data, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("failed to marshal product: %w", err)
	}

	err = c.client.Client.Set(ctx, key, data, time.Duration(expiration)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// Get retrieves a product from cache
func (c *productCache) Get(ctx context.Context, key string) (*model.Product, error) {
	data, err := c.client.Client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var product model.Product
	if err := json.Unmarshal([]byte(data), &product); err != nil {
		return nil, fmt.Errorf("failed to unmarshal product: %w", err)
	}

	return &product, nil
}

// Delete removes a product from cache
func (c *productCache) Delete(ctx context.Context, key string) error {
	return c.client.Client.Del(ctx, key).Err()
}

// SetList stores a list of products in cache
func (c *productCache) SetList(ctx context.Context, key string, products []*model.Product, expiration int) error {
	data, err := json.Marshal(products)
	if err != nil {
		return fmt.Errorf("failed to marshal products list: %w", err)
	}

	err = c.client.Client.Set(ctx, key, data, time.Duration(expiration)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("failed to set list cache: %w", err)
	}

	return nil
}

// GetList retrieves a list of products from cache
func (c *productCache) GetList(ctx context.Context, key string) ([]*model.Product, error) {
	data, err := c.client.Client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var products []*model.Product
	if err := json.Unmarshal([]byte(data), &products); err != nil {
		return nil, fmt.Errorf("failed to unmarshal products list: %w", err)
	}

	return products, nil
}

// InvalidatePattern removes all cache keys matching a pattern
func (c *productCache) InvalidatePattern(ctx context.Context, pattern string) error {
	keys, err := c.client.Client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.client.Client.Del(ctx, keys...).Err()
	}

	return nil
}

// CacheProductList caches a product list with pagination info
func (c *productCache) CacheProductList(ctx context.Context, key string, products []*model.Product, total int64, expiration int) error {
	data := map[string]interface{}{
		"products":  products,
		"total":     total,
		"cached_at": time.Now().Unix(),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal product list data: %w", err)
	}

	return c.client.Client.Set(ctx, key, jsonData, time.Duration(expiration)*time.Second).Err()
}

// GetCachedProductList retrieves a cached product list with metadata
func (c *productCache) GetCachedProductList(ctx context.Context, key string) ([]*model.Product, int64, error) {
	data, err := c.client.Client.Get(ctx, key).Result()
	if err != nil {
		return nil, 0, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, 0, fmt.Errorf("failed to unmarshal cached list: %w", err)
	}

	// Extract products
	productsData, ok := result["products"]
	if !ok {
		return nil, 0, fmt.Errorf("products data not found in cache")
	}

	productsJSON, err := json.Marshal(productsData)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal products data: %w", err)
	}

	var products []*model.Product
	if err := json.Unmarshal(productsJSON, &products); err != nil {
		return nil, 0, fmt.Errorf("failed to unmarshal products: %w", err)
	}

	// Extract total
	total, ok := result["total"].(float64)
	if !ok {
		return nil, 0, fmt.Errorf("total not found in cache")
	}

	return products, int64(total), nil
}
