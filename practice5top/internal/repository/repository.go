package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"practice5/internal/models"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// allowedColumns whitelist for safe dynamic filtering & sorting
var allowedColumns = map[string]string{
	"id":         "u.id::text",
	"name":       "u.name",
	"email":      "u.email",
	"gender":     "u.gender",
	"birth_date": "u.birth_date::text",
}

var allowedSortColumns = map[string]bool{
	"id":         true,
	"name":       true,
	"email":      true,
	"gender":     true,
	"birth_date": true,
}

// GetPaginatedUsers returns paginated, filtered and sorted users.
func (r *Repository) GetPaginatedUsers(f models.UserFilter) (models.PaginatedResponse, error) {
	args := []interface{}{}
	argIdx := 1

	whereClause, args, argIdx := buildWhereClause(f, args, argIdx)

	// Validate & build ORDER BY
	orderBy := "u.id" // default
	orderDir := "ASC"
	if f.OrderBy != "" {
		if _, ok := allowedSortColumns[f.OrderBy]; ok {
			orderBy = "u." + f.OrderBy
		}
	}
	if strings.ToUpper(f.OrderDir) == "DESC" {
		orderDir = "DESC"
	}

	// Count query
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM users u %s`, whereClause)
	var totalCount int
	err := r.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return models.PaginatedResponse{}, fmt.Errorf("count query: %w", err)
	}

	// Paginate args
	offset := (f.Page - 1) * f.PageSize
	args = append(args, f.PageSize, offset)

	dataQuery := fmt.Sprintf(`
		SELECT u.id, u.name, u.email, u.gender, u.birth_date
		FROM users u
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d`,
		whereClause,
		orderBy, orderDir,
		argIdx, argIdx+1,
	)

	rows, err := r.db.Query(dataQuery, args...)
	if err != nil {
		return models.PaginatedResponse{}, fmt.Errorf("data query: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return models.PaginatedResponse{}, err
		}
		users = append(users, u)
	}

	return models.PaginatedResponse{
		Data:       users,
		TotalCount: totalCount,
		Page:       f.Page,
		PageSize:   f.PageSize,
	}, nil
}

// buildWhereClause constructs a safe WHERE clause from filter fields.
func buildWhereClause(f models.UserFilter, args []interface{}, argIdx int) (string, []interface{}, int) {
	conditions := []string{}

	addFilter := func(field string, val *string) {
		if val == nil || *val == "" {
			return
		}
		col := allowedColumns[field]
		if field == "name" || field == "email" {
			args = append(args, "%"+*val+"%")
			conditions = append(conditions, fmt.Sprintf("%s ILIKE $%d", col, argIdx))
		} else {
			args = append(args, *val)
			conditions = append(conditions, fmt.Sprintf("%s = $%d", col, argIdx))
		}
		argIdx++
	}

	// Convert UUID string pointer for id
	addFilter("id", f.ID)
	addFilter("name", f.Name)
	addFilter("email", f.Email)
	addFilter("gender", f.Gender)
	addFilter("birth_date", f.BirthDate)

	if len(conditions) == 0 {
		return "", args, argIdx
	}
	return "WHERE " + strings.Join(conditions, " AND "), args, argIdx
}

// GetCommonFriends returns common friends of two users using a single JOIN query (no N+1).
func (r *Repository) GetCommonFriends(userID1, userID2 string) ([]models.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.gender, u.birth_date
		FROM users u
		JOIN user_friends uf1 ON uf1.friend_id = u.id AND uf1.user_id = $1
		JOIN user_friends uf2 ON uf2.friend_id = u.id AND uf2.user_id = $2
		ORDER BY u.name
	`

	rows, err := r.db.Query(query, userID1, userID2)
	if err != nil {
		return nil, fmt.Errorf("common friends query: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
