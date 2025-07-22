package service

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"go-microservice-boilerplate/internal/services/user/model"
	"go-microservice-boilerplate/internal/services/user/repository"
)

type userService struct {
	repo  repository.UserRepository
	cache repository.UserCache
}

func NewUserService(repo repository.UserRepository, cache repository.UserCache) UserService {
	return &userService{
		repo:  repo,
		cache: cache,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	// Check if user already exists
	existingUser, err := s.repo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Cache the user
	cacheKey := fmt.Sprintf("user:%s", user.ID.Hex())
	s.cache.Set(ctx, cacheKey, user, 3600) // 1 hour

	return user, nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*model.User, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("user:%s", id)
	if user, err := s.cache.Get(ctx, cacheKey); err == nil && user != nil {
		return user, nil
	}

	// Get from database
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Cache the result
	s.cache.Set(ctx, cacheKey, user, 3600)

	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, id string, req *model.UpdateUserRequest) (*model.User, error) {
	// Get existing user
	user, err := s.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := s.repo.Update(ctx, id, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Update cache
	cacheKey := fmt.Sprintf("user:%s", id)
	s.cache.Set(ctx, cacheKey, user, 3600)

	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Remove from cache
	cacheKey := fmt.Sprintf("user:%s", id)
	s.cache.Delete(ctx, cacheKey)

	return nil
}

func (s *userService) ListUsers(ctx context.Context, page, limit int, search string) ([]*model.User, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	users, total, err := s.repo.List(ctx, page, limit, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}
