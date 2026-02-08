package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"task-api/internal/models"
	"task-api/pkg/storage"
)

type TaskHandler struct{}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{}
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, models.ErrorResponse{Error: message})
}

// GET /tasks
func GetTasks(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	
	if idStr != "" {
		// Get single task by ID
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid id")
			return
		}

		// Find task in array
		for _, task := range storage.Tasks {
			if task.ID == id {
				respondJSON(w, http.StatusOK, task)
				return
			}
		}

		// Task not found
		respondError(w, http.StatusNotFound, "task not found")
		return
	}

	// Get all tasks
	respondJSON(w, http.StatusOK, storage.Tasks)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate title
	title := strings.TrimSpace(req.Title)
	if title == "" {
		respondError(w, http.StatusBadRequest, "invalid title")
		return
	}

	task := models.Task{
		ID:    storage.NextID,
		Title: title,
		Done:  false,
	}
	
	storage.Tasks = append(storage.Tasks, task)
	storage.NextID++

	respondJSON(w, http.StatusCreated, task)
}

// PATCH /tasks?id=X
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "id is required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req models.UpdateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	for i := range storage.Tasks {
		if storage.Tasks[i].ID == id {
			storage.Tasks[i].Done = req.Done
			respondJSON(w, http.StatusOK, models.SuccessResponse{Updated: true})
			return
		}
	}

	// Task not found
	respondError(w, http.StatusNotFound, "task not found")
}
