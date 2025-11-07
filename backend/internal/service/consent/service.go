// internal/service/consent/service.go
package consent

import (
	"context"
	"log/slog"
	"multibank/backend/internal/domain"
	"multibank/backend/internal/logger"
	ob "multibank/backend/internal/service/openbanking"
	"strings"
	"time"
)

type ConsentRepo interface {
	Create(ctx context.Context, c *domain.AccountConsent) (int64, error)
	UpdateAfterCheck(ctx context.Context, id int64, upd *domain.AccountConsent) error
	GetByID(ctx context.Context, id int64) (domain.AccountConsent, error)
	ListByUser(ctx context.Context, userID int64, bankID *int64) ([]domain.AccountConsent, error)
	DeleteByID(ctx context.Context, id int64) error
}

type BankService interface {
	GetBankByID(ctx context.Context, id int64) (domain.Bank, error)
	GetBankByCode(ctx context.Context, code string) (domain.Bank, error)
	// token for the request to the bank
	GetOrRefreshToken(ctx context.Context, bankID int64) (string, time.Time, error)
}

type OBConsentClient interface {
	RequestConsent(bank domain.Bank, clientID string, perms []domain.Permission, bearer string) (*ob.ConsentRequestResp, error)
	GetConsent(bank domain.Bank, requestOrConsentID string, xFapi string) (*ob.ConsentViewWrapper, error)
}

type Service struct {
	log    *slog.Logger
	repo   ConsentRepo
	banks  BankService
	client OBConsentClient

	defaultPerms  []domain.Permission
	reqBankCode   string
	reqBankName   string
	defaultReason string
}

func New(log *slog.Logger, repo ConsentRepo, banks BankService, client OBConsentClient,
	defaultPerms []domain.Permission, reqBankCode, reqBankName, defaultReason string) *Service {
	return &Service{log: log, repo: repo, banks: banks, client: client,
		defaultPerms: defaultPerms, reqBankCode: reqBankCode, reqBankName: reqBankName, defaultReason: defaultReason}
}

type CreateInput struct {
	UserID   int64
	BankCode string // (оставим по id; если хочешь по имени — скажу где поменять)
	ClientID string // comes from FRONTEND
	// Permissions []domain.Permission // always default
}

func (s *Service) Request(ctx context.Context, in CreateInput) (int64, error) {
	const op = "service.consent.Request"

	log := s.log.With(
		slog.String("op", op),
		slog.String("bank_code", in.BankCode),
	)

	log.Info("requesting a new consent")

	bank, err := s.banks.GetBankByCode(ctx, in.BankCode)
	if err != nil {
		log.Warn("failed to get bank", logger.Err(err))
		return 0, err
	}

	log = log.With(slog.Int64("bank_id", bank.ID))

	token, _, err := s.banks.GetOrRefreshToken(ctx, bank.ID)
	if err != nil {
		log.Warn("failed to get bank access token", slog.Int64("bank_id", bank.ID), logger.Err(err))
		return 0, err
	}

	// фиксированный набор разрешений
	perms := s.defaultPerms

	resp, err := s.client.RequestConsent(bank, in.ClientID, perms, token)
	if err != nil {
		log.Warn("failed to request consent", logger.Err(err))
		return 0, err
	}

	now := time.Now()

	status := normalizeRequestStatus(resp.Status, resp.AutoApproved)

	var (
		creation  *time.Time
		updated   *time.Time
		expire    *time.Time
		consentID *string = resp.ConsentID
	)

	// Если автоодобрено — подтянем детальный вид, чтобы заполнить даты.
	if resp.AutoApproved != nil && *resp.AutoApproved {
		key := resp.RequestID
		if consentID != nil && *consentID != "" {
			key = *consentID
		}
		if v, err := s.client.GetConsent(bank, key, s.reqBankCode); err == nil {
			// у банка уже правильные CamelCase статусы
			status = domain.ConsentStatus(v.Data.Status) // "Authorized"
			creation = &v.Data.CreationDateTime
			updated = &v.Data.StatusUpdateDateTime
			expire = &v.Data.ExpirationDateTime
			// иногда consentID выдают только в GET
			if v.Data.ConsentID != "" {
				cid := v.Data.ConsentID
				consentID = &cid
			}
		} else {
			log.Warn("auto-approved but failed to fetch detailed consent", logger.Err(err))
		}
	}

	c := &domain.AccountConsent{
		UserID:             in.UserID,
		BankID:             bank.ID,
		RequestID:          resp.RequestID,
		ConsentID:          consentID,
		Status:             status,
		AutoApproved:       resp.AutoApproved,
		ClientID:           in.ClientID,
		Permissions:        perms,
		Reason:             s.defaultReason,
		RequestingBank:     s.reqBankCode,
		RequestingBankName: s.reqBankName,
		CreatedAt:          now,
		UpdatedAt:          now,

		CreationDateTime:     creation,
		StatusUpdateDateTime: updated,
		ExpirationDateTime:   expire,
	}
	return s.repo.Create(ctx, c)
}

func normalizeRequestStatus(raw string, auto *bool) domain.ConsentStatus {
	if auto != nil && *auto {
		return domain.Authorised
	}
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "approved", "authorised", "authorized":
		return domain.Authorised
	case "pending", "awaitingauthorization", "awaitingauthorisation":
		return domain.AwaitingAuthorisation
	case "rejected":
		return domain.Rejected
	case "revoked":
		return domain.Revoked
	default:
		return domain.ConsentStatus(raw)
	}
}

func (s *Service) Refresh(ctx context.Context, id int64) (domain.AccountConsent, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.AccountConsent{}, err
	}

	bank, err := s.banks.GetBankByID(ctx, c.BankID)
	if err != nil {
		return domain.AccountConsent{}, err
	}

	key := c.RequestID
	if c.ConsentID != nil && *c.ConsentID != "" {
		key = *c.ConsentID
	}

	v, err := s.client.GetConsent(bank, key, s.reqBankCode)
	if err != nil {
		return domain.AccountConsent{}, err
	}

	upd := domain.AccountConsent{
		Status:               domain.ConsentStatus(v.Data.Status),
		CreationDateTime:     &v.Data.CreationDateTime,
		StatusUpdateDateTime: &v.Data.StatusUpdateDateTime,
		ExpirationDateTime:   &v.Data.ExpirationDateTime,
		AutoApproved:         nil, // не меняем
	}
	cid := v.Data.ConsentID
	upd.ConsentID = &cid

	if err := s.repo.UpdateAfterCheck(ctx, c.ID, &upd); err != nil {
		return domain.AccountConsent{}, err
	}
	return s.repo.GetByID(ctx, c.ID)
}

func (s *Service) Get(ctx context.Context, id int64) (domain.AccountConsent, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListMine(ctx context.Context, userID int64, bankID *int64) ([]domain.AccountConsent, error) {
	return s.repo.ListByUser(ctx, userID, bankID)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	// здесь можно дополнительно позвать DELETE в банк (если есть).
	return s.repo.DeleteByID(ctx, id)
}
