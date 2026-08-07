package booking

import (
	"context"
	"dach-admin/internal/application"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"

	"github.com/google/uuid"
)

type Service struct {
	repo repositories.BookingRepository
}

func NewService(repo repositories.BookingRepository) *Service { return &Service{repo: repo} }
func (s *Service) Create(ctx context.Context, item *domain.Booking) error {
	if err := application.Required("pickup_address", item.PickupAddress); err != nil {
		return err
	}
	if err := application.Required("dropoff_address", item.DropoffAddress); err != nil {
		return err
	}
	if item.PickupDate.IsZero() {
		return domain.ValidationError{Field: "pickup_date", Message: "is required"}
	}
	if item.Status == "" {
		item.Status = domain.BookingStatusPending
	}
	if err := validStatus(item.Status); err != nil {
		return err
	}
	if err := application.NonNegative("total_price", item.TotalPrice); err != nil {
		return err
	}
	return s.repo.Create(ctx, item)
}
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *Service) List(ctx context.Context, filter repositories.BookingFilter) (repositories.ListResult[domain.Booking], error) {
	return s.repo.List(ctx, filter)
}
func (s *Service) Update(ctx context.Context, item *domain.Booking) error {
	if err := validStatus(item.Status); err != nil {
		return err
	}
	return s.repo.Update(ctx, item)
}
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error { return s.repo.Delete(ctx, id) }
func (s *Service) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.BookingStatus, actorID *uuid.UUID) error {
	if err := validStatus(status); err != nil {
		return err
	}
	return s.repo.UpdateStatus(ctx, id, status, actorID)
}
func (s *Service) AssignDriver(ctx context.Context, bookingID uuid.UUID, driverID uuid.UUID, actorID *uuid.UUID) error {
	return s.repo.AssignDriver(ctx, repositories.AssignDriverParams{BookingID: bookingID, DriverID: driverID, ActorID: actorID})
}
func validStatus(v domain.BookingStatus) error {
	return application.In("status", v, domain.BookingStatusPending, domain.BookingStatusConfirmed, domain.BookingStatusAssigned, domain.BookingStatusPickedUp, domain.BookingStatusInTransit, domain.BookingStatusDelivered, domain.BookingStatusCancelled)
}
