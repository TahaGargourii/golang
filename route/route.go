package route

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

func NewRouter(rdb *redis.Client, ctx context.Context) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/add/{word}/{definition}", HandleAddEntry(rdb, ctx)).Methods("POST")
	router.HandleFunc("/get/{word}", HandleGetDefinition(rdb, ctx)).Methods("GET")
	router.HandleFunc("/remove/{word}", HandleRemoveEntry(rdb, ctx)).Methods("DELETE")

	return router
}

func HandleAddEntry(rdb *redis.Client, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		word := vars["word"]
		definition := vars["definition"]

		err := rdb.Set(ctx, word, definition, 0).Err()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func HandleGetDefinition(rdb *redis.Client, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		word := vars["word"]

		definition, err := rdb.Get(ctx, word).Result()
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Definition of '%s': %s\n", word, definition)
	}
}

func HandleRemoveEntry(rdb *redis.Client, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		word := vars["word"]

		err := rdb.Del(ctx, word).Err()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
