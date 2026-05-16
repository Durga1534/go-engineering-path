package middlewares

import (
	"context"
	"expense-tracker-api/internal/auth"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: Token missing", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return auth.SecretKey, nil
		})
		if err != nil {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", int(claims["user_id"].(float64)))
		next.ServeHTTP(w, r.WithContext(ctx))

	}
}
