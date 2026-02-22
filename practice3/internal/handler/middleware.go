package handler

import (
	"log"
	"net/http"
	"time"
)

const ApiKeyHeader = "X-API-KEY"

// ValidAPIKey is the key that must be sent in X-API-KEY (change in production).
var ValidAPIKey = "secret-api-key"

// LoggingMiddleware logs every request: timestamp, method, endpoint.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		method := r.Method
		endpoint := r.URL.Path
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s (elapsed: %v)", start.Format(time.RFC3339), method, endpoint, time.Since(start))
	})
}

// AuthMiddleware returns 401 if X-API-KEY is missing or invalid.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get(ApiKeyHeader)
		if key == "" {
			WriteError(w, http.StatusUnauthorized, "missing X-API-KEY header")
			return
		}
		if key != ValidAPIKey {
			WriteError(w, http.StatusUnauthorized, "invalid X-API-KEY")
			return
		}
		next.ServeHTTP(w, r)
	})
}
