package domain

import "time"

type QuoteStatus string

const (
	QuoteStatusDraft    QuoteStatus = "DRAFT"
	QuoteStatusSent     QuoteStatus = "SENT"
	QuoteStatusAccepted QuoteStatus = "ACCEPTED"
	QuoteStatusRejected QuoteStatus = "REJECTED"
	QuoteStatusExpired  QuoteStatus = "EXPIRED"
)

type Quote struct {
	ID             ID
	CustomerID     ID
	ServiceID      *ID
	PickupAddress  string
	DropoffAddress string
	Status         QuoteStatus
	ValidUntil     *time.Time
	TotalPrice     *float64
	Notes          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
