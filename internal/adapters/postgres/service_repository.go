package postgres

import (
	"context"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type ServiceRepository struct{ db *sql.DB }

func NewServiceRepository(db *sql.DB) *ServiceRepository { return &ServiceRepository{db: db} }
func (r *ServiceRepository) Create(ctx context.Context, s *domain.Service) error {
	return mapError(r.db.QueryRowContext(ctx, `INSERT INTO services (name, description, icon, status) VALUES ($1, nullif($2,''), nullif($3,''), $4) RETURNING id, created_at, updated_at`, s.Name, s.Description, s.Icon, s.Status).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt))
}
func (r *ServiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Service, error) {
	return scanService(r.db.QueryRowContext(ctx, `SELECT id, name, coalesce(description,''), coalesce(icon,''), status, created_at, updated_at FROM services WHERE id=$1`, id))
}
func (r *ServiceRepository) List(ctx context.Context, f repositories.ServiceFilter) (repositories.ListResult[domain.Service], error) {
	limit, offset := f.LimitOffset()
	where := "TRUE"
	args := []any{limit, offset}
	if f.Status != "" {
		where = "status = $3"
		args = append(args, f.Status)
	}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`SELECT id, name, coalesce(description,''), coalesce(icon,''), status, created_at, updated_at, count(*) OVER() FROM services WHERE %s ORDER BY created_at DESC LIMIT $1 OFFSET $2`, where), args...)
	if err != nil {
		return repositories.ListResult[domain.Service]{}, mapError(err)
	}
	defer rows.Close()
	res := repositories.ListResult[domain.Service]{Items: []domain.Service{}}
	for rows.Next() {
		s := domain.Service{}
		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.Icon, &s.Status, &s.CreatedAt, &s.UpdatedAt, &res.Total); err != nil {
			return res, mapError(err)
		}
		res.Items = append(res.Items, s)
	}
	return res, mapError(rows.Err())
}
func (r *ServiceRepository) Update(ctx context.Context, s *domain.Service) error {
	return mapError(r.db.QueryRowContext(ctx, `UPDATE services SET name=$2, description=nullif($3,''), icon=nullif($4,''), status=$5, updated_at=now() WHERE id=$1 RETURNING updated_at`, s.ID, s.Name, s.Description, s.Icon, s.Status).Scan(&s.UpdatedAt))
}
func (r *ServiceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return execOne(ctx, r.db, `DELETE FROM services WHERE id=$1`, id)
}
func (r *ServiceRepository) AddCoverage(ctx context.Context, serviceID uuid.UUID, coverageID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO service_coverage (service_id, coverage_zone_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, serviceID, coverageID)
	return mapError(err)
}
func (r *ServiceRepository) RemoveCoverage(ctx context.Context, serviceID uuid.UUID, coverageID uuid.UUID) error {
	return execOne(ctx, r.db, `DELETE FROM service_coverage WHERE service_id=$1 AND coverage_zone_id=$2`, serviceID, coverageID)
}
func scanService(row interface{ Scan(dest ...any) error }) (*domain.Service, error) {
	s := domain.Service{}
	if err := row.Scan(&s.ID, &s.Name, &s.Description, &s.Icon, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
		return nil, mapError(err)
	}
	return &s, nil
}
