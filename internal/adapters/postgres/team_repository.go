package postgres

import (
	"context"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

type TeamRepository struct{ db *sql.DB }

func NewTeamRepository(db *sql.DB) *TeamRepository { return &TeamRepository{db: db} }
func (r *TeamRepository) Create(ctx context.Context, m *domain.TeamMember) error {
	return mapError(r.db.QueryRowContext(ctx, `INSERT INTO team (name,email,password_hash,role,status) VALUES ($1,$2,$3,$4,$5) RETURNING id,created_at,updated_at`, m.Name, m.Email, m.PasswordHash, m.Role, m.Status).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt))
}
func (r *TeamRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.TeamMember, error) {
	return scanTeam(r.db.QueryRowContext(ctx, `SELECT id,name,email,password_hash,role,status,created_at,updated_at FROM team WHERE id=$1`, id))
}
func (r *TeamRepository) GetByEmail(ctx context.Context, email string) (*domain.TeamMember, error) {
	return scanTeam(r.db.QueryRowContext(ctx, `SELECT id,name,email,password_hash,role,status,created_at,updated_at FROM team WHERE email=$1`, email))
}
func (r *TeamRepository) List(ctx context.Context, f repositories.TeamFilter) (repositories.ListResult[domain.TeamMember], error) {
	limit, offset := f.LimitOffset()
	where := "TRUE"
	args := []any{limit, offset}
	if f.Role != "" {
		where += fmt.Sprintf(" AND role=$%d", len(args)+1)
		args = append(args, f.Role)
	}
	if f.Status != "" {
		where += fmt.Sprintf(" AND status=$%d", len(args)+1)
		args = append(args, f.Status)
	}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`SELECT id,name,email,password_hash,role,status,created_at,updated_at,count(*) OVER() FROM team WHERE %s ORDER BY created_at DESC LIMIT $1 OFFSET $2`, where), args...)
	if err != nil {
		return repositories.ListResult[domain.TeamMember]{}, mapError(err)
	}
	defer rows.Close()
	res := repositories.ListResult[domain.TeamMember]{Items: []domain.TeamMember{}}
	for rows.Next() {
		m := domain.TeamMember{}
		if err := rows.Scan(&m.ID, &m.Name, &m.Email, &m.PasswordHash, &m.Role, &m.Status, &m.CreatedAt, &m.UpdatedAt, &res.Total); err != nil {
			return res, mapError(err)
		}
		m.PasswordHash = ""
		res.Items = append(res.Items, m)
	}
	return res, mapError(rows.Err())
}
func (r *TeamRepository) Update(ctx context.Context, m *domain.TeamMember) error {
	return mapError(r.db.QueryRowContext(ctx, `UPDATE team SET name=$2,email=$3,role=$4,status=$5,updated_at=now() WHERE id=$1 RETURNING updated_at`, m.ID, m.Name, m.Email, m.Role, m.Status).Scan(&m.UpdatedAt))
}
func (r *TeamRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return execOne(ctx, r.db, `DELETE FROM team WHERE id=$1`, id)
}
func (r *TeamRepository) Deactivate(ctx context.Context, id uuid.UUID, actorID *uuid.UUID) error {
	return execOne(ctx, r.db, `UPDATE team SET status=$2,updated_at=now() WHERE id=$1`, id, domain.TeamStatusDisabled)
}
func scanTeam(row interface{ Scan(dest ...any) error }) (*domain.TeamMember, error) {
	m := domain.TeamMember{}
	if err := row.Scan(&m.ID, &m.Name, &m.Email, &m.PasswordHash, &m.Role, &m.Status, &m.CreatedAt, &m.UpdatedAt); err != nil {
		return nil, mapError(err)
	}
	return &m, nil
}
