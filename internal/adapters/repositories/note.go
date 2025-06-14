package repositories

import (
	"context"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/infrastructure/database/queries"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NoteRepo struct {
	pool    *pgxpool.Pool
	queries *queries.Queries
}

func NewNoteRepo(db *queries.Queries, pool *pgxpool.Pool) *NoteRepo {
	return &NoteRepo{
		queries: db,
		pool:    pool,
	}
}

func (nr *NoteRepo) CreateNote(ctx context.Context, note entity.Note) (int, error) {
	return nr.queries.CreateNote(ctx, queries.CreateNoteParams{
		Title:       note.Title,
		Description: note.Description,
		Status:      note.Status,
		AuthorID:    note.AuthorID,
	})
}

func (nr *NoteRepo) DeleteNote(ctx context.Context, noteID int, authorID int) error {
	return nr.queries.DeleteNote(ctx, queries.DeleteNoteParams{
		ID:       noteID,
		AuthorID: authorID,
	})
}

func (nr *NoteRepo) GetNotes(ctx context.Context, authorID int) ([]entity.Note, error) {
	notesWithoutFormat, err := nr.queries.GetNotes(ctx, authorID)
	if err != nil {
		return nil, err
	}

	notes := make([]entity.Note, 0)

	for _, val := range notesWithoutFormat {
		notes = append(notes, entity.Note{
			NoteID:      val.ID,
			Title:       val.Title,
			Description: val.Description,
			Status:      val.Status,
			AuthorID:    val.AuthorID,
		})
	}

	return notes, nil
}

func (nr *NoteRepo) GetNote(ctx context.Context, noteID int, authorID int) (entity.Note, error) {
	noteWithoutFormat, err := nr.queries.GetNote(ctx, queries.GetNoteParams{
		ID:       noteID,
		AuthorID: authorID,
	})

	if err != nil {
		return entity.Note{}, err
	}

	note := entity.Note{
		NoteID:      noteWithoutFormat.ID,
		Title:       noteWithoutFormat.Title,
		Description: noteWithoutFormat.Description,
		Status:      noteWithoutFormat.Status,
		AuthorID:    noteWithoutFormat.AuthorID,
	}

	return note, nil
}

func (nr *NoteRepo) UpdateNote(ctx context.Context, note entity.Note) error {
	return nr.queries.UpdateNote(ctx, queries.UpdateNoteParams{
		ID:          note.NoteID,
		Title:       note.Title,
		Description: note.Description,
		Status:      note.Description,
		AuthorID:    note.AuthorID,
	})
}
