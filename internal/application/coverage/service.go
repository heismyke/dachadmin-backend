package coverage

import (
	"context"
	"dach-admin/internal/application"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"

	"github.com/google/uuid"
)

type Service struct {
	repo repositories.CoverageRepository
}

func NewService(repo repositories.CoverageRepository) *Service { return &Service{repo: repo} }
func (s *Service) Create(ctx context.Context, item *domain.CoverageZone) error {
	if err := application.Required("name", item.Name); err != nil {
		return err
	}
	if item.Status == "" {
		item.Status = domain.CoverageZoneStatusActive
	}
	return s.repo.Create(ctx, item)
}
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.CoverageZone, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *Service) List(ctx context.Context, filter repositories.CoverageFilter) (repositories.ListResult[domain.CoverageZone], error) {
	return s.repo.List(ctx, filter)
}
func (s *Service) Update(ctx context.Context, item *domain.CoverageZone) error {
	if err := application.Required("name", item.Name); err != nil {
		return err
	}
	return s.repo.Update(ctx, item)
}
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error { return s.repo.Delete(ctx, id) }
