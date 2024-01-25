package route

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationRoutes(t *testing.T) {
	// Create a Redis client for testing
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // default DB
	})
	ctx := context.Background()

	// Create the router using the test Redis client
	router := NewRouter(rdb, ctx)

	// Test adding an entry
	t.Run("AddEntry", func(t *testing.T) {
		word := "testWord"
		definition := "testDefinition"
		url := fmt.Sprintf("/add/%s/%s", word, definition)

		req, err := http.NewRequest("POST", url, nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	// Test getting a definition
	t.Run("GetDefinition", func(t *testing.T) {
		word := "testWord"
		url := fmt.Sprintf("/get/%s", word)

		req, err := http.NewRequest("GET", url, nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Definition of")
	})

	// Test removing an entry
	t.Run("RemoveEntry", func(t *testing.T) {
		word := "testWord"
		url := fmt.Sprintf("/remove/%s", word)

		req, err := http.NewRequest("DELETE", url, nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
