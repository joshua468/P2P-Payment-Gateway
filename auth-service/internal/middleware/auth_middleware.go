package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if strings.HasPrefix(tokenString, "Bearer ") {
				tokenString = tokenString[len("Bearer "):]
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
