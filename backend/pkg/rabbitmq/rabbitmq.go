package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Client wraps the AMQP connection. Consumers and publishers are layered on
// top of it in later iterations; this iteration only wires the connection,
// health checks and graceful shutdown.
type Client struct {
	conn *amqp.Connection
	uri  string
}

// Connect establishes a RabbitMQ connection and verifies it by opening and
// closing a channel.
func Connect(ctx context.Context, uri string) (*Client, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("dial rabbitmq: %w", err)
	}
	c := &Client{conn: conn, uri: uri}
	if err := c.Check(ctx); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return c, nil
}

// Name identifies this dependency in readiness checks.
func (c *Client) Name() string {
	return "rabbitmq"
}

// Check implements the health checker interface used by /readyz.
func (c *Client) Check(ctx context.Context) error {
	if c.conn == nil || c.conn.IsClosed() {
		return fmt.Errorf("rabbitmq connection is closed")
	}
	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("rabbitmq channel: %w", err)
	}
	defer ch.Close()
	return nil
}

// Channel returns a fresh AMQP channel for publishing or consuming.
func (c *Client) Channel() (*amqp.Channel, error) {
	return c.conn.Channel()
}

// IsClosed reports whether the underlying connection is closed.
func (c *Client) IsClosed() bool {
	return c.conn == nil || c.conn.IsClosed()
}

// Close gracefully closes the connection.
func (c *Client) Close() error {
	if c.conn == nil || c.conn.IsClosed() {
		return nil
	}
	return c.conn.Close()
}
