package user

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"go-microservice-boilerplate/internal/config"
	"go-microservice-boilerplate/internal/database"
	"go-microservice-boilerplate/internal/proto/user"
	"go-microservice-boilerplate/internal/services/user/handler"
	"go-microservice-boilerplate/internal/services/user/repository"
	"go-microservice-boilerplate/internal/services/user/service"
	"go-microservice-boilerplate/internal/utils/logger"
)

type Server struct {
	config      *config.Config
	grpcServer  *grpc.Server
	userService service.UserService
}

func NewServer(cfg *config.Config, mongodb *database.MongoDB, redis *database.Redis) *Server {
	// Initialize repositories
	userRepo := repository.NewMongoUserRepository(mongodb)
	userCache := repository.NewRedisUserCache(redis)

	// Initialize service
	userService := service.NewUserService(userRepo, userCache)

	// Initialize gRPC server
	grpcServer := grpc.NewServer()

	// Register handlers
	userHandler := handler.NewUserGRPCHandler(userService)
	user.RegisterUserServiceServer(grpcServer, userHandler)

	// Enable reflection for grpcurl/grpc clients
	reflection.Register(grpcServer)

	return &Server{
		config:      cfg,
		grpcServer:  grpcServer,
		userService: userService,
	}
}

func (s *Server) Start() error {
	port := s.config.Services.User.Port
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	logger.Infof("User service starting on port %s", port)

	if err := s.grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	logger.Info("Shutting down User service...")
	s.grpcServer.GracefulStop()
}
