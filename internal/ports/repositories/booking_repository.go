package repositories

import (
	"context"
	"dach-admin/internal/domain"

	"github.com/google/uuid"
)

type BookingFilter struct {
	PageFilter
	Status     string
	CustomerID *uuid.UUID
}

type AssignDriverParams struct {
	BookingID uuid.UUID
	DriverID  uuid.UUID
	ActorID   *uuid.UUID
}

type BookingRepository interface {
	Create(ctx context.Context, item *domain.Booking) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error)
	List(ctx context.Context, filter BookingFilter) (ListResult[domain.Booking], error)
	Update(ctx context.Context, item *domain.Booking) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.BookingStatus, actorID *uuid.UUID) error
	AssignDriver(ctx context.Context, params AssignDriverParams) error
}
