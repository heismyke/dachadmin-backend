package contactrequest

import (
	"context"
	"dach-admin/internal/application"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"

	"github.com/google/uuid"
)

type Service struct {
	repo repositories.ContactRequestRepository
}

func NewService(repo repositories.ContactRequestRepository) *Service { return &Service{repo: repo} }
func (s *Service) Create(ctx context.Context, item *domain.ContactRequest) error {
	if err := application.Required("full_name", item.FullName); err != nil {
		return err
	}
	if err := application.Email("email", item.Email); err != nil {
		return err
	}
	if err := application.Required("message", item.Message); err != nil {
		return err
	}
	if item.Status == "" {
		item.Status = domain.ContactRequestStatusNew
	}
	return s.repo.Create(ctx, item)
}
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.ContactRequest, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *Service) List(ctx context.Context, filter repositories.ContactRequestFilter) (repositories.ListResult[domain.ContactRequest], error) {
	return s.repo.List(ctx, filter)
}
func (s *Service) Update(ctx context.Context, item *domain.ContactRequest) error {
	if err := application.Required("message", item.Message); err != nil {
		return err
	}
	return s.repo.Update(ctx, item)
}
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error { return s.repo.Delete(ctx, id) }
func (s *Service) Resolve(ctx context.Context, id uuid.UUID, actorID *uuid.UUID) error {
	return s.repo.Resolve(ctx, id, actorID)
}
