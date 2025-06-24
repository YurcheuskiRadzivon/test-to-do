package repositories

import (
	"context"
	"errors"
	"log"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/infrastructure/database/queries"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	ErrGetUsers      = "FAILED_TO_GET_USERS"
	ErrGetUser       = "FAILED_TO_GET_USER"
	ErrCreateUser    = "FAILED_CREATING_USER"
	ErrGetUserParams = "FAILED_TO_GET_USER_PARAMS_FOR_LOGIN"
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

func (ur *UserRepo) CreateUser(ctx context.Context, tx pgx.Tx, user entity.User) (int, error) {
	var id int
	var err error
	switch tx {
	case nil:
		id, err = ur.queries.CreateUser(ctx, queries.CreateUserParams{
			Username: user.Username,
			Password: user.Password,
			Email:    user.Email,
		})
		if err != nil {
			log.Printf("Failed to create user: %v", err)
			return 0, errors.New(ErrCreateUser)
		}
	default:
		id, err = ur.queries.WithTx(tx).CreateUser(ctx, queries.CreateUserParams{
			Username: user.Username,
			Password: user.Password,
			Email:    user.Email,
		})
		if err != nil {
			log.Printf("Failed to create user: %v", err)
			return 0, errors.New(ErrCreateUser)
		}
	}

	return id, nil
}

func (ur *UserRepo) DeleteUser(ctx context.Context, tx pgx.Tx, userID int) error {
	if tx != nil {
		return ur.queries.WithTx(tx).DeleteUser(ctx, userID)
	}
	return ur.queries.DeleteUser(ctx, userID)
}

func (ur *UserRepo) GetUser(ctx context.Context, tx pgx.Tx, userID int) (string, string, error) {
	var userFromDB queries.GetUserRow
	var err error

	switch tx {
	case nil:
		userFromDB, err = ur.queries.GetUser(ctx, userID)
		if err != nil {
			log.Printf("Failed to get user: %v", err)
			return "", "", errors.New(ErrGetUser)
		}
	default:
		userFromDB, err = ur.queries.WithTx(tx).GetUser(ctx, userID)
		if err != nil {
			log.Printf("Failed to get user: %v", err)
			return "", "", errors.New(ErrGetUser)
		}
	}

	return userFromDB.Username, userFromDB.Email, nil
}

func (ur *UserRepo) GetUsers(ctx context.Context, tx pgx.Tx) ([]entity.User, error) {
	var usersFromDB []queries.User
	var err error

	switch tx {
	case nil:
		usersFromDB, err = ur.queries.GetUsers(ctx)
		if err != nil {
			log.Printf("Failed to get users: %v", err)
			return nil, errors.New(ErrGetUsers)
		}
	default:
		usersFromDB, err = ur.queries.WithTx(tx).GetUsers(ctx)
		if err != nil {
			log.Printf("Failed to get users: %v", err)
			return nil, errors.New(ErrGetUsers)
		}
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

func (ur *UserRepo) UpdateUser(ctx context.Context, tx pgx.Tx, user entity.User) error {
	if tx != nil {
		return ur.queries.WithTx(tx).UpdateUser(ctx, queries.UpdateUserParams{
			ID:       user.UserID,
			Username: user.Username,
			Password: user.Password,
			Email:    user.Email,
		})
	}
	return ur.queries.UpdateUser(ctx, queries.UpdateUserParams{
		ID:       user.UserID,
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	})
}

func (ur *UserRepo) GetUserLoginParams(ctx context.Context, tx pgx.Tx, username string) (int, string, error) {
	var loginParams queries.GetUserLoginParamsRow
	var err error
	switch tx {
	case nil:
		loginParams, err = ur.queries.GetUserLoginParams(ctx, username)
		if err != nil {
			log.Printf("Failed to get user params: %v", err)
			return -1, "", errors.New(ErrGetUserParams)
		}
	default:
		loginParams, err = ur.queries.WithTx(tx).GetUserLoginParams(ctx, username)
		if err != nil {
			log.Printf("Failed to get user params: %v", err)
			return -1, "", errors.New(ErrGetUserParams)
		}
	}
	return loginParams.ID, loginParams.Password, nil
}

func (ur *UserRepo) UserExistsByID(ctx context.Context, tx pgx.Tx, id int) (bool, error) {
	if tx != nil {
		return ur.queries.WithTx(tx).UserExistsByID(ctx, id)
	}
	return ur.queries.UserExistsByID(ctx, id)
}
