package postgres

import (
	"context"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type CustomerRepository struct{ db *sql.DB }

func NewCustomerRepository(db *sql.DB) *CustomerRepository { return &CustomerRepository{db: db} }

func (r *CustomerRepository) Create(ctx context.Context, c *domain.Customer) error {
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO customers (full_name, email, phone, company_name, address, city, country)
		VALUES ($1, $2, nullif($3, ''), nullif($4, ''), nullif($5, ''), nullif($6, ''), nullif($7, ''))
		RETURNING id, created_at, updated_at`,
		c.FullName, c.Email, c.Phone, c.CompanyName, c.Address, c.City, c.Country)
	return mapError(row.Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt))
}

func (r *CustomerRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, full_name, email, coalesce(phone, ''), coalesce(company_name, ''), coalesce(address, ''), coalesce(city, ''), coalesce(country, ''), created_at, updated_at
		FROM customers
		WHERE id = $1`, id)
	c := domain.Customer{}
	if err := row.Scan(&c.ID, &c.FullName, &c.Email, &c.Phone, &c.CompanyName, &c.Address, &c.City, &c.Country, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, mapError(err)
	}
	return &c, nil
}

func (r *CustomerRepository) List(ctx context.Context, filter repositories.CustomerFilter) (repositories.ListResult[domain.Customer], error) {
	limit, offset := filter.LimitOffset()
	where := "TRUE"
	args := []any{limit, offset}
	if filter.Search != "" {
		where = "(full_name ILIKE $3 OR email ILIKE $3 OR company_name ILIKE $3)"
		args = append(args, "%"+filter.Search+"%")
	}
	query := fmt.Sprintf(`
		SELECT id, full_name, email, coalesce(phone, ''), coalesce(company_name, ''), coalesce(address, ''), coalesce(city, ''), coalesce(country, ''), created_at, updated_at, count(*) OVER()
		FROM customers
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`, where)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return repositories.ListResult[domain.Customer]{}, mapError(err)
	}
	defer rows.Close()
	result := repositories.ListResult[domain.Customer]{Items: []domain.Customer{}}
	for rows.Next() {
		c := domain.Customer{}
		if err := rows.Scan(&c.ID, &c.FullName, &c.Email, &c.Phone, &c.CompanyName, &c.Address, &c.City, &c.Country, &c.CreatedAt, &c.UpdatedAt, &result.Total); err != nil {
			return result, mapError(err)
		}
		result.Items = append(result.Items, c)
	}
	return result, mapError(rows.Err())
}

func (r *CustomerRepository) Update(ctx context.Context, c *domain.Customer) error {
	row := r.db.QueryRowContext(ctx, `
		UPDATE customers
		SET full_name = $2, email = $3, phone = nullif($4, ''), company_name = nullif($5, ''), address = nullif($6, ''), city = nullif($7, ''), country = nullif($8, ''), updated_at = now()
		WHERE id = $1
		RETURNING updated_at`,
		c.ID, c.FullName, c.Email, c.Phone, c.CompanyName, c.Address, c.City, c.Country)
	return mapError(row.Scan(&c.UpdatedAt))
}

func (r *CustomerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM customers WHERE id = $1`, id)
	if err != nil {
		return mapError(err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}
