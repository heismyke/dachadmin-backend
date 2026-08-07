package repositories

import (
	"context"
	"dach-admin/internal/domain"

	"github.com/google/uuid"
)

type CustomerFilter struct {
	PageFilter
	Search string
}

type CustomerRepository interface {
	Create(ctx context.Context, customer *domain.Customer) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error)
	List(ctx context.Context, filter CustomerFilter) (ListResult[domain.Customer], error)
	Update(ctx context.Context, customer *domain.Customer) error
	Delete(ctx context.Context, id uuid.UUID) error
}
