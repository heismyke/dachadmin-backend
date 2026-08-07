package quote

import (
	"context"
	"dach-admin/internal/application"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"

	"github.com/google/uuid"
)

type Service struct{ repo repositories.QuoteRepository }

func NewService(repo repositories.QuoteRepository) *Service { return &Service{repo: repo} }
func (s *Service) Create(ctx context.Context, item *domain.Quote) error {
	if err := application.Required("pickup_address", item.PickupAddress); err != nil {
		return err
	}
	if err := application.Required("dropoff_address", item.DropoffAddress); err != nil {
		return err
	}
	if item.Status == "" {
		item.Status = domain.QuoteStatusDraft
	}
	if err := valid(item.Status); err != nil {
		return err
	}
	if err := application.NonNegative("total_price", item.TotalPrice); err != nil {
		return err
	}
	return s.repo.Create(ctx, item)
}
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.Quote, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *Service) List(ctx context.Context, filter repositories.QuoteFilter) (repositories.ListResult[domain.Quote], error) {
	return s.repo.List(ctx, filter)
}
func (s *Service) Update(ctx context.Context, item *domain.Quote) error {
	if err := valid(item.Status); err != nil {
		return err
	}
	return s.repo.Update(ctx, item)
}
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error { return s.repo.Delete(ctx, id) }
func valid(v domain.QuoteStatus) error {
	return application.In("status", v, domain.QuoteStatusDraft, domain.QuoteStatusSent, domain.QuoteStatusAccepted, domain.QuoteStatusRejected, domain.QuoteStatusExpired)
}
