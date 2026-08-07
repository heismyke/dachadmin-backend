package customer

import (
	"context"
	"errors"
	"testing"

	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"

	"github.com/google/uuid"
)

type fakeCustomerRepo struct{ created *domain.Customer }

func (f *fakeCustomerRepo) Create(ctx context.Context, c *domain.Customer) error {
	c.ID = uuid.New()
	f.created = c
	return nil
}
func (f *fakeCustomerRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	return nil, domain.ErrNotFound
}
func (f *fakeCustomerRepo) List(ctx context.Context, filter repositories.CustomerFilter) (repositories.ListResult[domain.Customer], error) {
	return repositories.ListResult[domain.Customer]{}, nil
}
func (f *fakeCustomerRepo) Update(ctx context.Context, c *domain.Customer) error { return nil }
func (f *fakeCustomerRepo) Delete(ctx context.Context, id uuid.UUID) error       { return nil }

func TestCreateCustomerValidatesEmail(t *testing.T) {
	svc := NewService(&fakeCustomerRepo{}, nil)
	err := svc.Create(context.Background(), &domain.Customer{FullName: "John Doe", Email: "bad"})
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestCreateCustomerPersistsValidCustomer(t *testing.T) {
	repo := &fakeCustomerRepo{}
	svc := NewService(repo, nil)
	err := svc.Create(context.Background(), &domain.Customer{FullName: "John Doe", Email: "john@example.com"})
	if err != nil {
		t.Fatalf("create customer: %v", err)
	}
	if repo.created == nil || repo.created.ID == uuid.Nil {
		t.Fatalf("expected customer to be persisted with id")
	}
}
