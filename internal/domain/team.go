package domain

import "time"

type TeamRole string
type TeamStatus string

const (
	RoleSuperAdmin      TeamRole = "SUPER_ADMIN"
	RoleAdmin           TeamRole = "ADMIN"
	RoleDispatcher      TeamRole = "DISPATCHER"
	RoleCustomerSupport TeamRole = "CUSTOMER_SUPPORT"
	RoleFinance         TeamRole = "FINANCE"

	TeamStatusActive   TeamStatus = "ACTIVE"
	TeamStatusDisabled TeamStatus = "DISABLED"
)

type TeamMember struct {
	ID           ID
	Name         string
	Email        string
	PasswordHash string
	Role         TeamRole
	Status       TeamStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
