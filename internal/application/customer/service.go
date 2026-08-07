package customer

import (
	"context"
	"dach-admin/internal/application"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"

	"github.com/google/uuid"
)

type Service struct {
	repo  repositories.CustomerRepository
	audit repositories.AuditLogger
}

func NewService(repo repositories.CustomerRepository, audit repositories.AuditLogger) *Service {
	return &Service{repo: repo, audit: audit}
}

func (s *Service) Create(ctx context.Context, item *domain.Customer) error {
	if err := validate(item); err != nil {
		return err
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return err
	}
	return audit(ctx, s.audit, "CREATE", "CUSTOMER", &item.ID)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, filter repositories.CustomerFilter) (repositories.ListResult[domain.Customer], error) {
	return s.repo.List(ctx, filter)
}

func (s *Service) Update(ctx context.Context, item *domain.Customer) error {
	if err := validate(item); err != nil {
		return err
	}
	if err := s.repo.Update(ctx, item); err != nil {
		return err
	}
	return audit(ctx, s.audit, "UPDATE", "CUSTOMER", &item.ID)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	return audit(ctx, s.audit, "DELETE", "CUSTOMER", &id)
}

func validate(item *domain.Customer) error {
	if err := application.Required("full_name", item.FullName); err != nil {
		return err
	}
	return application.Email("email", item.Email)
}

func audit(ctx context.Context, logger repositories.AuditLogger, action string, entity string, entityID *uuid.UUID) error {
	if logger == nil {
		return nil
	}
	return logger.Log(ctx, domain.AuditLog{Action: action, Entity: entity, EntityID: entityID})
}
