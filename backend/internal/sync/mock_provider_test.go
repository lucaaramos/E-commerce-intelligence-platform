package sync

import (
	"context"
	"errors"
	"testing"

	"ecommerce-intelligence-platform/backend/internal/product"
)

// mockProvider implements product.Provider for testing purposes.
type mockProvider struct {
	products       []product.ProductDTO
	productBySKU   map[string]product.ProductDTO
	shouldFail     bool
	failAfter      int
	callCount      int
}

func newMockProvider(products ...product.ProductDTO) *mockProvider {
	m := &mockProvider{
		products:      products,
		productBySKU:  make(map[string]product.ProductDTO),
		shouldFail:    false,
		failAfter:     999,
	}
	for _, p := range products {
		m.productBySKU[p.SKU] = p
	}
	return m
}

func (m *mockProvider) ProviderID() string {
	return "mock"
}

func (m *mockProvider) FetchProducts(ctx context.Context) ([]product.ProductDTO, error) {
	m.callCount++
	if m.shouldFail && m.callCount >= m.failAfter {
		return nil, errors.New("simulated failure after " + string(rune('0' + m.callCount)))
	}
	// Return a copy to prevent mutation
	result := make([]product.ProductDTO, len(m.products))
	copy(result, m.products)
	return result, nil
}

func (m *mockProvider) FetchProductBySKU(ctx context.Context, sku string) (*product.ProductDTO, error) {
	m.callCount++
	product, exists := m.productBySKU[sku]
	if !exists {
		return nil, errors.New("product not found")
	}
	// Return a pointer to a copy
	p := product
	return &p, nil
}

func TestWorkerFetchProducts(t *testing.T) {
	mock := newMockProvider(
		product.ProductDTO{ID: 1, SKU: "TEST-001", Name: "Test Product", Price: 10.00, Stock: 50},
		product.ProductDTO{ID: 2, SKU: "TEST-002", Name: "Test Product 2", Price: 20.00, Stock: 100},
	)

	products, err := mock.FetchProducts(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(products))
	}
	if products[0].SKU != "TEST-001" {
		t.Fatalf("expected first product SKU TEST-001, got %s", products[0].SKU)
	}
	if products[1].Price != 20.00 {
		t.Fatalf("expected second product price 20.00, got %f", products[1].Price)
	}
}

func TestWorkerFetchProductBySKU(t *testing.T) {
	mock := newMockProvider(
		product.ProductDTO{ID: 1, SKU: "TEST-001", Name: "Test Product", Price: 10.00, Stock: 50},
	)

	product, err := mock.FetchProductBySKU(context.Background(), "TEST-001")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if product == nil {
		t.Fatal("expected non-nil product")
	}
	if product.SKU != "TEST-001" {
		t.Fatalf("expected SKU TEST-001, got %s", product.SKU)
	}

	// Test not found
	_, err = mock.FetchProductBySKU(context.Background(), "NONEXISTENT")
	if err == nil {
		t.Fatal("expected error for non-existent SKU")
	}
}

func TestWorkerProviderID(t *testing.T) {
	mock := newMockProvider()
	if mock.ProviderID() != "mock" {
		t.Fatalf("expected ProviderID 'mock', got '%s'", mock.ProviderID())
	}
}