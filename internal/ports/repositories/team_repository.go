package repositories

import (
	"context"
	"dach-admin/internal/domain"

	"github.com/google/uuid"
)

type TeamFilter struct {
	PageFilter
	Role   string
	Status string
}

type TeamRepository interface {
	Create(ctx context.Context, item *domain.TeamMember) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.TeamMember, error)
	GetByEmail(ctx context.Context, email string) (*domain.TeamMember, error)
	List(ctx context.Context, filter TeamFilter) (ListResult[domain.TeamMember], error)
	Update(ctx context.Context, item *domain.TeamMember) error
	Delete(ctx context.Context, id uuid.UUID) error
	Deactivate(ctx context.Context, id uuid.UUID, actorID *uuid.UUID) error
}
