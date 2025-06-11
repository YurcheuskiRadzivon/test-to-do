package service

import (
	"context"
	"errors"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	ports "github.com/YurcheuskiRadzivon/test-to-do/internal/core/ports/repositories"
)

const (
	statusSuccessfully = "SUCCESS"
	statusInProgress   = "IN_PROGRESS"
	statusNotStart     = "NOT_START"
)

var (
	ErrInvalidStatusFormat = errors.New("INVALID_STATUS_FORMAT")
	ErrInvalidIDFormat     = errors.New("INVALID_ID_FORMAT")
	ErrInvalidTitleFormat  = errors.New("INVALID_Title_FORMAT")
)

type NoteService struct {
	repo ports.NoteRepository
}

func NewNoteService(repo ports.NoteRepository) *NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) CreateNote(ctx context.Context, note entity.Note) error {
	return s.repo.CreateNote(ctx, note)
}

func (s *NoteService) GetNote(ctx context.Context, noteID int, authorID int) (entity.Note, error) {
	return s.repo.GetNote(ctx, noteID, authorID)
}

func (s *NoteService) GetNotes(ctx context.Context, authorID int) ([]entity.Note, error) {
	return s.repo.GetNotes(ctx, authorID)
}

func (s *NoteService) UpdateNote(ctx context.Context, note entity.Note) error {
	if CheckStatus(note.Status) == false {
		return ErrInvalidStatusFormat
	}
	if note.NoteID <= 0 {
		return ErrInvalidIDFormat
	}
	if note.Title == "" {
		return ErrInvalidTitleFormat
	}
	return s.repo.UpdateNote(ctx, note)
}

func (s *NoteService) DeleteNote(ctx context.Context, noteID int, authorID int) error {
	if noteID <= 0 {
		return ErrInvalidIDFormat
	}
	return s.repo.DeleteNote(ctx, noteID, authorID)
}

func CheckStatus(status string) bool {
	switch status {
	case statusInProgress, statusNotStart, statusSuccessfully:
		return true
	default:
		return false
	}
}
