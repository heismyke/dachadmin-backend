package repositories

import (
	"context"
	"dach-admin/internal/domain"
)

type DashboardRepository interface {
	GetStats(ctx context.Context) (domain.DashboardStats, error)
}
