package rabbitmq

import (
	"context"
	"testing"
)

func TestConnectInvalidURIReturnsError(t *testing.T) {
	_, err := Connect(context.Background(), "://not-a-valid-uri")
	if err == nil {
		t.Fatal("expected error for invalid URI, got nil")
	}
}

func TestConnectUnreachableReturnsError(t *testing.T) {
	ctx := context.Background()
	_, err := Connect(ctx, "amqp://guest:guest@127.0.0.1:1/")
	if err == nil {
		t.Fatal("expected error for unreachable broker, got nil")
	}
}

func TestClosedConnectionIsClosed(t *testing.T) {
	c := &Client{conn: nil}
	if !c.IsClosed() {
		t.Fatal("expected nil connection to report closed")
	}
	if err := c.Close(); err != nil {
		t.Fatalf("expected Close on nil connection to be a no-op, got %v", err)
	}
	if err := c.Check(context.Background()); err == nil {
		t.Fatal("expected Check on closed client to fail")
	}
}
