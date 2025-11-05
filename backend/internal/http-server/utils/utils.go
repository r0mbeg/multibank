// internal/http-server/utils/utils.go

package utils

import (
	"encoding/json"
	"net/http"
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
