package service

import (
	"context"
	"errors"
	"log"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	ports "github.com/YurcheuskiRadzivon/test-to-do/internal/core/ports/repositories"
)

const (
	contentType = "Content-Type"

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
	uow ports.UnitOfWork
}

func NewNoteService(uow ports.UnitOfWork) *NoteService {
	return &NoteService{uow: uow}
}

func (ns *NoteService) CreateNoteWithFilesWithTx(
	ctx context.Context,
	note entity.Note,
	uriList []string,
	filesContentType []string,
	userID int,
) error {
	tx, err := ns.uow.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	NoteRepository := ns.uow.NoteRepository(tx)
	FileMetaRepository := ns.uow.FileMetaRepository(tx)

	noteID, err := NoteRepository.CreateNote(ctx, note)
	if err != nil {
		log.Printf("Failed to create note: %v", err)
		return errors.New(ErrCreateNote)
	}

	for i := range uriList {
		err := FileMetaRepository.CreateFileMeta(ctx, entity.FileMeta{
			ContentType: filesContentType[i],
			OwnerType:   entity.OwnerNote,
			OwnerID:     noteID,
			UserID:      userID,
			URI:         uriList[i],
		})
		if err != nil {
			log.Printf("Failed to create fileMeta: %v", err)
			return errors.New(ErrCreateNote)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return errors.New("FAILED_TO_END_CREATING_NOTE")
	}

	return nil
}

func (ns *NoteService) GetNote(ctx context.Context, noteID int, authorID int) (entity.Note, error) {
	NoteRepository := ns.uow.NoteRepository(nil)
	note, err := NoteRepository.GetNote(ctx, noteID, authorID)
	if err != nil {
		return entity.Note{}, err
	}
	return note, nil
}

func (ns *NoteService) GetNotes(ctx context.Context, authorID int) ([]entity.Note, error) {
	NoteRepository := ns.uow.NoteRepository(nil)
	notes, err := NoteRepository.GetNotes(ctx, authorID)
	if err != nil {
		return []entity.Note{}, err
	}
	return notes, nil
}

func (ns *NoteService) UpdateNote(ctx context.Context, note entity.Note) error {
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
	NoteRepository := ns.uow.NoteRepository(nil)
	if err := NoteRepository.UpdateNote(ctx, note); err != nil {
		log.Printf("Failed to update note: %v", err)
		return errors.New(ErrUpdateNote)
	}
	return nil
}

func (ns *NoteService) DeleteNote(ctx context.Context, noteID int, authorID int) error {
	if noteID <= 0 {
		return errors.New(ErrInvalidIDFormat)
	}
	NoteRepository := ns.uow.NoteRepository(nil)
	if err := NoteRepository.DeleteNote(ctx, noteID, authorID); err != nil {
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
