package router

import (
	"dach-admin/internal/adapters/http/handler"
	"dach-admin/internal/adapters/http/middleware"
	"dach-admin/internal/domain"
	"net/http"
)

func New(h handler.Handlers, verifier middleware.Verifier) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /ready", h.Ready)
	mux.HandleFunc("POST /api/v1/auth/login", h.Login)
	mux.Handle("GET /api/v1/auth/me", middleware.Chain(http.HandlerFunc(h.Me), middleware.Authenticate(verifier), middleware.Authorize(domain.RoleAdmin, domain.RoleDispatcher, domain.RoleCustomerSupport, domain.RoleFinance)))

	admin := func(next http.HandlerFunc) http.Handler {
		return middleware.Chain(http.HandlerFunc(next), middleware.Authenticate(verifier), middleware.Authorize(domain.RoleAdmin))
	}
	dispatch := func(next http.HandlerFunc) http.Handler {
		return middleware.Chain(http.HandlerFunc(next), middleware.Authenticate(verifier), middleware.Authorize(domain.RoleAdmin, domain.RoleDispatcher, domain.RoleCustomerSupport, domain.RoleFinance))
	}
	support := func(next http.HandlerFunc) http.Handler {
		return middleware.Chain(http.HandlerFunc(next), middleware.Authenticate(verifier), middleware.Authorize(domain.RoleAdmin, domain.RoleCustomerSupport))
	}
	finance := func(next http.HandlerFunc) http.Handler {
		return middleware.Chain(http.HandlerFunc(next), middleware.Authenticate(verifier), middleware.Authorize(domain.RoleAdmin, domain.RoleFinance))
	}

	mux.Handle("GET /api/v1/dashboard", dispatch(h.DashboardStats))

	mux.Handle("GET /api/v1/customers", support(h.ListCustomers))
	mux.Handle("POST /api/v1/customers", support(h.CreateCustomer))
	mux.Handle("GET /api/v1/customers/{id}", support(h.GetCustomer))
	mux.Handle("PUT /api/v1/customers/{id}", support(h.UpdateCustomer))
	mux.Handle("DELETE /api/v1/customers/{id}", admin(h.DeleteCustomer))
	mux.Handle("GET /api/v1/services", admin(h.ListServices))
	mux.Handle("POST /api/v1/services", admin(h.CreateService))
	mux.Handle("GET /api/v1/services/{id}", admin(h.GetService))
	mux.Handle("PUT /api/v1/services/{id}", admin(h.UpdateService))
	mux.Handle("DELETE /api/v1/services/{id}", admin(h.DeleteService))
	mux.Handle("POST /api/v1/services/{id}/coverage/{coverageID}", admin(h.AddServiceCoverage))
	mux.Handle("DELETE /api/v1/services/{id}/coverage/{coverageID}", admin(h.RemoveServiceCoverage))
	mux.Handle("GET /api/v1/quotes", finance(h.ListQuotes))
	mux.Handle("POST /api/v1/quotes", finance(h.CreateQuote))
	mux.Handle("GET /api/v1/quotes/{id}", finance(h.GetQuote))
	mux.Handle("PUT /api/v1/quotes/{id}", finance(h.UpdateQuote))
	mux.Handle("DELETE /api/v1/quotes/{id}", admin(h.DeleteQuote))
	mux.Handle("GET /api/v1/bookings", dispatch(h.ListBookings))
	mux.Handle("POST /api/v1/bookings", dispatch(h.CreateBooking))
	mux.Handle("GET /api/v1/bookings/{id}", dispatch(h.GetBooking))
	mux.Handle("PUT /api/v1/bookings/{id}", dispatch(h.UpdateBooking))
	mux.Handle("DELETE /api/v1/bookings/{id}", admin(h.DeleteBooking))
	mux.Handle("PATCH /api/v1/bookings/{id}/status", dispatch(h.UpdateBookingStatus))
	mux.Handle("POST /api/v1/bookings/{id}/assign-driver", dispatch(h.AssignDriver))
	mux.Handle("GET /api/v1/drivers", dispatch(h.ListDrivers))
	mux.Handle("POST /api/v1/drivers", dispatch(h.CreateDriver))
	mux.Handle("GET /api/v1/drivers/{id}", dispatch(h.GetDriver))
	mux.Handle("PUT /api/v1/drivers/{id}", dispatch(h.UpdateDriver))
	mux.Handle("DELETE /api/v1/drivers/{id}", admin(h.DeleteDriver))
	mux.Handle("PATCH /api/v1/drivers/{id}/status", dispatch(h.UpdateDriverStatus))
	mux.Handle("GET /api/v1/deliveries", dispatch(h.ListDeliveries))
	mux.Handle("POST /api/v1/deliveries", dispatch(h.CreateDelivery))
	mux.Handle("GET /api/v1/deliveries/{id}", dispatch(h.GetDelivery))
	mux.Handle("PUT /api/v1/deliveries/{id}", dispatch(h.UpdateDelivery))
	mux.Handle("DELETE /api/v1/deliveries/{id}", admin(h.DeleteDelivery))
	mux.Handle("PATCH /api/v1/deliveries/{id}/location", dispatch(h.UpdateDeliveryLocation))
	mux.Handle("PATCH /api/v1/deliveries/{id}/status", dispatch(h.UpdateDeliveryStatus))
	mux.Handle("GET /api/v1/service-pricing", finance(h.ListServicePricing))
	mux.Handle("POST /api/v1/service-pricing", finance(h.CreateServicePricing))
	mux.Handle("GET /api/v1/service-pricing/{id}", finance(h.GetServicePricing))
	mux.Handle("PUT /api/v1/service-pricing/{id}", finance(h.UpdateServicePricing))
	mux.Handle("DELETE /api/v1/service-pricing/{id}", admin(h.DeleteServicePricing))
	mux.Handle("GET /api/v1/pricing-rules", finance(h.ListPricingRules))
	mux.Handle("POST /api/v1/pricing-rules", finance(h.CreatePricingRule))
	mux.Handle("GET /api/v1/pricing-rules/{id}", finance(h.GetPricingRule))
	mux.Handle("PUT /api/v1/pricing-rules/{id}", finance(h.UpdatePricingRule))
	mux.Handle("DELETE /api/v1/pricing-rules/{id}", admin(h.DeletePricingRule))
	mux.Handle("GET /api/v1/coverage-zones", admin(h.ListCoverage))
	mux.Handle("POST /api/v1/coverage-zones", admin(h.CreateCoverage))
	mux.Handle("GET /api/v1/coverage-zones/{id}", admin(h.GetCoverage))
	mux.Handle("PUT /api/v1/coverage-zones/{id}", admin(h.UpdateCoverage))
	mux.Handle("DELETE /api/v1/coverage-zones/{id}", admin(h.DeleteCoverage))
	mux.Handle("GET /api/v1/reviews", support(h.ListReviews))
	mux.Handle("POST /api/v1/reviews", support(h.CreateReview))
	mux.Handle("GET /api/v1/reviews/{id}", support(h.GetReview))
	mux.Handle("PUT /api/v1/reviews/{id}", support(h.UpdateReview))
	mux.Handle("DELETE /api/v1/reviews/{id}", admin(h.DeleteReview))
	mux.Handle("PATCH /api/v1/reviews/{id}/publish", support(h.PublishReview))
	mux.Handle("GET /api/v1/contact-requests", support(h.ListContactRequests))
	mux.Handle("POST /api/v1/contact-requests", support(h.CreateContactRequest))
	mux.Handle("GET /api/v1/contact-requests/{id}", support(h.GetContactRequest))
	mux.Handle("PUT /api/v1/contact-requests/{id}", support(h.UpdateContactRequest))
	mux.Handle("DELETE /api/v1/contact-requests/{id}", admin(h.DeleteContactRequest))
	mux.Handle("PATCH /api/v1/contact-requests/{id}/resolve", support(h.ResolveContactRequest))
	mux.Handle("GET /api/v1/team", admin(h.ListTeamMembers))
	mux.Handle("POST /api/v1/team", admin(h.CreateTeamMember))
	mux.Handle("GET /api/v1/team/{id}", admin(h.GetTeamMember))
	mux.Handle("PUT /api/v1/team/{id}", admin(h.UpdateTeamMember))
	mux.Handle("DELETE /api/v1/team/{id}", admin(h.DeleteTeamMember))
	mux.Handle("PATCH /api/v1/team/{id}/deactivate", admin(h.DeactivateTeamMember))
	mux.Handle("GET /api/v1/audit-logs", admin(h.ListAuditLogs))
	mux.Handle("GET /api/v1/audit-logs/{id}", admin(h.GetAuditLog))

	return mux
}
