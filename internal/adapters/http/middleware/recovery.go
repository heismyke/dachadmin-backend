package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func Recovery(log *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if v := recover(); v != nil {
					log.Error("panic recovered", "panic", v)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]any{"code": "INTERNAL_ERROR", "message": "internal server error"}})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
