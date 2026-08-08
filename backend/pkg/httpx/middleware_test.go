package httpx

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func newTestEngine(middlewares ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(middlewares...)
	return engine
}

func doRequest(t *testing.T, engine http.Handler, method, path string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	return rec
}

func TestRequestIDGeneratesAndEchoes(t *testing.T) {
	engine := newTestEngine(RequestID())
	engine.GET("/", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	rec := doRequest(t, engine, http.MethodGet, "/", nil)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	id := rec.Header().Get("X-Request-ID")
	if id == "" {
		t.Fatal("expected X-Request-ID header to be set")
	}
}

func TestRequestIDPreservesIncoming(t *testing.T) {
	engine := newTestEngine(RequestID())
	engine.GET("/", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	rec := doRequest(t, engine, http.MethodGet, "/", map[string]string{"X-Request-ID": "client-id-1"})
	if got := rec.Header().Get("X-Request-ID"); got != "client-id-1" {
		t.Fatalf("expected incoming request ID preserved, got %q", got)
	}
}

func TestRecoveryReturnsJSONError(t *testing.T) {
	logger := newTestLogger()
	engine := newTestEngine(RequestID(), Recovery(logger))
	engine.GET("/panic", func(c *gin.Context) { panic("boom") })

	rec := doRequest(t, engine, http.MethodGet, "/panic", nil)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "INTERNAL_ERROR") {
		t.Fatalf("expected error code in body, got %s", rec.Body.String())
	}
}

func TestCORSSetsHeadersForAllowedOrigin(t *testing.T) {
	engine := newTestEngine(CORS([]string{"http://allowed.com"}))
	engine.GET("/", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	rec := doRequest(t, engine, http.MethodGet, "/", map[string]string{"Origin": "http://allowed.com"})
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://allowed.com" {
		t.Fatalf("expected allow-origin header, got %q", got)
	}
}

func TestCORSRejectsDisallowedOrigin(t *testing.T) {
	engine := newTestEngine(CORS([]string{"http://allowed.com"}))
	engine.GET("/", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	rec := doRequest(t, engine, http.MethodGet, "/", map[string]string{"Origin": "http://evil.com"})
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no allow-origin header for disallowed origin, got %q", got)
	}
}

func TestCORSPreflight(t *testing.T) {
	engine := newTestEngine(CORS([]string{"http://allowed.com"}))
	rec := doRequest(t, engine, http.MethodOptions, "/", map[string]string{
		"Origin":                        "http://allowed.com",
		"Access-Control-Request-Method": "POST",
	})
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 preflight, got %d", rec.Code)
	}
}

func TestRateLimitRejectsOverLimit(t *testing.T) {
	engine := newTestEngine(RateLimit(2))
	engine.GET("/", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	headers := map[string]string{"X-Forwarded-For": "1.2.3.4"}
	for i := 0; i < 2; i++ {
		rec := doRequest(t, engine, http.MethodGet, "/", headers)
		if rec.Code != http.StatusNoContent {
			t.Fatalf("request %d: expected 204, got %d", i, rec.Code)
		}
	}
	rec := doRequest(t, engine, http.MethodGet, "/", headers)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "RATE_LIMITED") {
		t.Fatalf("expected RATE_LIMITED code, got %s", rec.Body.String())
	}
}

func TestRespondErrorIncludesRequestID(t *testing.T) {
	engine := newTestEngine(RequestID())
	engine.GET("/", func(c *gin.Context) { RespondError(c, http.StatusNotFound, "NOT_FOUND", "nope") })

	rec := doRequest(t, engine, http.MethodGet, "/", nil)
	body := rec.Body.String()
	if !strings.Contains(body, "NOT_FOUND") || !strings.Contains(body, "request_id") {
		t.Fatalf("unexpected error body: %s", body)
	}
}

func TestRespondOKEnvelope(t *testing.T) {
	engine := newTestEngine()
	engine.GET("/", func(c *gin.Context) { RespondOK(c, http.StatusOK, map[string]int{"n": 1}) })

	rec := doRequest(t, engine, http.MethodGet, "/", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"data"`) {
		t.Fatalf("expected data envelope, got %s", rec.Body.String())
	}
}

func TestLiveness(t *testing.T) {
	engine := newTestEngine()
	engine.GET("/healthz", Liveness())
	rec := doRequest(t, engine, http.MethodGet, "/healthz", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

type fakeChecker struct {
	name string
	err  error
}

func (f fakeChecker) Name() string              { return f.name }
func (f fakeChecker) Check(ctx context.Context) error { return f.err }

func TestReadinessAllOK(t *testing.T) {
	engine := newTestEngine()
	engine.GET("/readyz", Readiness(fakeChecker{name: "mongo"}, fakeChecker{name: "rabbit"}))

	rec := doRequest(t, engine, http.MethodGet, "/readyz", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestReadinessFailure(t *testing.T) {
	engine := newTestEngine()
	engine.GET("/readyz", Readiness(
		fakeChecker{name: "mongo"},
		fakeChecker{name: "rabbit", err: io.ErrUnexpectedEOF},
	))

	rec := doRequest(t, engine, http.MethodGet, "/readyz", nil)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "unavailable") {
		t.Fatalf("expected unavailable status, got %s", rec.Body.String())
	}
}
