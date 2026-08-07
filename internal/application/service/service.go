package service

import (
	"context"
	"dach-admin/internal/application"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"

	"github.com/google/uuid"
)

type Service struct {
	repo repositories.ServiceRepository
}

func NewService(repo repositories.ServiceRepository) *Service { return &Service{repo: repo} }
func (s *Service) Create(ctx context.Context, item *domain.Service) error {
	if err := application.Required("name", item.Name); err != nil {
		return err
	}
	if item.Status == "" {
		item.Status = domain.ServiceStatusActive
	}
	if err := application.In("status", item.Status, domain.ServiceStatusActive, domain.ServiceStatusInactive); err != nil {
		return err
	}
	return s.repo.Create(ctx, item)
}
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.Service, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *Service) List(ctx context.Context, filter repositories.ServiceFilter) (repositories.ListResult[domain.Service], error) {
	return s.repo.List(ctx, filter)
}
func (s *Service) Update(ctx context.Context, item *domain.Service) error {
	if err := application.Required("name", item.Name); err != nil {
		return err
	}
	return s.repo.Update(ctx, item)
}
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error { return s.repo.Delete(ctx, id) }
func (s *Service) AddCoverage(ctx context.Context, serviceID uuid.UUID, coverageID uuid.UUID) error {
	return s.repo.AddCoverage(ctx, serviceID, coverageID)
}
func (s *Service) RemoveCoverage(ctx context.Context, serviceID uuid.UUID, coverageID uuid.UUID) error {
	return s.repo.RemoveCoverage(ctx, serviceID, coverageID)
}
