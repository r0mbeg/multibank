package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"multibank/backend/internal/domain"
	"time"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo { return &UserRepo{db: db} }

const tsLayout = "2006-01-02 15:04:05"

// helpers for date parsing
func parseTS(s string) (time.Time, error) { return time.Parse(tsLayout, s) }
func nowUTC() string                      { return time.Now().UTC().Format(tsLayout) }

func (r *UserRepo) Create(ctx context.Context, u domain.User) (int64, error) {
	res, err := r.db.ExecContext(ctx, `
INSERT INTO users (email, first_name, last_name, patronymic, birthdate, password_hash, created_at, updated_at)
VALUES (?,     ?,          ?,        ?,          ?,         ?,             datetime('now'), datetime('now'))
`, u.Email, u.FirstName, u.LastName, u.Patronymic, u.BirthDate, u.PasswordHash)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (domain.User, error) {
	var u domain.User
	var created, updated string
	err := r.db.QueryRowContext(ctx, `
SELECT id, email, first_name, last_name, patronymic, birthdate, password_hash, created_at, updated_at
FROM users WHERE id = ?`, id).Scan(
		&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Patronymic, &u.BirthDate, &u.PasswordHash, &created, &updated,
	)
	if err != nil {
		return domain.User{}, err
	}
	if t, err := parseTS(created); err == nil {
		u.CreatedAt = t
	}
	if t, err := parseTS(updated); err == nil {
		u.UpdatedAt = t
	}
	return u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var u domain.User
	var created, updated string
	err := r.db.QueryRowContext(ctx, `
SELECT id, email, first_name, last_name, patronymic, birthdate, password_hash, created_at, updated_at
FROM users WHERE email = ?`, email).Scan(
		&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Patronymic, &u.BirthDate, &u.PasswordHash, &created, &updated,
	)
	if err != nil {
		return domain.User{}, err
	}
	if t, err := parseTS(created); err == nil {
		u.CreatedAt = t
	}
	if t, err := parseTS(updated); err == nil {
		u.UpdatedAt = t
	}
	return u, nil
}

func (r *UserRepo) UpdateNames(ctx context.Context, id int64, first, last, patr string) error {
	res, err := r.db.ExecContext(ctx, `
UPDATE users
SET first_name = ?, last_name = ?, patronymic = ?, updated_at = datetime('now')
WHERE id = ?`, first, last, patr, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("not found")
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("not found")
	}
	return nil
}

func (r *UserRepo) List(ctx context.Context, limit, offset int) ([]domain.User, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT id, email, first_name, last_name, patronymic, birthdate, password_hash, created_at, updated_at
FROM users
ORDER BY id
LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.User
	for rows.Next() {
		var u domain.User
		var created, updated string
		if err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Patronymic, &u.BirthDate, &u.PasswordHash, &created, &updated); err != nil {
			return nil, err
		}
		if t, err := parseTS(created); err == nil {
			u.CreatedAt = t
		}
		if t, err := parseTS(updated); err == nil {
			u.UpdatedAt = t
		}
		out = append(out, u)
	}
	return out, rows.Err()
}
