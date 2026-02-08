package main

import (
	"log"
	"net/http"

	"task-api/internal/handlers"
	"task-api/internal/middleware"
)

func main() {
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		middleware.Logger(middleware.APIKeyAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
				case "GET":
					handlers.GetTasks(w, r)
				case "POST":
					handlers.CreateTask(w, r)
				case "PATCH":
					handlers.UpdateTask(w, r)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}))).ServeHTTP(w, r)
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
