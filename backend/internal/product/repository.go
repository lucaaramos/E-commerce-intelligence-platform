package product

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const collectionName = "products"

type Repository interface {
	Create(ctx context.Context, p *Product) error
	GetByID(ctx context.Context, id bson.ObjectID) (*Product, error)
	GetBySKU(ctx context.Context, sku string) (*Product, error)
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id bson.ObjectID) error
	List(ctx context.Context, filter ProductFilter) ([]*Product, int64, error)
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
		{Keys: bson.D{{Key: "sku", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "deletedAt", Value: 1}}},
	}
	_, err := r.coll.Indexes().CreateMany(ctx, indexes)
	return err
}

func (r *mongoRepository) Create(ctx context.Context, p *Product) error {
	if p.ID.IsZero() {
		p.ID = bson.NewObjectID()
	}
	p.CalculateDerivedFields()
	_, err := r.coll.InsertOne(ctx, p)
	if mongo.IsDuplicateKeyError(err) {
		return ErrSKUExists
	}
	return err
}

func (r *mongoRepository) GetByID(ctx context.Context, id bson.ObjectID) (*Product, error) {
	var p Product
	err := r.coll.FindOne(ctx, bson.M{"_id": id, "deletedAt": bson.M{"$exists": false}}).Decode(&p)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	return &p, err
}

func (r *mongoRepository) GetBySKU(ctx context.Context, sku string) (*Product, error) {
	var p Product
	err := r.coll.FindOne(ctx, bson.M{"sku": sku, "deletedAt": bson.M{"$exists": false}}).Decode(&p)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	return &p, err
}

func (r *mongoRepository) Update(ctx context.Context, p *Product) error {
	p.CalculateDerivedFields()
	p.UpdatedAt = time.Now()
	res, err := r.coll.ReplaceOne(ctx, bson.M{"_id": p.ID, "deletedAt": bson.M{"$exists": false}}, p)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *mongoRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	now := time.Now()
	res, err := r.coll.UpdateOne(ctx,
		bson.M{"_id": id, "deletedAt": bson.M{"$exists": false}},
		bson.M{"$set": bson.M{"deletedAt": now, "status": StatusArchived}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *mongoRepository) List(ctx context.Context, filter ProductFilter) ([]*Product, int64, error) {
	query := bson.M{"deletedAt": bson.M{"$exists": false}}
	if filter.SKU != "" {
		query["sku"] = filter.SKU
	}
	if filter.Status != nil {
		query["status"] = *filter.Status
	}

	opts := options.Find()
	if filter.Limit > 0 {
		opts.SetLimit(int64(filter.Limit))
	}
	if filter.Offset > 0 {
		opts.SetSkip(int64(filter.Offset))
	}
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.coll.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var products []*Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, 0, err
	}

	total, err := r.coll.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}