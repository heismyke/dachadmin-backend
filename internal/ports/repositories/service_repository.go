package repositories

import (
	"context"
	"dach-admin/internal/domain"

	"github.com/google/uuid"
)

type ServiceFilter struct {
	PageFilter
	Status string
}

type ServiceRepository interface {
	Create(ctx context.Context, item *domain.Service) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Service, error)
	List(ctx context.Context, filter ServiceFilter) (ListResult[domain.Service], error)
	Update(ctx context.Context, item *domain.Service) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddCoverage(ctx context.Context, serviceID uuid.UUID, coverageID uuid.UUID) error
	RemoveCoverage(ctx context.Context, serviceID uuid.UUID, coverageID uuid.UUID) error
}
