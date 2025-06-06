package repositories

import (
	"context"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user entity.User) (int, error)
	DeleteUser(ctx context.Context, userID int) error
	GetUser(ctx context.Context, userID int) (string, string, error)
	GetUsers(ctx context.Context) ([]entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) error
	GetUserLoginParams(ctx context.Context, username string) (int, string, error)
	UserExistsByID(ctx context.Context, id int) (bool, error)
}
