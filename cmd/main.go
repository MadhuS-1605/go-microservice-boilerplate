package main

import (
	"fmt"
	"log"
	"os"

	"go-microservice-boilerplate/internal/config"
	"go-microservice-boilerplate/internal/database"
	"go-microservice-boilerplate/internal/services/gateway"
	"go-microservice-boilerplate/internal/services/product"
	"go-microservice-boilerplate/internal/services/user"
	"go-microservice-boilerplate/internal/utils/logger"
)

// @title Go Microservice Boilerplate API
// @version 1.0
// @description Microservice Boilerplate API for user and product management
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Parse command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/main.go <service>")
		fmt.Println("Available services: web, user, product")
		os.Exit(1)
	}

	service := os.Args[1]

	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize logger
	logger.Init(cfg.LogLevel)

	// Initialize databases
	mongodb, err := database.NewMongoDB(cfg.MongoDB)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer mongodb.Disconnect()

	redisClient, err := database.NewRedis(cfg.Redis)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redisClient.Close()

	// Run the specified service
	switch service {
	case "web", "gateway":
		runGateway(cfg)
	case "user":
		runUserService(cfg, *mongodb, *redisClient)
	case "product":
		runProductService(cfg, *mongodb, *redisClient)
	default:
		fmt.Printf("Unknown service: %s\n", service)
		fmt.Println("Available services: web, user, product")
		os.Exit(1)
	}
}

func runGateway(cfg *config.Config) {
	logger.Info("Starting Gateway Service...")

	gatewayServer := gateway.NewServer(cfg)
	if err := gatewayServer.Start(); err != nil {
		log.Fatal("Failed to start gateway service:", err)
	}
}

func runUserService(cfg *config.Config, mongodb database.MongoDB, redis database.Redis) {
	logger.Info("Starting User Service...")

	userServer := user.NewServer(cfg, &mongodb, &redis)
	if err := userServer.Start(); err != nil {
		log.Fatal("Failed to start user service:", err)
	}
}

func runProductService(cfg *config.Config, mongodb database.MongoDB, redis database.Redis) {
	logger.Info("Starting Product Service...")

	productServer := product.NewServer(cfg, &mongodb, &redis)
	if err := productServer.Start(); err != nil {
		log.Fatal("Failed to start product service:", err)
	}
}
