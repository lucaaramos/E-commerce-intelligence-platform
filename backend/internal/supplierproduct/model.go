package supplierproduct

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type SupplierProduct struct {
	ID              bson.ObjectID `bson:"_id,omitempty" json:"id"`
	ProductID       bson.ObjectID `bson:"productId" json:"productId" binding:"required"`
	SupplierID      bson.ObjectID `bson:"supplierId" json:"supplierId" binding:"required"`
	SupplierSKU     string        `bson:"supplierSku" json:"supplierSku"`
	Cost            float64       `bson:"cost" json:"cost" binding:"required,gt=0"`
	ShippingCost    float64       `bson:"shippingCost" json:"shippingCost"`
	LeadTimeDays    int           `bson:"leadTimeDays" json:"leadTimeDays"`
	MinOrderQty     int           `bson:"minOrderQty" json:"minOrderQty"`
	URL             string        `bson:"url" json:"url"`
	IsActive        bool          `bson:"isActive" json:"isActive"`
	LastSyncedAt    *time.Time    `bson:"lastSyncedAt,omitempty" json:"lastSyncedAt,omitempty"`
	CreatedAt       time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time     `bson:"updatedAt" json:"updatedAt"`
}

func (sp *SupplierProduct) TotalCost() float64 {
	return sp.Cost + sp.ShippingCost
}

type CreateSupplierProductInput struct {
	ProductID    bson.ObjectID `json:"productId" binding:"required"`
	SupplierID   bson.ObjectID `json:"supplierId" binding:"required"`
	SupplierSKU  string        `json:"supplierSku"`
	Cost         float64       `json:"cost" binding:"required,gt=0"`
	ShippingCost float64       `json:"shippingCost"`
	LeadTimeDays int           `json:"leadTimeDays"`
	MinOrderQty  int           `json:"minOrderQty"`
	URL          string        `json:"url"`
}

type UpdateSupplierProductInput struct {
	SupplierSKU  *string  `json:"supplierSku"`
	Cost         *float64 `json:"cost"`
	ShippingCost *float64 `json:"shippingCost"`
	LeadTimeDays *int     `json:"leadTimeDays"`
	MinOrderQty  *int     `json:"minOrderQty"`
	URL          *string  `json:"url"`
	IsActive     *bool    `json:"isActive"`
}

type SupplierProductFilter struct {
	ProductID  bson.ObjectID
	SupplierID bson.ObjectID
	IsActive   *bool
}
