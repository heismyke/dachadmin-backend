package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSAllowsConfiguredFrontendOrigin(t *testing.T) {
	handler := CORS([]string{"http://localhost:5184/"})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/customers", nil)
	req.Header.Set("Origin", "http://localhost:5184")
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusNoContent {
		t.Fatalf("status = %d", res.Code)
	}
	if got := res.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5184" {
		t.Fatalf("Access-Control-Allow-Origin = %q", got)
	}
	if got := res.Header().Get("Access-Control-Allow-Headers"); got == "" {
		t.Fatalf("expected Access-Control-Allow-Headers")
	}
}
