package service

import (
	"context"
	"errors"
	"testovoe/internal/domain"
	"testovoe/internal/repository"
)

var ErrUserNotFound = errors.New("пользователь не найден")
var ErrEmptyFields = errors.New("имя пользователя или email не могут быть пустыми")

type UserServiceInterface interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)
	UpdateUserByID(ctx context.Context, id int64, user *domain.User) error
	DeleteUserByID(ctx context.Context, id int64) error
}

type UserService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User) error {
	if user.Name == "" || user.Email == "" {
		return ErrEmptyFields
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) UpdateUserByID(ctx context.Context, id int64, user *domain.User) error {
	if user.Name == "" || user.Email == "" {
		return ErrEmptyFields
	}
	return s.repo.UpdateUserByID(ctx, id, user)
}

func (s *UserService) DeleteUserByID(ctx context.Context, id int64) error {
	return s.repo.DeleteUserByID(ctx, id)
}
