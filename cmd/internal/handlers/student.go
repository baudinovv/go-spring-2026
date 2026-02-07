package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"strconv"
)

// Create student
type Student struct {
	ID int `json:"id"`;
	Name string `json:"name"`;
	Year int `json:"year"`;
}

var Students []Student = make([]Student, 0);

type CreateStudentDTO struct {
	Name string `json:"name"`
	Year int `json:"year"`;
}

// reqd request body from user
// take some fields to the DTO
// Get student by id
func CreateStudent(w http.ResponseWriter, r *http.Request) {
	var studentToCreate CreateStudentDTO
	if err := json.NewDecoder(r.Body).Decode(&studentToCreate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	student := Student{
		Name: studentToCreate.Name,
		Year: studentToCreate.Year,
		ID: len(Students) + 1,
	}
	
	Students = append(Students, student)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(student)

}

func GetStudentById(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	
	// Validate URL structure (expecting /students/{id})
	if len(parts) < 3 || parts[len(parts)-1] == "" {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}
	
	// Extract ID from URL
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid student ID format", http.StatusBadRequest)
		return
	}
	
	// Find student (assuming you have a students slice or database)
	var student *Student
	for i := range Students {
		if Students[i].ID == id {
			student = &Students[i]
			break
		}
	}
	
	if student == nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	
	// Return student as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}