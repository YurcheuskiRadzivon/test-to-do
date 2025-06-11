package service

import (
	"context"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	ports "github.com/YurcheuskiRadzivon/test-to-do/internal/core/ports/repositories"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user entity.User) (int, error) {
	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) GetUser(ctx context.Context, userID int) (string, string, error) {
	return s.repo.GetUser(ctx, userID)
}

func (s *UserService) GetUsers(ctx context.Context) ([]entity.User, error) {
	return s.repo.GetUsers(ctx)
}

func (s *UserService) UpdateUser(ctx context.Context, user entity.User) error {
	return s.repo.UpdateUser(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, userID int) error {
	return s.repo.DeleteUser(ctx, userID)
}

func (s *UserService) GetUserLoginParams(ctx context.Context, username string) (int, string, error) {
	return s.repo.GetUserLoginParams(ctx, username)
}

func (s *UserService) UserExistsByID(ctx context.Context, userID int) (bool, error) {
	return s.repo.UserExistsByID(ctx, userID)
}
