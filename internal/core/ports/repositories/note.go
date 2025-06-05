package repositories

import (
	"context"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
)

type NoteRepository interface {
	CreateNote(ctx context.Context, note entity.Note) error
	DeleteNote(ctx context.Context, noteID int, authorID int) error
	GetNote(ctx context.Context, noteID int, authorID int) (entity.Note, error)
	GetNotes(ctx context.Context, authorID int) ([]entity.Note, error)
	UpdateNote(ctx context.Context, note entity.Note) error
}
