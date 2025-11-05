// tests/auth_e2e_test.go
package tests

import (
	"net/http"
	"strings"
	"testing"

	"multibank/backend/tests/suite"
	testutils "multibank/backend/tests/utils"

	"github.com/stretchr/testify/require"
)

func TestHTTP_AuthFlow(t *testing.T) {
	st := suite.New(t)
	defer st.Cancel()

	user := testutils.NewFakeUser()

	t.Run("try to get /me being unauthorized", func(t *testing.T) {
		testutils.GetWOBody(t, st, "/me").
			ExpectStatus(t, http.StatusUnauthorized)
	})

	t.Run("try to get /banks being unauthorized", func(t *testing.T) {
		testutils.GetWOBody(t, st, "/banks").
			ExpectStatus(t, http.StatusUnauthorized)
	})

	t.Run("login with wrong password -> 401", func(t *testing.T) {
		req := map[string]string{
			"email":    user.Email,
			"password": "wrong-pass-123",
		}
		testutils.
			PostWithBody(t, st, "/auth/login", req).
			ExpectStatus(t, http.StatusUnauthorized)
	})

	var token string
	t.Run("register new user -> 201 + token", func(t *testing.T) {
		tr := testutils.
			PostWithBody(t, st, "/auth/register", user).
			ExpectStatus(t, http.StatusCreated).
			DecodeTokenResponse(t)

		require.NotEmpty(t, tr.AccessToken)
		token = tr.AccessToken
	})

	t.Run("register same user again -> 409", func(t *testing.T) {
		testutils.
			PostWithBody(t, st, "/auth/register", user).
			ExpectStatus(t, http.StatusConflict)
	})

	t.Run("login with correct credentials -> 200 + token", func(t *testing.T) {
		req := map[string]string{
			"email":    user.Email,
			"password": user.Password,
		}
		testutils.
			PostWithBody(t, st, "/auth/login", req).
			ExpectStatus(t, http.StatusOK).
			DecodeTokenResponse(t)
	})

	t.Run("get /me being authorized -> 200 and correct payload", func(t *testing.T) {
		// Локальный тип, чтобы не тянуть dto в тесты
		type meResp struct {
			ID         int64  `json:"id"`
			Email      string `json:"email"`
			FirstName  string `json:"first_name"`
			LastName   string `json:"last_name"`
			Patronymic string `json:"patronymic"`
			BirthDate  string `json:"birthdate"`
			IsAdmin    bool   `json:"is_admin"`
			// created_at / updated_at можно тоже проверить при желании
		}

		resp := testutils.
			GetWithAuth(t, st, "/me", token).
			ExpectStatus(t, http.StatusOK).Resp

		got := testutils.DecodeJSON[meResp](t, resp)

		require.Greater(t, got.ID, int64(0))
		require.Equal(t, strings.ToLower(user.Email), got.Email) // у тебя в сервисе email приводится к lower
		require.Equal(t, user.FirstName, got.FirstName)
		require.Equal(t, user.LastName, got.LastName)
		require.Equal(t, user.Patronymic, got.Patronymic)
		require.Equal(t, user.BirthDate, got.BirthDate)
		require.False(t, got.IsAdmin) // по умолчанию false, если не менял
	})

	t.Run("get /banks being authorized -> 200", func(t *testing.T) {
		resp := testutils.
			GetWithAuth(t, st, "/banks", token).
			ExpectStatus(t, http.StatusOK).Resp

		// Если хочешь проверить схему — декодируй массив
		type bankResp struct {
			ID         int64  `json:"id"`
			Name       string `json:"name"`
			Code       string `json:"code"`
			APIBaseURL string `json:"api_base_url"`
			IsEnabled  bool   `json:"is_enabled"`
		}
		_ = testutils.DecodeJSON[[]bankResp](t, resp) // массив может быть пустым — это ок
	})

}
