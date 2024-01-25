package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"context"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"

	"golang/middleware"
)

var mySigningKey = []byte("secret")
var ctx = context.Background()

func generateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user"] = "Taha Gargouri"
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

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	router := mux.NewRouter()

	router.Use(middleware.LoggingMiddleware)

	router.Use(middleware.AuthMiddleware)

	var wg sync.WaitGroup

	router.HandleFunc("/add/{word}/{definition}", func(w http.ResponseWriter, r *http.Request) {
		wg.Add(1)
		actionAdd(rdb, w, r, &wg)
	})

	router.HandleFunc("/get/{word}", func(w http.ResponseWriter, r *http.Request) {
		actionDefine(rdb, w, r)
	})

	router.HandleFunc("/remove/{word}", func(w http.ResponseWriter, r *http.Request) {
		wg.Add(1)
		actionRemove(rdb, w, r, &wg)
	})

	router.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		actionList(rdb, w, r)
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

func actionAdd(rdb *redis.Client, w http.ResponseWriter, r *http.Request, wg *sync.WaitGroup) {
	defer wg.Done()

	vars := mux.Vars(r)
	word := vars["word"]
	definition := vars["definition"]

	err := rdb.Set(ctx, word, definition, 0).Err()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error adding word: %s", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Word '%s' added to the dictionary.\n", word)
}

func actionRemove(rdb *redis.Client, w http.ResponseWriter, r *http.Request, wg *sync.WaitGroup) {
	defer wg.Done()

	vars := mux.Vars(r)
	word := vars["word"]

	err := rdb.Del(ctx, word).Err()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error removing word: %s", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Word '%s' removed from the dictionary.\n", word)
}

func actionDefine(rdb *redis.Client, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	word := vars["word"]

	definition, err := rdb.Get(ctx, word).Result()
	if err == redis.Nil {
		http.Error(w, fmt.Sprintf("Word '%s' does not exist", word), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Definition of '%s': %s\n", word, definition)
}

func actionList(rdb *redis.Client, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This operation is not supported. Redis does not support listing all keys in a performant way!")
}
