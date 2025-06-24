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
	ErrGetNotes = "FAILED_TO_GET_NOTES"
	ErrGetNote  = "FAILED_TO_GET_NOTE"
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

func (nr *NoteRepo) WithTx(tx pgx.Tx) *NoteRepo {
	return &NoteRepo{
		queries: queries.New(tx),
		pool:    nr.pool,
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

func (nr *NoteRepo) CreateNoteWithTx(ctx context.Context, tx pgx.Tx, note entity.Note) (int, error) {
	return nr.WithTx(tx).queries.CreateNote(ctx, queries.CreateNoteParams{
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
		log.Printf("Failed to get notes: %v", err)
		return nil, errors.New(ErrGetNotes)
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
		log.Printf("Failed to get note: %v", err)
		return entity.Note{}, errors.New(ErrGetNote)
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
