package service

import (
	"errors"
	"practice8/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// ── GetUserByID ──────────────────────────────────────────────────────────────

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Bakytzhan Agai"}
	mockRepo.EXPECT().GetUserByID(1).Return(user, nil)

	result, err := svc.GetUserByID(1)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

// ── CreateUser ───────────────────────────────────────────────────────────────

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	user := &repository.User{ID: 2, Name: "Bakytzhan Agai"}
	mockRepo.EXPECT().CreateUser(user).Return(nil)

	err := svc.CreateUser(user)
	assert.NoError(t, err)
}

// ── RegisterUser ─────────────────────────────────────────────────────────────

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	existing := &repository.User{ID: 1, Name: "Existing", Email: "test@test.com"}
	mockRepo.EXPECT().GetByEmail("test@test.com").Return(existing, nil)

	err := svc.RegisterUser(&repository.User{Name: "New"}, "test@test.com")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestRegisterUser_NewUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	newUser := &repository.User{Name: "New", Email: "new@test.com"}
	mockRepo.EXPECT().GetByEmail("new@test.com").Return(nil, nil)
	mockRepo.EXPECT().CreateUser(newUser).Return(nil)

	err := svc.RegisterUser(newUser, "new@test.com")
	assert.NoError(t, err)
}

func TestRegisterUser_RepositoryErrorOnCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	newUser := &repository.User{Name: "New", Email: "new@test.com"}
	mockRepo.EXPECT().GetByEmail("new@test.com").Return(nil, nil)
	mockRepo.EXPECT().CreateUser(newUser).Return(errors.New("db write error"))

	err := svc.RegisterUser(newUser, "new@test.com")
	assert.Error(t, err)
}

func TestRegisterUser_GetByEmailRepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	mockRepo.EXPECT().GetByEmail("fail@test.com").Return(nil, errors.New("db error"))

	err := svc.RegisterUser(&repository.User{Name: "New"}, "fail@test.com")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting user")
}

// ── UpdateUserName ────────────────────────────────────────────────────────────

func TestUpdateUserName_EmptyName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	err := svc.UpdateUserName(1, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name cannot be empty")
}

func TestUpdateUserName_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	mockRepo.EXPECT().GetUserByID(99).Return(nil, errors.New("user not found"))

	err := svc.UpdateUserName(99, "NewName")
	assert.Error(t, err)
}

func TestUpdateUserName_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	user := &repository.User{ID: 2, Name: "OldName"}
	mockRepo.EXPECT().GetUserByID(2).Return(user, nil)
	// verify name was changed before UpdateUser is called
	mockRepo.EXPECT().UpdateUser(gomock.Any()).DoAndReturn(func(u *repository.User) error {
		assert.Equal(t, "NewName", u.Name)
		return nil
	})

	err := svc.UpdateUserName(2, "NewName")
	assert.NoError(t, err)
}

func TestUpdateUserName_UpdateUserFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	user := &repository.User{ID: 2, Name: "OldName"}
	mockRepo.EXPECT().GetUserByID(2).Return(user, nil)
	mockRepo.EXPECT().UpdateUser(gomock.Any()).DoAndReturn(func(u *repository.User) error {
		assert.Equal(t, "NewName", u.Name)
		return errors.New("update failed")
	})

	err := svc.UpdateUserName(2, "NewName")
	assert.Error(t, err)
}

// ── DeleteUser ────────────────────────────────────────────────────────────────

func TestDeleteUser_AttemptToDeleteAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	err := svc.DeleteUser(1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not allowed to delete admin")
}

func TestDeleteUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	// verify the correct id is passed to repo
	mockRepo.EXPECT().DeleteUser(5).DoAndReturn(func(id int) error {
		assert.Equal(t, 5, id)
		return nil
	})

	err := svc.DeleteUser(5)
	assert.NoError(t, err)
}

func TestDeleteUser_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	mockRepo.EXPECT().DeleteUser(5).Return(errors.New("db delete error"))

	err := svc.DeleteUser(5)
	assert.Error(t, err)
}
