package domain

type AuthClaims struct {
	Subject string   `json:"sub"`
	Email   string   `json:"email"`
	Role    TeamRole `json:"role"`
	Exp     int64    `json:"exp"`
}
