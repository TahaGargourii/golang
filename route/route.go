package route

import (
	"fmt"
	"goProject/dictionary"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(d *dictionary.Dictionary) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/add/{word}/{definition}", HandleAddEntry(d)).Methods("POST")
	router.HandleFunc("/get/{word}", HandleGetDefinition(d)).Methods("GET")
	router.HandleFunc("/remove/{word}", HandleRemoveEntry(d)).Methods("DELETE")

	return router
}

func HandleAddEntry(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		word := vars["word"]
		definition := vars["definition"]

		err := d.Add(word, definition)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func HandleGetDefinition(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		word := vars["word"]

		entry, err := d.Get(word)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Definition of '%s': %s\n", word, entry.String())
	}
}

func HandleRemoveEntry(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		word := vars["word"]

		err := d.Remove(word)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
