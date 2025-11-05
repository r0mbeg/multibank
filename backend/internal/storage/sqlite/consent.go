package sqlite

import (
	"context"
	"database/sql"
	"multibank/backend/internal/domain"
	"time"
)

type ConsentRepo struct {
	db *sql.DB
}


func NewConsentRepo(db *sql.DB) *ConsentRepo {
	return &ConsentRepo{db: db}
}

func (s *ConsentRepo) GetConsentByID(ctx context.Context, id int64) (domain.FullConsent, error) {

	row := s.db.QueryRowContext(ctx, `
		SELECT consent_id, user_id, bank_id, created_at, expiration_date_time, status
		FROM consents WHERE if = ?
		`, id)
	var c domain.FullConsent
	if err := row.Scan(&c.ID, &c.UserID, &c.BankID, &c.CreationDateTime, &c.ExpirationDateTime, &c.Status); err != nil {
		return domain.FullConsent{}, err
	}
	return c, nil
}

func (s *ConsentRepo) UpsertConsent(ctx context.Context, c domain.FullConsent) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO consents(consent_id, user_id, bank_id, created_at, expiration_date_time, status)
		VALUES(?,?,?,?,?,?)
		ON CONFLICT(bank_id) DO UPDATE SET
		  access_token = excluded.access_token,
		  expires_at   = excluded.expires_at,
		  updated_at   = datetime('now')`,
		c.ID, c.UserID, c.BankID,
		c.CreationDateTime.UTC().Format(time.RFC3339),
		c.ExpirationDateTime.UTC().Format(time.RFC3339),
		c.Status)
	return err
}

// TODO: удалить по ид - по идее это временный объект и их нужно постоянно удалять, но пока это не особо нужно