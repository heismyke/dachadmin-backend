package handler

import (
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"
	"net/http"
)

func (h Handlers) CreateQuote(w http.ResponseWriter, r *http.Request) {
	var req QuoteRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	q := domain.Quote{CustomerID: req.CustomerID, ServiceID: req.ServiceID, PickupAddress: req.PickupAddress, DropoffAddress: req.DropoffAddress, Status: req.Status, ValidUntil: req.ValidUntil, TotalPrice: req.TotalPrice, Notes: req.Notes}
	if err := h.Quotes.Create(r.Context(), &q); err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusCreated, q)
}
func (h Handlers) GetQuote(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	q, err := h.Quotes.GetByID(r.Context(), rid)
	if err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, q)
}
func (h Handlers) ListQuotes(w http.ResponseWriter, r *http.Request) {
	pf := page(r)
	res, err := h.Quotes.List(r.Context(), repositories.QuoteFilter{PageFilter: pf, Status: r.URL.Query().Get("status")})
	if err != nil {
		writeError(w, err)
		return
	}
	writeList(w, res, pf.Page, pf.PageSize)
}
func (h Handlers) UpdateQuote(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var req QuoteRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	q := domain.Quote{ID: rid, CustomerID: req.CustomerID, ServiceID: req.ServiceID, PickupAddress: req.PickupAddress, DropoffAddress: req.DropoffAddress, Status: req.Status, ValidUntil: req.ValidUntil, TotalPrice: req.TotalPrice, Notes: req.Notes}
	if err := h.Quotes.Update(r.Context(), &q); err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, q)
}
func (h Handlers) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.Quotes.Delete(r.Context(), rid); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h Handlers) CreateDriver(w http.ResponseWriter, r *http.Request) {
	var req DriverRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	d := domain.Driver{FullName: req.FullName, Email: req.Email, Phone: req.Phone, LicenseNumber: req.LicenseNumber, VehicleType: req.VehicleType, VehiclePlate: req.VehiclePlate, Status: req.Status}
	if err := h.Drivers.Create(r.Context(), &d); err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusCreated, d)
}
func (h Handlers) GetDriver(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	d, err := h.Drivers.GetByID(r.Context(), rid)
	if err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, d)
}
func (h Handlers) ListDrivers(w http.ResponseWriter, r *http.Request) {
	pf := page(r)
	res, err := h.Drivers.List(r.Context(), repositories.DriverFilter{PageFilter: pf, Status: r.URL.Query().Get("status")})
	if err != nil {
		writeError(w, err)
		return
	}
	writeList(w, res, pf.Page, pf.PageSize)
}
func (h Handlers) UpdateDriver(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var req DriverRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	d := domain.Driver{ID: rid, FullName: req.FullName, Email: req.Email, Phone: req.Phone, LicenseNumber: req.LicenseNumber, VehicleType: req.VehicleType, VehiclePlate: req.VehiclePlate, Status: req.Status}
	if err := h.Drivers.Update(r.Context(), &d); err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, d)
}
func (h Handlers) DeleteDriver(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.Drivers.Delete(r.Context(), rid); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h Handlers) UpdateDriverStatus(w http.ResponseWriter, r *http.Request) {
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
	if err := h.Drivers.UpdateStatus(r.Context(), rid, domain.DriverStatus(req.Status), actorID(r)); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h Handlers) CreateDelivery(w http.ResponseWriter, r *http.Request) {
	var req DeliveryRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	d := domain.Delivery{BookingID: req.BookingID, DriverID: req.DriverID, Status: req.Status, CurrentLatitude: req.CurrentLatitude, CurrentLongitude: req.CurrentLongitude, StartedAt: req.StartedAt, CompletedAt: req.CompletedAt}
	if err := h.Deliveries.Create(r.Context(), &d); err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusCreated, d)
}
func (h Handlers) GetDelivery(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	d, err := h.Deliveries.GetByID(r.Context(), rid)
	if err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, d)
}
func (h Handlers) ListDeliveries(w http.ResponseWriter, r *http.Request) {
	pf := page(r)
	res, err := h.Deliveries.List(r.Context(), repositories.DeliveryFilter{PageFilter: pf, Status: r.URL.Query().Get("status")})
	if err != nil {
		writeError(w, err)
		return
	}
	writeList(w, res, pf.Page, pf.PageSize)
}
func (h Handlers) UpdateDelivery(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var req DeliveryRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	d := domain.Delivery{ID: rid, BookingID: req.BookingID, DriverID: req.DriverID, Status: req.Status, CurrentLatitude: req.CurrentLatitude, CurrentLongitude: req.CurrentLongitude, StartedAt: req.StartedAt, CompletedAt: req.CompletedAt}
	if err := h.Deliveries.Update(r.Context(), &d); err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, d)
}
func (h Handlers) DeleteDelivery(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.Deliveries.Delete(r.Context(), rid); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h Handlers) UpdateDeliveryLocation(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var req LocationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	if err := h.Deliveries.UpdateLocation(r.Context(), rid, req.Latitude, req.Longitude, actorID(r)); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h Handlers) UpdateDeliveryStatus(w http.ResponseWriter, r *http.Request) {
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
	if err := h.Deliveries.UpdateStatus(r.Context(), rid, domain.DeliveryStatus(req.Status), actorID(r)); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
