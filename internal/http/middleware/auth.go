package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"yatdl/internal/auth"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Auth(jwt *auth.Jwt) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Println("missing authorization header")
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Println("invalid authorization header")
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			userID, err := jwt.ValidateJWT(parts[1])
			if err != nil {
				log.Printf("failed to validate JWT: %v", err)
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
