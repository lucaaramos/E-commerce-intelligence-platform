package product

import "context"

type ProductProvider interface {
	FetchProducts(ctx context.Context) ([]ProductDTO, error)
	FetchProductBySKU(ctx context.Context, sku string) (*ProductDTO, error)
	ProviderID() string
}