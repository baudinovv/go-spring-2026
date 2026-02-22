package modules

import "time"

type User struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// CreateUserInput is used for creating a new user (ID and timestamps are generated).
type CreateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UpdateUserInput is used for updating an existing user (partial update).
type UpdateUserInput struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}
