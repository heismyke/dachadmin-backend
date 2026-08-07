package repositories

import (
	"context"
	"dach-admin/internal/domain"

	"github.com/google/uuid"
)

type DeliveryFilter struct {
	PageFilter
	Status string
}

type DeliveryRepository interface {
	Create(ctx context.Context, item *domain.Delivery) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Delivery, error)
	List(ctx context.Context, filter DeliveryFilter) (ListResult[domain.Delivery], error)
	Update(ctx context.Context, item *domain.Delivery) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateLocation(ctx context.Context, id uuid.UUID, lat float64, lon float64, actorID *uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.DeliveryStatus, actorID *uuid.UUID) error
}
