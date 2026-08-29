package transport

import (
	"context"
	"net/http"
	"portfolio/internal/user"
	"strings"
)

type ctxKey int

const userIDKey ctxKey = 0

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authz := r.Header.Get("Authorization")
		tokenStr, ok := strings.CutPrefix(authz, "Bearer ")
		if !ok {
			writeError(w, http.StatusUnauthorized, "authorization header format must be Bearer {token}")
			return
		}
		claims, err := user.VerifyToken(tokenStr)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid token")
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIDFromContext(ctx context.Context) (int32, bool) {
	id, ok := ctx.Value(userIDKey).(int32)
	return id, ok
}
