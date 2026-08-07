package review

import (
	"context"
	"dach-admin/internal/application"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"

	"github.com/google/uuid"
)

type Service struct{ repo repositories.ReviewRepository }

func NewService(repo repositories.ReviewRepository) *Service { return &Service{repo: repo} }
func (s *Service) Create(ctx context.Context, item *domain.Review) error {
	if err := validReview(item); err != nil {
		return err
	}
	if item.Status == "" {
		item.Status = domain.ReviewStatusPending
	}
	return s.repo.Create(ctx, item)
}
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.Review, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *Service) List(ctx context.Context, filter repositories.ReviewFilter) (repositories.ListResult[domain.Review], error) {
	return s.repo.List(ctx, filter)
}
func (s *Service) Update(ctx context.Context, item *domain.Review) error {
	if err := validReview(item); err != nil {
		return err
	}
	return s.repo.Update(ctx, item)
}
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error { return s.repo.Delete(ctx, id) }
func (s *Service) Publish(ctx context.Context, id uuid.UUID, actorID *uuid.UUID) error {
	return s.repo.Publish(ctx, id, actorID)
}
func validReview(item *domain.Review) error {
	if item.Rating < 1 || item.Rating > 5 {
		return domain.ValidationError{Field: "rating", Message: "must be between 1 and 5"}
	}
	if item.Status != "" {
		return application.In("status", item.Status, domain.ReviewStatusPending, domain.ReviewStatusPublished, domain.ReviewStatusHidden)
	}
	return nil
}
