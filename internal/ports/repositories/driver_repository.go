package repositories

import (
	"context"
	"dach-admin/internal/domain"

	"github.com/google/uuid"
)

type DriverFilter struct {
	PageFilter
	Status string
}

type DriverRepository interface {
	Create(ctx context.Context, item *domain.Driver) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Driver, error)
	List(ctx context.Context, filter DriverFilter) (ListResult[domain.Driver], error)
	Update(ctx context.Context, item *domain.Driver) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.DriverStatus, actorID *uuid.UUID) error
}
