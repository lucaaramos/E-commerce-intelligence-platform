package dummyjson

import (
	"context"
	"fmt"
	"log/slog"

	"ecommerce-intelligence-platform/backend/internal/product"
)

type DummyJSONProvider struct {
	client *Client
	logger *slog.Logger
}

func NewDummyJSONProvider(client *Client, logger *slog.Logger) *DummyJSONProvider {
	if client == nil {
		client = NewClient()
	}
	return &DummyJSONProvider{
		client: client,
		logger: logger,
	}
}

func (p *DummyJSONProvider) ProviderID() string {
	return "dummyjson"
}

func (p *DummyJSONProvider) HealthCheck(ctx context.Context) error {
	_, _, err := p.client.ListProducts(ctx, 1, 0)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	p.logger.Info("dummyjson health check passed")
	return nil
}

func (p *DummyJSONProvider) FetchProducts(ctx context.Context) ([]product.ProductDTO, error) {
	products, total, err := p.client.ListProducts(ctx, 0, 0)
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

func (p *DummyJSONProvider) FetchProductBySKU(ctx context.Context, sku string) (*product.ProductDTO, error) {
	products, err := p.FetchProducts(ctx)
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