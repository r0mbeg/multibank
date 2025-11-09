// internal/storage/sqlite/sqlite.go

package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

	// 1) make sure that the parent directory exists
	dir := filepath.Dir(storagePath)
	if dir == "." || dir == "/" {
		// the relative file in the current directory is ok
	} else {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("%s: mkdir %q: %w", op, dir, err)
		}
	}

	// 2) If the storagePath points to a directory => error
	if fi, err := os.Stat(storagePath); err == nil && fi.IsDir() {
		return nil, fmt.Errorf("%s: path %q is a directory, expected file", op, storagePath)
	}

	// 3) touch DB file, to get an understandable rights/path error before opening sqlite
	// O_CREATE will not recreate an existing file
	f, err := os.OpenFile(storagePath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("%s: create/open %q: %w", op, storagePath, err)
	}
	_ = f.Close()

	// 4) connection string (busy_timeout, foreign_keys и т.п.)
	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(DELETE)&_pragma=foreign_keys(ON)", storagePath)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: open: %w", op, err)
	}

	// one writer to DB
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(0)

	// 5) conn check
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("%s: ping: %w", op, translateSQLiteOpenErr(err))
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error { return s.db.Close() }
func (s *Storage) DB() *sql.DB  { return s.db }

// translateSQLiteOpenErr makes the message clearer for the most common cases
func translateSQLiteOpenErr(err error) error {
	// modernc.org/sqlite usually returns "unable to open database file: out of memory (14)"
	// but in fact the problem is about the dir
	if strings.Contains(err.Error(), "unable to open database file") {
		return fmt.Errorf("unable to open database file (check the existence of the directory and write permissions): %w", err)
	}
	return err
}

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

	// recommended products
	if _, err = tx.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS recommended_products (
    product_id   TEXT NOT NULL,
    bank_code    TEXT NOT NULL,
    product_type TEXT NOT NULL,
    created_at   TEXT NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (product_id, bank_code, product_type)
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
CREATE TABLE IF NOT EXISTS account_consents (
  id                     INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id                INTEGER NOT NULL,
  bank_id                INTEGER NOT NULL,
  request_id             TEXT    NOT NULL UNIQUE,
  consent_id             TEXT    UNIQUE,
  status                 TEXT    NOT NULL,     -- Authorised | AwaitingAuthorisation | Rejected | Revoked
  auto_approved          INTEGER,              -- NULL | 0 | 1
  permissions_json       TEXT    NOT NULL,     -- JSON-массив строк

  reason                 TEXT    NOT NULL,
  requesting_bank        TEXT    NOT NULL,
  requesting_bank_name   TEXT    NOT NULL,

  creation_datetime      TEXT,
  status_update_datetime TEXT,
  expiration_datetime    TEXT,

  client_id              TEXT    NOT NULL, -- e.g. team014-1

  created_at             TEXT    NOT NULL DEFAULT (datetime('now')),
  updated_at             TEXT    NOT NULL DEFAULT (datetime('now')),

  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (bank_id) REFERENCES banks(id)
);
`); err != nil {
		return err
	}

	// view with bank code
	if _, err = tx.ExecContext(ctx, `
CREATE VIEW IF NOT EXISTS account_consents_view AS
SELECT
  c.*,
  b.code AS bank_code
FROM account_consents c
JOIN banks b ON b.id = c.bank_id;
`); err != nil {
		return err
	}

	return tx.Commit()
}
