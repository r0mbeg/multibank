// internal/storage/sqlite/sqlite.go

package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {

	const op = "storage.sqlite.New"

	if storagePath == "" {
		return nil, fmt.Errorf("%s: empty storage path", op)
	}

	// connection params (timeout, wal mode, foreign keys)
	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(DELETE)&_pragma=foreign_keys(ON)", storagePath)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: open: %w", op, err)
	}

	// single writer
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	// connection check (creates database)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("%s: ping: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error { return s.db.Close() }
func (s *Storage) DB() *sql.DB  { return s.db }

// Migrate - idempotent schema initialization/migration.
func (s *Storage) Migrate(ctx context.Context) error {
	// TODO add op
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// versioning
	if _, err = tx.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY
);`); err != nil {
		return err
	}

	// zero migration (idempotent IF NOT EXISTS)
	// users
	if _, err = tx.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS users (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    email         TEXT    NOT NULL UNIQUE,
    first_name    TEXT    NOT NULL,
    last_name     TEXT    NOT NULL,
    patronymic    TEXT    NOT NULL,
    birthdate     TEXT    NOT NULL CHECK (length(birthdate) = 10), -- YYYY-MM-DD
    password_hash TEXT    NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    created_at    TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at    TEXT    NOT NULL DEFAULT (datetime('now'))
);
`); err != nil {
		return err
	}

	// banks
	if _, err = tx.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS banks(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    name         TEXT    NOT NULL,
    code         TEXT    NOT NULL UNIQUE,
    api_base_url TEXT    NOT NULL,
    login        TEXT    NOT NULL,  -- client_id
    password     TEXT    NOT NULL,  -- client_secret
    is_enabled   INTEGER NOT NULL DEFAULT 1,
    created_at   TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at   TEXT    NOT NULL DEFAULT (datetime('now'))
);
`); err != nil {
		return err
	}

	// bank tokens
	if _, err = tx.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS bank_tokens(
    bank_id      INTEGER NOT NULL UNIQUE,
    access_token TEXT    NOT NULL,
    expires_at   TEXT    NOT NULL, -- RFC3339 или datetime('...')
    created_at   TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at   TEXT    NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY(bank_id) REFERENCES banks(id) ON DELETE CASCADE
);
`); err != nil {
		return err
	}

	// consents
	if _, err = tx.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS bank_tokens(
id         INTEGER  NOT NULL,
bank_id    INTEGER  NOT NULL,
created_at TEXT     NOT NULL,
expiration_date_time TEXT NOT NULL,
status     TEXT     NOT NULL
FOREIGN KEY (bank_id) REFERENCES banks(id) ON DELETE CASCADE,
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
`); err != nil {
		return err
	}

	return tx.Commit()
}
