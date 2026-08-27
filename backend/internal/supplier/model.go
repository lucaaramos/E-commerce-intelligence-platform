package supplier

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Supplier struct {
	ID                bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name              string        `bson:"name" json:"name" binding:"required"`
	Country           string        `bson:"country" json:"country"`
	IntegrationMethod string        `bson:"integrationMethod" json:"integrationMethod"`
	ContactEmail      string        `bson:"contactEmail" json:"contactEmail"`
	CreatedAt         time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt         time.Time     `bson:"updatedAt" json:"updatedAt"`
	DeletedAt         *time.Time    `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}

type CreateSupplierInput struct {
	Name              string `json:"name" binding:"required"`
	Country           string `json:"country"`
	IntegrationMethod string `json:"integrationMethod"`
	ContactEmail      string `json:"contactEmail"`
}

type UpdateSupplierInput struct {
	Name              *string `json:"name"`
	Country           *string `json:"country"`
	IntegrationMethod *string `json:"integrationMethod"`
	ContactEmail      *string `json:"contactEmail"`
}

type SupplierFilter struct {
	Name    string
	Limit   int
	Offset  int
}
