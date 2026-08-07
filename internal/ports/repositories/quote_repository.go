package repositories

import (
	"context"
	"dach-admin/internal/domain"

	"github.com/google/uuid"
)

type QuoteFilter struct {
	PageFilter
	CustomerID *uuid.UUID
	Status     string
}

type QuoteRepository interface {
	Create(ctx context.Context, item *domain.Quote) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Quote, error)
	List(ctx context.Context, filter QuoteFilter) (ListResult[domain.Quote], error)
	Update(ctx context.Context, item *domain.Quote) error
	Delete(ctx context.Context, id uuid.UUID) error
}
