package dashboard

import (
	"context"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"
)

type Service struct {
	repo repositories.DashboardRepository
}

func NewService(repo repositories.DashboardRepository) *Service { return &Service{repo: repo} }
func (s *Service) GetStats(ctx context.Context) (domain.DashboardStats, error) {
	return s.repo.GetStats(ctx)
}
