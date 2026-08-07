package domain

import "time"

type AuditLog struct {
	ID        ID
	UserID    *ID
	Action    string
	Entity    string
	EntityID  *ID
	Changes   []byte
	IPAddress string
	CreatedAt time.Time
}
