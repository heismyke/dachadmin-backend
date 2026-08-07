package pricing

import (
	"context"
	"dach-admin/internal/application"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"

	"github.com/google/uuid"
)

type ServicePricingService struct {
	repo repositories.ServicePricingRepository
}

func NewServicePricingService(repo repositories.ServicePricingRepository) *ServicePricingService {
	return &ServicePricingService{repo: repo}
}
func (s *ServicePricingService) Create(ctx context.Context, item *domain.ServicePricing) error {
	if item.Currency == "" {
		item.Currency = "GBP"
	}
	if err := validPrice(item); err != nil {
		return err
	}
	return s.repo.Create(ctx, item)
}
func (s *ServicePricingService) GetByID(ctx context.Context, id uuid.UUID) (*domain.ServicePricing, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *ServicePricingService) List(ctx context.Context, filter repositories.ServicePricingFilter) (repositories.ListResult[domain.ServicePricing], error) {
	return s.repo.List(ctx, filter)
}
func (s *ServicePricingService) Update(ctx context.Context, item *domain.ServicePricing) error {
	if err := validPrice(item); err != nil {
		return err
	}
	return s.repo.Update(ctx, item)
}
func (s *ServicePricingService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
func validPrice(item *domain.ServicePricing) error {
	if err := application.In("price_type", item.PriceType, domain.PriceTypeBase, domain.PriceTypePerKM, domain.PriceTypePerMile, domain.PriceTypePerKG, domain.PriceTypeFlat); err != nil {
		return err
	}
	if err := application.Positive("price", item.Price); err != nil {
		return err
	}
	return application.DateRange(item.ValidFrom, item.ValidTo)
}

type PricingRuleService struct {
	repo repositories.PricingRuleRepository
}

func NewPricingRuleService(repo repositories.PricingRuleRepository) *PricingRuleService {
	return &PricingRuleService{repo: repo}
}
func (s *PricingRuleService) Create(ctx context.Context, item *domain.PricingRule) error {
	if err := validRule(item); err != nil {
		return err
	}
	return s.repo.Create(ctx, item)
}
func (s *PricingRuleService) GetByID(ctx context.Context, id uuid.UUID) (*domain.PricingRule, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *PricingRuleService) List(ctx context.Context, filter repositories.PricingRuleFilter) (repositories.ListResult[domain.PricingRule], error) {
	return s.repo.List(ctx, filter)
}
func (s *PricingRuleService) Update(ctx context.Context, item *domain.PricingRule) error {
	if err := validRule(item); err != nil {
		return err
	}
	return s.repo.Update(ctx, item)
}
func (s *PricingRuleService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
func validRule(item *domain.PricingRule) error {
	if err := application.Required("name", item.Name); err != nil {
		return err
	}
	if err := application.Required("rule_type", item.RuleType); err != nil {
		return err
	}
	if err := application.In("adjustment_type", item.AdjustmentType, domain.AdjustmentTypeFixed, domain.AdjustmentTypePercentage); err != nil {
		return err
	}
	return application.Positive("adjustment_value", item.AdjustmentValue)
}
