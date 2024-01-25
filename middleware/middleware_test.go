package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"goProject/middleware"

	"github.com/dgrijalva/jwt-go"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	tokenString, err := generateValidToken()
	if err != nil {
		t.Fatalf("Error generating valid token: %v", err)
	}

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", tokenString)

	rr := httptest.NewRecorder()

	handler := middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %v, but got %v", http.StatusUnauthorized, rr.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	invalidTokenString := "invalid_token"

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", invalidTokenString)

	rr := httptest.NewRecorder()

	handler := middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with an invalid token")
	}))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %v, but got %v", http.StatusUnauthorized, rr.Code)
	}
}

func generateValidToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "John Doe"
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString([]byte("eyJhbGcI5yririzvyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE3MDYxMjI1MjgsInVzZXIiOiJKb2huIERvZSJ9.Ni7icx3noQB_N18y6lkF-FA0qV4yEcCkrjwmjj42tzY"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
