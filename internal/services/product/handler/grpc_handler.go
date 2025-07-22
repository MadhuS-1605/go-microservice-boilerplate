package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go-microservice-boilerplate/internal/proto/common"
	"go-microservice-boilerplate/internal/proto/product"
	"go-microservice-boilerplate/internal/services/product/model"
	"go-microservice-boilerplate/internal/services/product/service"
)

type ProductGRPCHandler struct {
	product.UnimplementedProductServiceServer
	productService service.ProductService
}

func NewProductGRPCHandler(productService service.ProductService) *ProductGRPCHandler {
	return &ProductGRPCHandler{
		productService: productService,
	}
}

func (h *ProductGRPCHandler) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.ProductResponse, error) {
	createReq := &model.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    req.Quantity,
		Category:    req.Category,
		SKU:         req.Sku,
	}

	productModel, err := h.productService.CreateProduct(ctx, createReq)
	if err != nil {
		return &product.ProductResponse{
			Status: &common.StatusResponse{
				Code:    int32(codes.Internal),
				Message: err.Error(),
				Success: false,
			},
		}, status.Error(codes.Internal, err.Error())
	}

	return &product.ProductResponse{
		Product: h.modelToProto(productModel),
		Status: &common.StatusResponse{
			Code:    int32(codes.OK),
			Message: "Product created successfully",
			Success: true,
		},
	}, nil
}

func (h *ProductGRPCHandler) GetProduct(ctx context.Context, req *product.GetProductRequest) (*product.ProductResponse, error) {
	productModel, err := h.productService.GetProduct(ctx, req.Id)
	if err != nil {
		return &product.ProductResponse{
			Status: &common.StatusResponse{
				Code:    int32(codes.NotFound),
				Message: err.Error(),
				Success: false,
			},
		}, status.Error(codes.NotFound, err.Error())
	}

	return &product.ProductResponse{
		Product: h.modelToProto(productModel),
		Status: &common.StatusResponse{
			Code:    int32(codes.OK),
			Message: "Product retrieved successfully",
			Success: true,
		},
	}, nil
}

func (h *ProductGRPCHandler) UpdateProduct(ctx context.Context, req *product.UpdateProductRequest) (*product.ProductResponse, error) {
	updateReq := &model.UpdateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    req.Quantity,
		Category:    req.Category,
	}

	productModel, err := h.productService.UpdateProduct(ctx, req.Id, updateReq)
	if err != nil {
		return &product.ProductResponse{
			Status: &common.StatusResponse{
				Code:    int32(codes.Internal),
				Message: err.Error(),
				Success: false,
			},
		}, status.Error(codes.Internal, err.Error())
	}

	return &product.ProductResponse{
		Product: h.modelToProto(productModel),
		Status: &common.StatusResponse{
			Code:    int32(codes.OK),
			Message: "Product updated successfully",
			Success: true,
		},
	}, nil
}

func (h *ProductGRPCHandler) DeleteProduct(ctx context.Context, req *product.DeleteProductRequest) (*common.StatusResponse, error) {
	err := h.productService.DeleteProduct(ctx, req.Id)
	if err != nil {
		return &common.StatusResponse{
			Code:    int32(codes.Internal),
			Message: err.Error(),
			Success: false,
		}, status.Error(codes.Internal, err.Error())
	}

	return &common.StatusResponse{
		Code:    int32(codes.OK),
		Message: "Product deleted successfully",
		Success: true,
	}, nil
}

func (h *ProductGRPCHandler) ListProducts(ctx context.Context, req *product.ListProductsRequest) (*product.ListProductsResponse, error) {
	page := int(req.Page)
	limit := int(req.Limit)
	search := req.Search
	category := req.Category

	products, total, err := h.productService.ListProducts(ctx, page, limit, search, category)
	if err != nil {
		return &product.ListProductsResponse{
			Status: &common.StatusResponse{
				Code:    int32(codes.Internal),
				Message: err.Error(),
				Success: false,
			},
		}, status.Error(codes.Internal, err.Error())
	}

	protoProducts := make([]*product.Product, len(products))
	for i, p := range products {
		protoProducts[i] = h.modelToProto(p)
	}

	return &product.ListProductsResponse{
		Products: protoProducts,
		Total:    int32(total),
		Status: &common.StatusResponse{
			Code:    int32(codes.OK),
			Message: "Products retrieved successfully",
			Success: true,
		},
	}, nil
}

func (h *ProductGRPCHandler) modelToProto(p *model.Product) *product.Product {
	return &product.Product{
		Id:          p.ID.Hex(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Quantity:    p.Quantity,
		Category:    p.Category,
		Sku:         p.SKU,
		CreatedAt:   p.CreatedAt.Unix(),
		UpdatedAt:   p.UpdatedAt.Unix(),
	}
}
