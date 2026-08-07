package middleware

import (
	"net/http"
	"strings"
)

func CORS(allowedOrigins []string) Middleware {
	allowed := map[string]bool{}
	for _, origin := range allowedOrigins {
		allowed[normalizeOrigin(origin)] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := normalizeOrigin(r.Header.Get("Origin"))
			if origin != "" && (allowed["*"] || allowed[origin]) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				if allowed["*"] {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				}
				w.Header().Set("Vary", "Origin")
			}
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Max-Age", "600")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func normalizeOrigin(origin string) string {
	return strings.TrimRight(strings.TrimSpace(origin), "/")
}
