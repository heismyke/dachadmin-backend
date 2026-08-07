package domain

import "time"

type DriverStatus string

const (
	DriverStatusAvailable DriverStatus = "AVAILABLE"
	DriverStatusBusy      DriverStatus = "BUSY"
	DriverStatusOffline   DriverStatus = "OFFLINE"
	DriverStatusSuspended DriverStatus = "SUSPENDED"
)

type Driver struct {
	ID            ID
	FullName      string
	Email         string
	Phone         string
	LicenseNumber string
	VehicleType   string
	VehiclePlate  string
	Status        DriverStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
