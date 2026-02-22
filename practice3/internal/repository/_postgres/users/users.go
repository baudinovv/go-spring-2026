package users

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"practice3/internal/repository/_postgres"
	"practice3/pkg/modules"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type Repository struct {
	db               *_postgres.Dialect
	executionTimeout time.Duration
}

func NewUserRepository(db *_postgres.Dialect) *Repository {
	return &Repository{
		db:               db,
		executionTimeout: 5 * time.Second,
	}
}

func (r *Repository) GetUsers() ([]modules.User, error) {
	var list []modules.User
	query := `SELECT id, name, email, created_at, updated_at FROM users`
	err := r.db.DB.Select(&list, query)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) GetUserByID(id int) (*modules.User, error) {
	var u modules.User
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.DB.Get(&u, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *Repository) CreateUser(in modules.CreateUserInput) (int, error) {
	query := `INSERT INTO users (name, email, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING id`
	var id int
	err := r.db.DB.QueryRow(query, in.Name, in.Email).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) UpdateUser(id int, in modules.UpdateUserInput) error {
	// Build dynamic update: only set non-nil fields
	args := []interface{}{id}
	updates := []string{"updated_at = NOW()"}
	pos := 2
	if in.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", pos))
		args = append(args, *in.Name)
		pos++
	}
	if in.Email != nil {
		updates = append(updates, fmt.Sprintf("email = $%d", pos))
		args = append(args, *in.Email)
		pos++
	}
	if len(updates) <= 1 {
		return nil // nothing to update
	}
	query := "UPDATE users SET " + strings.Join(updates, ", ") + " WHERE id = $1"
	result, err := r.db.DB.Exec(query, args...)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *Repository) DeleteUser(id int) (int64, error) {
	result, err := r.db.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
