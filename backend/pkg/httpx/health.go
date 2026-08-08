package httpx

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthChecker reports whether a dependency is healthy.
type HealthChecker interface {
	Check(ctx context.Context) error
}

type healthResponse struct {
	Status  string                 `json:"status"`
	Checks  map[string]checkResult `json:"checks,omitempty"`
	RequestID string               `json:"request_id,omitempty"`
}

type checkResult struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// Liveness returns the liveness probe handler.
func Liveness() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// Readiness probes all registered dependencies and returns 200 only when all
// of them are reachable.
func Readiness(checks ...HealthChecker) gin.HandlerFunc {
	return func(c *gin.Context) {
		results := make(map[string]checkResult, len(checks))
		allOK := true
		for _, check := range checks {
			name := dependencyName(check)
			status := "ok"
			errMsg := ""
			if err := check.Check(c.Request.Context()); err != nil {
				status = "error"
				errMsg = err.Error()
				allOK = false
			}
			results[name] = checkResult{Status: status, Error: errMsg}
		}

		status := http.StatusOK
		body := healthResponse{Status: "ok", Checks: results, RequestID: GetRequestID(c)}
		if !allOK {
			status = http.StatusServiceUnavailable
			body.Status = "unavailable"
		}
		c.JSON(status, body)
	}
}

func dependencyName(check HealthChecker) string {
	if n, ok := check.(interface{ Name() string }); ok {
		return n.Name()
	}
	return "dependency"
}
