package usecase

import (
		repository "practice3/internal/_repository"
	"practice3/pkg/modules"
)

type UserUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) GetUsers() ([]modules.User, error) {
	return u.repo.GetUsers()
}

func (u *UserUsecase) GetUserByID(id int) (*modules.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *UserUsecase) CreateUser(in modules.CreateUserInput) (int, error) {
	return u.repo.CreateUser(in)
}

func (u *UserUsecase) UpdateUser(id int, in modules.UpdateUserInput) error {
	return u.repo.UpdateUser(id, in)
}

func (u *UserUsecase) DeleteUser(id int) (int64, error) {
	return u.repo.DeleteUser(id)
}
