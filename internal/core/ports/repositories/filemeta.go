package repositories

import (
	"context"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
)

type FileMetaRepository interface {
	CreateFileMeta(ctx context.Context, fileMeta entity.FileMeta) error
	DeleteFileMetaByID(ctx context.Context, id int) error
	DeleteFileMetaByNoteID(ctx context.Context, ownerType entity.OwnerType, ownerID int) error
	FileMetasExistsByIDAndUserID(ctx context.Context, id int, userID int) (bool, error)
	GetFileMetaIDByID(ctx context.Context, ownerType entity.OwnerType, ownerID int) ([]int, error)
	GetFileMetaByID(ctx context.Context, id int) (entity.FileMeta, error)
	GetFileMetaURI(ctx context.Context, id int) (string, error)
	GetFileMetas(ctx context.Context) ([]entity.FileMeta, error)
	GetFileMetasIDByUserID(ctx context.Context, userID int) ([]int, error)
}
