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

func (r *ConsentRepo) Create(ctx context.Context, c *domain.AccountConsent) (int64, error) {
	const q = `
INSERT INTO account_consents
(user_id, bank_id, request_id, consent_id, status, auto_approved, permissions_json, reason, requesting_bank, requesting_bank_name,
 bank_status, bank_creation_datetime, bank_status_update_datetime, bank_expiration_datetime)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	perms, _ := json.Marshal(c.Permissions)
	res, err := r.db.ExecContext(ctx, q,
		c.UserID, c.BankID, c.RequestID, c.ConsentID, string(c.Status), c.AutoApproved,
		string(perms), c.Reason, c.RequestingBank, c.RequestingBankName,
		c.BankStatus, sqliteutils.ToISO(c.BankCreationDateTime), sqliteutils.ToISO(c.BankStatusUpdateTime), sqliteutils.ToISO(c.BankExpirationTime),
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
    bank_status = ?,
    bank_creation_datetime = ?,
    bank_status_update_datetime = ?,
    bank_expiration_datetime = ?,
    updated_at = datetime('now')
WHERE id = ?`
	_, err := r.db.ExecContext(ctx, q,
		upd.ConsentID,
		string(upd.Status),
		upd.AutoApproved,
		upd.BankStatus,
		sqliteutils.ToISO(upd.BankCreationDateTime),
		sqliteutils.ToISO(upd.BankStatusUpdateTime),
		sqliteutils.ToISO(upd.BankExpirationTime),
		id,
	)
	return err
}

func (r *ConsentRepo) GetByID(ctx context.Context, id int64) (domain.AccountConsent, error) {
	const q = `SELECT id,user_id,bank_id,request_id,consent_id,status,auto_approved,permissions_json,reason,requesting_bank,requesting_bank_name,
    bank_status,bank_creation_datetime,bank_status_update_datetime,bank_expiration_datetime,created_at,updated_at
    FROM account_consents WHERE id=?`

	var (
		c                                    domain.AccountConsent
		perms                                string
		autoApproved                         *int64
		consentID                            *string
		bankCreation, bankUpdate, bankExpiry *string
		createdAtStr, updatedAtStr           *string
	)

	row := r.db.QueryRowContext(ctx, q, id)
	if err := row.Scan(
		&c.ID, &c.UserID, &c.BankID, &c.RequestID, &consentID, &c.Status, &autoApproved, &perms,
		&c.Reason, &c.RequestingBank, &c.RequestingBankName,
		&c.BankStatus, &bankCreation, &bankUpdate, &bankExpiry,
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
	_ = json.Unmarshal([]byte(perms), &c.Permissions)

	c.BankCreationDateTime = sqliteutils.FromISO(bankCreation)
	c.BankStatusUpdateTime = sqliteutils.FromISO(bankUpdate)
	c.BankExpirationTime = sqliteutils.FromISO(bankExpiry)

	if t := sqliteutils.FromISO(createdAtStr); t != nil {
		c.CreatedAt = *t
	}
	if t := sqliteutils.FromISO(updatedAtStr); t != nil {
		c.UpdatedAt = *t
	}

	return c, nil
}

func (r *ConsentRepo) ListByUser(ctx context.Context, userID int64, bankID *int64) ([]domain.AccountConsent, error) {
	q := `SELECT id FROM account_consents WHERE user_id=?`
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
	var out []domain.AccountConsent
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		c, err := r.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}

func (r *ConsentRepo) DeleteByID(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM account_consents WHERE id=?`, id)
	return err
}
