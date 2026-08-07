package bootstrap

import (
	"context"
	"dach-admin/internal/domain"
	"dach-admin/internal/infrastructure/config"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func EnsureAdmin(ctx context.Context, db *sql.DB, cfg config.BootstrapAdminConfig) error {
	if cfg.Email == "" && cfg.Password == "" {
		return nil
	}
	if cfg.Email == "" || cfg.Password == "" {
		return errors.New("both BOOTSTRAP_ADMIN_EMAIL and BOOTSTRAP_ADMIN_PASSWORD are required")
	}
	if len(cfg.Password) < 8 {
		return errors.New("BOOTSTRAP_ADMIN_PASSWORD must be at least 8 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, `
		INSERT INTO team (name, email, password_hash, role, status)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (email) DO UPDATE
		SET name = EXCLUDED.name,
			password_hash = EXCLUDED.password_hash,
			role = EXCLUDED.role,
			status = EXCLUDED.status,
			updated_at = now()`,
		cfg.Name,
		cfg.Email,
		string(hash),
		domain.RoleSuperAdmin,
		domain.TeamStatusActive,
	)
	return err
}
