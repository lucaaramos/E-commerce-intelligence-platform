package dummyjson

import (
	"context"
	"errors"
	"io"
	"testing"

	"ecommerce-intelligence-platform/backend/internal/product"
	"log/slog"
)

// mockClient implements Client for testing purposes (no actual HTTP requests).
type mockClient struct {
	products      []ProductDTO
	productBySKU  map[string]ProductDTO
	shouldFail    bool
	callCount     int
}

func newMockClient(products ...ProductDTO) *mockClient {
	m := &mockClient{
		products:      products,
		productBySKU:  make(map[string]ProductDTO),
		shouldFail:    false,
		callCount:     0,
	}
	for _, p := range products {
		m.productBySKU[p.SKU] = p
	}
	return m
}

func (m *mockClient) ListProducts(ctx context.Context, limit, skip int) ([]ProductDTO, int, error) {
	m.callCount++
	// Return all products regardless of limit/skip for simplicity in tests
	result := make([]ProductDTO, len(m.products))
	for i, p := range m.products {
		result[i] = p
	}
	total := len(m.products)
	if m.shouldFail && m.callCount >= 2 {
		return nil, 0, errors.New("simulated failure")
	}
	return result, total, nil
}

func (m *mockClient) GetProductByID(ctx context.Context, id int) (*ProductDTO, error) {
	m.callCount++
	for _, p := range m.products {
		if p.ID == id {
			p := p
			return &p, nil
		}
	}
	return nil, product.ErrNotFound
}

func (m *mockClient) SearchProducts(ctx context.Context, query string, limit, skip int) ([]ProductDTO, int, error) {
	m.callCount++
	// Simple search: return products whose title or description contains query
	var result []ProductDTO
	for _, p := range m.products {
		if len(result) >= limit && limit > 0 {
			break
		}
		result = append(result, p)
	}
	total := len(m.products)
	return result, total, nil
}

func (m *mockClient) HealthCheck(ctx context.Context) error {
	if m.shouldFail {
		return errors.New("simulated health check failure")
	}
	_, _, err := m.ListProducts(ctx, 1, 0)
	return err
}

type mockProvider struct {
	client *mockClient
	logger *slog.Logger
}

func newMockProvider(client *mockClient) *mockProvider {
	if client == nil {
		client = newMockClient()
	}
	return &mockProvider{
		client: client,
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

func (m *mockProvider) ProviderID() string {
	return "dummyjson"
}

func (m *mockProvider) HealthCheck(ctx context.Context) error {
	return m.client.HealthCheck(ctx)
}

func (m *mockProvider) FetchProducts(ctx context.Context) ([]product.ProductDTO, error) {
	products, total, err := m.client.ListProducts(ctx, 0, 0)
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return nil, nil
	}

	result := make([]product.ProductDTO, len(products))
	for i, p := range products {
		result[i] = product.ProductDTO{
			ID:        p.ID,
			SKU:       p.SKU,
			Name:      p.Name,
			Price:     p.Price,
			Stock:     p.Stock,
			UpdatedAt: p.UpdatedAt,
			Category:  p.Category,
			Rating:    p.Rating,
			Brand:     p.Brand,
		}
	}
	return result, nil
}

func (m *mockProvider) FetchProductBySKU(ctx context.Context, sku string) (*product.ProductDTO, error) {
	products, err := m.FetchProducts(ctx)
	if err != nil {
		return nil, err
	}

	for _, p := range products {
		if p.SKU == sku {
			return &p, nil
		}
	}

	return nil, product.ErrNotFound
}

func TestDummyJSONProviderFetchProducts(t *testing.T) {
	mockClient := newMockClient(
		ProductDTO{ID: 1, SKU: "BEA-ESS-ESS-001", Name: "Essence Mascara", Price: 10.00, Stock: 50, Category: "beauty", Rating: 4.5, Brand: "Essence"},
		ProductDTO{ID: 2, SKU: "FRA-CAL-CAL-006", Name: "Eau de Parfum", Price: 15.00, Stock: 30, Category: "fragrance", Rating: 4.0, Brand: "Generic"},
	)

	provider := newMockProvider(mockClient)
	products, err := provider.FetchProducts(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(products))
	}
	if products[0].SKU != "BEA-ESS-ESS-001" {
		t.Fatalf("expected SKU BEA-ESS-ESS-001, got %s", products[0].SKU)
	}
	if products[0].Price != 10.00 {
		t.Fatalf("expected price 10.00, got %f", products[0].Price)
	}
	if products[1].Name != "Eau de Parfum" {
		t.Fatalf("expected name Eau de Parfum, got %s", products[1].Name)
	}
}

func TestDummyJSONProviderFetchProductBySKU(t *testing.T) {
	mockClient := newMockClient(
		ProductDTO{ID: 1, SKU: "BEA-ESS-ESS-001", Name: "Essence Mascara", Price: 10.00, Stock: 50, Category: "beauty", Rating: 4.5, Brand: "Essence"},
	)

	provider := newMockProvider(mockClient)
	product, err := provider.FetchProductBySKU(context.Background(), "BEA-ESS-ESS-001")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if product == nil {
		t.Fatal("expected non-nil product")
	}
	if product.SKU != "BEA-ESS-ESS-001" {
		t.Fatalf("expected SKU BEA-ESS-ESS-001, got %s", product.SKU)
	}
	if product.Price != 10.00 {
		t.Fatalf("expected price 10.00, got %f", product.Price)
	}
}

func TestDummyJSONProviderFetchProductBySKUNotFound(t *testing.T) {
	mockClient := newMockClient(
		ProductDTO{ID: 1, SKU: "BEA-ESS-ESS-001", Name: "Essence Mascara", Price: 10.00, Stock: 50, Category: "beauty", Rating: 4.5, Brand: "Essence"},
	)

	provider := newMockProvider(mockClient)
	_, err := provider.FetchProductBySKU(context.Background(), "NONEXISTENT")
	if err == nil {
		t.Fatal("expected error for non-existent SKU")
	}
	if !errors.Is(err, product.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got: %v", err)
	}
}

func TestDummyJSONProviderHealthCheck(t *testing.T) {
	mockClient := newMockClient()
	provider := newMockProvider(mockClient)
	err := provider.HealthCheck(context.Background())
	if err != nil {
		t.Fatalf("expected no error for healthy provider, got: %v", err)
	}
}

func TestDummyJSONProviderHealthCheckFailed(t *testing.T) {
	mockClient := &mockClient{shouldFail: true}
	provider := newMockProvider(mockClient)
	err := provider.HealthCheck(context.Background())
	if err == nil {
		t.Fatal("expected error for failed health check")
	}
}

func TestDummyJSONProviderProviderID(t *testing.T) {
	mockClient := newMockClient()
	provider := newMockProvider(mockClient)
	if provider.ProviderID() != "dummyjson" {
		t.Fatalf("expected ProviderID 'dummyjson', got '%s'", provider.ProviderID())
	}
}