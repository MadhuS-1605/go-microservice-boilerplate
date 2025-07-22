package gateway

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "go-microservice-boilerplate/docs"

	"go-microservice-boilerplate/internal/config"
	"go-microservice-boilerplate/internal/middleware"
	"go-microservice-boilerplate/internal/services/gateway/client"
	"go-microservice-boilerplate/internal/services/gateway/handler"
	"go-microservice-boilerplate/internal/utils/logger"
)

type Server struct {
	config        *config.Config
	httpServer    *http.Server
	userClient    *client.UserClient
	productClient *client.ProductClient
}

func NewServer(cfg *config.Config) *Server {
	// Initialize gRPC clients
	userClient, err := client.NewUserClient(cfg)
	if err != nil {
		logger.Fatalf("Failed to create user client: %v", err)
	}

	productClient, err := client.NewProductClient(cfg)
	if err != nil {
		logger.Fatalf("Failed to create product client: %v", err)
	}

	// Initialize Gin router
	router := gin.New()

	// Add middleware
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	// Swagger documentation with authentication
	if cfg.Swagger.Enabled {
		swaggerGroup := router.Group("/swagger")

		// Apply authentication if enabled
		if cfg.Swagger.Auth.Enabled {
			swaggerGroup.Use(middleware.SwaggerAuth(cfg.Swagger.Auth.Username, cfg.Swagger.Auth.Password))
			logger.Infof("Swagger authentication enabled. Username: %s", cfg.Swagger.Auth.Username)
		}

		swaggerGroup.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Initialize handlers
	gatewayHandler := handler.NewGatewayHandler(userClient, productClient)
	gatewayHandler.RegisterRoutes(router)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    ":" + cfg.Services.Gateway.Port,
		Handler: router,
	}

	return &Server{
		config:        cfg,
		httpServer:    httpServer,
		userClient:    userClient,
		productClient: productClient,
	}
}

func (s *Server) Start() error {
	logger.Infof("Gateway server starting on port %s", s.config.Services.Gateway.Port)
	logger.Infof("Swagger documentation available at http://localhost:%s/swagger/index.html", s.config.Services.Gateway.Port)

	// Start server in a goroutine
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start gateway server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Gateway server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("gateway server forced to shutdown: %w", err)
	}

	// Close gRPC clients
	s.userClient.Close()
	s.productClient.Close()

	logger.Info("Gateway server exited")
	return nil
}
