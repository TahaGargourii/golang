package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"estiam/dictionary"
)

type TokenManager interface {
	VerifyToken(token string) bool
}

type InMemoryTokenManager struct {
	tokens map[string]time.Time
}

func NewInMemoryTokenManager() *InMemoryTokenManager {
	return &InMemoryTokenManager{tokens: make(map[string]time.Time)}
}

func (m *InMemoryTokenManager) VerifyToken(token string) bool {
	_, exists := m.tokens[token]
	if !exists {
		return false
	}

	// Check if the token has expired
	expirationTime := m.tokens[token]
	if time.Now().After(expirationTime) {
		return false
	}

	// Token is valid, renew its expiration time
	m.tokens[token] = time.Now().Add(time.Hour)
	return true
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf("[%s] %s %s %s %d", start.Format("2006-01-02T15:04:05.000Z"), r.Method, r.RequestURI, r.RemoteAddr, r.StatusCode)
	})
}

func TokenValidationMiddleware(tokenManager TokenManager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		tokenPrefix := "Bearer "

		if !strings.HasPrefix(token, tokenPrefix) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Invalid token format")
			return
		}
		tokenString := strings.TrimPrefix(token, tokenPrefix)
		if !tokenManager.VerifyToken(tokenString) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Invalid token")
			return
		}
		next.ServeHTTP(w, r)
	})
}
