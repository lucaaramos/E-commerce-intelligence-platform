package product

type ProductDTO struct {
	ID           any     `json:"id"`
	SKU          string  `json:"sku"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	Stock        int     `json:"stock"`
	UpdatedAt    string  `json:"updatedAt"`
	Category     string  `json:"category"`
	Rating       float64 `json:"rating"`
	Brand        string  `json:"brand"`
	SellingPrice float64 `json:"sellingPrice"`
	RealCost     float64 `json:"realCost"`
	Margin       float64 `json:"margin"`
	Status       string  `json:"status"`
}
