package mongodb

import (
	"context"
	"testing"
)

func TestConnectInvalidURIReturnsError(t *testing.T) {
	_, err := Connect(context.Background(), "://not-a-valid-uri", "test")
	if err == nil {
		t.Fatal("expected error for invalid URI, got nil")
	}
}

func TestConnectEmptyURIReturnsError(t *testing.T) {
	_, err := Connect(context.Background(), "", "test")
	if err == nil {
		t.Fatal("expected error for empty URI, got nil")
	}
}
