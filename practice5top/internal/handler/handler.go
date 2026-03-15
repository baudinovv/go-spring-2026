package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"practice5/internal/models"
	"practice5/internal/repository"
)

type Handler struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// GET /users?page=1&page_size=10&order_by=name&order_dir=asc&name=john&gender=male
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}

	f := models.UserFilter{
		OrderBy:  q.Get("order_by"),
		OrderDir: strings.ToUpper(q.Get("order_dir")),
		Page:     page,
		PageSize: pageSize,
	}

	if v := q.Get("id"); v != "" {
		f.ID = &v
	}
	if v := q.Get("name"); v != "" {
		f.Name = &v
	}
	if v := q.Get("email"); v != "" {
		f.Email = &v
	}
	if v := q.Get("gender"); v != "" {
		f.Gender = &v
	}
	if v := q.Get("birth_date"); v != "" {
		f.BirthDate = &v
	}

	resp, err := h.repo.GetPaginatedUsers(f)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// GET /users/common-friends?user1=<uuid>&user2=<uuid>
func (h *Handler) GetCommonFriends(w http.ResponseWriter, r *http.Request) {
	user1 := r.URL.Query().Get("user1")
	user2 := r.URL.Query().Get("user2")

	if user1 == "" || user2 == "" {
		writeError(w, http.StatusBadRequest, "user1 and user2 query params are required")
		return
	}
	if user1 == user2 {
		writeError(w, http.StatusBadRequest, "user1 and user2 must be different")
		return
	}

	friends, err := h.repo.GetCommonFriends(user1, user2)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user1":          user1,
		"user2":          user2,
		"common_friends": friends,
		"count":          len(friends),
	})
}
