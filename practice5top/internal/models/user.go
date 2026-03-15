package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
}

type PaginatedResponse struct {
	Data       []User `json:"data"`
	TotalCount int    `json:"totalCount"`
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
}

// UserFilter holds dynamic filter + sort params
type UserFilter struct {
	ID        *string
	Name      *string
	Email     *string
	Gender    *string
	BirthDate *string // "YYYY-MM-DD"
	OrderBy   string
	OrderDir  string // "ASC" or "DESC"
	Page      int
	PageSize  int
}
