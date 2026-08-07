package services

import (
	"context"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"

	"github.com/google/uuid"
)

type CustomerService interface {
	Create(ctx context.Context, item *domain.Customer) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error)
	List(ctx context.Context, filter repositories.CustomerFilter) (repositories.ListResult[domain.Customer], error)
	Update(ctx context.Context, item *domain.Customer) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ServiceService interface {
	Create(ctx context.Context, item *domain.Service) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Service, error)
	List(ctx context.Context, filter repositories.ServiceFilter) (repositories.ListResult[domain.Service], error)
	Update(ctx context.Context, item *domain.Service) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddCoverage(ctx context.Context, serviceID uuid.UUID, coverageID uuid.UUID) error
	RemoveCoverage(ctx context.Context, serviceID uuid.UUID, coverageID uuid.UUID) error
}

type QuoteService interface {
	Create(ctx context.Context, item *domain.Quote) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Quote, error)
	List(ctx context.Context, filter repositories.QuoteFilter) (repositories.ListResult[domain.Quote], error)
	Update(ctx context.Context, item *domain.Quote) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type BookingService interface {
	Create(ctx context.Context, item *domain.Booking) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error)
	List(ctx context.Context, filter repositories.BookingFilter) (repositories.ListResult[domain.Booking], error)
	Update(ctx context.Context, item *domain.Booking) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.BookingStatus, actorID *uuid.UUID) error
	AssignDriver(ctx context.Context, bookingID uuid.UUID, driverID uuid.UUID, actorID *uuid.UUID) error
}

type DriverService interface {
	Create(ctx context.Context, item *domain.Driver) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Driver, error)
	List(ctx context.Context, filter repositories.DriverFilter) (repositories.ListResult[domain.Driver], error)
	Update(ctx context.Context, item *domain.Driver) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.DriverStatus, actorID *uuid.UUID) error
}

type DeliveryService interface {
	Create(ctx context.Context, item *domain.Delivery) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Delivery, error)
	List(ctx context.Context, filter repositories.DeliveryFilter) (repositories.ListResult[domain.Delivery], error)
	Update(ctx context.Context, item *domain.Delivery) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateLocation(ctx context.Context, id uuid.UUID, lat float64, lon float64, actorID *uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.DeliveryStatus, actorID *uuid.UUID) error
}

type ServicePricingService interface {
	Create(ctx context.Context, item *domain.ServicePricing) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.ServicePricing, error)
	List(ctx context.Context, filter repositories.ServicePricingFilter) (repositories.ListResult[domain.ServicePricing], error)
	Update(ctx context.Context, item *domain.ServicePricing) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PricingRuleService interface {
	Create(ctx context.Context, item *domain.PricingRule) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.PricingRule, error)
	List(ctx context.Context, filter repositories.PricingRuleFilter) (repositories.ListResult[domain.PricingRule], error)
	Update(ctx context.Context, item *domain.PricingRule) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CoverageService interface {
	Create(ctx context.Context, item *domain.CoverageZone) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.CoverageZone, error)
	List(ctx context.Context, filter repositories.CoverageFilter) (repositories.ListResult[domain.CoverageZone], error)
	Update(ctx context.Context, item *domain.CoverageZone) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ReviewService interface {
	Create(ctx context.Context, item *domain.Review) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Review, error)
	List(ctx context.Context, filter repositories.ReviewFilter) (repositories.ListResult[domain.Review], error)
	Update(ctx context.Context, item *domain.Review) error
	Delete(ctx context.Context, id uuid.UUID) error
	Publish(ctx context.Context, id uuid.UUID, actorID *uuid.UUID) error
}

type ContactRequestService interface {
	Create(ctx context.Context, item *domain.ContactRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.ContactRequest, error)
	List(ctx context.Context, filter repositories.ContactRequestFilter) (repositories.ListResult[domain.ContactRequest], error)
	Update(ctx context.Context, item *domain.ContactRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	Resolve(ctx context.Context, id uuid.UUID, actorID *uuid.UUID) error
}

type TeamService interface {
	Create(ctx context.Context, item *domain.TeamMember, plainPassword string) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.TeamMember, error)
	List(ctx context.Context, filter repositories.TeamFilter) (repositories.ListResult[domain.TeamMember], error)
	Update(ctx context.Context, item *domain.TeamMember) error
	Delete(ctx context.Context, id uuid.UUID) error
	Deactivate(ctx context.Context, id uuid.UUID, actorID *uuid.UUID) error
}

type AuditLogService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.AuditLog, error)
	List(ctx context.Context, filter repositories.AuditLogFilter) (repositories.ListResult[domain.AuditLog], error)
}

type DashboardService interface {
	GetStats(ctx context.Context) (domain.DashboardStats, error)
}

type AuthService interface {
	Login(ctx context.Context, email string, password string) (string, *domain.TeamMember, error)
}
