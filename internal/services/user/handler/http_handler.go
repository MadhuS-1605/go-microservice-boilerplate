package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-microservice-boilerplate/internal/services/user/model"
	"go-microservice-boilerplate/internal/services/user/service"
	"go-microservice-boilerplate/internal/utils/response"
)

type UserHTTPHandler struct {
	userService service.UserService
}

func NewUserHTTPHandler(userService service.UserService) *UserHTTPHandler {
	return &UserHTTPHandler{
		userService: userService,
	}
}

func (h *UserHTTPHandler) RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.POST("", h.CreateUser)
		users.GET("/:id", h.GetUser)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
		users.GET("", h.ListUsers)
	}
}

func (h *UserHTTPHandler) CreateUser(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create user", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "User created successfully", user)
}

func (h *UserHTTPHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User retrieved successfully", user)
}

func (h *UserHTTPHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), id, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update user", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User updated successfully", user)
}

func (h *UserHTTPHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	err := h.userService.DeleteUser(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete user", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User deleted successfully", nil)
}

func (h *UserHTTPHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	users, total, err := h.userService.ListUsers(c.Request.Context(), page, limit, search)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list users", err.Error())
		return
	}

	result := gin.H{
		"users": users,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	}

	response.Success(c, http.StatusOK, "Users retrieved successfully", result)
}
