package tests

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

func TestHTTP_Register_Login(t *testing.T) {
	st := suite.New(t)
	defer st.Cancel()

	gofakeit.Seed(time.Now().UnixNano())

	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, true, false, 12)

	// register
	bodyReg, _ := json.Marshal(map[string]any{
		"email":      email,
		"first_name": gofakeit.FirstName(),
		"last_name":  gofakeit.LastName(),
		"patronymic": gofakeit.LetterN(1),
		"birthdate":  gofakeit.Date().Format("2006-01-02"),
		"password":   pass,
	})
	req, _ := http.NewRequest(http.MethodPost, st.BaseURL+"/auth/register", bytes.NewReader(bodyReg))
	req.Header.Set("Content-Type", "application/json")
	resp, err := st.Client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var regResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&regResp))
	require.NotEmpty(t, regResp.AccessToken)

	// login
	bodyLogin, _ := json.Marshal(map[string]any{
		"email":    email,
		"password": pass,
	})
	req, _ = http.NewRequest(http.MethodPost, st.BaseURL+"/auth/login", bytes.NewReader(bodyLogin))
	req.Header.Set("Content-Type", "application/json")
	resp, err = st.Client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResp struct {
		AccessToken string `json:"access_token"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&loginResp))
	require.NotEmpty(t, loginResp.AccessToken)
}
