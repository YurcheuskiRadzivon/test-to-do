package repositories

import (
	"context"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/infrastructure/database/queries"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FileMetaRepo struct {
	pool    *pgxpool.Pool
	queries *queries.Queries
}

func NewFileMetaRepo(db *queries.Queries, pool *pgxpool.Pool) *FileMetaRepo {
	return &FileMetaRepo{
		queries: db,
		pool:    pool,
	}
}

func (fmr *FileMetaRepo) CreateFileMeta(ctx context.Context, fileMeta entity.FileMeta) error {
	return fmr.queries.CreateFileMeta(ctx, queries.CreateFileMetaParams{
		ContentType: fileMeta.ContentType,
		OwnerType:   string(fileMeta.OwnerType),
		OwnerID:     fileMeta.OwnerID,
		UserID:      fileMeta.UserID,
		Uri:         fileMeta.URI,
	})
}
func (fmr *FileMetaRepo) DeleteFileMetaByID(ctx context.Context, id int) error {
	return fmr.queries.DeleteFileMetaByID(ctx, id)
}
func (fmr *FileMetaRepo) DeleteFileMetaByNoteID(ctx context.Context, ownerType entity.OwnerType, ownerID int) error {
	return fmr.queries.DeleteFileMetaByNoteID(ctx, queries.DeleteFileMetaByNoteIDParams{
		OwnerType: string(ownerType),
		OwnerID:   ownerID,
	})
}
func (fmr *FileMetaRepo) FileMetasExistsByIDAndUserID(ctx context.Context, id int, userID int) (bool, error) {
	return fmr.queries.FileMetasExistsByIDAndUserID(ctx, queries.FileMetasExistsByIDAndUserIDParams{
		ID:     id,
		UserID: userID,
	})
}
func (fmr *FileMetaRepo) GetFileMetaByID(ctx context.Context, id int) (entity.FileMeta, error) {
	fileMeta, err := fmr.queries.GetFileMetaByID(ctx, id)
	if err != nil {
		return entity.FileMeta{}, err
	}
	return entity.FileMeta{
		ContentType: fileMeta.ContentType,
		OwnerType:   entity.OwnerType(fileMeta.OwnerType),
		OwnerID:     fileMeta.OwnerID,
		UserID:      fileMeta.UserID,
		URI:         fileMeta.Uri,
	}, nil
}
func (fmr *FileMetaRepo) GetFileMetaURI(ctx context.Context, id int) (string, error) {
	return fmr.queries.GetFileMetaURI(ctx, id)
}

func (fmr *FileMetaRepo) GetFileMetaIDByID(ctx context.Context, ownerType entity.OwnerType, ownerID int) ([]int, error) {
	return fmr.queries.GetFileMetaIDByID(ctx, queries.GetFileMetaIDByIDParams{
		OwnerType: string(ownerType),
		OwnerID:   ownerID,
	})
}

func (fmr *FileMetaRepo) GetFileMetas(ctx context.Context) ([]entity.FileMeta, error) {
	res, err := fmr.queries.GetFileMetas(ctx)
	if err != nil {
		return nil, err
	}

	fileMetas := make([]entity.FileMeta, 0)

	for _, fileMeta := range res {
		fileMetas = append(fileMetas, entity.FileMeta{
			FileID:      fileMeta.ID,
			ContentType: fileMeta.ContentType,
			OwnerType:   entity.OwnerType(fileMeta.OwnerType),
			OwnerID:     fileMeta.ID,
			UserID:      fileMeta.UserID,
			URI:         fileMeta.Uri,
		})
	}

	return fileMetas, nil
}
