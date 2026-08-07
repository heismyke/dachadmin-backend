package driver

import (
	"context"
	"dach-admin/internal/application"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"

	"github.com/google/uuid"
)

type Service struct{ repo repositories.DriverRepository }

func NewService(repo repositories.DriverRepository) *Service { return &Service{repo: repo} }
func (s *Service) Create(ctx context.Context, item *domain.Driver) error {
	if err := application.Required("full_name", item.FullName); err != nil {
		return err
	}
	if err := application.Required("license_number", item.LicenseNumber); err != nil {
		return err
	}
	if item.Status == "" {
		item.Status = domain.DriverStatusAvailable
	}
	if err := valid(item.Status); err != nil {
		return err
	}
	return s.repo.Create(ctx, item)
}
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.Driver, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *Service) List(ctx context.Context, filter repositories.DriverFilter) (repositories.ListResult[domain.Driver], error) {
	return s.repo.List(ctx, filter)
}
func (s *Service) Update(ctx context.Context, item *domain.Driver) error {
	if err := application.Required("full_name", item.FullName); err != nil {
		return err
	}
	if err := valid(item.Status); err != nil {
		return err
	}
	return s.repo.Update(ctx, item)
}
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error { return s.repo.Delete(ctx, id) }
func (s *Service) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.DriverStatus, actorID *uuid.UUID) error {
	if err := valid(status); err != nil {
		return err
	}
	return s.repo.UpdateStatus(ctx, id, status, actorID)
}
func valid(v domain.DriverStatus) error {
	return application.In("status", v, domain.DriverStatusAvailable, domain.DriverStatusBusy, domain.DriverStatusOffline, domain.DriverStatusSuspended)
}
