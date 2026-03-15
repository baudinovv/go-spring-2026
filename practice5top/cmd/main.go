package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"practice5/db"
	"practice5/internal/handler"
	"practice5/internal/repository"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:admin@localhost:5432/postgres?sslmode=disable"
	}

	database, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}
	log.Println("Connected to database")

	// Run migrations automatically on startup
	if err := db.RunMigrations(database,
		"migrations/001_create_tables.sql",
		"migrations/002_seed.sql",
	); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	repo := repository.New(database)
	h := handler.New(repo)

	mux := http.NewServeMux()

	// // Endpoint 1: GET /users  — paginated, filtered, sorted
	// mux.HandleFunc("GET /users", h.GetUsers)

	// // Endpoint 2: GET /users/common-friends?user1=<uuid>&user2=<uuid>
	// mux.HandleFunc("GET /users/common-friends", h.GetCommonFriends)

	mux.HandleFunc("/users/common-friends", h.GetCommonFriends) // more specific first!
	mux.HandleFunc("/users", h.GetUsers)

	addr := ":8080"
	fmt.Printf("Server running on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
