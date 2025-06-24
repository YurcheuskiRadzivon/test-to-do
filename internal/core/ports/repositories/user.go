package repositories

import (
	"context"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	CreateUser(ctx context.Context, tx pgx.Tx, user entity.User) (int, error)
	DeleteUser(ctx context.Context, tx pgx.Tx, userID int) error
	GetUser(ctx context.Context, tx pgx.Tx, userID int) (string, string, error)
	GetUsers(ctx context.Context, tx pgx.Tx) ([]entity.User, error)
	UpdateUser(ctx context.Context, tx pgx.Tx, user entity.User) error
	GetUserLoginParams(ctx context.Context, tx pgx.Tx, username string) (int, string, error)
	UserExistsByID(ctx context.Context, tx pgx.Tx, id int) (bool, error)
}
