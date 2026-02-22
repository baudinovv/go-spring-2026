package repository

import (
	"practice3/internal/repository/_postgres"
	"practice3/internal/repository/_postgres/users"
	"practice3/pkg/modules"
)

type UserRepository interface {
	GetUsers() ([]modules.User, error)
	GetUserByID(id int) (*modules.User, error)
	CreateUser(in modules.CreateUserInput) (int, error)
	UpdateUser(id int, in modules.UpdateUserInput) error
	DeleteUser(id int) (int64, error)
}

type Repositories struct {
	UserRepository UserRepository
}

func NewRepositories(db *_postgres.Dialect) *Repositories {
	return &Repositories{
		UserRepository: users.NewUserRepository(db),
	}
}

func (r *Repositories) GetUsers() ([]modules.User, error) {
	return r.UserRepository.GetUsers()
}

func (r *Repositories) GetUserByID(id int) (*modules.User, error) {
	return r.UserRepository.GetUserByID(id)
}

func (r *Repositories) CreateUser(in modules.CreateUserInput) (int, error) {
	return r.UserRepository.CreateUser(in)
}

func (r *Repositories) UpdateUser(id int, in modules.UpdateUserInput) error {
	return r.UserRepository.UpdateUser(id, in)
}

func (r *Repositories) DeleteUser(id int) (int64, error) {
	return r.UserRepository.DeleteUser(id)
}
