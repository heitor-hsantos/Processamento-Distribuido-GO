package middleware

import (
	"context"
	"net/http"
	"strings"

	"JungleGaming-test/internal/infrastructure/auth"
)

type contextKey string

const ProviderIDKey contextKey = "providerId"

func Auth(validator *auth.OIDCValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "missing or invalid authorization header", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			providerID, err := validator.ValidateToken(r.Context(), tokenStr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ProviderIDKey, providerID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
