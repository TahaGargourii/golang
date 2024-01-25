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
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()

	router := NewRouter(rdb, ctx)

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
