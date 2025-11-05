// tests/testutils/utils.go

package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"

	"multibank/backend/tests/suite"
)

const PassDefaultLen = 10

type ResponseWrapper struct {
	Resp *http.Response
}

func PostWithBody(t *testing.T, s *suite.Suite, path string, body any) *ResponseWrapper {
	data, err := json.Marshal(body)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, s.BaseURL+path, bytes.NewReader(data))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	require.NoError(t, err)

	return &ResponseWrapper{Resp: resp}
}

// GetWOBody doest HTTP GET-request without and returns ResponseWrapper.
func GetWOBody(t *testing.T, s *suite.Suite, path string, headers ...map[string]string) *ResponseWrapper {
	req, err := http.NewRequest(http.MethodGet, s.BaseURL+path, nil)
	require.NoError(t, err)

	req.Header.Set("Accept", "application/json")

	// Если переданы доп. заголовки (например, Authorization)
	if len(headers) > 0 {
		for k, v := range headers[0] {
			req.Header.Set(k, v)
		}
	}

	resp, err := s.Client.Do(req)
	require.NoError(t, err)

	return &ResponseWrapper{Resp: resp}
}

// GetWithAuth doest GET with the header Authorization: Bearer <token>
func GetWithAuth(t *testing.T, s *suite.Suite, path, token string) *ResponseWrapper {
	req, err := http.NewRequest(http.MethodGet, s.BaseURL+path, nil)
	require.NoError(t, err)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.Client.Do(req)
	require.NoError(t, err)
	return &ResponseWrapper{Resp: resp}
}

func (r *ResponseWrapper) ExpectStatus(t *testing.T, code int) *ResponseWrapper {
	require.Equal(t, code, r.Resp.StatusCode)
	return r
}

func (r *ResponseWrapper) DecodeTokenResponse(t *testing.T) TokenResponse {
	defer r.Resp.Body.Close()
	var out TokenResponse
	require.NoError(t, json.NewDecoder(r.Resp.Body).Decode(&out))
	require.NotEmpty(t, out.AccessToken)
	return out
}

func DecodeJSON[T any](t *testing.T, resp *http.Response) T {
	defer resp.Body.Close()
	var out T
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&out))
	return out
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// ==== Fake data

type FakeUser struct {
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Patronymic string `json:"patronymic"`
	BirthDate  string `json:"birthdate"`
	Password   string `json:"password"`
}

func NewFakeUser() FakeUser {
	gofakeit.Seed(time.Now().UnixNano())
	return FakeUser{
		Email:      gofakeit.Email(),
		FirstName:  gofakeit.FirstName(),
		LastName:   gofakeit.LastName(),
		Patronymic: gofakeit.LetterN(1),
		BirthDate:  gofakeit.Date().Format("2006-01-02"),
		Password:   RandomFakePassword(),
	}
}

func RandomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, PassDefaultLen)
}
