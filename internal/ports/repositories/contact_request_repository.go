package repositories

import (
	"context"
	"dach-admin/internal/domain"

	"github.com/google/uuid"
)

type ContactRequestFilter struct {
	PageFilter
	Status string
}

type ContactRequestRepository interface {
	Create(ctx context.Context, item *domain.ContactRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.ContactRequest, error)
	List(ctx context.Context, filter ContactRequestFilter) (ListResult[domain.ContactRequest], error)
	Update(ctx context.Context, item *domain.ContactRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	Resolve(ctx context.Context, id uuid.UUID, actorID *uuid.UUID) error
}
