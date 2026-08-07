package repositories

import (
	"context"
	"dach-admin/internal/domain"

	"github.com/google/uuid"
)

type ServicePricingFilter struct {
	PageFilter
	ServiceID *uuid.UUID
}

type PricingRuleFilter struct {
	PageFilter
	Active *bool
}

type ServicePricingRepository interface {
	Create(ctx context.Context, item *domain.ServicePricing) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.ServicePricing, error)
	List(ctx context.Context, filter ServicePricingFilter) (ListResult[domain.ServicePricing], error)
	Update(ctx context.Context, item *domain.ServicePricing) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PricingRuleRepository interface {
	Create(ctx context.Context, item *domain.PricingRule) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.PricingRule, error)
	List(ctx context.Context, filter PricingRuleFilter) (ListResult[domain.PricingRule], error)
	Update(ctx context.Context, item *domain.PricingRule) error
	Delete(ctx context.Context, id uuid.UUID) error
}
