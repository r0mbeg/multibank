// internal/service/consent/service.go
package consent

import (
	"context"
	"log/slog"
	"multibank/backend/internal/domain"
	"multibank/backend/internal/logger"
	ob "multibank/backend/internal/service/openbanking"
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
	UserID      int64
	BankID      int64
	ClientID    string              // clien_id from FRONTEND
	Permissions []domain.Permission // perms can be overrided
}

func (s *Service) Request(ctx context.Context, in CreateInput) (int64, error) {

	const op = "service.consent.Request"

	log := s.log.With(slog.String("op", op))

	bank, err := s.banks.GetBankByID(ctx, in.BankID)
	if err != nil {
		log.Warn("failed to get bank",
			slog.Int64("bank_id", in.BankID),
			logger.Err(err),
		)
		return 0, err
	}
	token, _, err := s.banks.GetOrRefreshToken(ctx, in.BankID)
	if err != nil {
		log.Warn("failed to get bank access token",
			slog.Int64("bank_id", in.BankID),
			logger.Err(err),
		)
		return 0, err
	}

	perms := in.Permissions
	if len(perms) == 0 {
		perms = s.defaultPerms
	}

	resp, err := s.client.RequestConsent(bank, in.ClientID, perms, token)
	if err != nil {
		log.Warn("failed to request consent", logger.Err(err))
		return 0, err
	}

	now := time.Now()
	status := domain.ConsentStatus(resp.Status)

	c := &domain.AccountConsent{
		UserID:             in.UserID,
		BankID:             in.BankID,
		RequestID:          resp.RequestID,
		ConsentID:          resp.ConsentID,
		Status:             status,
		AutoApproved:       resp.AutoApproved,
		Permissions:        perms,
		Reason:             s.defaultReason,
		RequestingBank:     s.reqBankCode,
		RequestingBankName: s.reqBankName,
		CreatedAt:          now,
		UpdatedAt:          now,
		BankStatus:         &resp.Status,
	}
	return s.repo.Create(ctx, c)
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
	v, err := s.client.GetConsent(bank, key, s.reqBankCode) // x-fapi-interaction-id
	if err != nil {
		return domain.AccountConsent{}, err
	}

	our := domain.AccountConsent{
		Status: domain.ConsentStatus(v.Data.Status),
	}
	consentID := v.Data.ConsentID
	our.ConsentID = &consentID
	our.BankStatus = &v.Data.Status
	our.BankCreationDateTime = &v.Data.CreationDateTime
	our.BankStatusUpdateTime = &v.Data.StatusUpdateDateTime
	our.BankExpirationTime = &v.Data.ExpirationDateTime

	if err := s.repo.UpdateAfterCheck(ctx, c.ID, &our); err != nil {
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
