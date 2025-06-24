package repositories

import (
	"context"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	"github.com/jackc/pgx/v5"
)

type NoteRepository interface {
	CreateNote(ctx context.Context, tx pgx.Tx, note entity.Note) (int, error)
	DeleteNote(ctx context.Context, tx pgx.Tx, noteID int, authorID int) error
	GetNote(ctx context.Context, tx pgx.Tx, noteID int, authorID int) (entity.Note, error)
	GetNotes(ctx context.Context, tx pgx.Tx, authorID int) ([]entity.Note, error)
	UpdateNote(ctx context.Context, tx pgx.Tx, note entity.Note) error
}
