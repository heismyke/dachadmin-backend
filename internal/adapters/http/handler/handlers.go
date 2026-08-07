package handler

import (
	"dach-admin/internal/adapters/http/middleware"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"
	"dach-admin/internal/ports/services"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type Handlers struct {
	Customers       services.CustomerService
	Services        services.ServiceService
	Quotes          services.QuoteService
	Bookings        services.BookingService
	Drivers         services.DriverService
	Deliveries      services.DeliveryService
	ServicePricing  services.ServicePricingService
	PricingRules    services.PricingRuleService
	Coverage        services.CoverageService
	Reviews         services.ReviewService
	ContactRequests services.ContactRequestService
	Team            services.TeamService
	AuditLogs       services.AuditLogService
	Auth            services.AuthService
	Dashboard       services.DashboardService
	DB              *sql.DB
}

func id(r *http.Request, name string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(r.PathValue(name))
	if err != nil {
		return uuid.Nil, domain.ErrValidation
	}
	return parsed, nil
}

func actorID(r *http.Request) *uuid.UUID { return nil }

func (h Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	token, member, err := h.Auth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, envelope{"data": envelope{"token": token, "team_member": teamResponse(*member)}})
}

func (h Handlers) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.ClaimsKey).(domain.AuthClaims)
	if !ok {
		writeError(w, domain.ErrUnauthenticated)
		return
	}
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		writeError(w, domain.ErrUnauthenticated)
		return
	}
	member, err := h.Team.GetByID(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, teamResponse(*member))
}

func (h Handlers) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, envelope{"status": "ok"})
}
func (h Handlers) Ready(w http.ResponseWriter, r *http.Request) {
	if err := h.DB.PingContext(r.Context()); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, envelope{"status": "ok"})
}
func (h Handlers) DashboardStats(w http.ResponseWriter, r *http.Request) {
	data, err := h.Dashboard.GetStats(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, data)
}

func (h Handlers) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var req CustomerRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	c := domain.Customer{FullName: req.FullName, Email: req.Email, Phone: req.Phone, CompanyName: req.CompanyName, Address: req.Address, City: req.City, Country: req.Country}
	if err := h.Customers.Create(r.Context(), &c); err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusCreated, c)
}
func (h Handlers) GetCustomer(w http.ResponseWriter, r *http.Request) {
	id, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	c, err := h.Customers.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, c)
}
func (h Handlers) ListCustomers(w http.ResponseWriter, r *http.Request) {
	pf := page(r)
	res, err := h.Customers.List(r.Context(), repositories.CustomerFilter{PageFilter: pf, Search: r.URL.Query().Get("search")})
	if err != nil {
		writeError(w, err)
		return
	}
	writeList(w, res, pf.Page, pf.PageSize)
}
func (h Handlers) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var req CustomerRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	c := domain.Customer{ID: rid, FullName: req.FullName, Email: req.Email, Phone: req.Phone, CompanyName: req.CompanyName, Address: req.Address, City: req.City, Country: req.Country}
	if err := h.Customers.Update(r.Context(), &c); err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, c)
}
func (h Handlers) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.Customers.Delete(r.Context(), rid); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h Handlers) CreateService(w http.ResponseWriter, r *http.Request) {
	var req ServiceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	item := domain.Service{Name: req.Name, Description: req.Description, Icon: req.Icon, Status: req.Status}
	if err := h.Services.Create(r.Context(), &item); err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusCreated, item)
}
func (h Handlers) GetService(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	item, err := h.Services.GetByID(r.Context(), rid)
	if err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, item)
}
func (h Handlers) ListServices(w http.ResponseWriter, r *http.Request) {
	pf := page(r)
	res, err := h.Services.List(r.Context(), repositories.ServiceFilter{PageFilter: pf, Status: r.URL.Query().Get("status")})
	if err != nil {
		writeError(w, err)
		return
	}
	writeList(w, res, pf.Page, pf.PageSize)
}
func (h Handlers) UpdateService(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var req ServiceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	item := domain.Service{ID: rid, Name: req.Name, Description: req.Description, Icon: req.Icon, Status: req.Status}
	if err := h.Services.Update(r.Context(), &item); err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, item)
}
func (h Handlers) DeleteService(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.Services.Delete(r.Context(), rid); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h Handlers) AddServiceCoverage(w http.ResponseWriter, r *http.Request) {
	sid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	cid, err := id(r, "coverageID")
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.Services.AddCoverage(r.Context(), sid, cid); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h Handlers) RemoveServiceCoverage(w http.ResponseWriter, r *http.Request) {
	sid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	cid, err := id(r, "coverageID")
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.Services.RemoveCoverage(r.Context(), sid, cid); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h Handlers) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var req BookingRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	b := domain.Booking{CustomerID: req.CustomerID, ServiceID: req.ServiceID, QuoteID: req.QuoteID, PickupAddress: req.PickupAddress, DropoffAddress: req.DropoffAddress, PickupDate: req.PickupDate, DeliveryDate: req.DeliveryDate, Status: req.Status, TotalPrice: req.TotalPrice, Notes: req.Notes}
	if err := h.Bookings.Create(r.Context(), &b); err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusCreated, b)
}
func (h Handlers) GetBooking(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	b, err := h.Bookings.GetByID(r.Context(), rid)
	if err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, b)
}
func (h Handlers) ListBookings(w http.ResponseWriter, r *http.Request) {
	pf := page(r)
	res, err := h.Bookings.List(r.Context(), repositories.BookingFilter{PageFilter: pf, Status: r.URL.Query().Get("status")})
	if err != nil {
		writeError(w, err)
		return
	}
	writeList(w, res, pf.Page, pf.PageSize)
}
func (h Handlers) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var req BookingRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	b := domain.Booking{ID: rid, CustomerID: req.CustomerID, ServiceID: req.ServiceID, QuoteID: req.QuoteID, PickupAddress: req.PickupAddress, DropoffAddress: req.DropoffAddress, PickupDate: req.PickupDate, DeliveryDate: req.DeliveryDate, Status: req.Status, TotalPrice: req.TotalPrice, Notes: req.Notes}
	if err := h.Bookings.Update(r.Context(), &b); err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, b)
}
func (h Handlers) DeleteBooking(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.Bookings.Delete(r.Context(), rid); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h Handlers) UpdateBookingStatus(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var req StatusRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	if err := h.Bookings.UpdateStatus(r.Context(), rid, domain.BookingStatus(req.Status), actorID(r)); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h Handlers) AssignDriver(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var req AssignDriverRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	if err := h.Bookings.AssignDriver(r.Context(), rid, req.DriverID, actorID(r)); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h Handlers) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	pf := page(r)
	var entityID *uuid.UUID
	if raw := r.URL.Query().Get("entity_id"); raw != "" {
		parsed, err := uuid.Parse(raw)
		if err != nil {
			writeError(w, domain.ErrValidation)
			return
		}
		entityID = &parsed
	}
	res, err := h.AuditLogs.List(r.Context(), repositories.AuditLogFilter{PageFilter: pf, Entity: r.URL.Query().Get("entity"), EntityID: entityID})
	if err != nil {
		writeError(w, err)
		return
	}
	writeList(w, res, pf.Page, pf.PageSize)
}
func (h Handlers) GetAuditLog(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	a, err := h.AuditLogs.GetByID(r.Context(), rid)
	if err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, a)
}

func boolPtrQuery(r *http.Request, key string) *bool {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return nil
	}
	v, err := strconv.ParseBool(raw)
	if err != nil {
		return nil
	}
	return &v
}
