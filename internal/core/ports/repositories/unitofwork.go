package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type UnitOfWork interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)

	NoteRepository(tx pgx.Tx) NoteRepository
	UserRepository(tx pgx.Tx) UserRepository
	FileMetaRepository(tx pgx.Tx) FileMetaRepository
}
