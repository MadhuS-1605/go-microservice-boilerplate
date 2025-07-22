package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-microservice-boilerplate/internal/proto/product"
	"go-microservice-boilerplate/internal/proto/user"
	"go-microservice-boilerplate/internal/services/gateway/client"
	"go-microservice-boilerplate/internal/utils/response"
)

type GatewayHandler struct {
	userClient    *client.UserClient
	productClient *client.ProductClient
}

func NewGatewayHandler(userClient *client.UserClient, productClient *client.ProductClient) *GatewayHandler {
	return &GatewayHandler{
		userClient:    userClient,
		productClient: productClient,
	}
}

func (h *GatewayHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")

	// Health check
	api.GET("/health", h.HealthCheck)

	// User routes
	users := api.Group("/users")
	{
		users.POST("", h.CreateUser)
		users.GET("/:id", h.GetUser)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
		users.GET("", h.ListUsers)
	}

	// Product routes
	products := api.Group("/products")
	{
		products.POST("", h.CreateProduct)
		products.GET("/:id", h.GetProduct)
		products.PUT("/:id", h.UpdateProduct)
		products.DELETE("/:id", h.DeleteProduct)
		products.GET("", h.ListProducts)
	}
}

// HealthCheck godoc
// @Summary Health Check
// @Description Check the health status of the gateway service
// @Tags Health
// @Produce json
// @Router /health [get]
func (h *GatewayHandler) HealthCheck(c *gin.Context) {
	response.Success(c, http.StatusOK, "Gateway is healthy", gin.H{
		"service": "gateway",
		"status":  "ok",
	})
}

// CreateUser godoc
// @Summary Create User
// @Description Create a new user account
// @Tags Users
// @Accept json
// @Produce json
// @Router /users [post]
func (h *GatewayHandler) CreateUser(c *gin.Context) {
	var req user.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	resp, err := h.userClient.CreateUser(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create user", err.Error())
		return
	}

	if !resp.Status.Success {
		response.Error(c, int(resp.Status.Code), resp.Status.Message, nil)
		return
	}

	response.Success(c, http.StatusCreated, resp.Status.Message, resp.User)
}

// GetUser godoc
// @Summary Get User
// @Description Get user by ID
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Router /users/{id} [get]
func (h *GatewayHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	req := &user.GetUserRequest{Id: id}
	resp, err := h.userClient.GetUser(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found", err.Error())
		return
	}

	if !resp.Status.Success {
		response.Error(c, int(resp.Status.Code), resp.Status.Message, nil)
		return
	}

	response.Success(c, http.StatusOK, resp.Status.Message, resp.User)
}

// UpdateUser godoc
// @Summary Update User
// @Description Update user information
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Router /users/{id} [put]
func (h *GatewayHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req user.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	req.Id = id
	resp, err := h.userClient.UpdateUser(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update user", err.Error())
		return
	}

	if !resp.Status.Success {
		response.Error(c, int(resp.Status.Code), resp.Status.Message, nil)
		return
	}

	response.Success(c, http.StatusOK, resp.Status.Message, resp.User)
}

// DeleteUser godoc
// @Summary Delete User
// @Description Delete user by ID
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Router /users/{id} [delete]
func (h *GatewayHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	req := &user.DeleteUserRequest{Id: id}
	resp, err := h.userClient.DeleteUser(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete user", err.Error())
		return
	}

	if !resp.Success {
		response.Error(c, int(resp.Code), resp.Message, nil)
		return
	}

	response.Success(c, http.StatusOK, resp.Message, nil)
}

// ListUsers godoc
// @Summary List Users
// @Description Get paginated list of users
// @Tags Users
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param search query string false "Search term"
// @Router /users [get]
func (h *GatewayHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	req := &user.ListUsersRequest{
		Page:   int32(page),
		Limit:  int32(limit),
		Search: search,
	}

	resp, err := h.userClient.ListUsers(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list users", err.Error())
		return
	}

	if !resp.Status.Success {
		response.Error(c, int(resp.Status.Code), resp.Status.Message, nil)
		return
	}

	result := gin.H{
		"users": resp.Users,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       resp.Total,
			"total_pages": (resp.Total + int32(limit) - 1) / int32(limit),
		},
	}

	response.Success(c, http.StatusOK, resp.Status.Message, result)
}

// CreateProduct godoc
// @Summary Create Product
// @Description Create a new product
// @Tags Products
// @Accept json
// @Produce json
// @Router /products [post]
func (h *GatewayHandler) CreateProduct(c *gin.Context) {
	var req product.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	resp, err := h.productClient.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create product", err.Error())
		return
	}

	if !resp.Status.Success {
		response.Error(c, int(resp.Status.Code), resp.Status.Message, nil)
		return
	}

	response.Success(c, http.StatusCreated, resp.Status.Message, resp.Product)
}

// GetProduct godoc
// @Summary Get Product
// @Description Get product by ID
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Router /products/{id} [get]
func (h *GatewayHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")

	req := &product.GetProductRequest{Id: id}
	resp, err := h.productClient.GetProduct(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Product not found", err.Error())
		return
	}

	if !resp.Status.Success {
		response.Error(c, int(resp.Status.Code), resp.Status.Message, nil)
		return
	}

	response.Success(c, http.StatusOK, resp.Status.Message, resp.Product)
}

// UpdateProduct godoc
// @Summary Update Product
// @Description Update product information
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Router /products/{id} [put]
func (h *GatewayHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var req product.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	req.Id = id
	resp, err := h.productClient.UpdateProduct(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update product", err.Error())
		return
	}

	if !resp.Status.Success {
		response.Error(c, int(resp.Status.Code), resp.Status.Message, nil)
		return
	}

	response.Success(c, http.StatusOK, resp.Status.Message, resp.Product)
}

// DeleteProduct godoc
// @Summary Delete Product
// @Description Delete product by ID
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Router /products/{id} [delete]
func (h *GatewayHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	req := &product.DeleteProductRequest{Id: id}
	resp, err := h.productClient.DeleteProduct(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete product", err.Error())
		return
	}

	if !resp.Success {
		response.Error(c, int(resp.Code), resp.Message, nil)
		return
	}

	response.Success(c, http.StatusOK, resp.Message, nil)
}

// ListProducts godoc
// @Summary List Products
// @Description Get paginated list of products
// @Tags Products
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Param category query string false "Category filter"
// @Router /products [get]
func (h *GatewayHandler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	category := c.Query("category")

	req := &product.ListProductsRequest{
		Page:     int32(page),
		Limit:    int32(limit),
		Search:   search,
		Category: category,
	}

	resp, err := h.productClient.ListProducts(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list products", err.Error())
		return
	}

	if !resp.Status.Success {
		response.Error(c, int(resp.Status.Code), resp.Status.Message, nil)
		return
	}

	result := gin.H{
		"products": resp.Products,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       resp.Total,
			"total_pages": (resp.Total + int32(limit) - 1) / int32(limit),
		},
	}

	response.Success(c, http.StatusOK, resp.Status.Message, result)
}
