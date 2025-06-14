package repositories

import (
	"context"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/infrastructure/database/queries"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool    *pgxpool.Pool
	queries *queries.Queries
}

func NewUserRepo(db *queries.Queries, pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		queries: db,
		pool:    pool,
	}
}

func (ur *UserRepo) CreateUser(ctx context.Context, user entity.User) (int, error) {
	id, err := ur.queries.CreateUser(ctx, queries.CreateUserParams{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (ur *UserRepo) DeleteUser(ctx context.Context, userID int) error {
	return ur.queries.DeleteUser(ctx, userID)
}

func (ur *UserRepo) GetUser(ctx context.Context, userID int) (string, string, error) {
	userFromDB, err := ur.queries.GetUser(ctx, userID)
	if err != nil {
		return "", "", err
	}

	return userFromDB.Username, userFromDB.Email, nil
}

func (ur *UserRepo) GetUsers(ctx context.Context) ([]entity.User, error) {
	usersFromDB, err := ur.queries.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	var users []entity.User
	for _, u := range usersFromDB {
		users = append(users, entity.User{
			UserID:   u.ID,
			Username: u.Username,
			Password: u.Password,
			Email:    u.Email,
		})
	}

	return users, nil
}

func (ur *UserRepo) UpdateUser(ctx context.Context, user entity.User) error {
	return ur.queries.UpdateUser(ctx, queries.UpdateUserParams{
		ID:       user.UserID,
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	})
}

func (ur *UserRepo) GetUserLoginParams(ctx context.Context, username string) (int, string, error) {
	loginParams, err := ur.queries.GetUserLoginParams(ctx, username)
	if err != nil {
		return -1, "", err
	}
	return loginParams.ID, loginParams.Password, nil
}

func (ur *UserRepo) UserExistsByID(ctx context.Context, id int) (bool, error) {
	return ur.queries.UserExistsByID(ctx, id)
}
