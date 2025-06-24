package service

import (
	"context"
	"errors"
	"log"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	ports "github.com/YurcheuskiRadzivon/test-to-do/internal/core/ports/repositories"
)

const (
	ErrInvalidOwnerType = "INVALID_OWNER_TYPE"
	ErrGetFileMeta      = "FAILED_TO_GET_FILE_META"
)

type FileMetaService struct {
	uow ports.UnitOfWork
}

func NewFileMetaService(uow ports.UnitOfWork) *FileMetaService {
	return &FileMetaService{uow: uow}
}

func (fms *FileMetaService) CreateFileMeta(ctx context.Context, fileMeta entity.FileMeta) error {
	fileMetaRepository := fms.uow.FileMetaRepository(nil)
	return fileMetaRepository.CreateFileMeta(ctx, fileMeta)
}

func (fms *FileMetaService) DeleteFileMetaByID(ctx context.Context, id int) error {
	fileMetaRepository := fms.uow.FileMetaRepository(nil)
	return fileMetaRepository.DeleteFileMetaByID(ctx, id)
}

func (fms *FileMetaService) DeleteFileMetaByNoteID(ctx context.Context, ownerType string, ownerID int) error {
	if ownerType != string(entity.OwnerNote) {
		log.Printf("Failed owner type: %v - have type", ownerType)
		return errors.New(ErrInvalidOwnerType)
	}

	fileMetaRepository := fms.uow.FileMetaRepository(nil)

	return fileMetaRepository.DeleteFileMetaByNoteID(ctx, entity.OwnerNote, ownerID)
}

func (fms *FileMetaService) FileMetasExistsByIDAndUserID(ctx context.Context, id int, userID int) (bool, error) {
	fileMetaRepository := fms.uow.FileMetaRepository(nil)
	return fileMetaRepository.FileMetasExistsByIDAndUserID(ctx, id, userID)
}

func (fms *FileMetaService) GetFileMetaIDByID(ctx context.Context, ownerType string, ownerID int) ([]int, error) {
	if ownerType != string(entity.OwnerNote) {
		log.Printf("Failed owner type: %v - have type", ownerType)
		return nil, errors.New(ErrInvalidOwnerType)
	}

	fileMetaRepository := fms.uow.FileMetaRepository(nil)

	filemetas, err := fileMetaRepository.GetFileMetaIDByID(ctx, entity.OwnerNote, ownerID)
	if err != nil {
		log.Printf("Failed to get meta id by id: %v", err)
		return nil, errors.New(ErrGetFileMeta)
	}
	return filemetas, nil
}

func (fms *FileMetaService) GetFileMetaByID(ctx context.Context, id int) (entity.FileMeta, error) {
	fileMetaRepository := fms.uow.FileMetaRepository(nil)

	fileMeta, err := fileMetaRepository.GetFileMetaByID(ctx, id)
	if err != nil {
		return entity.FileMeta{}, err
	}
	return fileMeta, nil
}

func (fms *FileMetaService) GetFileMetaURI(ctx context.Context, id int) (string, error) {
	fileMetaRepository := fms.uow.FileMetaRepository(nil)
	return fileMetaRepository.GetFileMetaURI(ctx, id)
}

func (fms *FileMetaService) GetFileMetas(ctx context.Context) ([]entity.FileMeta, error) {
	fileMetaRepository := fms.uow.FileMetaRepository(nil)
	fileMetas, err := fileMetaRepository.GetFileMetas(ctx)
	if err != nil {
		return nil, err
	}

	return fileMetas, nil
}

func (fms *FileMetaService) GetFileMetasIDByUserID(ctx context.Context, userID int) ([]int, error) {
	fileMetaRepository := fms.uow.FileMetaRepository(nil)
	filemetas, err := fileMetaRepository.GetFileMetasIDByUserID(ctx, userID)
	if err != nil {
		log.Printf("Failed to get meta id by id: %v", err)
		return nil, errors.New(ErrGetFileMeta)
	}
	return filemetas, nil
}
