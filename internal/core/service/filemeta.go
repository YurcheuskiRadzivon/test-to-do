package service

import (
	"context"
	"errors"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	ports "github.com/YurcheuskiRadzivon/test-to-do/internal/core/ports/repositories"
)

type FileMetaService struct {
	repo ports.FileMetaRepository
}

func NewFileMetaService(repo ports.FileMetaRepository) *FileMetaService {
	return &FileMetaService{repo: repo}
}

func (fms *FileMetaService) CreateFileMeta(ctx context.Context, fileMeta entity.FileMeta) error {
	return fms.repo.CreateFileMeta(ctx, fileMeta)
}
func (fms *FileMetaService) DeleteFileMetaByID(ctx context.Context, id int) error {
	return fms.repo.DeleteFileMetaByID(ctx, id)
}
func (fms *FileMetaService) DeleteFileMetaByNoteID(ctx context.Context, ownerType string, ownerID int) error {
	if ownerType != string(entity.OwnerNote) {
		return errors.New("INVALID_OWNER_TYPE")
	}
	return fms.repo.DeleteFileMetaByNoteID(ctx, entity.OwnerNote, ownerID)
}
func (fms *FileMetaService) FileMetasExistsByIDAndUserID(ctx context.Context, id int, userID int) (bool, error) {
	return fms.repo.FileMetasExistsByIDAndUserID(ctx, id, userID)
}
func (fms *FileMetaService) GetFileMetaIDByID(ctx context.Context, ownerType string, ownerID int) ([]int, error) {
	if ownerType != string(entity.OwnerNote) {
		return nil, errors.New("INVALID_OWNER_TYPE")
	}
	return fms.repo.GetFileMetaIDByID(ctx, entity.OwnerNote, ownerID)
}
func (fms *FileMetaService) GetFileMetaByID(ctx context.Context, id int) (entity.FileMeta, error) {
	fileMeta, err := fms.repo.GetFileMetaByID(ctx, id)
	if err != nil {
		return entity.FileMeta{}, err
	}
	return fileMeta, nil
}
func (fms *FileMetaService) GetFileMetaURI(ctx context.Context, id int) (string, error) {
	return fms.repo.GetFileMetaURI(ctx, id)
}
func (fms *FileMetaService) GetFileMetas(ctx context.Context) ([]entity.FileMeta, error) {
	fileMetas, err := fms.repo.GetFileMetas(ctx)
	if err != nil {
		return nil, err
	}

	return fileMetas, nil
}
