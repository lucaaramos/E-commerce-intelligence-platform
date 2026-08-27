package pricehistory

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type PriceHistory struct {
	ID               bson.ObjectID `bson:"_id,omitempty" json:"id"`
	SupplierProductID bson.ObjectID `bson:"supplierProductId" json:"supplierProductId"`
	ProductID        bson.ObjectID `bson:"productId" json:"productId"`
	SupplierID       bson.ObjectID `bson:"supplierId" json:"supplierId"`
	Price            float64       `bson:"price" json:"price"`
	ShippingCost     float64       `bson:"shippingCost" json:"shippingCost"`
	Source           string        `bson:"source" json:"source"`
	Timestamp        time.Time     `bson:"timestamp" json:"timestamp"`
}

type PriceHistoryFilter struct {
	SupplierProductID bson.ObjectID
	ProductID         bson.ObjectID
	SupplierID        bson.ObjectID
	Since             *time.Time
	Limit             int
}

type PriceChange struct {
	ProductID  bson.ObjectID `json:"productId"`
	SupplierID bson.ObjectID `json:"supplierId"`
	OldPrice   float64       `json:"oldPrice"`
	NewPrice   float64       `json:"newPrice"`
	Change     float64       `json:"change"`
	ChangePct  float64       `json:"changePct"`
	Timestamp  time.Time     `json:"timestamp"`
}
