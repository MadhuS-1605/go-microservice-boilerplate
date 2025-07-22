package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-microservice-boilerplate/internal/services/product/model"
	"go-microservice-boilerplate/internal/services/product/service"
	"go-microservice-boilerplate/internal/utils/response"
)

type ProductHTTPHandler struct {
	productService service.ProductService
}

func (h *ProductHTTPHandler) RegisterRoutes(router *gin.RouterGroup) {
	products := router.Group("/products")
	{
		products.POST("", h.CreateProduct)
		products.GET("/:id", h.GetProduct)
		products.PUT("/:id", h.UpdateProduct)
		products.DELETE("/:id", h.DeleteProduct)
		products.GET("", h.ListProducts)
	}
}

func (h *ProductHTTPHandler) CreateProduct(c *gin.Context) {
	var req model.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	product, err := h.productService.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create product", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Product created successfully", product)
}

func (h *ProductHTTPHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")

	product, err := h.productService.GetProduct(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Product not found", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Product retrieved successfully", product)
}

func (h *ProductHTTPHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var req model.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	product, err := h.productService.UpdateProduct(c.Request.Context(), id, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update product", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Product updated successfully", product)
}

func (h *ProductHTTPHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	err := h.productService.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete product", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Product deleted successfully", nil)
}

func (h *ProductHTTPHandler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	category := c.Query("category")

	products, total, err := h.productService.ListProducts(c.Request.Context(), page, limit, search, category)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list products", err.Error())
		return
	}

	result := gin.H{
		"products": products,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	}

	response.Success(c, http.StatusOK, "Products retrieved successfully", result)
}
