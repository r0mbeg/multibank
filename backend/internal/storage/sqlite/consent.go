// internal/storage/sqlite/consent.go

package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"multibank/backend/internal/domain"
	sqliteutils "multibank/backend/internal/storage/sqlite/utils"
)

type ConsentRepo struct {
	db *sql.DB
}

func NewConsentRepo(db *sql.DB) *ConsentRepo { return &ConsentRepo{db: db} }

// single consentCols - not hardcode
const consentCols = `
id,user_id,bank_id,request_id,consent_id,status,auto_approved,permissions_json,reason,requesting_bank,requesting_bank_name,
creation_datetime,status_update_datetime,expiration_datetime,client_id,
created_at,updated_at
`

// rowScanner allows you to scan both *sql.Row and *sql.Rows
// resolves a problem with deadlock in ListByUser
type rowScanner interface{ Scan(dest ...any) error }

// scanConsent parses DB row and returns domain.AccountConsent
func scanConsent(rs rowScanner) (domain.AccountConsent, error) {
	var (
		c                               domain.AccountConsent
		perms                           string
		autoApproved                    *int64
		consentID                       *string
		creation, statusUpd, expiration *string
		createdAtStr, updatedAtStr      *string
	)

	if err := rs.Scan(
		&c.ID, &c.UserID, &c.BankID, &c.RequestID, &consentID, &c.Status, &autoApproved, &perms,
		&c.Reason, &c.RequestingBank, &c.RequestingBankName,
		&creation, &statusUpd, &expiration, &c.ClientID,
		&createdAtStr, &updatedAtStr,
	); err != nil {
		return domain.AccountConsent{}, err
	}

	if consentID != nil {
		c.ConsentID = consentID
	}
	if autoApproved != nil {
		v := (*autoApproved) != 0
		c.AutoApproved = &v
	}
	if perms != "" {
		_ = json.Unmarshal([]byte(perms), &c.Permissions)
	}

	c.CreationDateTime = sqliteutils.FromISO(creation)
	c.StatusUpdateDateTime = sqliteutils.FromISO(statusUpd)
	c.ExpirationDateTime = sqliteutils.FromISO(expiration)

	if t := sqliteutils.FromISO(createdAtStr); t != nil {
		c.CreatedAt = *t
	}
	if t := sqliteutils.FromISO(updatedAtStr); t != nil {
		c.UpdatedAt = *t
	}

	return c, nil
}

func (r *ConsentRepo) Create(ctx context.Context, c *domain.AccountConsent) (int64, error) {
	const q = `
INSERT INTO account_consents
(user_id, bank_id, request_id, consent_id, status, auto_approved, permissions_json, reason, requesting_bank, requesting_bank_name,
 creation_datetime, status_update_datetime, expiration_datetime, client_login)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	perms, _ := json.Marshal(c.Permissions)
	res, err := r.db.ExecContext(ctx, q,
		c.UserID, c.BankID, c.RequestID, c.ConsentID, string(c.Status), c.AutoApproved,
		string(perms), c.Reason, c.RequestingBank, c.RequestingBankName,
		sqliteutils.ToISO(c.CreationDateTime), sqliteutils.ToISO(c.StatusUpdateDateTime), sqliteutils.ToISO(c.ExpirationDateTime),
		c.ClientID,
	)
	if err != nil {
		return 0, fmt.Errorf("consent create: %w", err)
	}
	return res.LastInsertId()
}

func (r *ConsentRepo) UpdateAfterCheck(ctx context.Context, id int64, upd *domain.AccountConsent) error {
	const q = `
UPDATE account_consents
SET consent_id = COALESCE(?, consent_id),
    status = ?,
    auto_approved = COALESCE(?, auto_approved),
    creation_datetime = COALESCE(?, creation_datetime),
    status_update_datetime = COALESCE(?, status_update_datetime),
    expiration_datetime = COALESCE(?, expiration_datetime),
    updated_at = datetime('now')
WHERE id = ?`
	_, err := r.db.ExecContext(ctx, q,
		upd.ConsentID,
		string(upd.Status),
		upd.AutoApproved,
		sqliteutils.ToISO(upd.CreationDateTime),
		sqliteutils.ToISO(upd.StatusUpdateDateTime),
		sqliteutils.ToISO(upd.ExpirationDateTime),
		id,
	)
	return err
}

func (r *ConsentRepo) GetByID(ctx context.Context, id int64) (domain.AccountConsent, error) {
	q := `SELECT ` + consentCols + ` FROM account_consents WHERE id=?`
	row := r.db.QueryRowContext(ctx, q, id)
	return scanConsent(row)
}

func (r *ConsentRepo) ListByUser(ctx context.Context, userID int64, bankID *int64) ([]domain.AccountConsent, error) {
	q := `SELECT ` + consentCols + ` FROM account_consents WHERE user_id=?`
	args := []any{userID}
	if bankID != nil {
		q += ` AND bank_id=?`
		args = append(args, *bankID)
	}
	q += ` ORDER BY id DESC`

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.AccountConsent, 0, 16)
	for rows.Next() {
		c, err := scanConsent(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *ConsentRepo) DeleteByID(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM account_consents WHERE id=?`, id)
	return err
}
