package config

import (
	"os"
	"strconv"
)

func Load() (*Config, error) {
	return &Config{
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGODB_DATABASE", "microservices_db"),
			Timeout:  getEnvInt("MONGODB_TIMEOUT", 30),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Services: ServicesConfig{
			Gateway: GatewayConfig{
				Port: getEnv("GATEWAY_PORT", "8080"),
			},
			User: ServiceConfig{
				Port: getEnv("USER_SERVICE_PORT", "50051"),
				Host: getEnv("USER_SERVICE_HOST", "localhost"),
			},
			Product: ServiceConfig{
				Port: getEnv("PRODUCT_SERVICE_PORT", "50052"),
				Host: getEnv("PRODUCT_SERVICE_HOST", "localhost"),
			},
		},
		Swagger: SwaggerConfig{ // Add swagger config
			Enabled: getBoolEnv("SWAGGER_ENABLED", true),
			Auth: SwaggerAuth{
				Enabled:  getBoolEnv("SWAGGER_AUTH_ENABLED", true),
				Username: getEnv("SWAGGER_USERNAME", "admin"),
				Password: getEnv("SWAGGER_PASSWORD", "boilerplate@123"),
			},
			Title:       getEnv("SWAGGER_TITLE", "Go Microservice API"),
			Version:     getEnv("SWAGGER_VERSION", "1.0.0"),
			Description: getEnv("SWAGGER_DESCRIPTION", "Microservice API for user and product management"),
		},
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		JWTSecret: getEnv("JWT_SECRET", "boilerplate@123"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
