package postgres

import (
	"context"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

type ReviewRepository struct{ db *sql.DB }

func NewReviewRepository(db *sql.DB) *ReviewRepository { return &ReviewRepository{db: db} }
func (r *ReviewRepository) Create(ctx context.Context, v *domain.Review) error {
	return mapError(r.db.QueryRowContext(ctx, `INSERT INTO reviews (customer_id,booking_id,rating,comment,status) VALUES ($1,$2,$3,nullif($4,''),$5) RETURNING id,created_at,updated_at`, v.CustomerID, v.BookingID, v.Rating, v.Comment, v.Status).Scan(&v.ID, &v.CreatedAt, &v.UpdatedAt))
}
func (r *ReviewRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Review, error) {
	return scanReview(r.db.QueryRowContext(ctx, `SELECT id,customer_id,booking_id,rating,coalesce(comment,''),status,created_at,updated_at FROM reviews WHERE id=$1`, id))
}
func (r *ReviewRepository) List(ctx context.Context, f repositories.ReviewFilter) (repositories.ListResult[domain.Review], error) {
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
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`SELECT id,customer_id,booking_id,rating,coalesce(comment,''),status,created_at,updated_at,count(*) OVER() FROM reviews WHERE %s ORDER BY created_at DESC LIMIT $1 OFFSET $2`, where), args...)
	if err != nil {
		return repositories.ListResult[domain.Review]{}, mapError(err)
	}
	defer rows.Close()
	res := repositories.ListResult[domain.Review]{Items: []domain.Review{}}
	for rows.Next() {
		v := domain.Review{}
		if err := rows.Scan(&v.ID, &v.CustomerID, &v.BookingID, &v.Rating, &v.Comment, &v.Status, &v.CreatedAt, &v.UpdatedAt, &res.Total); err != nil {
			return res, mapError(err)
		}
		res.Items = append(res.Items, v)
	}
	return res, mapError(rows.Err())
}
func (r *ReviewRepository) Update(ctx context.Context, v *domain.Review) error {
	return mapError(r.db.QueryRowContext(ctx, `UPDATE reviews SET customer_id=$2,booking_id=$3,rating=$4,comment=nullif($5,''),status=$6,updated_at=now() WHERE id=$1 RETURNING updated_at`, v.ID, v.CustomerID, v.BookingID, v.Rating, v.Comment, v.Status).Scan(&v.UpdatedAt))
}
func (r *ReviewRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return execOne(ctx, r.db, `DELETE FROM reviews WHERE id=$1`, id)
}
func (r *ReviewRepository) Publish(ctx context.Context, id uuid.UUID, actorID *uuid.UUID) error {
	return execOne(ctx, r.db, `UPDATE reviews SET status=$2,updated_at=now() WHERE id=$1`, id, domain.ReviewStatusPublished)
}
func scanReview(row interface{ Scan(dest ...any) error }) (*domain.Review, error) {
	v := domain.Review{}
	if err := row.Scan(&v.ID, &v.CustomerID, &v.BookingID, &v.Rating, &v.Comment, &v.Status, &v.CreatedAt, &v.UpdatedAt); err != nil {
		return nil, mapError(err)
	}
	return &v, nil
}
