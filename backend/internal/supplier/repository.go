package supplier

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	ErrNotFound = errors.New("supplier not found")
	ErrNameExists = errors.New("supplier with this name already exists")
)

const collectionName = "suppliers"

type Repository interface {
	Create(ctx context.Context, s *Supplier) error
	GetByID(ctx context.Context, id bson.ObjectID) (*Supplier, error)
	GetByName(ctx context.Context, name string) (*Supplier, error)
	Update(ctx context.Context, s *Supplier) error
	Delete(ctx context.Context, id bson.ObjectID) error
	List(ctx context.Context, filter SupplierFilter) ([]*Supplier, int64, error)
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
		{Keys: bson.D{{Key: "name", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "deletedAt", Value: 1}}},
	}
	_, err := r.coll.Indexes().CreateMany(ctx, indexes)
	return err
}

func (r *mongoRepository) Create(ctx context.Context, s *Supplier) error {
	if s.ID.IsZero() {
		s.ID = bson.NewObjectID()
	}
	_, err := r.coll.InsertOne(ctx, s)
	if mongo.IsDuplicateKeyError(err) {
		return ErrNameExists
	}
	return err
}

func (r *mongoRepository) GetByID(ctx context.Context, id bson.ObjectID) (*Supplier, error) {
	var s Supplier
	err := r.coll.FindOne(ctx, bson.M{"_id": id, "deletedAt": bson.M{"$exists": false}}).Decode(&s)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	return &s, err
}

func (r *mongoRepository) GetByName(ctx context.Context, name string) (*Supplier, error) {
	var s Supplier
	err := r.coll.FindOne(ctx, bson.M{"name": name, "deletedAt": bson.M{"$exists": false}}).Decode(&s)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	return &s, err
}

func (r *mongoRepository) Update(ctx context.Context, s *Supplier) error {
	s.UpdatedAt = time.Now()
	res, err := r.coll.ReplaceOne(ctx, bson.M{"_id": s.ID, "deletedAt": bson.M{"$exists": false}}, s)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrNameExists
		}
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
		bson.M{"$set": bson.M{"deletedAt": now}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *mongoRepository) List(ctx context.Context, filter SupplierFilter) ([]*Supplier, int64, error) {
	query := bson.M{"deletedAt": bson.M{"$exists": false}}
	if filter.Name != "" {
		query["name"] = bson.M{"$regex": filter.Name, "$options": "i"}
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

	var suppliers []*Supplier
	if err := cursor.All(ctx, &suppliers); err != nil {
		return nil, 0, err
	}

	total, err := r.coll.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	return suppliers, total, nil
}
