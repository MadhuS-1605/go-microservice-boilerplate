package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go-microservice-boilerplate/internal/config"
	"go-microservice-boilerplate/internal/proto/common"
	"go-microservice-boilerplate/internal/proto/user"
)

type UserClient struct {
	conn   *grpc.ClientConn
	client user.UserServiceClient
}

func NewUserClient(cfg *config.Config) (*UserClient, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Services.User.Host, cfg.Services.User.Port)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	client := user.NewUserServiceClient(conn)

	return &UserClient{
		conn:   conn,
		client: client,
	}, nil
}

func (c *UserClient) Close() error {
	return c.conn.Close()
}

func (c *UserClient) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.UserResponse, error) {
	return c.client.CreateUser(ctx, req)
}

func (c *UserClient) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.UserResponse, error) {
	return c.client.GetUser(ctx, req)
}

func (c *UserClient) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UserResponse, error) {
	return c.client.UpdateUser(ctx, req)
}

func (c *UserClient) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*common.StatusResponse, error) {
	return c.client.DeleteUser(ctx, req)
}

func (c *UserClient) ListUsers(ctx context.Context, req *user.ListUsersRequest) (*user.ListUsersResponse, error) {
	return c.client.ListUsers(ctx, req)
}
