package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"

	"goProject/dictionary"
	"goProject/middleware"
)

const filePath = "dictionary.txt"

var mySigningKey = []byte("secret")

func generateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user"] = "John Doe"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", fmt.Errorf("Something Went Wrong: %s", err.Error())
	}

	return tokenString, nil
}

func main() {

	fmt.Println("My Simple JWT Creation Program")
	tokenString, err := generateJWT()
	if err != nil {
		fmt.Println("Error generating token string")
	} else {
		fmt.Println("Generated Token String: ", tokenString)
	}

	d := dictionary.New(filePath)
	defer d.Close()

	router := mux.NewRouter()

	router.Use(middleware.LoggingMiddleware)

	router.Use(middleware.AuthMiddleware)

	var wg sync.WaitGroup

	router.HandleFunc("/add/{word}/{definition}", func(w http.ResponseWriter, r *http.Request) {
		wg.Add(1)
		actionAdd(d, w, r, &wg)
	})

	router.HandleFunc("/get/{word}", func(w http.ResponseWriter, r *http.Request) {
		actionDefine(d, w, r)
	})

	router.HandleFunc("/remove/{word}", func(w http.ResponseWriter, r *http.Request) {
		wg.Add(1)
		actionRemove(d, w, r, &wg)
	})

	router.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		actionList(d, w, r)
	})

	router.HandleFunc("/exit", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Exiting program.")
		wg.Wait()
		os.Exit(0)
	})

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
	wg.Wait()
}

func actionAdd(d *dictionary.Dictionary, w http.ResponseWriter, r *http.Request, wg *sync.WaitGroup) {
	defer wg.Done()

	vars := mux.Vars(r)
	word := vars["word"]
	definition := vars["definition"]

	err := d.Add(word, definition)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error adding word: %s", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Word '%s' added to the dictionary.\n", word)
}

func actionRemove(d *dictionary.Dictionary, w http.ResponseWriter, r *http.Request, wg *sync.WaitGroup) {
	defer wg.Done()

	vars := mux.Vars(r)
	word := vars["word"]

	err := d.Remove(word)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error removing word: %s", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Word '%s' removed from the dictionary.\n", word)
}

func actionDefine(d *dictionary.Dictionary, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	word := vars["word"]

	entry, err := d.Get(word)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Definition of '%s': %s\n", word, entry.String())
}

func actionList(d *dictionary.Dictionary, w http.ResponseWriter, r *http.Request) {
	words, entries, err := d.List()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing words: %s", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Words in the dictionary:")
	for _, word := range words {
		fmt.Fprintf(w, "%s: %s\n", word, entries[word].String())
	}
}
