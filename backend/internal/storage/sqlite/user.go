// internal/storage/sqlite/user.go
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"multibank/backend/internal/domain"
	"multibank/backend/internal/storage"
	sqliteutils "multibank/backend/internal/storage/sqlite/utils"

	"modernc.org/sqlite"               // type of Error
	sqlitelib "modernc.org/sqlite/lib" // code constants
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) Create(ctx context.Context, u domain.User) (int64, error) {
	const op = "storage.sqlite.user.Create"

	res, err := r.db.ExecContext(ctx, `
INSERT INTO users (email, first_name, last_name, patronymic, birthdate, password_hash, is_admin, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?,?, datetime('now'), datetime('now'))`,
		u.Email, u.FirstName, u.LastName, u.Patronymic, u.BirthDate, u.PasswordHash, u.IsAdmin,
	)

	if err != nil {
		var se *sqlite.Error
		if errors.As(err, &se) {
			if se.Code() == sqlitelib.SQLITE_CONSTRAINT_UNIQUE {
				return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (domain.User, error) {

	const op = "storage.sqlite.user.GetByID"

	var u domain.User
	var created, updated string
	err := r.db.QueryRowContext(ctx, `
SELECT id, email, first_name, last_name, patronymic, birthdate, password_hash, is_admin, created_at, updated_at
FROM users WHERE id = ?`, id).Scan(
		&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Patronymic, &u.BirthDate, &u.PasswordHash, &u.IsAdmin, &created, &updated,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, fmt.Errorf("%s : %w", op, storage.ErrUserNotFound)
		}

		return domain.User{}, fmt.Errorf("%s : %w", op, err)
	}

	if t, err := sqliteutils.ParseTS(created); err == nil {
		u.CreatedAt = t
	}
	if t, err := sqliteutils.ParseTS(updated); err == nil {
		u.UpdatedAt = t
	}
	return u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {

	const op = "storage.sqlite.user.GetByEmail"

	var u domain.User
	var created, updated string
	err := r.db.QueryRowContext(ctx, `
SELECT id, email, first_name, last_name, patronymic, birthdate, password_hash, is_admin, created_at, updated_at
FROM users WHERE email = ?`, email).Scan(
		&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Patronymic, &u.BirthDate, &u.PasswordHash, &u.IsAdmin, &created, &updated,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, fmt.Errorf("%s : %w", op, storage.ErrUserNotFound)
		}

		return domain.User{}, fmt.Errorf("%s : %w", op, err)
	}
	if t, err := sqliteutils.ParseTS(created); err == nil {
		u.CreatedAt = t
	}
	if t, err := sqliteutils.ParseTS(updated); err == nil {
		u.UpdatedAt = t
	}
	return u, nil
}
