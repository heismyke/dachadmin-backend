package postgres

import (
	"context"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type DriverRepository struct{ db *sql.DB }

func NewDriverRepository(db *sql.DB) *DriverRepository { return &DriverRepository{db: db} }
func (r *DriverRepository) Create(ctx context.Context, d *domain.Driver) error {
	return mapError(r.db.QueryRowContext(ctx, `INSERT INTO drivers (full_name,email,phone,license_number,vehicle_type,vehicle_plate,status) VALUES ($1,nullif($2,''),nullif($3,''),$4,nullif($5,''),nullif($6,''),$7) RETURNING id, created_at, updated_at`, d.FullName, d.Email, d.Phone, d.LicenseNumber, d.VehicleType, d.VehiclePlate, d.Status).Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt))
}
func (r *DriverRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Driver, error) {
	return scanDriver(r.db.QueryRowContext(ctx, `SELECT id, full_name, coalesce(email,''), coalesce(phone,''), license_number, coalesce(vehicle_type,''), coalesce(vehicle_plate,''), status, created_at, updated_at FROM drivers WHERE id=$1`, id))
}
func (r *DriverRepository) List(ctx context.Context, f repositories.DriverFilter) (repositories.ListResult[domain.Driver], error) {
	limit, offset := f.LimitOffset()
	where := "TRUE"
	args := []any{limit, offset}
	if f.Status != "" {
		where = "status=$3"
		args = append(args, f.Status)
	}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`SELECT id, full_name, coalesce(email,''), coalesce(phone,''), license_number, coalesce(vehicle_type,''), coalesce(vehicle_plate,''), status, created_at, updated_at, count(*) OVER() FROM drivers WHERE %s ORDER BY created_at DESC LIMIT $1 OFFSET $2`, where), args...)
	if err != nil {
		return repositories.ListResult[domain.Driver]{}, mapError(err)
	}
	defer rows.Close()
	res := repositories.ListResult[domain.Driver]{Items: []domain.Driver{}}
	for rows.Next() {
		d := domain.Driver{}
		if err := rows.Scan(&d.ID, &d.FullName, &d.Email, &d.Phone, &d.LicenseNumber, &d.VehicleType, &d.VehiclePlate, &d.Status, &d.CreatedAt, &d.UpdatedAt, &res.Total); err != nil {
			return res, mapError(err)
		}
		res.Items = append(res.Items, d)
	}
	return res, mapError(rows.Err())
}
func (r *DriverRepository) Update(ctx context.Context, d *domain.Driver) error {
	return mapError(r.db.QueryRowContext(ctx, `UPDATE drivers SET full_name=$2,email=nullif($3,''),phone=nullif($4,''),license_number=$5,vehicle_type=nullif($6,''),vehicle_plate=nullif($7,''),status=$8,updated_at=now() WHERE id=$1 RETURNING updated_at`, d.ID, d.FullName, d.Email, d.Phone, d.LicenseNumber, d.VehicleType, d.VehiclePlate, d.Status).Scan(&d.UpdatedAt))
}
func (r *DriverRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return execOne(ctx, r.db, `DELETE FROM drivers WHERE id=$1`, id)
}
func (r *DriverRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.DriverStatus, actorID *uuid.UUID) error {
	return execOne(ctx, r.db, `UPDATE drivers SET status=$2, updated_at=now() WHERE id=$1`, id, status)
}
func scanDriver(row interface{ Scan(dest ...any) error }) (*domain.Driver, error) {
	d := domain.Driver{}
	if err := row.Scan(&d.ID, &d.FullName, &d.Email, &d.Phone, &d.LicenseNumber, &d.VehicleType, &d.VehiclePlate, &d.Status, &d.CreatedAt, &d.UpdatedAt); err != nil {
		return nil, mapError(err)
	}
	return &d, nil
}
