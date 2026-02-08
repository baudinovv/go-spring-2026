package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"task-api/internal/models"
)

const (
	APIKey = "secret12345"
)

// APIKeyAuth validates the X-API-KEY header
func APIKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")

		if apiKey == "" || apiKey != APIKey {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "unauthorized"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Logger logs HTTP method and path for every request
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timestamp := time.Now().Format("2006-01-02T15:04:05")

		log.Printf("%s %s %s Request received", timestamp, r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}
