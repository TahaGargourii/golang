package middleware

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("secret")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if isValidToken(tokenString) {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}

func isValidToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return mySigningKey, nil
	})

	if err != nil {
		return false
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true
	}

	return false
}
