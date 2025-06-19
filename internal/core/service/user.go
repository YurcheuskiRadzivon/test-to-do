package service

import (
	"context"
	"errors"
	"log"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	ports "github.com/YurcheuskiRadzivon/test-to-do/internal/core/ports/repositories"
)

const (
	ErrUpdateUser = "FAILED_UPDATE_USER"
	ErrDeleteUser = "FAILED_DELETE_USER"
	ErrExistUser  = "FAILED_EXISTING_USER"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user entity.User) (int, error) {
	userID, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (s *UserService) GetUser(ctx context.Context, userID int) (string, string, error) {
	username, email, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return "", "", err
	}
	return username, email, nil
}

func (s *UserService) GetUsers(ctx context.Context) ([]entity.User, error) {
	users, err := s.repo.GetUsers(ctx)
	if err != nil {
		return []entity.User{}, err
	}
	return users, nil
}

func (s *UserService) UpdateUser(ctx context.Context, user entity.User) error {
	err := s.repo.UpdateUser(ctx, user)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		return errors.New(ErrUpdateUser)
	}
	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID int) error {
	err := s.repo.DeleteUser(ctx, userID)
	if err != nil {
		log.Printf("Failed to delete user: %v", err)
		return errors.New(ErrUpdateUser)
	}
	return nil
}

func (s *UserService) GetUserLoginParams(ctx context.Context, username string) (int, string, error) {
	id, password, err := s.repo.GetUserLoginParams(ctx, username)
	if err != nil {
		return 0, "", err
	}
	return id, password, nil
}

func (s *UserService) UserExistsByID(ctx context.Context, userID int) (bool, error) {
	exist, err := s.repo.UserExistsByID(ctx, userID)
	if err != nil {
		log.Printf("Failed to check of existing user: %v", err)
		return false, errors.New(ErrExistUser)
	}
	return exist, nil
}
