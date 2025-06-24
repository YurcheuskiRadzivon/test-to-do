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
	uow ports.UnitOfWork
}

func NewUserService(uow ports.UnitOfWork) *UserService {
	return &UserService{uow: uow}
}

func (us *UserService) CreateUser(ctx context.Context, user entity.User) (int, error) {
	userRepository := us.uow.UserRepository(nil)
	userID, err := userRepository.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (us *UserService) GetUser(ctx context.Context, userID int) (string, string, error) {
	userRepository := us.uow.UserRepository(nil)
	username, email, err := userRepository.GetUser(ctx, userID)
	if err != nil {
		return "", "", err
	}
	return username, email, nil
}

func (us *UserService) GetUsers(ctx context.Context) ([]entity.User, error) {
	userRepository := us.uow.UserRepository(nil)
	users, err := userRepository.GetUsers(ctx)
	if err != nil {
		return []entity.User{}, err
	}
	return users, nil
}

func (us *UserService) UpdateUser(ctx context.Context, user entity.User) error {
	userRepository := us.uow.UserRepository(nil)
	err := userRepository.UpdateUser(ctx, user)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		return errors.New(ErrUpdateUser)
	}
	return nil
}

func (us *UserService) DeleteUser(ctx context.Context, userID int) error {
	userRepository := us.uow.UserRepository(nil)
	err := userRepository.DeleteUser(ctx, userID)
	if err != nil {
		log.Printf("Failed to delete user: %v", err)
		return errors.New(ErrUpdateUser)
	}
	return nil
}

func (us *UserService) GetUserLoginParams(ctx context.Context, username string) (int, string, error) {
	userRepository := us.uow.UserRepository(nil)
	id, password, err := userRepository.GetUserLoginParams(ctx, username)
	if err != nil {
		return 0, "", err
	}
	return id, password, nil
}

func (us *UserService) UserExistsByID(ctx context.Context, userID int) (bool, error) {
	userRepository := us.uow.UserRepository(nil)
	exist, err := userRepository.UserExistsByID(ctx, userID)
	if err != nil {
		log.Printf("Failed to check of existing user: %v", err)
		return false, errors.New(ErrExistUser)
	}
	return exist, nil
}
