package middleware

import (
	"context"
	"dach-admin/internal/domain"
	"strings"

	"net/http"
)

type Verifier interface {
	Verify(token string) (domain.AuthClaims, error)
}

type authKey string

const ClaimsKey authKey = "claims"

func Authenticate(verifier Verifier) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			token, ok := strings.CutPrefix(auth, "Bearer ")
			if !ok || token == "" {
				errorJSON(w, http.StatusUnauthorized, "UNAUTHENTICATED", "authentication required")
				return
			}
			claims, err := verifier.Verify(token)
			if err != nil {
				errorJSON(w, http.StatusUnauthorized, "UNAUTHENTICATED", "authentication required")
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ClaimsKey, claims)))
		})
	}
}

func Authorize(roles ...domain.TeamRole) Middleware {
	allowed := map[domain.TeamRole]bool{domain.RoleSuperAdmin: true}
	for _, role := range roles {
		allowed[role] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(ClaimsKey).(domain.AuthClaims)
			if !ok {
				errorJSON(w, http.StatusUnauthorized, "UNAUTHENTICATED", "authentication required")
				return
			}
			if !allowed[claims.Role] {
				errorJSON(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func errorJSON(w http.ResponseWriter, status int, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(`{"error":{"code":"` + code + `","message":"` + message + `"}}`))
}
