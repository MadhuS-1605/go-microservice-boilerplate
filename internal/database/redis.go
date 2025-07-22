package database

import (
	"github.com/go-redis/redis/v8"
	"go-microservice-boilerplate/internal/config"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis(config config.RedisConfig) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		return nil, err
	}

	return &Redis{Client: client}, nil
}

func (r *Redis) Close() error {
	return r.Client.Close()
}
