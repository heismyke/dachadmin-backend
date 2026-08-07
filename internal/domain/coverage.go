package domain

import "time"

type CoverageZoneStatus string

const (
	CoverageZoneStatusActive   CoverageZoneStatus = "ACTIVE"
	CoverageZoneStatusInactive CoverageZoneStatus = "INACTIVE"
)

type CoverageZone struct {
	ID          ID
	Name        string
	Description string
	Status      CoverageZoneStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
