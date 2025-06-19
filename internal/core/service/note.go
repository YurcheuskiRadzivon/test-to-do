package service

import (
	"context"
	"errors"
	"log"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	ports "github.com/YurcheuskiRadzivon/test-to-do/internal/core/ports/repositories"
)

const (
	statusSuccessfully = "SUCCESS"
	statusInProgress   = "IN_PROGRESS"
	statusNotStart     = "NOT_START"

	ErrCreateNote          = "FAILED_CREATE_NOTE"
	ErrInvalidStatusFormat = "INVALID_STATUS_FORMAT"
	ErrInvalidIDFormat     = "INVALID_ID_FORMAT"
	ErrInvalidTitleFormat  = "INVALID_TITLE_FORMAT"
	ErrUpdateNote          = "FAILED_UPDATE_NOTE"
	ErrDeleteNote          = "FAILED_DELETE_NOTE"
)

type NoteService struct {
	repo ports.NoteRepository
}

func NewNoteService(repo ports.NoteRepository) *NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) CreateNote(ctx context.Context, note entity.Note) (int, error) {
	id, err := s.repo.CreateNote(ctx, note)
	if err != nil {
		log.Printf("Failed to create note: %v", err)
		return 0, errors.New(ErrCreateNote)
	}
	return id, nil
}

func (s *NoteService) GetNote(ctx context.Context, noteID int, authorID int) (entity.Note, error) {
	note, err := s.repo.GetNote(ctx, noteID, authorID)
	if err != nil {
		return entity.Note{}, err
	}
	return note, nil
}

func (s *NoteService) GetNotes(ctx context.Context, authorID int) ([]entity.Note, error) {
	notes, err := s.repo.GetNotes(ctx, authorID)
	if err != nil {
		return []entity.Note{}, err
	}
	return notes, nil
}

func (s *NoteService) UpdateNote(ctx context.Context, note entity.Note) error {
	if CheckStatus(note.Status) == false {
		log.Printf("Failed to update note: %v - check status", CheckStatus(note.Status))
		return errors.New(ErrInvalidStatusFormat)
	}
	if note.NoteID <= 0 {
		log.Printf("Failed to update note: %v - noteID", note.NoteID)
		return errors.New(ErrInvalidIDFormat)
	}
	if note.Title == "" {
		log.Printf("Failed to update note: %v - note title", note.Title)
		return errors.New(ErrInvalidTitleFormat)
	}

	if err := s.repo.UpdateNote(ctx, note); err != nil {
		log.Printf("Failed to update note: %v", err)
		return errors.New(ErrUpdateNote)
	}
	return nil
}

func (s *NoteService) DeleteNote(ctx context.Context, noteID int, authorID int) error {
	if noteID <= 0 {
		return errors.New(ErrInvalidIDFormat)
	}
	if err := s.repo.DeleteNote(ctx, noteID, authorID); err != nil {
		log.Printf("Failed to delete note: %v", err)
		return errors.New(ErrDeleteNote)
	}
	return nil
}

func CheckStatus(status string) bool {
	switch status {
	case statusInProgress, statusNotStart, statusSuccessfully:
		return true
	default:
		return false
	}
}
