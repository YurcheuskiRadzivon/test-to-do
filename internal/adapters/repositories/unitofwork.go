package repositories

import (
	"context"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/ports/repositories"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/infrastructure/database/queries"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UOW struct {
	pool *pgxpool.Pool
}

func NewUOW(pool *pgxpool.Pool) *UOW {
	return &UOW{
		pool: pool,
	}
}

func (u *UOW) BeginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (u *UOW) NoteRepository(tx pgx.Tx) repositories.NoteRepository {
	if tx != nil {
		noteRepoWithTx := NewNoteRepo(
			queries.New(tx),
			u.pool,
		)
		return noteRepoWithTx
	}
	noteRepoWithoutTx := NewNoteRepo(
		queries.New(u.pool),
		u.pool,
	)
	return noteRepoWithoutTx
}
func (u *UOW) UserRepository(tx pgx.Tx) repositories.UserRepository {
	if tx != nil {
		userRepoWithTx := NewUserRepo(
			queries.New(tx),
			u.pool,
		)
		return userRepoWithTx
	}
	userRepoWithoutTx := NewUserRepo(
		queries.New(u.pool),
		u.pool,
	)
	return userRepoWithoutTx
}

func (u *UOW) FileMetaRepository(tx pgx.Tx) repositories.FileMetaRepository {
	if tx != nil {
		fileMetaRepoWithTx := NewFileMetaRepo(
			queries.New(tx),
			u.pool,
		)
		return fileMetaRepoWithTx
	}
	fileMetaRepoWithoutTx := NewFileMetaRepo(
		queries.New(u.pool),
		u.pool,
	)
	return fileMetaRepoWithoutTx
}
