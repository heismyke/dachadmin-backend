package domain

import "time"

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "PENDING"
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
	BookingStatusAssigned  BookingStatus = "ASSIGNED"
	BookingStatusPickedUp  BookingStatus = "PICKED_UP"
	BookingStatusInTransit BookingStatus = "IN_TRANSIT"
	BookingStatusDelivered BookingStatus = "DELIVERED"
	BookingStatusCancelled BookingStatus = "CANCELLED"
)

type Booking struct {
	ID             ID
	CustomerID     ID
	ServiceID      ID
	QuoteID        *ID
	PickupAddress  string
	DropoffAddress string
	PickupDate     time.Time
	DeliveryDate   *time.Time
	Status         BookingStatus
	TotalPrice     *float64
	Notes          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
