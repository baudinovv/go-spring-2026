package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"practice3/internal/repository/_postgres/users"
	"practice3/internal/usecase"
	"practice3/pkg/modules"
)

type UserHandler struct {
	uc *usecase.UserUsecase
}

func NewUserHandler(uc *usecase.UserUsecase) *UserHandler {
	return &UserHandler{uc: uc}
}

// ListUsers handles GET /users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/users" && r.URL.Path != "/users/" {
		return
	}
	list, err := h.uc.GetUsers()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if list == nil {
		list = []modules.User{}
	}
	WriteJSON(w, http.StatusOK, list)
}

// GetUserByID handles GET /users/{id}
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, ok := userIDFromPath(r.URL.Path, w)
	if !ok {
		return
	}
	user, err := h.uc.GetUserByID(id)
	if err != nil {
		if err == users.ErrUserNotFound {
			WriteError(w, http.StatusNotFound, "user not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, user)
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/users" && r.URL.Path != "/users/" {
		return
	}
	var in modules.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if in.Name == "" || in.Email == "" {
		WriteError(w, http.StatusBadRequest, "name and email are required")
		return
	}
	id, err := h.uc.CreateUser(in)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, map[string]int{"id": id})
}

// UpdateUser handles PATCH /users/{id} or PUT /users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, ok := userIDFromPath(r.URL.Path, w)
	if !ok {
		return
	}
	var in modules.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if err := h.uc.UpdateUser(id, in); err != nil {
		if err == users.ErrUserNotFound {
			WriteError(w, http.StatusNotFound, "user not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"message": "user updated successfully"})
}

// DeleteUser handles DELETE /users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, ok := userIDFromPath(r.URL.Path, w)
	if !ok {
		return
	}
	rows, err := h.uc.DeleteUser(id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if rows == 0 {
		WriteError(w, http.StatusNotFound, "user not found")
		return
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{"message": "user deleted successfully", "rows_affected": rows})
}

func userIDFromPath(path string, w http.ResponseWriter) (int, bool) {
	path = strings.TrimSuffix(path, "/")
	if !strings.HasPrefix(path, "/users/") {
		WriteError(w, http.StatusNotFound, "not found")
		return 0, false
	}
	s := strings.TrimPrefix(path, "/users/")
	if s == "" {
		WriteError(w, http.StatusNotFound, "not found")
		return 0, false
	}
	id, err := strconv.Atoi(s)
	if err != nil || id < 1 {
		WriteError(w, http.StatusBadRequest, "invalid user id")
		return 0, false
	}
	return id, true
}
