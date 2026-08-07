package repositories

import (
	"context"
	"dach-admin/internal/domain"

	"github.com/google/uuid"
)

type AuditLogFilter struct {
	PageFilter
	Entity   string
	EntityID *uuid.UUID
	UserID   *uuid.UUID
}

type AuditLogger interface {
	Log(ctx context.Context, entry domain.AuditLog) error
}

type AuditLogRepository interface {
	AuditLogger
	GetByID(ctx context.Context, id uuid.UUID) (*domain.AuditLog, error)
	List(ctx context.Context, filter AuditLogFilter) (ListResult[domain.AuditLog], error)
}
