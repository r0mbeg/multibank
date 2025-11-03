package mw

import (
	"context"
	authjwt "multibank/backend/internal/service/auth/jwt"
	"net/http"
	"strings"
)

// ctxKey — типизированный ключ для контекста
type ctxKey int

const userIDKey ctxKey = iota

// WithUserID добавляет userID в контекст
func WithUserID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

// UserIDFromContext извлекает userID из контекста
func UserIDFromContext(ctx context.Context) (int64, bool) {
	v, ok := ctx.Value(userIDKey).(int64)
	return v, ok
}

// Auth — middleware, проверяющий Bearer-токен и добавляющий userID в контекст
func Auth(jwtMgr *authjwt.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				http.Error(w, `{"error":"missing or invalid token"}`, http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(header, "Bearer ")
			claims, err := jwtMgr.Parse(token)
			if err != nil {
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			ctx := WithUserID(r.Context(), claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
