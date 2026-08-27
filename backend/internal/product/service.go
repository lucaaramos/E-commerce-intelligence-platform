package product

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service interface {
	Create(ctx context.Context, input CreateProductInput) (*Product, error)
	GetByID(ctx context.Context, id bson.ObjectID) (*Product, error)
	GetBySKU(ctx context.Context, sku string) (*Product, error)
	Update(ctx context.Context, id bson.ObjectID, input UpdateProductInput) (*Product, error)
	Delete(ctx context.Context, id bson.ObjectID) error
	List(ctx context.Context, filter ProductFilter) ([]*Product, int64, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, input CreateProductInput) (*Product, error) {
	if input.SKU == "" || input.Name == "" || input.SellingPrice <= 0 {
		return nil, ErrInvalidInput
	}

	existing, err := s.repo.GetBySKU(ctx, input.SKU)
	if err == nil && existing != nil {
		return nil, ErrSKUExists
	}

	now := time.Now()
	p := &Product{
		SKU:          input.SKU,
		Name:         input.Name,
		Description:  input.Description,
		Status:       StatusActive,
		SellingPrice: input.SellingPrice,
		Variants:     input.Variants,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) GetByID(ctx context.Context, id bson.ObjectID) (*Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetBySKU(ctx context.Context, sku string) (*Product, error) {
	return s.repo.GetBySKU(ctx, sku)
}

func (s *service) Update(ctx context.Context, id bson.ObjectID, input UpdateProductInput) (*Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		p.Name = *input.Name
	}
	if input.Description != nil {
		p.Description = *input.Description
	}
	if input.Status != nil {
		p.Status = *input.Status
	}
	if input.SellingPrice != nil {
		if *input.SellingPrice <= 0 {
			return nil, ErrInvalidInput
		}
		p.SellingPrice = *input.SellingPrice
	}
	if input.Variants != nil {
		p.Variants = input.Variants
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) Delete(ctx context.Context, id bson.ObjectID) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) List(ctx context.Context, filter ProductFilter) ([]*Product, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	return s.repo.List(ctx, filter)
}