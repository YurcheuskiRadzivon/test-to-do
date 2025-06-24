package repositories

import (
	"context"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	"github.com/jackc/pgx/v5"
)

type FileMetaRepository interface {
	CreateFileMeta(ctx context.Context, tx pgx.Tx, fileMeta entity.FileMeta) error
	DeleteFileMetaByID(ctx context.Context, tx pgx.Tx, id int) error
	DeleteFileMetaByNoteID(ctx context.Context, tx pgx.Tx, ownerType entity.OwnerType, ownerID int) error
	FileMetasExistsByIDAndUserID(ctx context.Context, tx pgx.Tx, id int, userID int) (bool, error)
	GetFileMetaIDByID(ctx context.Context, tx pgx.Tx, ownerType entity.OwnerType, ownerID int) ([]int, error)
	GetFileMetaByID(ctx context.Context, tx pgx.Tx, id int) (entity.FileMeta, error)
	GetFileMetaURI(ctx context.Context, tx pgx.Tx, id int) (string, error)
	GetFileMetas(ctx context.Context, tx pgx.Tx) ([]entity.FileMeta, error)
	GetFileMetasIDByUserID(ctx context.Context, tx pgx.Tx, userID int) ([]int, error)
}
