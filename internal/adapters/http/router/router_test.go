package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"dach-admin/internal/adapters/http/handler"
	"dach-admin/internal/domain"
)

type fakeVerifier struct{}

func (fakeVerifier) Verify(token string) (domain.AuthClaims, error) {
	return domain.AuthClaims{Subject: "00000000-0000-0000-0000-000000000001", Email: "admin@example.com", Role: domain.RoleAdmin}, nil
}

func TestProtectedRoutesAreRegistered(t *testing.T) {
	app := New(handler.Handlers{}, fakeVerifier{})
	routes := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/v1/dashboard"},
		{http.MethodGet, "/api/v1/customers"},
		{http.MethodGet, "/api/v1/bookings"},
		{http.MethodPatch, "/api/v1/bookings/00000000-0000-0000-0000-000000000001/status"},
		{http.MethodPost, "/api/v1/bookings/00000000-0000-0000-0000-000000000001/assign-driver"},
		{http.MethodGet, "/api/v1/drivers"},
		{http.MethodPatch, "/api/v1/drivers/00000000-0000-0000-0000-000000000001/status"},
		{http.MethodGet, "/api/v1/deliveries"},
		{http.MethodPatch, "/api/v1/deliveries/00000000-0000-0000-0000-000000000001/location"},
		{http.MethodGet, "/api/v1/services"},
		{http.MethodGet, "/api/v1/service-pricing"},
		{http.MethodGet, "/api/v1/coverage-zones"},
		{http.MethodGet, "/api/v1/pricing-rules"},
		{http.MethodGet, "/api/v1/reviews"},
		{http.MethodGet, "/api/v1/contact-requests"},
		{http.MethodGet, "/api/v1/team"},
		{http.MethodGet, "/api/v1/audit-logs"},
	}
	for _, route := range routes {
		req := httptest.NewRequest(route.method, route.path, nil)
		res := httptest.NewRecorder()
		app.ServeHTTP(res, req)
		if res.Code == http.StatusNotFound || res.Code == http.StatusMethodNotAllowed {
			t.Fatalf("%s %s was not registered, status %d", route.method, route.path, res.Code)
		}
	}
}
