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
	return &DummyJSONProvider{
		client: client,
		logger: logger,
	}
}

func (p *DummyJSONProvider) FetchProducts(ctx context.Context) ([]product.ProductDTO, error) {
	if p == nil || p.client == nil {
		return nil, fmt.Errorf("dummyjson client is nil")
	}

	products, _, err := p.client.ListProducts(ctx, 0, 0)
	if err != nil {
		return nil, err
	}

	mapped := make([]product.ProductDTO, 0, len(products))
	for _, item := range products {
		mapped = append(mapped, product.ProductDTO{
			ID:        item.ID,
			SKU:       item.SKU,
			Name:      item.Name,
			Price:     item.Price,
			Stock:     item.Stock,
			UpdatedAt: item.UpdatedAt,
			Category:  item.Category,
			Rating:    item.Rating,
			Brand:     item.Brand,
		})
	}

	return mapped, nil
}

func (p *DummyJSONProvider) FetchProductBySKU(ctx context.Context, sku string) (*product.ProductDTO, error) {
	if p == nil || p.client == nil {
		return nil, fmt.Errorf("dummyjson client is nil")
	}

	products, _, err := p.client.ListProducts(ctx, 0, 0)
	if err != nil {
		return nil, err
	}

	for _, item := range products {
		if item.SKU == sku {
			mapped := &product.ProductDTO{
				ID:        item.ID,
				SKU:       item.SKU,
				Name:      item.Name,
				Price:     item.Price,
				Stock:     item.Stock,
				UpdatedAt: item.UpdatedAt,
				Category:  item.Category,
				Rating:    item.Rating,
				Brand:     item.Brand,
			}
			return mapped, nil
		}
	}

	return nil, fmt.Errorf("product with sku %q not found", sku)
}

func (p *DummyJSONProvider) ProviderID() string {
	return "dummyjson"
}
