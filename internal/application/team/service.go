package team

import (
	"context"
	"dach-admin/internal/application"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct{ repo repositories.TeamRepository }

func NewService(repo repositories.TeamRepository) *Service { return &Service{repo: repo} }
func (s *Service) Create(ctx context.Context, item *domain.TeamMember, plainPassword string) error {
	if err := application.Required("name", item.Name); err != nil {
		return err
	}
	if err := application.Email("email", item.Email); err != nil {
		return err
	}
	if len(plainPassword) < 8 {
		return domain.ValidationError{Field: "password", Message: "must be at least 8 characters"}
	}
	if item.Status == "" {
		item.Status = domain.TeamStatusActive
	}
	if err := valid(item); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	item.PasswordHash = string(hash)
	return s.repo.Create(ctx, item)
}
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.TeamMember, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *Service) List(ctx context.Context, filter repositories.TeamFilter) (repositories.ListResult[domain.TeamMember], error) {
	return s.repo.List(ctx, filter)
}
func (s *Service) Update(ctx context.Context, item *domain.TeamMember) error {
	if err := application.Required("name", item.Name); err != nil {
		return err
	}
	if err := valid(item); err != nil {
		return err
	}
	return s.repo.Update(ctx, item)
}
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error { return s.repo.Delete(ctx, id) }
func (s *Service) Deactivate(ctx context.Context, id uuid.UUID, actorID *uuid.UUID) error {
	return s.repo.Deactivate(ctx, id, actorID)
}
func valid(item *domain.TeamMember) error {
	if err := application.In("role", item.Role, domain.RoleSuperAdmin, domain.RoleAdmin, domain.RoleDispatcher, domain.RoleCustomerSupport, domain.RoleFinance); err != nil {
		return err
	}
	return application.In("status", item.Status, domain.TeamStatusActive, domain.TeamStatusDisabled)
}
