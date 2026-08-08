package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Env     string
	Host    string
	Port    int
	MongoDB MongoDBConfig
	RabbitMQ RabbitMQConfig
	Auth    AuthConfig
	HTTP    HTTPConfig
	Log     LogConfig
}

type MongoDBConfig struct {
	URI string
	DB  string
}

type RabbitMQConfig struct {
	URI string
}

type AuthConfig struct {
	JWTSecret     string
	AdminEmail    string
	AdminPassword string
}

type HTTPConfig struct {
	CORSAllowedOrigins    []string
	RateLimitPerMinute    int
	ShutdownTimeout       time.Duration
	ReadHeaderTimeout     time.Duration
}

type LogConfig struct {
	Level  string
	Format string
}

func Load() Config {
	return Config{
		Env:  getEnv("APP_ENV", "development"),
		Host: getEnv("APP_HOST", "0.0.0.0"),
		Port: getEnvInt("APP_PORT", 8080),
		MongoDB: MongoDBConfig{
			URI: getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			DB:  getEnv("MONGODB_DB", "ecommerce_intel"),
		},
		RabbitMQ: RabbitMQConfig{
			URI: getEnv("RABBITMQ_URI", "amqp://guest:guest@localhost:5672/"),
		},
		Auth: AuthConfig{
			JWTSecret:     getEnv("JWT_SECRET", "dev-secret-change-me"),
			AdminEmail:    getEnv("ADMIN_EMAIL", "admin@example.com"),
			AdminPassword: getEnv("ADMIN_PASSWORD", "admin123"),
		},
		HTTP: HTTPConfig{
			CORSAllowedOrigins: getEnvSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:5173"}),
			RateLimitPerMinute: getEnvInt("RATE_LIMIT_PER_MINUTE", 120),
			ShutdownTimeout:    getEnvDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
			ReadHeaderTimeout:  getEnvDuration("READ_HEADER_TIMEOUT", 5*time.Second),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "text"),
		},
	}
}

func (c Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}

func getEnvSlice(key string, fallback []string) []string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	parts := strings.Split(v, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	if len(result) == 0 {
		return fallback
	}
	return result
}
