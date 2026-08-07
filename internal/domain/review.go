package domain

import "time"

type ReviewStatus string

const (
	ReviewStatusPending   ReviewStatus = "PENDING"
	ReviewStatusPublished ReviewStatus = "PUBLISHED"
	ReviewStatusHidden    ReviewStatus = "HIDDEN"
)

type Review struct {
	ID         ID
	CustomerID ID
	BookingID  ID
	Rating     int
	Comment    string
	Status     ReviewStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
