package product

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"go-microservice-boilerplate/internal/config"
	"go-microservice-boilerplate/internal/database"
	"go-microservice-boilerplate/internal/proto/product"
	"go-microservice-boilerplate/internal/services/product/handler"
	"go-microservice-boilerplate/internal/services/product/repository"
	"go-microservice-boilerplate/internal/services/product/service"
	"go-microservice-boilerplate/internal/utils/logger"
)

type Server struct {
	config         *config.Config
	grpcServer     *grpc.Server
	productService service.ProductService
}

func NewServer(cfg *config.Config, mongodb *database.MongoDB, redis *database.Redis) *Server {
	// Initialize repositories
	productRepo := repository.NewMongoProductRepository(mongodb)
	productCache := repository.NewRedisProductCache(redis)

	// Initialize service
	productService := service.NewProductService(productRepo, productCache)

	// Initialize gRPC server
	grpcServer := grpc.NewServer()

	// Register handlers
	productHandler := handler.NewProductGRPCHandler(productService)
	product.RegisterProductServiceServer(grpcServer, productHandler)

	// Enable reflection for grpcurl/grpc clients
	reflection.Register(grpcServer)

	return &Server{
		config:         cfg,
		grpcServer:     grpcServer,
		productService: productService,
	}
}

func (s *Server) Start() error {
	port := s.config.Services.Product.Port
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	logger.Infof("Product service starting on port %s", port)

	if err := s.grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	logger.Info("Shutting down Product service...")
	s.grpcServer.GracefulStop()
}
