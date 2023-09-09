package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net/http"
	"os"
)

const (
	shortURLLen = 6
)

var (
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),      // Redis server address
		Password: os.Getenv("REDIS_PASSWORD"), // Your Redis password
		DB:       0,                           // Use default DB
	})
	ctx = context.Background()
)

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/", expandHandler)
	fmt.Println("Server is running on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
}
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Service is healthy!")
}
func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusBadRequest)
		return
	}

	longURL := r.FormValue("url")
	if longURL == "" {
		http.Error(w, "Missing url parameter", http.StatusBadRequest)
		return
	}
	shortURL := generateShortURL(longURL)
	err := rdb.HSet(ctx, "longURL", shortURL, longURL).Err()
	if err != nil {
		fmt.Printf(err.Error())
		http.Error(w, "Failed to save to Redis", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s/%s", os.Getenv("BASE_URL"), shortURL)
}

func expandHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[len("/"):]
	fmt.Println("接受参数", shortURL)
	// Using Redis HGET to retrieve the long URL associated with shortURL
	longURL, err := rdb.HGet(ctx, "longURL", shortURL).Result()
	fmt.Println(longURL)
	if err == redis.Nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch from Redis", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, longURL, http.StatusSeeOther)
}

func generateShortURL(url string) string {
	hash := sha256.Sum256([]byte(url))
	return base64.URLEncoding.EncodeToString(hash[:])[:shortURLLen]
}
