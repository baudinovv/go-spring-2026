package handler

import (
	"net/http"
	"strings"

	"practice3/internal/usecase"
)

// NewRouter returns an http.Handler that routes to user handlers and health.
// Middleware (logging, auth) should be wrapped around this.
func NewRouter(uc *usecase.UserUsecase) http.Handler {
	uh := NewUserHandler(uc)
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" && r.URL.Path != "/health/" {
			WriteError(w, http.StatusNotFound, "not found")
			return
		}
		if r.Method != http.MethodGet {
			WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimSuffix(r.URL.Path, "/")
		if path != "/users" {
			WriteError(w, http.StatusNotFound, "not found")
			return
		}
		switch r.Method {
			case http.MethodGet:
				uh.ListUsers(w, r)
			case http.MethodPost:
				uh.CreateUser(w, r)
			default:
				WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			uh.GetUserByID(w, r)
		case http.MethodPatch, http.MethodPut:
			uh.UpdateUser(w, r)
		case http.MethodDelete:
			uh.DeleteUser(w, r)
		default:
			WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	return mux
}
