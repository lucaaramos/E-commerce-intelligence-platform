package httpx

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"ecommerce-intelligence-platform/backend/pkg/config"
)

// New builds the gin engine with the application middleware chain, health
// endpoints and the versioned API group. Module routes are registered by
// calling RegisterRoutes on the returned engine.
func New(cfg *config.Config, logger *slog.Logger) *gin.Engine {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()
	engine.Use(Middlewares(cfg, logger)...)
	engine.GET("/healthz", Liveness())

	// Versioned API group; module routes register against this.
	engine.Group("/api/v1")

	return engine
}

// RegisterReadiness attaches the /readyz endpoint with the given dependency
// checks.
func RegisterReadiness(engine *gin.Engine, checks ...HealthChecker) {
	engine.GET("/readyz", Readiness(checks...))
}
