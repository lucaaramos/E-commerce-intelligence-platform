package product

import "context"

type ProductProvider interface {
	FetchProducts(ctx context.Context) ([]Product, error)
	FetchProductBySKU(ctx context.Context, sku string) (*Product, error)
	ProviderID() string
}