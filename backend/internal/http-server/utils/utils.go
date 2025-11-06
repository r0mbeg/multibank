// internal/http-server/utils/utils.go

package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Local handlers' utils
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	type errResp struct {
		Error string `json:"error"`
	}
	WriteJSON(w, status, errResp{Error: msg})
}

// NormalizeURL normalizes usr
func NormalizeURL(raw string) (*url.URL, error) {
	if raw == "" {
		return nil, fmt.Errorf("empty url")
	}
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return nil, err
	}
	// Если схема не указана — добавим http по умолчанию
	if u.Scheme == "" {
		u, err = url.Parse("https://" + raw)
		if err != nil {
			return nil, err
		}
	}
	// Уберём лишний мусор в пути: хотим чтобы base оканчивался без хвостового слеша
	u.Path = strings.TrimRight(u.Path, "/")
	return u, nil
}

// MaskSecret hides secret (for logs)
func MaskSecret(u *url.URL) string {
	cp := *u
	q := cp.Query()
	if q.Has("client_secret") {
		q.Set("client_secret", "******")
	}
	cp.RawQuery = q.Encode()
	return cp.String()
}
