package main

import (
	"log"
	"net/http"

	"github.com/NurmagambetovBakytzhan/spring26/cmd/internal/handlers"
)

// get student by id, post create student
// net/http

// middlware - auth

func main() {

	// GET localhost:8080/hello/
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	}) 

	http.HandleFunc("/students", handlers.CreateStudent)

	log.Fatal(http.ListenAndServe(":8080", nil))
}


