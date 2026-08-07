package handler

import (
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"
	"net/http"
)

func (h Handlers) CreateServicePricing(w http.ResponseWriter, r *http.Request) {
	var req ServicePricingRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	p := domain.ServicePricing{ServiceID: req.ServiceID, PriceType: req.PriceType, Price: req.Price, Currency: req.Currency, ValidFrom: req.ValidFrom, ValidTo: req.ValidTo}
	if err := h.ServicePricing.Create(r.Context(), &p); err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusCreated, p)
}
func (h Handlers) GetServicePricing(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	p, err := h.ServicePricing.GetByID(r.Context(), rid)
	if err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, p)
}
func (h Handlers) ListServicePricing(w http.ResponseWriter, r *http.Request) {
	pf := page(r)
	res, err := h.ServicePricing.List(r.Context(), repositories.ServicePricingFilter{PageFilter: pf})
	if err != nil {
		writeError(w, err)
		return
	}
	writeList(w, res, pf.Page, pf.PageSize)
}
func (h Handlers) UpdateServicePricing(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var req ServicePricingRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	p := domain.ServicePricing{ID: rid, ServiceID: req.ServiceID, PriceType: req.PriceType, Price: req.Price, Currency: req.Currency, ValidFrom: req.ValidFrom, ValidTo: req.ValidTo}
	if err := h.ServicePricing.Update(r.Context(), &p); err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, p)
}
func (h Handlers) DeleteServicePricing(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.ServicePricing.Delete(r.Context(), rid); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h Handlers) CreatePricingRule(w http.ResponseWriter, r *http.Request) {
	var req PricingRuleRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	p := domain.PricingRule{Name: req.Name, Description: req.Description, RuleType: req.RuleType, Conditions: req.Conditions, AdjustmentType: req.AdjustmentType, AdjustmentValue: req.AdjustmentValue, IsActive: req.IsActive}
	if err := h.PricingRules.Create(r.Context(), &p); err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusCreated, p)
}
func (h Handlers) GetPricingRule(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	p, err := h.PricingRules.GetByID(r.Context(), rid)
	if err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, p)
}
func (h Handlers) ListPricingRules(w http.ResponseWriter, r *http.Request) {
	pf := page(r)
	res, err := h.PricingRules.List(r.Context(), repositories.PricingRuleFilter{PageFilter: pf, Active: boolPtrQuery(r, "active")})
	if err != nil {
		writeError(w, err)
		return
	}
	writeList(w, res, pf.Page, pf.PageSize)
}
func (h Handlers) UpdatePricingRule(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var req PricingRuleRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	p := domain.PricingRule{ID: rid, Name: req.Name, Description: req.Description, RuleType: req.RuleType, Conditions: req.Conditions, AdjustmentType: req.AdjustmentType, AdjustmentValue: req.AdjustmentValue, IsActive: req.IsActive}
	if err := h.PricingRules.Update(r.Context(), &p); err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, p)
}
func (h Handlers) DeletePricingRule(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.PricingRules.Delete(r.Context(), rid); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h Handlers) CreateCoverage(w http.ResponseWriter, r *http.Request) {
	var req CoverageRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	z := domain.CoverageZone{Name: req.Name, Description: req.Description, Status: req.Status}
	if err := h.Coverage.Create(r.Context(), &z); err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusCreated, z)
}
func (h Handlers) GetCoverage(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	z, err := h.Coverage.GetByID(r.Context(), rid)
	if err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, z)
}
func (h Handlers) ListCoverage(w http.ResponseWriter, r *http.Request) {
	pf := page(r)
	res, err := h.Coverage.List(r.Context(), repositories.CoverageFilter{PageFilter: pf, Status: r.URL.Query().Get("status")})
	if err != nil {
		writeError(w, err)
		return
	}
	writeList(w, res, pf.Page, pf.PageSize)
}
func (h Handlers) UpdateCoverage(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var req CoverageRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	z := domain.CoverageZone{ID: rid, Name: req.Name, Description: req.Description, Status: req.Status}
	if err := h.Coverage.Update(r.Context(), &z); err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, z)
}
func (h Handlers) DeleteCoverage(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.Coverage.Delete(r.Context(), rid); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h Handlers) CreateReview(w http.ResponseWriter, r *http.Request) {
	var req ReviewRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	v := domain.Review{CustomerID: req.CustomerID, BookingID: req.BookingID, Rating: req.Rating, Comment: req.Comment, Status: req.Status}
	if err := h.Reviews.Create(r.Context(), &v); err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusCreated, v)
}
func (h Handlers) GetReview(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	v, err := h.Reviews.GetByID(r.Context(), rid)
	if err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, v)
}
func (h Handlers) ListReviews(w http.ResponseWriter, r *http.Request) {
	pf := page(r)
	res, err := h.Reviews.List(r.Context(), repositories.ReviewFilter{PageFilter: pf, Status: r.URL.Query().Get("status")})
	if err != nil {
		writeError(w, err)
		return
	}
	writeList(w, res, pf.Page, pf.PageSize)
}
func (h Handlers) UpdateReview(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var req ReviewRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	v := domain.Review{ID: rid, CustomerID: req.CustomerID, BookingID: req.BookingID, Rating: req.Rating, Comment: req.Comment, Status: req.Status}
	if err := h.Reviews.Update(r.Context(), &v); err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, v)
}
func (h Handlers) DeleteReview(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.Reviews.Delete(r.Context(), rid); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h Handlers) PublishReview(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.Reviews.Publish(r.Context(), rid, actorID(r)); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
