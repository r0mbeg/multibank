// internal/storage/sqlite/bank.go

package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"multibank/backend/internal/domain"
	"multibank/backend/internal/storage"
	sqliteutils "multibank/backend/internal/storage/sqlite/utils"
	"time"
)

type BankRepo struct {
	db *sql.DB
}

func NewBankRepo(db *sql.DB) *BankRepo {
	return &BankRepo{db: db}
}

func (s *BankRepo) ListEnabledBanks(ctx context.Context) ([]domain.Bank, error) {
	const op = "storage.sqlite.bank.ListEnabledBanks"

	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, code, api_base_url, login, password, is_enabled, created_at, updated_at
		FROM banks
		WHERE is_enabled = 1
		ORDER BY name`)
	if err != nil {
		return []domain.Bank{}, fmt.Errorf("%s : %w", op, err)
	}
	defer rows.Close()

	var out []domain.Bank
	for rows.Next() {
		var b domain.Bank
		var created, updated string
		var en int
		err := rows.Scan(&b.ID, &b.Name, &b.Code, &b.APIBaseURL, &b.Login, &b.Password, &en, &created, &updated)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return []domain.Bank{}, fmt.Errorf("%s : %w", op, storage.ErrBanksNotFound)
			}

			return []domain.Bank{}, fmt.Errorf("%s : %w", op, err)
		}

		b.IsEnabled = en == 1

		if t, err := sqliteutils.ParseTS(created); err == nil {
			b.CreatedAt = t
		}

		if t, err := sqliteutils.ParseTS(updated); err == nil {
			b.UpdatedAt = t
		}

		out = append(out, b)
	}
	return out, rows.Err()
}

func (s *BankRepo) GetBankByID(ctx context.Context, id int64) (domain.Bank, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, name, code, api_base_url, login, password, is_enabled, created_at, updated_at
		FROM banks WHERE id = ?`, id)

	var b domain.Bank
	var en int
	var created, updated string

	if err := row.Scan(&b.ID, &b.Name, &b.Code, &b.APIBaseURL, &b.Login, &b.Password, &en, &created, &updated); err != nil {
		return domain.Bank{}, err
	}
	b.IsEnabled = en == 1

	if t, err := sqliteutils.ParseTS(created); err == nil {
		b.CreatedAt = t
	}
	if t, err := sqliteutils.ParseTS(updated); err == nil {
		b.UpdatedAt = t
	}
	return b, nil
}

func (s *BankRepo) GetBankByCode(ctx context.Context, code string) (domain.Bank, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, name, code, api_base_url, login, password, is_enabled, created_at, updated_at
		FROM banks WHERE code = ?`, code)

	var b domain.Bank
	var en int
	var created, updated string

	if err := row.Scan(&b.ID, &b.Name, &b.Code, &b.APIBaseURL, &b.Login, &b.Password, &en, &created, &updated); err != nil {
		return domain.Bank{}, err
	}
	b.IsEnabled = en == 1

	if t, err := sqliteutils.ParseTS(created); err == nil {
		b.CreatedAt = t
	}
	if t, err := sqliteutils.ParseTS(updated); err == nil {
		b.UpdatedAt = t
	}
	return b, nil
}

func (s *BankRepo) UpsertBankToken(ctx context.Context, t domain.BankToken) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO bank_tokens(bank_id, access_token, expires_at, updated_at)
		VALUES(?,?,?,?)
		ON CONFLICT(bank_id) DO UPDATE SET
		  access_token = excluded.access_token,
		  expires_at   = excluded.expires_at,
		  updated_at   = excluded.updated_at`,
		t.BankID,
		t.AccessToken,
		t.ExpiresAt.UTC().Format(sqliteutils.TsLayout),
		time.Now().UTC().Format(sqliteutils.TsLayout),
	)
	return err
}

func (s *BankRepo) GetBankToken(ctx context.Context, bankID int64) (domain.BankToken, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT bank_id, access_token, expires_at, created_at, updated_at
		FROM bank_tokens WHERE bank_id = ?`, bankID)

	var t domain.BankToken
	var expires, created, updated string

	if err := row.Scan(&t.BankID, &t.AccessToken, &expires, &created, &updated); err != nil {
		return domain.BankToken{}, err
	}
	if v, err := sqliteutils.ParseTS(expires); err == nil {
		t.ExpiresAt = v
	}
	if v, err := sqliteutils.ParseTS(created); err == nil {
		t.CreatedAt = v
	}
	if v, err := sqliteutils.ParseTS(updated); err == nil {
		t.UpdatedAt = v
	}
	return t, nil
}
