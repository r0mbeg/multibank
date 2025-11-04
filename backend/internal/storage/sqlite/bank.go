package sqlite

import (
	"context"
	"database/sql"
	"multibank/backend/internal/domain"
	"time"
)

type BankRepo struct {
	db *sql.DB
}

func NewBankRepo(db *sql.DB) *BankRepo {
	return &BankRepo{db: db}
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (s *BankRepo) ListEnabledBanks(ctx context.Context) ([]domain.Bank, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, code, api_base_url, login, password, is_enabled, created_at, updated_at
		FROM banks
		WHERE is_enabled = 1
		ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Bank
	for rows.Next() {
		var b domain.Bank
		var en int
		if err := rows.Scan(&b.ID, &b.Name, &b.Code, &b.APIBaseURL, &b.Login, &b.Password, &en, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		b.IsEnabled = en == 1
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
	if err := row.Scan(&b.ID, &b.Name, &b.Code, &b.APIBaseURL, &b.Login, &b.Password, &en, &b.CreatedAt, &b.UpdatedAt); err != nil {
		return domain.Bank{}, err
	}
	b.IsEnabled = en == 1
	return b, nil
}

func (s *BankRepo) UpsertBankToken(ctx context.Context, t domain.BankToken) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO bank_tokens(bank_id, access_token, expires_at)
		VALUES(?,?,?)
		ON CONFLICT(bank_id) DO UPDATE SET
		  access_token = excluded.access_token,
		  expires_at   = excluded.expires_at,
		  updated_at   = datetime('now')`,
		t.BankID, t.AccessToken, t.ExpiresAt.UTC().Format(time.RFC3339))
	return err
}

func (s *BankRepo) GetBankToken(ctx context.Context, bankID int64) (domain.BankToken, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT bank_id, access_token, expires_at, created_at, updated_at
		FROM bank_tokens WHERE bank_id = ?`, bankID)
	var t domain.BankToken
	if err := row.Scan(&t.BankID, &t.AccessToken, &t.ExpiresAt, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return domain.BankToken{}, err
	}
	return t, nil
}
