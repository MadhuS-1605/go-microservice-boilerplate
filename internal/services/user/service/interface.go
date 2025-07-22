package service

import (
	"context"
	"go-microservice-boilerplate/internal/services/user/model"
)

type UserService interface {
	CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
	GetUser(ctx context.Context, id string) (*model.User, error)
	UpdateUser(ctx context.Context, id string, req *model.UpdateUserRequest) (*model.User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, page, limit int, search string) ([]*model.User, int64, error)
}
