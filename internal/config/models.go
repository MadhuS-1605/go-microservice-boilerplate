package config

type Config struct {
	MongoDB   MongoDBConfig
	Redis     RedisConfig
	Services  ServicesConfig
	Swagger   SwaggerConfig
	LogLevel  string
	JWTSecret string
}

type MongoDBConfig struct {
	URI      string
	Database string
	Timeout  int
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type ServicesConfig struct {
	Gateway GatewayConfig
	User    ServiceConfig
	Product ServiceConfig
}

type GatewayConfig struct {
	Port string
}

type ServiceConfig struct {
	Port string
	Host string
}

type SwaggerConfig struct {
	Enabled     bool        `yaml:"enabled"`
	Auth        SwaggerAuth `yaml:"auth"`
	Title       string      `yaml:"title"`
	Version     string      `yaml:"version"`
	Description string      `yaml:"description"`
}

type SwaggerAuth struct {
	Enabled  bool   `yaml:"enabled"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
