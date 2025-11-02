package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"

	"multibank/backend/tests/suite"
)

const (
	passDefaultLen = 10
)

func TestHTTP_Register_Login(t *testing.T) {
	st := suite.New(t)
	defer st.Cancel()

	gofakeit.Seed(time.Now().UnixNano())

	email := gofakeit.Email()
	firstName := gofakeit.FirstName()
	lastName := gofakeit.LastName()
	patronymic := gofakeit.LetterN(1)
	birthdate := gofakeit.Date().Format("2006-01-02")
	pass := randomFakePassword()

	// register
	bodyReg, _ := json.Marshal(map[string]any{
		"email":      email,
		"first_name": firstName,
		"last_name":  lastName,
		"patronymic": patronymic,
		"birthdate":  birthdate,
		"password":   pass,
	})

	fmt.Printf("body in test: %s", bodyReg)

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

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
