package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-microservice-boilerplate/internal/database"
	"go-microservice-boilerplate/internal/services/user/model"
)

type redisUserCache struct {
	client *database.Redis
}

func NewRedisUserCache(redis *database.Redis) UserCache {
	return &redisUserCache{
		client: redis,
	}
}

func (c *redisUserCache) Set(ctx context.Context, key string, user *model.User, expiration int) error {
	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	err = c.client.Client.Set(ctx, key, data, time.Duration(expiration)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

func (c *redisUserCache) Get(ctx context.Context, key string) (*model.User, error) {
	data, err := c.client.Client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}

func (c *redisUserCache) Delete(ctx context.Context, key string) error {
	return c.client.Client.Del(ctx, key).Err()
}
