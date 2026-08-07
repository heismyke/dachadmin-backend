package postgres

import (
	"context"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type AuditLogRepository struct{ db *sql.DB }

func NewAuditLogRepository(db *sql.DB) *AuditLogRepository { return &AuditLogRepository{db: db} }

func (r *AuditLogRepository) Log(ctx context.Context, entry domain.AuditLog) error {
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO audit_logs (user_id, action, entity, entity_id, changes, ip_address)
		VALUES ($1, $2, $3, $4, $5, nullif($6, ''))
		RETURNING id, created_at`,
		entry.UserID, entry.Action, entry.Entity, entry.EntityID, nullBytes(entry.Changes), entry.IPAddress)
	return mapError(row.Scan(&entry.ID, &entry.CreatedAt))
}

func (r *AuditLogRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.AuditLog, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, action, entity, entity_id, coalesce(changes, '{}'::jsonb), coalesce(ip_address, ''), created_at
		FROM audit_logs
		WHERE id = $1`, id)
	return scanAudit(row)
}

func (r *AuditLogRepository) List(ctx context.Context, filter repositories.AuditLogFilter) (repositories.ListResult[domain.AuditLog], error) {
	limit, offset := filter.LimitOffset()
	where := "TRUE"
	args := []any{limit, offset}
	if filter.Entity != "" {
		where += fmt.Sprintf(" AND entity = $%d", len(args)+1)
		args = append(args, filter.Entity)
	}
	if filter.EntityID != nil {
		where += fmt.Sprintf(" AND entity_id = $%d", len(args)+1)
		args = append(args, *filter.EntityID)
	}
	if filter.UserID != nil {
		where += fmt.Sprintf(" AND user_id = $%d", len(args)+1)
		args = append(args, *filter.UserID)
	}
	query := fmt.Sprintf(`
		SELECT id, user_id, action, entity, entity_id, coalesce(changes, '{}'::jsonb), coalesce(ip_address, ''), created_at, count(*) OVER()
		FROM audit_logs
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`, where)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return repositories.ListResult[domain.AuditLog]{}, mapError(err)
	}
	defer rows.Close()
	result := repositories.ListResult[domain.AuditLog]{Items: []domain.AuditLog{}}
	for rows.Next() {
		a := domain.AuditLog{}
		if err := rows.Scan(&a.ID, &a.UserID, &a.Action, &a.Entity, &a.EntityID, &a.Changes, &a.IPAddress, &a.CreatedAt, &result.Total); err != nil {
			return result, mapError(err)
		}
		result.Items = append(result.Items, a)
	}
	return result, mapError(rows.Err())
}

func scanAudit(row interface{ Scan(dest ...any) error }) (*domain.AuditLog, error) {
	a := domain.AuditLog{}
	if err := row.Scan(&a.ID, &a.UserID, &a.Action, &a.Entity, &a.EntityID, &a.Changes, &a.IPAddress, &a.CreatedAt); err != nil {
		return nil, mapError(err)
	}
	return &a, nil
}

func nullBytes(b []byte) any {
	if len(b) == 0 {
		return nil
	}
	return b
}
