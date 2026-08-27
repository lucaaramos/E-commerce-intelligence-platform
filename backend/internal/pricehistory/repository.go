package pricehistory

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const collectionName = "price_history"

type Repository interface {
	Create(ctx context.Context, ph *PriceHistory) error
	List(ctx context.Context, filter PriceHistoryFilter) ([]*PriceHistory, error)
	GetLatest(ctx context.Context, supplierProductID bson.ObjectID) (*PriceHistory, error)
	GetPriceChange(ctx context.Context, supplierProductID bson.ObjectID) (*PriceChange, error)
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
		{Keys: bson.D{{Key: "supplierProductId", Value: 1}, {Key: "timestamp", Value: -1}}},
		{Keys: bson.D{{Key: "productId", Value: 1}, {Key: "timestamp", Value: -1}}},
		{Keys: bson.D{{Key: "supplierId", Value: 1}, {Key: "timestamp", Value: -1}}},
	}
	_, err := r.coll.Indexes().CreateMany(ctx, indexes)
	return err
}

func (r *mongoRepository) Create(ctx context.Context, ph *PriceHistory) error {
	ph.Timestamp = time.Now()
	_, err := r.coll.InsertOne(ctx, ph)
	return err
}

func (r *mongoRepository) List(ctx context.Context, filter PriceHistoryFilter) ([]*PriceHistory, error) {
	query := bson.M{}
	if !filter.SupplierProductID.IsZero() {
		query["supplierProductId"] = filter.SupplierProductID
	}
	if !filter.ProductID.IsZero() {
		query["productId"] = filter.ProductID
	}
	if !filter.SupplierID.IsZero() {
		query["supplierId"] = filter.SupplierID
	}
	if filter.Since != nil {
		query["timestamp"] = bson.M{"$gte": *filter.Since}
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.coll.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*PriceHistory
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *mongoRepository) GetLatest(ctx context.Context, supplierProductID bson.ObjectID) (*PriceHistory, error) {
	var ph PriceHistory
	err := r.coll.FindOne(ctx,
		bson.M{"supplierProductId": supplierProductID},
		options.FindOne().SetSort(bson.D{{Key: "timestamp", Value: -1}}),
	).Decode(&ph)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &ph, err
}

func (r *mongoRepository) GetPriceChange(ctx context.Context, supplierProductID bson.ObjectID) (*PriceChange, error) {
	cursor, err := r.coll.Find(ctx,
		bson.M{"supplierProductId": supplierProductID},
		options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(2),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var history []*PriceHistory
	if err := cursor.All(ctx, &history); err != nil {
		return nil, err
	}

	if len(history) < 2 {
		return nil, nil
	}

	old := history[1]
	new := history[0]
	change := new.Price - old.Price
	changePct := (change / old.Price) * 100

	return &PriceChange{
		ProductID:  new.ProductID,
		SupplierID: new.SupplierID,
		OldPrice:   old.Price,
		NewPrice:   new.Price,
		Change:     change,
		ChangePct:  changePct,
		Timestamp:  new.Timestamp,
	}, nil
}
