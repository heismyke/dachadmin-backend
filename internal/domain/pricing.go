package domain

import "time"

type PriceType string
type AdjustmentType string

const (
	PriceTypeBase    PriceType = "BASE"
	PriceTypePerKM   PriceType = "PER_KM"
	PriceTypePerMile PriceType = "PER_MILE"
	PriceTypePerKG   PriceType = "PER_KG"
	PriceTypeFlat    PriceType = "FLAT"

	AdjustmentTypeFixed      AdjustmentType = "FIXED"
	AdjustmentTypePercentage AdjustmentType = "PERCENTAGE"
)

type ServicePricing struct {
	ID        ID
	ServiceID ID
	PriceType PriceType
	Price     float64
	Currency  string
	ValidFrom *time.Time
	ValidTo   *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PricingRule struct {
	ID              ID
	Name            string
	Description     string
	RuleType        string
	Conditions      []byte
	AdjustmentType  AdjustmentType
	AdjustmentValue float64
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
