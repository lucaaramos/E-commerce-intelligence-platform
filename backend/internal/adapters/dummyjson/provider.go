package dummyjson

import (
	"context"
	"fmt"
	"log/slog"

	"ecommerce-intelligence-platform/backend/internal/product"
)

// type ProviderProduct struct {
// 	ID        int     `json:"id"`
// 	SKU       string  `json:"sku"`
// 	Name      string  `json:"title"`
// 	Price     float64 `json:"price"`
// 	Stock     int     `json:"stock"`
// 	UpdatedAt string  `json:"updatedAt"`
// 	Category  string  `json:"category"`
// 	Rating    float64 `json:"rating"`
// 	Brand     string  `json:"brand"`
// }

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

func (p *DummyJSONProvider) FetchProducts(ctx context.Context) ([]product.Product, error) {
	if p == nil || p.client == nil {
		return nil, fmt.Errorf("dummyjson client is nil")
	}

	products, _, err := p.client.ListProducts(ctx, 0, 0)
	if err != nil {
		return nil, err
	}

	mapped := make([]product.Product, 0, len(products))
	for _, item := range products {
		mapped = append(mapped, product.Product{
			SKU:          item.SKU,
			Name:         item.Name,
			SellingPrice: item.Price,
		})
	}

	return mapped, nil
}

func (p *DummyJSONProvider) FetchProductBySKU(ctx context.Context, sku string) (*product.Product, error) {
	if p == nil || p.client == nil {
		return nil, fmt.Errorf("dummyjson client is nil")
	}

	products, _, err := p.client.ListProducts(ctx, 0, 0)
	if err != nil {
		return nil, err
	}

	for _, item := range products {
		if item.SKU == sku {
			mapped := &product.Product{
				SKU:          item.SKU,
				Name:         item.Name,
				SellingPrice: item.Price,
			}
			return mapped, nil
		}
	}

	return nil, fmt.Errorf("product with sku %q not found", sku)
}

func (p *DummyJSONProvider) ProviderID() string {
	return "dummyjson"
}
