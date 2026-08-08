package httpx

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimit implements a simple in-memory fixed-window limiter per client IP.
// Suitable for a single-instance internal platform; replace with a shared
// store if the app is ever scaled horizontally.
func RateLimit(requestsPerMinute int) gin.HandlerFunc {
	if requestsPerMinute <= 0 {
		requestsPerMinute = 120
	}
	limiter := newFixedWindowLimiter(requestsPerMinute, time.Minute)
	return func(c *gin.Context) {
		if !limiter.Allow(c.ClientIP()) {
			RespondError(c, http.StatusTooManyRequests, "RATE_LIMITED", "too many requests, slow down")
			return
		}
		c.Next()
	}
}

type fixedWindowLimiter struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	windows map[string]*windowState
}

type windowState struct {
	count int
	reset time.Time
}

func newFixedWindowLimiter(limit int, window time.Duration) *fixedWindowLimiter {
	return &fixedWindowLimiter{
		limit:   limit,
		window:  window,
		windows: make(map[string]*windowState),
	}
}

func (l *fixedWindowLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	state, ok := l.windows[key]
	if !ok || now.After(state.reset) {
		l.windows[key] = &windowState{count: 1, reset: now.Add(l.window)}
		return true
	}
	if state.count >= l.limit {
		return false
	}
	state.count++
	return true
}
