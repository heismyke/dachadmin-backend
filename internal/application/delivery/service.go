package delivery

import (
	"context"
	"dach-admin/internal/application"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"

	"github.com/google/uuid"
)

type Service struct {
	repo repositories.DeliveryRepository
}

func NewService(repo repositories.DeliveryRepository) *Service { return &Service{repo: repo} }
func (s *Service) Create(ctx context.Context, item *domain.Delivery) error {
	if item.Status == "" {
		item.Status = domain.DeliveryStatusAssigned
	}
	if err := valid(item.Status); err != nil {
		return err
	}
	return s.repo.Create(ctx, item)
}
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.Delivery, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *Service) List(ctx context.Context, filter repositories.DeliveryFilter) (repositories.ListResult[domain.Delivery], error) {
	return s.repo.List(ctx, filter)
}
func (s *Service) Update(ctx context.Context, item *domain.Delivery) error {
	if err := valid(item.Status); err != nil {
		return err
	}
	return s.repo.Update(ctx, item)
}
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error { return s.repo.Delete(ctx, id) }
func (s *Service) UpdateLocation(ctx context.Context, id uuid.UUID, lat float64, lon float64, actorID *uuid.UUID) error {
	if lat < -90 || lat > 90 {
		return domain.ValidationError{Field: "latitude", Message: "must be between -90 and 90"}
	}
	if lon < -180 || lon > 180 {
		return domain.ValidationError{Field: "longitude", Message: "must be between -180 and 180"}
	}
	return s.repo.UpdateLocation(ctx, id, lat, lon, actorID)
}
func (s *Service) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.DeliveryStatus, actorID *uuid.UUID) error {
	if err := valid(status); err != nil {
		return err
	}
	return s.repo.UpdateStatus(ctx, id, status, actorID)
}
func valid(v domain.DeliveryStatus) error {
	return application.In("status", v, domain.DeliveryStatusAssigned, domain.DeliveryStatusPickupPending, domain.DeliveryStatusPickedUp, domain.DeliveryStatusInTransit, domain.DeliveryStatusDelivered, domain.DeliveryStatusFailed)
}
