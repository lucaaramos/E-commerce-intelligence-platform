package product

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Product struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"id"`
	SKU          string        `bson:"sku" json:"sku" binding:"required"`
	Name         string        `bson:"name" json:"name" binding:"required"`
	Description  string        `bson:"description" json:"description"`
	Status       Status        `bson:"status" json:"status"`
	SellingPrice float64       `bson:"sellingPrice" json:"sellingPrice"`
	RealCost     float64       `bson:"realCost" json:"realCost"`
	Margin       float64       `bson:"margin" json:"margin"`
	Variants     []Variant     `bson:"variants" json:"variants"`
	CreatedAt    time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time     `bson:"updatedAt" json:"updatedAt"`
	DeletedAt    *time.Time    `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
	Price        float64       `bson:"price" json:"price"`
}

type Variant struct {
	SKU      string  `bson:"sku" json:"sku"`
	Name     string  `bson:"name" json:"name"`
	Cost     float64 `bson:"cost" json:"cost"`
	Stock    int     `bson:"stock" json:"stock"`
	Position int     `bson:"position" json:"position"`
}

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusArchived Status = "archived"
)

type CreateProductInput struct {
	SKU          string    `json:"sku" binding:"required"`
	Name         string    `json:"name" binding:"required"`
	Description  string    `json:"description"`
	SellingPrice float64   `json:"sellingPrice" binding:"required,gt=0"`
	Variants     []Variant `json:"variants"`
}

type UpdateProductInput struct {
	Name         *string   `json:"name"`
	Description  *string   `json:"description"`
	Status       *Status   `json:"status"`
	SellingPrice *float64  `json:"sellingPrice"`
	Variants     []Variant `json:"variants"`
}

type ProductFilter struct {
	SKU    string
	Status *Status
	Limit  int
	Offset int
}

func (p *Product) CalculateDerivedFields() {
	var totalCost float64
	for _, v := range p.Variants {
		totalCost += v.Cost
	}
	p.RealCost = totalCost
	if p.SellingPrice > 0 && p.RealCost > 0 {
		p.Margin = ((p.SellingPrice - p.RealCost) / p.SellingPrice) * 100
	}
}

func (p *Product) ToDTO() *ProductDTO {
	p.CalculateDerivedFields()
	return &ProductDTO{
		ID:           p.ID.Hex(),
		SKU:          p.SKU,
		Name:         p.Name,
		Description:  p.Description,
		SellingPrice: p.SellingPrice,
		RealCost:     p.RealCost,
		Margin:       p.Margin,
		Status:       string(p.Status),
	}
}
