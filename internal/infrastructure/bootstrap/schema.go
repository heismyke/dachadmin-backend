package bootstrap

import (
	"context"
	"database/sql"
	"os"
)

func EnsureSchema(ctx context.Context, db *sql.DB) error {
	var exists bool
	if err := db.QueryRowContext(ctx, `SELECT to_regclass('public.team') IS NOT NULL`).Scan(&exists); err != nil {
		return err
	}
	if exists {
		return nil
	}

	schema, err := os.ReadFile("migrations/001_initial_schema.up.sql")
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, string(schema))
	return err
}
