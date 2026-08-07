package repositories

import (
	"context"
	"dach-admin/internal/domain"

	"github.com/google/uuid"
)

type ReviewFilter struct {
	PageFilter
	Status     string
	CustomerID *uuid.UUID
}

type ReviewRepository interface {
	Create(ctx context.Context, item *domain.Review) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Review, error)
	List(ctx context.Context, filter ReviewFilter) (ListResult[domain.Review], error)
	Update(ctx context.Context, item *domain.Review) error
	Delete(ctx context.Context, id uuid.UUID) error
	Publish(ctx context.Context, id uuid.UUID, actorID *uuid.UUID) error
}
