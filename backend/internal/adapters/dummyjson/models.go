package dummyjson

type ProductDTO struct {
	ID        int     `json:"id"`
	SKU       string  `json:"sku"`
	Name      string  `json:"title"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
	UpdatedAt string  `json:"updatedAt"`
	Category  string  `json:"category"`
	Rating    float64 `json:"rating"`
	Brand     string  `json:"brand"`
}

type ProviderProduct = ProductDTO
