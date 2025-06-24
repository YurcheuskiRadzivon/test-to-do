package service

import (
	"context"
	"errors"
	"log"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/managers/transaction"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	ports "github.com/YurcheuskiRadzivon/test-to-do/internal/core/ports/repositories"
)

const (
	ErrUpdateUser = "FAILED_UPDATE_USER"
	ErrDeleteUser = "FAILED_DELETE_USER"
	ErrExistUser  = "FAILED_EXISTING_USER"
)

type UserService struct {
	repoU     ports.UserRepository
	txManager transaction.TransactionManager
}

func NewUserService(repoU ports.UserRepository, txManager transaction.TransactionManager) *UserService {
	return &UserService{
		repoU:     repoU,
		txManager: txManager,
	}
}

func (us *UserService) CreateUser(ctx context.Context, user entity.User) (int, error) {
	userID, err := us.repoU.CreateUser(ctx, nil, user)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (us *UserService) GetUser(ctx context.Context, userID int) (string, string, error) {
	username, email, err := us.repoU.GetUser(ctx, nil, userID)
	if err != nil {
		return "", "", err
	}
	return username, email, nil
}

func (us *UserService) GetUsers(ctx context.Context) ([]entity.User, error) {
	users, err := us.repoU.GetUsers(ctx, nil)
	if err != nil {
		return []entity.User{}, err
	}
	return users, nil
}

func (us *UserService) UpdateUser(ctx context.Context, user entity.User) error {
	err := us.repoU.UpdateUser(ctx, nil, user)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		return errors.New(ErrUpdateUser)
	}
	return nil
}

func (us *UserService) DeleteUser(ctx context.Context, userID int) error {
	err := us.repoU.DeleteUser(ctx, nil, userID)
	if err != nil {
		log.Printf("Failed to delete user: %v", err)
		return errors.New(ErrUpdateUser)
	}
	return nil
}

func (us *UserService) GetUserLoginParams(ctx context.Context, username string) (int, string, error) {
	id, password, err := us.repoU.GetUserLoginParams(ctx, nil, username)
	if err != nil {
		return 0, "", err
	}
	return id, password, nil
}

func (us *UserService) UserExistsByID(ctx context.Context, userID int) (bool, error) {
	exist, err := us.repoU.UserExistsByID(ctx, nil, userID)
	if err != nil {
		log.Printf("Failed to check of existing user: %v", err)
		return false, errors.New(ErrExistUser)
	}
	return exist, nil
}
