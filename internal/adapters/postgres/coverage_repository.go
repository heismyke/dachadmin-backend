package postgres

import (
	"context"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

type CoverageRepository struct{ db *sql.DB }

func NewCoverageRepository(db *sql.DB) *CoverageRepository { return &CoverageRepository{db: db} }
func (r *CoverageRepository) Create(ctx context.Context, z *domain.CoverageZone) error {
	return mapError(r.db.QueryRowContext(ctx, `INSERT INTO coverage_zones (name,description,status) VALUES ($1,nullif($2,''),$3) RETURNING id,created_at,updated_at`, z.Name, z.Description, z.Status).Scan(&z.ID, &z.CreatedAt, &z.UpdatedAt))
}
func (r *CoverageRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.CoverageZone, error) {
	return scanCoverage(r.db.QueryRowContext(ctx, `SELECT id,name,coalesce(description,''),status,created_at,updated_at FROM coverage_zones WHERE id=$1`, id))
}
func (r *CoverageRepository) List(ctx context.Context, f repositories.CoverageFilter) (repositories.ListResult[domain.CoverageZone], error) {
	limit, offset := f.LimitOffset()
	where := "TRUE"
	args := []any{limit, offset}
	if f.Status != "" {
		where = "status=$3"
		args = append(args, f.Status)
	}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`SELECT id,name,coalesce(description,''),status,created_at,updated_at,count(*) OVER() FROM coverage_zones WHERE %s ORDER BY created_at DESC LIMIT $1 OFFSET $2`, where), args...)
	if err != nil {
		return repositories.ListResult[domain.CoverageZone]{}, mapError(err)
	}
	defer rows.Close()
	res := repositories.ListResult[domain.CoverageZone]{Items: []domain.CoverageZone{}}
	for rows.Next() {
		z := domain.CoverageZone{}
		if err := rows.Scan(&z.ID, &z.Name, &z.Description, &z.Status, &z.CreatedAt, &z.UpdatedAt, &res.Total); err != nil {
			return res, mapError(err)
		}
		res.Items = append(res.Items, z)
	}
	return res, mapError(rows.Err())
}
func (r *CoverageRepository) Update(ctx context.Context, z *domain.CoverageZone) error {
	return mapError(r.db.QueryRowContext(ctx, `UPDATE coverage_zones SET name=$2,description=nullif($3,''),status=$4,updated_at=now() WHERE id=$1 RETURNING updated_at`, z.ID, z.Name, z.Description, z.Status).Scan(&z.UpdatedAt))
}
func (r *CoverageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return execOne(ctx, r.db, `DELETE FROM coverage_zones WHERE id=$1`, id)
}
func scanCoverage(row interface{ Scan(dest ...any) error }) (*domain.CoverageZone, error) {
	z := domain.CoverageZone{}
	if err := row.Scan(&z.ID, &z.Name, &z.Description, &z.Status, &z.CreatedAt, &z.UpdatedAt); err != nil {
		return nil, mapError(err)
	}
	return &z, nil
}
