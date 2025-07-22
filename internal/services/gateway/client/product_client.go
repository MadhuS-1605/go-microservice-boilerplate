package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go-microservice-boilerplate/internal/config"
	"go-microservice-boilerplate/internal/proto/common"
	"go-microservice-boilerplate/internal/proto/product"
)

type ProductClient struct {
	conn   *grpc.ClientConn
	client product.ProductServiceClient
}

func NewProductClient(cfg *config.Config) (*ProductClient, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Services.Product.Host, cfg.Services.Product.Port)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}

	client := product.NewProductServiceClient(conn)

	return &ProductClient{
		conn:   conn,
		client: client,
	}, nil
}

func (c *ProductClient) Close() error {
	return c.conn.Close()
}

func (c *ProductClient) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.ProductResponse, error) {
	return c.client.CreateProduct(ctx, req)
}

func (c *ProductClient) GetProduct(ctx context.Context, req *product.GetProductRequest) (*product.ProductResponse, error) {
	return c.client.GetProduct(ctx, req)
}

func (c *ProductClient) UpdateProduct(ctx context.Context, req *product.UpdateProductRequest) (*product.ProductResponse, error) {
	return c.client.UpdateProduct(ctx, req)
}

func (c *ProductClient) DeleteProduct(ctx context.Context, req *product.DeleteProductRequest) (*common.StatusResponse, error) {
	return c.client.DeleteProduct(ctx, req)
}

func (c *ProductClient) ListProducts(ctx context.Context, req *product.ListProductsRequest) (*product.ListProductsResponse, error) {
	return c.client.ListProducts(ctx, req)
}
