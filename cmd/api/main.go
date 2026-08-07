package main

import (
	"context"
	"dach-admin/internal/adapters/http/handler"
	"dach-admin/internal/adapters/http/middleware"
	"dach-admin/internal/adapters/http/router"
	"dach-admin/internal/adapters/postgres"
	authapp "dach-admin/internal/application/auth"
	bookingapp "dach-admin/internal/application/booking"
	contactapp "dach-admin/internal/application/contactrequest"
	coverageapp "dach-admin/internal/application/coverage"
	customerapp "dach-admin/internal/application/customer"
	dashboardapp "dach-admin/internal/application/dashboard"
	deliveryapp "dach-admin/internal/application/delivery"
	driverapp "dach-admin/internal/application/driver"
	pricingapp "dach-admin/internal/application/pricing"
	quoteapp "dach-admin/internal/application/quote"
	reviewapp "dach-admin/internal/application/review"
	serviceapp "dach-admin/internal/application/service"
	teamapp "dach-admin/internal/application/team"
	infraauth "dach-admin/internal/infrastructure/auth"
	"dach-admin/internal/infrastructure/bootstrap"
	"dach-admin/internal/infrastructure/config"
	"dach-admin/internal/infrastructure/database"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := database.OpenPostgres(ctx, cfg.DB)
	if err != nil {
		log.Error("open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := bootstrap.EnsureSchema(ctx, db); err != nil {
		log.Error("bootstrap schema", "error", err)
		os.Exit(1)
	}
	if err := bootstrap.EnsureAdmin(ctx, db, cfg.BootstrapAdmin); err != nil {
		log.Error("bootstrap admin", "error", err)
		os.Exit(1)
	}

	jwt := infraauth.NewJWT(cfg.JWT.Secret, cfg.JWT.TTL)
	auditRepo := postgres.NewAuditLogRepository(db)
	customerRepo := postgres.NewCustomerRepository(db)
	serviceRepo := postgres.NewServiceRepository(db)
	quoteRepo := postgres.NewQuoteRepository(db)
	bookingRepo := postgres.NewBookingRepository(db)
	driverRepo := postgres.NewDriverRepository(db)
	deliveryRepo := postgres.NewDeliveryRepository(db)
	servicePricingRepo := postgres.NewServicePricingRepository(db)
	pricingRuleRepo := postgres.NewPricingRuleRepository(db)
	coverageRepo := postgres.NewCoverageRepository(db)
	reviewRepo := postgres.NewReviewRepository(db)
	contactRepo := postgres.NewContactRequestRepository(db)
	teamRepo := postgres.NewTeamRepository(db)
	dashboardRepo := postgres.NewDashboardRepository(db)

	h := handler.Handlers{
		Customers:       customerapp.NewService(customerRepo, auditRepo),
		Services:        serviceapp.NewService(serviceRepo),
		Quotes:          quoteapp.NewService(quoteRepo),
		Bookings:        bookingapp.NewService(bookingRepo),
		Drivers:         driverapp.NewService(driverRepo),
		Deliveries:      deliveryapp.NewService(deliveryRepo),
		ServicePricing:  pricingapp.NewServicePricingService(servicePricingRepo),
		PricingRules:    pricingapp.NewPricingRuleService(pricingRuleRepo),
		Coverage:        coverageapp.NewService(coverageRepo),
		Reviews:         reviewapp.NewService(reviewRepo),
		ContactRequests: contactapp.NewService(contactRepo),
		Team:            teamapp.NewService(teamRepo),
		AuditLogs:       auditRepo,
		Auth:            authapp.NewService(teamRepo, jwt),
		Dashboard:       dashboardapp.NewService(dashboardRepo),
		DB:              db,
	}

	app := middleware.Chain(
		router.New(h, jwt),
		middleware.RequestID,
		middleware.CORS(cfg.AllowedOrigins),
		middleware.Recovery(log),
		middleware.Logging(log),
	)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           app,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Info("server listening", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("server shutdown", "error", err)
	}
}
