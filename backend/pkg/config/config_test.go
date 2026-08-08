package config

import (
	"testing"
	"time"
)

func setEnv(t *testing.T, key, value string) {
	t.Helper()
	t.Setenv(key, value)
}

func clearEnv(t *testing.T, key string) {
	t.Helper()
	t.Setenv(key, "")
}

func TestLoadDefaults(t *testing.T) {
	clearEnv(t, "APP_ENV")
	clearEnv(t, "APP_PORT")
	clearEnv(t, "MONGODB_URI")
	clearEnv(t, "RATE_LIMIT_PER_MINUTE")

	cfg := Load()

	if cfg.Env != "development" {
		t.Errorf("expected Env development, got %q", cfg.Env)
	}
	if cfg.Port != 8080 {
		t.Errorf("expected Port 8080, got %d", cfg.Port)
	}
	if cfg.Host != "0.0.0.0" {
		t.Errorf("expected Host 0.0.0.0, got %q", cfg.Host)
	}
	if cfg.MongoDB.URI != "mongodb://localhost:27017" {
		t.Errorf("unexpected MongoDB URI: %q", cfg.MongoDB.URI)
	}
	if cfg.MongoDB.DB != "ecommerce_intel" {
		t.Errorf("unexpected MongoDB DB: %q", cfg.MongoDB.DB)
	}
	if cfg.HTTP.RateLimitPerMinute != 120 {
		t.Errorf("expected RateLimitPerMinute 120, got %d", cfg.HTTP.RateLimitPerMinute)
	}
	if cfg.HTTP.ShutdownTimeout != 10*time.Second {
		t.Errorf("expected ShutdownTimeout 10s, got %s", cfg.HTTP.ShutdownTimeout)
	}
	if cfg.Addr() != "0.0.0.0:8080" {
		t.Errorf("unexpected Addr: %q", cfg.Addr())
	}
}

func TestLoadFromEnv(t *testing.T) {
	setEnv(t, "APP_ENV", "production")
	setEnv(t, "APP_HOST", "127.0.0.1")
	setEnv(t, "APP_PORT", "9000")
	setEnv(t, "MONGODB_URI", "mongodb://mongo:27017")
	setEnv(t, "MONGODB_DB", "custom_db")
	setEnv(t, "RABBITMQ_URI", "amqp://user:pass@host:5672/")
	setEnv(t, "JWT_SECRET", "super-secret")
	setEnv(t, "ADMIN_EMAIL", "owner@shop.com")
	setEnv(t, "ADMIN_PASSWORD", "strongpass")
	setEnv(t, "RATE_LIMIT_PER_MINUTE", "5")
	setEnv(t, "SHUTDOWN_TIMEOUT", "30s")
	setEnv(t, "CORS_ALLOWED_ORIGINS", "http://a.com, http://b.com")
	setEnv(t, "LOG_LEVEL", "debug")
	setEnv(t, "LOG_FORMAT", "json")

	cfg := Load()

	if cfg.Env != "production" {
		t.Errorf("expected Env production, got %q", cfg.Env)
	}
	if cfg.Port != 9000 {
		t.Errorf("expected Port 9000, got %d", cfg.Port)
	}
	if cfg.Host != "127.0.0.1" {
		t.Errorf("expected Host 127.0.0.1, got %q", cfg.Host)
	}
	if cfg.MongoDB.URI != "mongodb://mongo:27017" {
		t.Errorf("unexpected MongoDB URI: %q", cfg.MongoDB.URI)
	}
	if cfg.MongoDB.DB != "custom_db" {
		t.Errorf("unexpected MongoDB DB: %q", cfg.MongoDB.DB)
	}
	if cfg.RabbitMQ.URI != "amqp://user:pass@host:5672/" {
		t.Errorf("unexpected RabbitMQ URI: %q", cfg.RabbitMQ.URI)
	}
	if cfg.Auth.JWTSecret != "super-secret" {
		t.Errorf("unexpected JWT secret: %q", cfg.Auth.JWTSecret)
	}
	if cfg.Auth.AdminEmail != "owner@shop.com" {
		t.Errorf("unexpected admin email: %q", cfg.Auth.AdminEmail)
	}
	if cfg.Auth.AdminPassword != "strongpass" {
		t.Errorf("unexpected admin password: %q", cfg.Auth.AdminPassword)
	}
	if cfg.HTTP.RateLimitPerMinute != 5 {
		t.Errorf("expected RateLimitPerMinute 5, got %d", cfg.HTTP.RateLimitPerMinute)
	}
	if cfg.HTTP.ShutdownTimeout != 30*time.Second {
		t.Errorf("expected ShutdownTimeout 30s, got %s", cfg.HTTP.ShutdownTimeout)
	}
	if len(cfg.HTTP.CORSAllowedOrigins) != 2 || cfg.HTTP.CORSAllowedOrigins[0] != "http://a.com" {
		t.Errorf("unexpected CORS origins: %v", cfg.HTTP.CORSAllowedOrigins)
	}
	if cfg.Log.Level != "debug" || cfg.Log.Format != "json" {
		t.Errorf("unexpected log config: %+v", cfg.Log)
	}
	if cfg.Addr() != "127.0.0.1:9000" {
		t.Errorf("unexpected Addr: %q", cfg.Addr())
	}
}

func TestLoadInvalidIntFallsBack(t *testing.T) {
	setEnv(t, "APP_PORT", "not-a-number")
	cfg := Load()
	if cfg.Port != 8080 {
		t.Errorf("expected fallback Port 8080, got %d", cfg.Port)
	}
}

func TestLoadInvalidDurationFallsBack(t *testing.T) {
	setEnv(t, "SHUTDOWN_TIMEOUT", "nope")
	cfg := Load()
	if cfg.HTTP.ShutdownTimeout != 10*time.Second {
		t.Errorf("expected fallback ShutdownTimeout, got %s", cfg.HTTP.ShutdownTimeout)
	}
}

func TestLoadEmptyCORSFallsBack(t *testing.T) {
	setEnv(t, "CORS_ALLOWED_ORIGINS", "")
	cfg := Load()
	if len(cfg.HTTP.CORSAllowedOrigins) != 1 || cfg.HTTP.CORSAllowedOrigins[0] != "http://localhost:5173" {
		t.Errorf("unexpected CORS fallback: %v", cfg.HTTP.CORSAllowedOrigins)
	}
}
