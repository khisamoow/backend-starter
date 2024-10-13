package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("SECRET_KEY"))

func TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Неверный токен", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
