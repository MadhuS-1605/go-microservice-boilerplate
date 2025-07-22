package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go-microservice-boilerplate/internal/proto/common"
	"go-microservice-boilerplate/internal/proto/user"
	"go-microservice-boilerplate/internal/services/user/model"
	"go-microservice-boilerplate/internal/services/user/service"
)

type UserGRPCHandler struct {
	user.UnimplementedUserServiceServer
	userService service.UserService
}

func NewUserGRPCHandler(userService service.UserService) *UserGRPCHandler {
	return &UserGRPCHandler{
		userService: userService,
	}
}

func (h *UserGRPCHandler) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.UserResponse, error) {
	createReq := &model.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	}

	userModel, err := h.userService.CreateUser(ctx, createReq)
	if err != nil {
		return &user.UserResponse{
			Status: &common.StatusResponse{
				Code:    int32(codes.Internal),
				Message: err.Error(),
				Success: false,
			},
		}, status.Error(codes.Internal, err.Error())
	}

	return &user.UserResponse{
		User: h.modelToProto(userModel),
		Status: &common.StatusResponse{
			Code:    int32(codes.OK),
			Message: "User created successfully",
			Success: true,
		},
	}, nil
}

func (h *UserGRPCHandler) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.UserResponse, error) {
	userModel, err := h.userService.GetUser(ctx, req.Id)
	if err != nil {
		return &user.UserResponse{
			Status: &common.StatusResponse{
				Code:    int32(codes.NotFound),
				Message: err.Error(),
				Success: false,
			},
		}, status.Error(codes.NotFound, err.Error())
	}

	return &user.UserResponse{
		User: h.modelToProto(userModel),
		Status: &common.StatusResponse{
			Code:    int32(codes.OK),
			Message: "User retrieved successfully",
			Success: true,
		},
	}, nil
}

func (h *UserGRPCHandler) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UserResponse, error) {
	updateReq := &model.UpdateUserRequest{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}

	userModel, err := h.userService.UpdateUser(ctx, req.Id, updateReq)
	if err != nil {
		return &user.UserResponse{
			Status: &common.StatusResponse{
				Code:    int32(codes.Internal),
				Message: err.Error(),
				Success: false,
			},
		}, status.Error(codes.Internal, err.Error())
	}

	return &user.UserResponse{
		User: h.modelToProto(userModel),
		Status: &common.StatusResponse{
			Code:    int32(codes.OK),
			Message: "User updated successfully",
			Success: true,
		},
	}, nil
}

func (h *UserGRPCHandler) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*common.StatusResponse, error) {
	err := h.userService.DeleteUser(ctx, req.Id)
	if err != nil {
		return &common.StatusResponse{
			Code:    int32(codes.Internal),
			Message: err.Error(),
			Success: false,
		}, status.Error(codes.Internal, err.Error())
	}

	return &common.StatusResponse{
		Code:    int32(codes.OK),
		Message: "User deleted successfully",
		Success: true,
	}, nil
}

func (h *UserGRPCHandler) ListUsers(ctx context.Context, req *user.ListUsersRequest) (*user.ListUsersResponse, error) {
	page := int(req.Page)
	limit := int(req.Limit)
	search := req.Search

	users, total, err := h.userService.ListUsers(ctx, page, limit, search)
	if err != nil {
		return &user.ListUsersResponse{
			Status: &common.StatusResponse{
				Code:    int32(codes.Internal),
				Message: err.Error(),
				Success: false,
			},
		}, status.Error(codes.Internal, err.Error())
	}

	protoUsers := make([]*user.User, len(users))
	for i, u := range users {
		protoUsers[i] = h.modelToProto(u)
	}

	return &user.ListUsersResponse{
		Users: protoUsers,
		Total: int32(total),
		Status: &common.StatusResponse{
			Code:    int32(codes.OK),
			Message: "Users retrieved successfully",
			Success: true,
		},
	}, nil
}

func (h *UserGRPCHandler) modelToProto(u *model.User) *user.User {
	return &user.User{
		Id:        u.ID.Hex(),
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		CreatedAt: u.CreatedAt.Unix(),
		UpdatedAt: u.UpdatedAt.Unix(),
	}
}
