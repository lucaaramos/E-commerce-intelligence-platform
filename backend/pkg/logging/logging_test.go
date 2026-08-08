package logging

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"ecommerce-intelligence-platform/backend/pkg/config"
)

func TestNewJSONFormat(t *testing.T) {
	logger := New(config.LogConfig{Level: "info", Format: "json"})
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestNewTextFormat(t *testing.T) {
	logger := New(config.LogConfig{Level: "debug", Format: "text"})
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestParseLevel(t *testing.T) {
	cases := []struct {
		in   string
		want slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"DEBUG", slog.LevelDebug},
		{"warn", slog.LevelWarn},
		{"error", slog.LevelError},
		{"", slog.LevelInfo},
		{"bogus", slog.LevelInfo},
		{"info", slog.LevelInfo},
	}
	for _, c := range cases {
		if got := parseLevel(c.in); got != c.want {
			t.Errorf("parseLevel(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestLoggingProducesOutput(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, nil)
	logger := slog.New(handler)
	logger.Info("hello", "key", "value")

	out := buf.String()
	if !strings.Contains(out, "hello") {
		t.Errorf("expected message in output, got %q", out)
	}
	if !strings.Contains(out, "value") {
		t.Errorf("expected attribute value in output, got %q", out)
	}
}
