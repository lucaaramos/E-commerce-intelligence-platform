package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// Client wraps the mongo client and database to keep the driver out of the
// application wiring and provide a single ping/close surface.
type Client struct {
	client *mongo.Client
	db     *mongo.Database
}

// Connect establishes a MongoDB connection and verifies it with a ping.
func Connect(ctx context.Context, uri, dbName string) (*Client, error) {
	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		_ = client.Disconnect(pingCtx)
		return nil, err
	}

	return &Client{
		client: client,
		db:     client.Database(dbName),
	}, nil
}

// DB returns the application database.
func (c *Client) DB() *mongo.Database {
	return c.db
}

// Name identifies this dependency in readiness checks.
func (c *Client) Name() string {
	return "mongodb"
}

// Check implements the health checker interface used by /readyz.
func (c *Client) Check(ctx context.Context) error {
	return c.client.Ping(ctx, readpref.Primary())
}

// Close gracefully releases the connection pool.
func (c *Client) Close(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}
