package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"

	"github.com/google/uuid"
)

type customerSvc struct{}

func (customerSvc) Create(ctx context.Context, item *domain.Customer) error { return nil }
func (customerSvc) GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	return nil, domain.ErrNotFound
}
func (customerSvc) List(ctx context.Context, filter repositories.CustomerFilter) (repositories.ListResult[domain.Customer], error) {
	return repositories.ListResult[domain.Customer]{Items: []domain.Customer{{ID: uuid.New(), FullName: "Jane Doe", Email: "jane@example.com"}}, Total: 1}, nil
}
func (customerSvc) Update(ctx context.Context, item *domain.Customer) error { return nil }
func (customerSvc) Delete(ctx context.Context, id uuid.UUID) error          { return nil }

func TestListCustomersReturnsDataEnvelope(t *testing.T) {
	h := Handlers{Customers: customerSvc{}}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/customers?page=1&page_size=20", nil)
	res := httptest.NewRecorder()
	h.ListCustomers(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d", res.Code)
	}
	body := res.Body.String()
	if body == "" || body[0:8] != `{"data":` {
		t.Fatalf("unexpected response body: %s", body)
	}
}
