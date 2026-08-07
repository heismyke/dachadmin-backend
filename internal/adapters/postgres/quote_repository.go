package postgres

import (
	"context"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

type QuoteRepository struct{ db *sql.DB }

func NewQuoteRepository(db *sql.DB) *QuoteRepository { return &QuoteRepository{db: db} }
func (r *QuoteRepository) Create(ctx context.Context, q *domain.Quote) error {
	return mapError(r.db.QueryRowContext(ctx, `INSERT INTO quotes (customer_id,service_id,pickup_address,dropoff_address,status,valid_until,total_price,notes) VALUES ($1,$2,$3,$4,$5,$6,$7,nullif($8,'')) RETURNING id,created_at,updated_at`, q.CustomerID, q.ServiceID, q.PickupAddress, q.DropoffAddress, q.Status, q.ValidUntil, q.TotalPrice, q.Notes).Scan(&q.ID, &q.CreatedAt, &q.UpdatedAt))
}
func (r *QuoteRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Quote, error) {
	return scanQuote(r.db.QueryRowContext(ctx, `SELECT id,customer_id,service_id,pickup_address,dropoff_address,status,valid_until,total_price::float8,coalesce(notes,''),created_at,updated_at FROM quotes WHERE id=$1`, id))
}
func (r *QuoteRepository) List(ctx context.Context, f repositories.QuoteFilter) (repositories.ListResult[domain.Quote], error) {
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
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`SELECT id,customer_id,service_id,pickup_address,dropoff_address,status,valid_until,total_price::float8,coalesce(notes,''),created_at,updated_at,count(*) OVER() FROM quotes WHERE %s ORDER BY created_at DESC LIMIT $1 OFFSET $2`, where), args...)
	if err != nil {
		return repositories.ListResult[domain.Quote]{}, mapError(err)
	}
	defer rows.Close()
	res := repositories.ListResult[domain.Quote]{Items: []domain.Quote{}}
	for rows.Next() {
		q := domain.Quote{}
		var total sql.NullFloat64
		if err := rows.Scan(&q.ID, &q.CustomerID, &q.ServiceID, &q.PickupAddress, &q.DropoffAddress, &q.Status, &q.ValidUntil, &total, &q.Notes, &q.CreatedAt, &q.UpdatedAt, &res.Total); err != nil {
			return res, mapError(err)
		}
		q.TotalPrice = floatPtr(total)
		res.Items = append(res.Items, q)
	}
	return res, mapError(rows.Err())
}
func (r *QuoteRepository) Update(ctx context.Context, q *domain.Quote) error {
	return mapError(r.db.QueryRowContext(ctx, `UPDATE quotes SET customer_id=$2,service_id=$3,pickup_address=$4,dropoff_address=$5,status=$6,valid_until=$7,total_price=$8,notes=nullif($9,''),updated_at=now() WHERE id=$1 RETURNING updated_at`, q.ID, q.CustomerID, q.ServiceID, q.PickupAddress, q.DropoffAddress, q.Status, q.ValidUntil, q.TotalPrice, q.Notes).Scan(&q.UpdatedAt))
}
func (r *QuoteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return execOne(ctx, r.db, `DELETE FROM quotes WHERE id=$1`, id)
}
func scanQuote(row interface{ Scan(dest ...any) error }) (*domain.Quote, error) {
	q := domain.Quote{}
	var total sql.NullFloat64
	if err := row.Scan(&q.ID, &q.CustomerID, &q.ServiceID, &q.PickupAddress, &q.DropoffAddress, &q.Status, &q.ValidUntil, &total, &q.Notes, &q.CreatedAt, &q.UpdatedAt); err != nil {
		return nil, mapError(err)
	}
	q.TotalPrice = floatPtr(total)
	return &q, nil
}
