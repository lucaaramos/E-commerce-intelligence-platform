package supplierproduct

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	ErrNotFound          = errors.New("supplier product not found")
	ErrAlreadyExists     = errors.New("supplier product already exists")
	ErrNoCheapestFound   = errors.New("no active supplier found for product")
)

const collectionName = "supplier_products"

type Repository interface {
	Create(ctx context.Context, sp *SupplierProduct) error
	GetByID(ctx context.Context, id bson.ObjectID) (*SupplierProduct, error)
	Update(ctx context.Context, sp *SupplierProduct) error
	Delete(ctx context.Context, id bson.ObjectID) error
	List(ctx context.Context, filter SupplierProductFilter) ([]*SupplierProduct, error)
	GetByProductID(ctx context.Context, productID bson.ObjectID) ([]*SupplierProduct, error)
	GetCheapest(ctx context.Context, productID bson.ObjectID) (*SupplierProduct, error)
	EnsureIndexes(ctx context.Context) error
}

type mongoRepository struct {
	coll *mongo.Collection
}

func NewRepository(db *mongo.Database) Repository {
	return &mongoRepository{coll: db.Collection(collectionName)}
}

func (r *mongoRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "productId", Value: 1}, {Key: "supplierId", Value: 1}}},
		{Keys: bson.D{{Key: "productId", Value: 1}}},
		{Keys: bson.D{{Key: "supplierId", Value: 1}}},
	}
	_, err := r.coll.Indexes().CreateMany(ctx, indexes)
	return err
}

func (r *mongoRepository) Create(ctx context.Context, sp *SupplierProduct) error {
	sp.IsActive = true
	sp.CreatedAt = time.Now()
	sp.UpdatedAt = time.Now()
	_, err := r.coll.InsertOne(ctx, sp)
	return err
}

func (r *mongoRepository) GetByID(ctx context.Context, id bson.ObjectID) (*SupplierProduct, error) {
	var sp SupplierProduct
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&sp)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	return &sp, err
}

func (r *mongoRepository) Update(ctx context.Context, sp *SupplierProduct) error {
	sp.UpdatedAt = time.Now()
	res, err := r.coll.ReplaceOne(ctx, bson.M{"_id": sp.ID}, sp)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *mongoRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *mongoRepository) List(ctx context.Context, filter SupplierProductFilter) ([]*SupplierProduct, error) {
	query := bson.M{}
	if !filter.ProductID.IsZero() {
		query["productId"] = filter.ProductID
	}
	if !filter.SupplierID.IsZero() {
		query["supplierId"] = filter.SupplierID
	}
	if filter.IsActive != nil {
		query["isActive"] = *filter.IsActive
	}

	opts := options.Find().SetSort(bson.D{{Key: "cost", Value: 1}})

	cursor, err := r.coll.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*SupplierProduct
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *mongoRepository) GetByProductID(ctx context.Context, productID bson.ObjectID) ([]*SupplierProduct, error) {
	return r.List(ctx, SupplierProductFilter{ProductID: productID, IsActive: boolPtr(true)})
}

func (r *mongoRepository) GetCheapest(ctx context.Context, productID bson.ObjectID) (*SupplierProduct, error) {
	var sp SupplierProduct
	err := r.coll.FindOne(ctx,
		bson.M{"productId": productID, "isActive": true},
		options.FindOne().SetSort(bson.D{{Key: "cost", Value: 1}}),
	).Decode(&sp)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNoCheapestFound
	}
	return &sp, err
}

func boolPtr(b bool) *bool { return &b }
