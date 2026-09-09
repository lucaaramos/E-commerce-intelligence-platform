package dummyjson

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseURL = "https://dummyjson.com"

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) ListProducts(ctx context.Context, limit, skip int) ([]ProductDTO, int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/products", nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	if limit > 0 {
		q.Add("limit", fmt.Sprintf("%d", limit))
	}
	q.Add("skip", fmt.Sprintf("%d", skip))
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("api returned status %d", resp.StatusCode)
	}

	type response struct {
		Products []ProductDTO `json:"products"`
		Total    int          `json:"total"`
		Skip     int          `json:"skip"`
		Limit    int          `json:"limit"`
	}

	var res response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return res.Products, res.Total, nil
}

func (c *Client) GetProductByID(ctx context.Context, id int) (*ProductDTO, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/products/%d", baseURL, id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api returned status %d", resp.StatusCode)
	}

	var product ProductDTO
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &product, nil
}

func (c *Client) SearchProducts(ctx context.Context, query string, limit, skip int) ([]ProductDTO, int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/products/search", nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Add("q", query)
	if limit > 0 {
		q.Add("limit", fmt.Sprintf("%d", limit))
	}
	q.Add("skip", fmt.Sprintf("%d", skip))
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("api returned status %d", resp.StatusCode)
	}

	type response struct {
		Products []ProductDTO `json:"products"`
		Total    int          `json:"total"`
	}

	var res response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return res.Products, res.Total, nil
}