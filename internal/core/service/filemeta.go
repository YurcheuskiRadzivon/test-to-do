package service

import (
	"context"
	"errors"
	"log"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/managers/transaction"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	ports "github.com/YurcheuskiRadzivon/test-to-do/internal/core/ports/repositories"
)

const (
	ErrInvalidOwnerType = "INVALID_OWNER_TYPE"
	ErrGetFileMeta      = "FAILED_TO_GET_FILE_META"
)

type FileMetaService struct {
	repoFM    ports.FileMetaRepository
	txManager transaction.TransactionManager
}

func NewFileMetaService(repoFM ports.FileMetaRepository, txManager transaction.TransactionManager) *FileMetaService {
	return &FileMetaService{
		repoFM:    repoFM,
		txManager: txManager,
	}
}

func (fms *FileMetaService) CreateFileMeta(ctx context.Context, fileMeta entity.FileMeta) error {
	return fms.repoFM.CreateFileMeta(ctx, nil, fileMeta)
}

func (fms *FileMetaService) DeleteFileMetaByID(ctx context.Context, id int) error {
	return fms.repoFM.DeleteFileMetaByID(ctx, nil, id)
}

func (fms *FileMetaService) DeleteFileMetaByNoteID(ctx context.Context, ownerType string, ownerID int) error {
	if ownerType != string(entity.OwnerNote) {
		log.Printf("Failed owner type: %v - have type", ownerType)
		return errors.New(ErrInvalidOwnerType)
	}

	return fms.repoFM.DeleteFileMetaByNoteID(ctx, nil, entity.OwnerNote, ownerID)
}

func (fms *FileMetaService) FileMetasExistsByIDAndUserID(ctx context.Context, id int, userID int) (bool, error) {
	return fms.repoFM.FileMetasExistsByIDAndUserID(ctx, nil, id, userID)
}

func (fms *FileMetaService) GetFileMetaIDByID(ctx context.Context, ownerType string, ownerID int) ([]int, error) {
	if ownerType != string(entity.OwnerNote) {
		log.Printf("Failed owner type: %v - have type", ownerType)
		return nil, errors.New(ErrInvalidOwnerType)
	}

	filemetas, err := fms.repoFM.GetFileMetaIDByID(ctx, nil, entity.OwnerNote, ownerID)
	if err != nil {
		log.Printf("Failed to get meta id by id: %v", err)
		return nil, errors.New(ErrGetFileMeta)
	}
	return filemetas, nil
}

func (fms *FileMetaService) GetFileMetaByID(ctx context.Context, id int) (entity.FileMeta, error) {
	fileMeta, err := fms.repoFM.GetFileMetaByID(ctx, nil, id)
	if err != nil {
		return entity.FileMeta{}, err
	}
	return fileMeta, nil
}

func (fms *FileMetaService) GetFileMetaURI(ctx context.Context, id int) (string, error) {
	return fms.repoFM.GetFileMetaURI(ctx, nil, id)
}

func (fms *FileMetaService) GetFileMetas(ctx context.Context) ([]entity.FileMeta, error) {
	fileMetas, err := fms.repoFM.GetFileMetas(ctx, nil)
	if err != nil {
		return nil, err
	}

	return fileMetas, nil
}

func (fms *FileMetaService) GetFileMetasIDByUserID(ctx context.Context, userID int) ([]int, error) {
	filemetas, err := fms.repoFM.GetFileMetasIDByUserID(ctx, nil, userID)
	if err != nil {
		log.Printf("Failed to get meta id by id: %v", err)
		return nil, errors.New(ErrGetFileMeta)
	}
	return filemetas, nil
}
