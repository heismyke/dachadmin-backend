package repositories

import (
	"context"
	"dach-admin/internal/domain"

	"github.com/google/uuid"
)

type CoverageFilter struct {
	PageFilter
	Status string
}

type CoverageRepository interface {
	Create(ctx context.Context, item *domain.CoverageZone) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.CoverageZone, error)
	List(ctx context.Context, filter CoverageFilter) (ListResult[domain.CoverageZone], error)
	Update(ctx context.Context, item *domain.CoverageZone) error
	Delete(ctx context.Context, id uuid.UUID) error
}
