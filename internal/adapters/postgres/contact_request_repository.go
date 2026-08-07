package postgres

import (
	"context"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

type ContactRequestRepository struct{ db *sql.DB }

func NewContactRequestRepository(db *sql.DB) *ContactRequestRepository {
	return &ContactRequestRepository{db: db}
}
func (r *ContactRequestRepository) Create(ctx context.Context, c *domain.ContactRequest) error {
	return mapError(r.db.QueryRowContext(ctx, `INSERT INTO contact_requests (full_name,email,phone,subject,message,status) VALUES ($1,$2,nullif($3,''),nullif($4,''),$5,$6) RETURNING id,created_at,updated_at`, c.FullName, c.Email, c.Phone, c.Subject, c.Message, c.Status).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt))
}
func (r *ContactRequestRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ContactRequest, error) {
	return scanContact(r.db.QueryRowContext(ctx, `SELECT id,full_name,email,coalesce(phone,''),coalesce(subject,''),message,status,created_at,updated_at FROM contact_requests WHERE id=$1`, id))
}
func (r *ContactRequestRepository) List(ctx context.Context, f repositories.ContactRequestFilter) (repositories.ListResult[domain.ContactRequest], error) {
	limit, offset := f.LimitOffset()
	where := "TRUE"
	args := []any{limit, offset}
	if f.Status != "" {
		where = "status=$3"
		args = append(args, f.Status)
	}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`SELECT id,full_name,email,coalesce(phone,''),coalesce(subject,''),message,status,created_at,updated_at,count(*) OVER() FROM contact_requests WHERE %s ORDER BY created_at DESC LIMIT $1 OFFSET $2`, where), args...)
	if err != nil {
		return repositories.ListResult[domain.ContactRequest]{}, mapError(err)
	}
	defer rows.Close()
	res := repositories.ListResult[domain.ContactRequest]{Items: []domain.ContactRequest{}}
	for rows.Next() {
		c := domain.ContactRequest{}
		if err := rows.Scan(&c.ID, &c.FullName, &c.Email, &c.Phone, &c.Subject, &c.Message, &c.Status, &c.CreatedAt, &c.UpdatedAt, &res.Total); err != nil {
			return res, mapError(err)
		}
		res.Items = append(res.Items, c)
	}
	return res, mapError(rows.Err())
}
func (r *ContactRequestRepository) Update(ctx context.Context, c *domain.ContactRequest) error {
	return mapError(r.db.QueryRowContext(ctx, `UPDATE contact_requests SET full_name=$2,email=$3,phone=nullif($4,''),subject=nullif($5,''),message=$6,status=$7,updated_at=now() WHERE id=$1 RETURNING updated_at`, c.ID, c.FullName, c.Email, c.Phone, c.Subject, c.Message, c.Status).Scan(&c.UpdatedAt))
}
func (r *ContactRequestRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return execOne(ctx, r.db, `DELETE FROM contact_requests WHERE id=$1`, id)
}
func (r *ContactRequestRepository) Resolve(ctx context.Context, id uuid.UUID, actorID *uuid.UUID) error {
	return execOne(ctx, r.db, `UPDATE contact_requests SET status=$2,updated_at=now() WHERE id=$1`, id, domain.ContactRequestStatusResolved)
}
func scanContact(row interface{ Scan(dest ...any) error }) (*domain.ContactRequest, error) {
	c := domain.ContactRequest{}
	if err := row.Scan(&c.ID, &c.FullName, &c.Email, &c.Phone, &c.Subject, &c.Message, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, mapError(err)
	}
	return &c, nil
}
