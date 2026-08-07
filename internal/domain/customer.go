package domain

import "time"

type Customer struct {
	ID          ID
	FullName    string
	Email       string
	Phone       string
	CompanyName string
	Address     string
	City        string
	Country     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
