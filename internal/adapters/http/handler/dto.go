package handler

import (
	"encoding/json"
	"time"

	"dach-admin/internal/domain"
)

type CustomerRequest struct {
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	City        string `json:"city"`
	Country     string `json:"country"`
}

type ServiceRequest struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Icon        string               `json:"icon"`
	Status      domain.ServiceStatus `json:"status"`
}

type QuoteRequest struct {
	CustomerID     domain.ID          `json:"customer_id"`
	ServiceID      *domain.ID         `json:"service_id"`
	PickupAddress  string             `json:"pickup_address"`
	DropoffAddress string             `json:"dropoff_address"`
	Status         domain.QuoteStatus `json:"status"`
	ValidUntil     *time.Time         `json:"valid_until"`
	TotalPrice     *float64           `json:"total_price"`
	Notes          string             `json:"notes"`
}

type BookingRequest struct {
	CustomerID     domain.ID            `json:"customer_id"`
	ServiceID      domain.ID            `json:"service_id"`
	QuoteID        *domain.ID           `json:"quote_id"`
	PickupAddress  string               `json:"pickup_address"`
	DropoffAddress string               `json:"dropoff_address"`
	PickupDate     time.Time            `json:"pickup_date"`
	DeliveryDate   *time.Time           `json:"delivery_date"`
	Status         domain.BookingStatus `json:"status"`
	TotalPrice     *float64             `json:"total_price"`
	Notes          string               `json:"notes"`
}

type DriverRequest struct {
	FullName      string              `json:"full_name"`
	Email         string              `json:"email"`
	Phone         string              `json:"phone"`
	LicenseNumber string              `json:"license_number"`
	VehicleType   string              `json:"vehicle_type"`
	VehiclePlate  string              `json:"vehicle_plate"`
	Status        domain.DriverStatus `json:"status"`
}

type DeliveryRequest struct {
	BookingID        domain.ID             `json:"booking_id"`
	DriverID         domain.ID             `json:"driver_id"`
	Status           domain.DeliveryStatus `json:"status"`
	CurrentLatitude  *float64              `json:"current_latitude"`
	CurrentLongitude *float64              `json:"current_longitude"`
	StartedAt        *time.Time            `json:"started_at"`
	CompletedAt      *time.Time            `json:"completed_at"`
}

type ServicePricingRequest struct {
	ServiceID domain.ID        `json:"service_id"`
	PriceType domain.PriceType `json:"price_type"`
	Price     float64          `json:"price"`
	Currency  string           `json:"currency"`
	ValidFrom *time.Time       `json:"valid_from"`
	ValidTo   *time.Time       `json:"valid_to"`
}

type PricingRuleRequest struct {
	Name            string                `json:"name"`
	Description     string                `json:"description"`
	RuleType        string                `json:"rule_type"`
	Conditions      json.RawMessage       `json:"conditions"`
	AdjustmentType  domain.AdjustmentType `json:"adjustment_type"`
	AdjustmentValue float64               `json:"adjustment_value"`
	IsActive        bool                  `json:"is_active"`
}

type CoverageRequest struct {
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Status      domain.CoverageZoneStatus `json:"status"`
}

type ReviewRequest struct {
	CustomerID domain.ID           `json:"customer_id"`
	BookingID  domain.ID           `json:"booking_id"`
	Rating     int                 `json:"rating"`
	Comment    string              `json:"comment"`
	Status     domain.ReviewStatus `json:"status"`
}

type ContactRequestRequest struct {
	FullName string                      `json:"full_name"`
	Email    string                      `json:"email"`
	Phone    string                      `json:"phone"`
	Subject  string                      `json:"subject"`
	Message  string                      `json:"message"`
	Status   domain.ContactRequestStatus `json:"status"`
}

type TeamRequest struct {
	Name     string            `json:"name"`
	Email    string            `json:"email"`
	Password string            `json:"password"`
	Role     domain.TeamRole   `json:"role"`
	Status   domain.TeamStatus `json:"status"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type StatusRequest struct {
	Status string `json:"status"`
}

type AssignDriverRequest struct {
	DriverID domain.ID `json:"driver_id"`
}

type LocationRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type TeamResponse struct {
	ID        domain.ID         `json:"id"`
	Name      string            `json:"name"`
	Email     string            `json:"email"`
	Role      domain.TeamRole   `json:"role"`
	Status    domain.TeamStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

func teamResponse(m domain.TeamMember) TeamResponse {
	return TeamResponse{ID: m.ID, Name: m.Name, Email: m.Email, Role: m.Role, Status: m.Status, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}
