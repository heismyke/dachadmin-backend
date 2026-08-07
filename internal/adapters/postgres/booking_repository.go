package postgres

import (
	"context"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type BookingRepository struct{ db *sql.DB }

func NewBookingRepository(db *sql.DB) *BookingRepository { return &BookingRepository{db: db} }
func (r *BookingRepository) Create(ctx context.Context, b *domain.Booking) error {
	return mapError(r.db.QueryRowContext(ctx, `INSERT INTO bookings (customer_id,service_id,quote_id,pickup_address,dropoff_address,pickup_date,delivery_date,status,total_price,notes) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,nullif($10,'')) RETURNING id,created_at,updated_at`, b.CustomerID, b.ServiceID, b.QuoteID, b.PickupAddress, b.DropoffAddress, b.PickupDate, b.DeliveryDate, b.Status, b.TotalPrice, b.Notes).Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt))
}
func (r *BookingRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	return scanBooking(r.db.QueryRowContext(ctx, `SELECT id,customer_id,service_id,quote_id,pickup_address,dropoff_address,pickup_date,delivery_date,status,total_price::float8,coalesce(notes,''),created_at,updated_at FROM bookings WHERE id=$1`, id))
}
func (r *BookingRepository) List(ctx context.Context, f repositories.BookingFilter) (repositories.ListResult[domain.Booking], error) {
	limit, offset := f.LimitOffset()
	where := "TRUE"
	args := []any{limit, offset}
	if f.Status != "" {
		where += fmt.Sprintf(" AND status=$%d", len(args)+1)
		args = append(args, f.Status)
	}
	if f.CustomerID != nil {
		where += fmt.Sprintf(" AND customer_id=$%d", len(args)+1)
		args = append(args, *f.CustomerID)
	}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`SELECT id,customer_id,service_id,quote_id,pickup_address,dropoff_address,pickup_date,delivery_date,status,total_price::float8,coalesce(notes,''),created_at,updated_at,count(*) OVER() FROM bookings WHERE %s ORDER BY created_at DESC LIMIT $1 OFFSET $2`, where), args...)
	if err != nil {
		return repositories.ListResult[domain.Booking]{}, mapError(err)
	}
	defer rows.Close()
	res := repositories.ListResult[domain.Booking]{Items: []domain.Booking{}}
	for rows.Next() {
		b := domain.Booking{}
		var total sql.NullFloat64
		if err := rows.Scan(&b.ID, &b.CustomerID, &b.ServiceID, &b.QuoteID, &b.PickupAddress, &b.DropoffAddress, &b.PickupDate, &b.DeliveryDate, &b.Status, &total, &b.Notes, &b.CreatedAt, &b.UpdatedAt, &res.Total); err != nil {
			return res, mapError(err)
		}
		b.TotalPrice = floatPtr(total)
		res.Items = append(res.Items, b)
	}
	return res, mapError(rows.Err())
}
func (r *BookingRepository) Update(ctx context.Context, b *domain.Booking) error {
	return mapError(r.db.QueryRowContext(ctx, `UPDATE bookings SET customer_id=$2,service_id=$3,quote_id=$4,pickup_address=$5,dropoff_address=$6,pickup_date=$7,delivery_date=$8,status=$9,total_price=$10,notes=nullif($11,''),updated_at=now() WHERE id=$1 RETURNING updated_at`, b.ID, b.CustomerID, b.ServiceID, b.QuoteID, b.PickupAddress, b.DropoffAddress, b.PickupDate, b.DeliveryDate, b.Status, b.TotalPrice, b.Notes).Scan(&b.UpdatedAt))
}
func (r *BookingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return execOne(ctx, r.db, `DELETE FROM bookings WHERE id=$1`, id)
}
func (r *BookingRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.BookingStatus, actorID *uuid.UUID) error {
	return execOne(ctx, r.db, `UPDATE bookings SET status=$2,updated_at=now() WHERE id=$1`, id, status)
}
func (r *BookingRepository) AssignDriver(ctx context.Context, p repositories.AssignDriverParams) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var bookingStatus domain.BookingStatus
	if err := tx.QueryRowContext(ctx, `SELECT status FROM bookings WHERE id=$1 FOR UPDATE`, p.BookingID).Scan(&bookingStatus); err != nil {
		return mapError(err)
	}
	if bookingStatus == domain.BookingStatusDelivered || bookingStatus == domain.BookingStatusCancelled {
		return domain.ErrConflict
	}
	var driverStatus domain.DriverStatus
	if err := tx.QueryRowContext(ctx, `SELECT status FROM drivers WHERE id=$1 FOR UPDATE`, p.DriverID).Scan(&driverStatus); err != nil {
		return mapError(err)
	}
	if driverStatus != domain.DriverStatusAvailable {
		return domain.ErrConflict
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO live_deliveries (booking_id,driver_id,status,started_at) VALUES ($1,$2,$3,now()) ON CONFLICT (booking_id) DO UPDATE SET driver_id=EXCLUDED.driver_id,status=EXCLUDED.status,updated_at=now()`, p.BookingID, p.DriverID, domain.DeliveryStatusAssigned); err != nil {
		return mapError(err)
	}
	if _, err := tx.ExecContext(ctx, `UPDATE bookings SET status=$2,updated_at=now() WHERE id=$1`, p.BookingID, domain.BookingStatusAssigned); err != nil {
		return mapError(err)
	}
	if _, err := tx.ExecContext(ctx, `UPDATE drivers SET status=$2,updated_at=now() WHERE id=$1`, p.DriverID, domain.DriverStatusBusy); err != nil {
		return mapError(err)
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO audit_logs (user_id,action,entity,entity_id,changes) VALUES ($1,$2,$3,$4,jsonb_build_object('driver_id',$5::text))`, p.ActorID, "ASSIGN_DRIVER", "BOOKING", p.BookingID, p.DriverID); err != nil {
		return mapError(err)
	}
	return mapError(tx.Commit())
}
func scanBooking(row interface{ Scan(dest ...any) error }) (*domain.Booking, error) {
	b := domain.Booking{}
	var total sql.NullFloat64
	if err := row.Scan(&b.ID, &b.CustomerID, &b.ServiceID, &b.QuoteID, &b.PickupAddress, &b.DropoffAddress, &b.PickupDate, &b.DeliveryDate, &b.Status, &total, &b.Notes, &b.CreatedAt, &b.UpdatedAt); err != nil {
		return nil, mapError(err)
	}
	b.TotalPrice = floatPtr(total)
	return &b, nil
}
