package postgres

import (
	"context"
	"dach-admin/internal/domain"
	"database/sql"
)

type DashboardRepository struct{ db *sql.DB }

func NewDashboardRepository(db *sql.DB) *DashboardRepository { return &DashboardRepository{db: db} }

func (r *DashboardRepository) GetStats(ctx context.Context) (domain.DashboardStats, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT
			(SELECT count(*) FROM customers),
			(SELECT count(*) FROM bookings),
			(SELECT count(*) FROM bookings WHERE status = 'PENDING'),
			(SELECT count(*) FROM live_deliveries WHERE status IN ('ASSIGNED','PICKUP_PENDING','PICKED_UP','IN_TRANSIT')),
			(SELECT count(*) FROM drivers WHERE status = 'AVAILABLE'),
			(SELECT count(*) FROM quotes WHERE status IN ('DRAFT','SENT')),
			(SELECT count(*) FROM contact_requests WHERE status = 'NEW'),
			(SELECT count(*) FROM reviews WHERE status = 'PENDING'),
			coalesce((SELECT sum(total_price)::float8 FROM bookings WHERE status = 'DELIVERED' AND created_at::date = current_date), 0)`)
	stats := domain.DashboardStats{}
	err := row.Scan(&stats.TotalCustomers, &stats.TotalBookings, &stats.PendingBookings, &stats.ActiveDeliveries, &stats.AvailableDrivers, &stats.PendingQuotes, &stats.NewContactRequests, &stats.PendingReviews, &stats.RevenueToday)
	return stats, mapError(err)
}
