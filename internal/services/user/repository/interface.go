package repository

import (
	"context"
	"go-microservice-boilerplate/internal/services/user/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, id string, user *model.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int, search string) ([]*model.User, int64, error)
}

type UserCache interface {
	Set(ctx context.Context, key string, user *model.User, expiration int) error
	Get(ctx context.Context, key string) (*model.User, error)
	Delete(ctx context.Context, key string) error
}
