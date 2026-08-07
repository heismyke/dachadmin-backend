package domain

import "time"

type DeliveryStatus string

const (
	DeliveryStatusAssigned      DeliveryStatus = "ASSIGNED"
	DeliveryStatusPickupPending DeliveryStatus = "PICKUP_PENDING"
	DeliveryStatusPickedUp      DeliveryStatus = "PICKED_UP"
	DeliveryStatusInTransit     DeliveryStatus = "IN_TRANSIT"
	DeliveryStatusDelivered     DeliveryStatus = "DELIVERED"
	DeliveryStatusFailed        DeliveryStatus = "FAILED"
)

type Delivery struct {
	ID               ID
	BookingID        ID
	DriverID         ID
	Status           DeliveryStatus
	CurrentLatitude  *float64
	CurrentLongitude *float64
	StartedAt        *time.Time
	CompletedAt      *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
