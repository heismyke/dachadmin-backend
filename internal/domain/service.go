package domain

import "time"

type ServiceStatus string

const (
	ServiceStatusActive   ServiceStatus = "ACTIVE"
	ServiceStatusInactive ServiceStatus = "INACTIVE"
)

type Service struct {
	ID          ID
	Name        string
	Description string
	Icon        string
	Status      ServiceStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
