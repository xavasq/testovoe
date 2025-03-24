package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"testovoe/internal/domain"
	"testovoe/internal/repository"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUserByID(ctx context.Context, id int64, user *domain.User) error {
	args := m.Called(ctx, id, user)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUserByID(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &domain.User{Name: "Test User", Email: "test@example.com"}
	mockRepo.On("CreateUser", mock.Anything, user).Return(nil)

	err := service.CreateUser(context.Background(), user)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_EmptyFields(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &domain.User{Name: "", Email: "test@example.com"}
	err := service.CreateUser(context.Background(), user)
	assert.Equal(t, ErrEmptyFields, err)

	user = &domain.User{Name: "Test User", Email: ""}
	err = service.CreateUser(context.Background(), user)
	assert.Equal(t, ErrEmptyFields, err)

	mockRepo.AssertNotCalled(t, "CreateUser")
}

func TestGetUserByID_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	expectedUser := &domain.User{ID: 1, Name: "Test User", Email: "test@example.com"}
	mockRepo.On("GetUserByID", mock.Anything, int64(1)).Return(expectedUser, nil)

	user, err := service.GetUserByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	mockRepo.On("GetUserByID", mock.Anything, int64(1)).Return((*domain.User)(nil), repository.ErrUserNotFound)

	user, err := service.GetUserByID(context.Background(), 1)
	assert.Equal(t, repository.ErrUserNotFound, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUserByID_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &domain.User{Name: "Updated User", Email: "updated@example.com"}
	mockRepo.On("UpdateUserByID", mock.Anything, int64(1), user).Return(nil)

	err := service.UpdateUserByID(context.Background(), 1, user)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUserByID_EmptyFields(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &domain.User{Name: "", Email: "updated@example.com"}
	err := service.UpdateUserByID(context.Background(), 1, user)
	assert.Equal(t, ErrEmptyFields, err)

	user = &domain.User{Name: "Updated User", Email: ""}
	err = service.UpdateUserByID(context.Background(), 1, user)
	assert.Equal(t, ErrEmptyFields, err)

	mockRepo.AssertNotCalled(t, "UpdateUserByID")
}

func TestDeleteUserByID_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	mockRepo.On("DeleteUserByID", mock.Anything, int64(1)).Return(nil)

	err := service.DeleteUserByID(context.Background(), 1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUserByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	mockRepo.On("DeleteUserByID", mock.Anything, int64(1)).Return(repository.ErrUserNotFound)

	err := service.DeleteUserByID(context.Background(), 1)
	assert.Equal(t, repository.ErrUserNotFound, err)
	mockRepo.AssertExpectations(t)
}
