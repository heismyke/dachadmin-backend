package postgres

import (
	"context"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

type DeliveryRepository struct{ db *sql.DB }

func NewDeliveryRepository(db *sql.DB) *DeliveryRepository { return &DeliveryRepository{db: db} }
func (r *DeliveryRepository) Create(ctx context.Context, d *domain.Delivery) error {
	return mapError(r.db.QueryRowContext(ctx, `INSERT INTO live_deliveries (booking_id,driver_id,status,current_latitude,current_longitude,started_at,completed_at) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id,created_at,updated_at`, d.BookingID, d.DriverID, d.Status, d.CurrentLatitude, d.CurrentLongitude, d.StartedAt, d.CompletedAt).Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt))
}
func (r *DeliveryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Delivery, error) {
	return scanDelivery(r.db.QueryRowContext(ctx, `SELECT id,booking_id,driver_id,status,current_latitude::float8,current_longitude::float8,started_at,completed_at,created_at,updated_at FROM live_deliveries WHERE id=$1`, id))
}
func (r *DeliveryRepository) List(ctx context.Context, f repositories.DeliveryFilter) (repositories.ListResult[domain.Delivery], error) {
	limit, offset := f.LimitOffset()
	where := "TRUE"
	args := []any{limit, offset}
	if f.Status != "" {
		where = "status=$3"
		args = append(args, f.Status)
	}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`SELECT id,booking_id,driver_id,status,current_latitude::float8,current_longitude::float8,started_at,completed_at,created_at,updated_at,count(*) OVER() FROM live_deliveries WHERE %s ORDER BY created_at DESC LIMIT $1 OFFSET $2`, where), args...)
	if err != nil {
		return repositories.ListResult[domain.Delivery]{}, mapError(err)
	}
	defer rows.Close()
	res := repositories.ListResult[domain.Delivery]{Items: []domain.Delivery{}}
	for rows.Next() {
		d := domain.Delivery{}
		var lat, lon sql.NullFloat64
		if err := rows.Scan(&d.ID, &d.BookingID, &d.DriverID, &d.Status, &lat, &lon, &d.StartedAt, &d.CompletedAt, &d.CreatedAt, &d.UpdatedAt, &res.Total); err != nil {
			return res, mapError(err)
		}
		d.CurrentLatitude = floatPtr(lat)
		d.CurrentLongitude = floatPtr(lon)
		res.Items = append(res.Items, d)
	}
	return res, mapError(rows.Err())
}
func (r *DeliveryRepository) Update(ctx context.Context, d *domain.Delivery) error {
	return mapError(r.db.QueryRowContext(ctx, `UPDATE live_deliveries SET booking_id=$2,driver_id=$3,status=$4,current_latitude=$5,current_longitude=$6,started_at=$7,completed_at=$8,updated_at=now() WHERE id=$1 RETURNING updated_at`, d.ID, d.BookingID, d.DriverID, d.Status, d.CurrentLatitude, d.CurrentLongitude, d.StartedAt, d.CompletedAt).Scan(&d.UpdatedAt))
}
func (r *DeliveryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return execOne(ctx, r.db, `DELETE FROM live_deliveries WHERE id=$1`, id)
}
func (r *DeliveryRepository) UpdateLocation(ctx context.Context, id uuid.UUID, lat float64, lon float64, actorID *uuid.UUID) error {
	return execOne(ctx, r.db, `UPDATE live_deliveries SET current_latitude=$2,current_longitude=$3,updated_at=now() WHERE id=$1`, id, lat, lon)
}
func (r *DeliveryRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.DeliveryStatus, actorID *uuid.UUID) error {
	return execOne(ctx, r.db, `UPDATE live_deliveries SET status=$2,completed_at=CASE WHEN $2='DELIVERED' THEN now() ELSE completed_at END,updated_at=now() WHERE id=$1`, id, status)
}
func scanDelivery(row interface{ Scan(dest ...any) error }) (*domain.Delivery, error) {
	d := domain.Delivery{}
	var lat, lon sql.NullFloat64
	if err := row.Scan(&d.ID, &d.BookingID, &d.DriverID, &d.Status, &lat, &lon, &d.StartedAt, &d.CompletedAt, &d.CreatedAt, &d.UpdatedAt); err != nil {
		return nil, mapError(err)
	}
	d.CurrentLatitude = floatPtr(lat)
	d.CurrentLongitude = floatPtr(lon)
	return &d, nil
}
