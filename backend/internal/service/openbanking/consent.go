package openbanking

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"multibank/backend/internal/domain"
	"multibank/backend/internal/http-server/dto"
	"net/http"
	"net/url"
	"path"
	"time"
)

type ConsentService struct {
	httpClient *http.Client
	log        *slog.Logger
}

type ConsentOperations interface {
	GetByID(
		Bank domain.Bank,
		ConsentID string,
	) (*domain.Consent, error)

	DeleteByIDDeleteByID(
		Bank domain.Bank,
		ConsentID string,
	) error

	Request(
		Bank domain.Bank,
		Permissions []domain.Permission,
		RequestingUser string,
		AccessToken domain.BankToken,
	) (*domain.Consent, error)
}

func CreatePath(Bank domain.Bank, ConsentID string) string {
	return path.Join(
		Bank.APIBaseURL,
		"account-consents",
		url.PathEscape(ConsentID),
	)
}

func (s *ConsentService) GetByID(Bank domain.Bank, ConsentID string) (*domain.Consent, error) {
	const operation = "service.openbanking.GetByID"

	log := s.log.With(
		slog.String("op", operation),
		slog.String("ConsentID", ConsentID),
	)

	log.Info("Retreiving consent")

	url := CreatePath(Bank, ConsentID)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode <= 200 && resp.StatusCode > 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("account-consents %d: %s", resp.StatusCode, string(body))
	}

	var consent dto.ConsentViewWrapper
	if err := json.NewDecoder(resp.Body).Decode(&consent); err != nil {
		return nil, err
	}

	// TODO какие-то проверки надо

	return &domain.Consent{
		ID:                   consent.Data.ID,
		Status:               consent.Data.Status,
		CreationDateTime:     consent.Data.CreationDateTime,
		StatusUpdateDateTime: consent.Data.StatusUpdateDateTime,
		Permissions:          consent.Data.Permissions,
		ExpirationDateTime:   consent.Data.ExpirationDateTime,
	}, nil
}

func (s *ConsentService) DeleteByID(Bank domain.Bank, ConsentID string) error {
	const operation = "service.openbanking.DeleteByID"

	log := s.log.With(
		slog.String("op", operation),
		slog.String("ConsentID", ConsentID),
	)

	log.Info("Deleting consent")

	url := CreatePath(Bank, ConsentID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode <= 200 && resp.StatusCode > 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("account-consents %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *ConsentService) Request(
	Bank domain.Bank,
	Permissions []domain.Permission,
	RequestingUser string,
	AccessToken domain.BankToken,
) (*domain.Consent, error) {
	const operation = "service.openbanking.Request"

	log := s.log.With(
		slog.String("op", operation),
		slog.Any("bank", Bank),
		slog.Any("permissions", Permissions),
		slog.String("requestingUser", RequestingUser),
	)

	log.Info("Creating a new request for consent")

	requestBody := &dto.ConsentRequest{
		ClientId:           RequestingUser,
		Permissions:        Permissions,
		Reason:             "just give me consent pretty please",
		RequestingBank:     Bank.Code,
		RequestingBankName: Bank.Name,
	}
	requestBodyJson, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	url := path.Join(
		Bank.APIBaseURL,
		"account-consents",
	)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyJson))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+AccessToken.AccessToken)
	req.Header.Set("X-Requesting-Bank", Bank.Code)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode <= 200 && resp.StatusCode > 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("account-consents %d: %s", resp.StatusCode, string(body))
	}

	var consentResponse dto.ConsentRequestResponse
	if err := json.NewDecoder(resp.Body).Decode(&consentResponse); err != nil {
		return nil, err
	}

	return &domain.Consent{
		ID:                   consentResponse.ConsentId,
		Status:               consentResponse.Status,
		CreationDateTime:     time.Now(),
		Permissions:          Permissions,
		StatusUpdateDateTime: time.Now(),
		ExpirationDateTime:   time.Now().Add(time.Duration(5) * time.Minute),
	}, nil
}
