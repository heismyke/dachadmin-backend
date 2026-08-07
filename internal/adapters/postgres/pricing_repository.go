package postgres

import (
	"context"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

type ServicePricingRepository struct{ db *sql.DB }

func NewServicePricingRepository(db *sql.DB) *ServicePricingRepository {
	return &ServicePricingRepository{db: db}
}
func (r *ServicePricingRepository) Create(ctx context.Context, p *domain.ServicePricing) error {
	return mapError(r.db.QueryRowContext(ctx, `INSERT INTO service_pricing (service_id,price_type,price,currency,valid_from,valid_to) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id,created_at,updated_at`, p.ServiceID, p.PriceType, p.Price, p.Currency, p.ValidFrom, p.ValidTo).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt))
}
func (r *ServicePricingRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ServicePricing, error) {
	return scanServicePricing(r.db.QueryRowContext(ctx, `SELECT id,service_id,price_type,price::float8,currency,valid_from,valid_to,created_at,updated_at FROM service_pricing WHERE id=$1`, id))
}
func (r *ServicePricingRepository) List(ctx context.Context, f repositories.ServicePricingFilter) (repositories.ListResult[domain.ServicePricing], error) {
	limit, offset := f.LimitOffset()
	where := "TRUE"
	args := []any{limit, offset}
	if f.ServiceID != nil {
		where = "service_id=$3"
		args = append(args, *f.ServiceID)
	}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`SELECT id,service_id,price_type,price::float8,currency,valid_from,valid_to,created_at,updated_at,count(*) OVER() FROM service_pricing WHERE %s ORDER BY created_at DESC LIMIT $1 OFFSET $2`, where), args...)
	if err != nil {
		return repositories.ListResult[domain.ServicePricing]{}, mapError(err)
	}
	defer rows.Close()
	res := repositories.ListResult[domain.ServicePricing]{Items: []domain.ServicePricing{}}
	for rows.Next() {
		p := domain.ServicePricing{}
		if err := rows.Scan(&p.ID, &p.ServiceID, &p.PriceType, &p.Price, &p.Currency, &p.ValidFrom, &p.ValidTo, &p.CreatedAt, &p.UpdatedAt, &res.Total); err != nil {
			return res, mapError(err)
		}
		res.Items = append(res.Items, p)
	}
	return res, mapError(rows.Err())
}
func (r *ServicePricingRepository) Update(ctx context.Context, p *domain.ServicePricing) error {
	return mapError(r.db.QueryRowContext(ctx, `UPDATE service_pricing SET service_id=$2,price_type=$3,price=$4,currency=$5,valid_from=$6,valid_to=$7,updated_at=now() WHERE id=$1 RETURNING updated_at`, p.ID, p.ServiceID, p.PriceType, p.Price, p.Currency, p.ValidFrom, p.ValidTo).Scan(&p.UpdatedAt))
}
func (r *ServicePricingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return execOne(ctx, r.db, `DELETE FROM service_pricing WHERE id=$1`, id)
}
func scanServicePricing(row interface{ Scan(dest ...any) error }) (*domain.ServicePricing, error) {
	p := domain.ServicePricing{}
	if err := row.Scan(&p.ID, &p.ServiceID, &p.PriceType, &p.Price, &p.Currency, &p.ValidFrom, &p.ValidTo, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, mapError(err)
	}
	return &p, nil
}

type PricingRuleRepository struct{ db *sql.DB }

func NewPricingRuleRepository(db *sql.DB) *PricingRuleRepository {
	return &PricingRuleRepository{db: db}
}
func (r *PricingRuleRepository) Create(ctx context.Context, p *domain.PricingRule) error {
	return mapError(r.db.QueryRowContext(ctx, `INSERT INTO pricing_rules (name,description,rule_type,conditions,adjustment_type,adjustment_value,is_active) VALUES ($1,nullif($2,''),$3,$4,$5,$6,$7) RETURNING id,created_at,updated_at`, p.Name, p.Description, p.RuleType, nullBytes(p.Conditions), p.AdjustmentType, p.AdjustmentValue, p.IsActive).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt))
}
func (r *PricingRuleRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.PricingRule, error) {
	return scanPricingRule(r.db.QueryRowContext(ctx, `SELECT id,name,coalesce(description,''),rule_type,coalesce(conditions,'{}'::jsonb),adjustment_type,adjustment_value::float8,is_active,created_at,updated_at FROM pricing_rules WHERE id=$1`, id))
}
func (r *PricingRuleRepository) List(ctx context.Context, f repositories.PricingRuleFilter) (repositories.ListResult[domain.PricingRule], error) {
	limit, offset := f.LimitOffset()
	where := "TRUE"
	args := []any{limit, offset}
	if f.Active != nil {
		where = "is_active=$3"
		args = append(args, *f.Active)
	}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`SELECT id,name,coalesce(description,''),rule_type,coalesce(conditions,'{}'::jsonb),adjustment_type,adjustment_value::float8,is_active,created_at,updated_at,count(*) OVER() FROM pricing_rules WHERE %s ORDER BY created_at DESC LIMIT $1 OFFSET $2`, where), args...)
	if err != nil {
		return repositories.ListResult[domain.PricingRule]{}, mapError(err)
	}
	defer rows.Close()
	res := repositories.ListResult[domain.PricingRule]{Items: []domain.PricingRule{}}
	for rows.Next() {
		p := domain.PricingRule{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.RuleType, &p.Conditions, &p.AdjustmentType, &p.AdjustmentValue, &p.IsActive, &p.CreatedAt, &p.UpdatedAt, &res.Total); err != nil {
			return res, mapError(err)
		}
		res.Items = append(res.Items, p)
	}
	return res, mapError(rows.Err())
}
func (r *PricingRuleRepository) Update(ctx context.Context, p *domain.PricingRule) error {
	return mapError(r.db.QueryRowContext(ctx, `UPDATE pricing_rules SET name=$2,description=nullif($3,''),rule_type=$4,conditions=$5,adjustment_type=$6,adjustment_value=$7,is_active=$8,updated_at=now() WHERE id=$1 RETURNING updated_at`, p.ID, p.Name, p.Description, p.RuleType, nullBytes(p.Conditions), p.AdjustmentType, p.AdjustmentValue, p.IsActive).Scan(&p.UpdatedAt))
}
func (r *PricingRuleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return execOne(ctx, r.db, `DELETE FROM pricing_rules WHERE id=$1`, id)
}
func scanPricingRule(row interface{ Scan(dest ...any) error }) (*domain.PricingRule, error) {
	p := domain.PricingRule{}
	if err := row.Scan(&p.ID, &p.Name, &p.Description, &p.RuleType, &p.Conditions, &p.AdjustmentType, &p.AdjustmentValue, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, mapError(err)
	}
	return &p, nil
}
