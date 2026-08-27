package supplier

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var (
	ErrInvalidInput = errors.New("invalid input")
)

type Service interface {
	Create(ctx context.Context, input CreateSupplierInput) (*Supplier, error)
	GetByID(ctx context.Context, id bson.ObjectID) (*Supplier, error)
	GetByName(ctx context.Context, name string) (*Supplier, error)
	Update(ctx context.Context, id bson.ObjectID, input UpdateSupplierInput) (*Supplier, error)
	Delete(ctx context.Context, id bson.ObjectID) error
	List(ctx context.Context, filter SupplierFilter) ([]*Supplier, int64, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, input CreateSupplierInput) (*Supplier, error) {
	if input.Name == "" {
		return nil, ErrInvalidInput
	}

	existing, err := s.repo.GetByName(ctx, input.Name)
	if err == nil && existing != nil {
		return nil, ErrNameExists
	}

	now := time.Now()
	sup := &Supplier{
		Name:              input.Name,
		Country:           input.Country,
		IntegrationMethod: input.IntegrationMethod,
		ContactEmail:      input.ContactEmail,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := s.repo.Create(ctx, sup); err != nil {
		return nil, err
	}
	return sup, nil
}

func (s *service) GetByID(ctx context.Context, id bson.ObjectID) (*Supplier, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetByName(ctx context.Context, name string) (*Supplier, error) {
	return s.repo.GetByName(ctx, name)
}

func (s *service) Update(ctx context.Context, id bson.ObjectID, input UpdateSupplierInput) (*Supplier, error) {
	sup, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		sup.Name = *input.Name
	}
	if input.Country != nil {
		sup.Country = *input.Country
	}
	if input.IntegrationMethod != nil {
		sup.IntegrationMethod = *input.IntegrationMethod
	}
	if input.ContactEmail != nil {
		sup.ContactEmail = *input.ContactEmail
	}

	if err := s.repo.Update(ctx, sup); err != nil {
		return nil, err
	}
	return sup, nil
}

func (s *service) Delete(ctx context.Context, id bson.ObjectID) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) List(ctx context.Context, filter SupplierFilter) ([]*Supplier, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	return s.repo.List(ctx, filter)
}
