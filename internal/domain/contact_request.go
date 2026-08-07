package domain

import "time"

type ContactRequestStatus string

const (
	ContactRequestStatusNew        ContactRequestStatus = "NEW"
	ContactRequestStatusInProgress ContactRequestStatus = "IN_PROGRESS"
	ContactRequestStatusResolved   ContactRequestStatus = "RESOLVED"
	ContactRequestStatusClosed     ContactRequestStatus = "CLOSED"
)

type ContactRequest struct {
	ID        ID
	FullName  string
	Email     string
	Phone     string
	Subject   string
	Message   string
	Status    ContactRequestStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}
