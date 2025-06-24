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
	ErrGetFileMetaByID = "FAILED_TO GET_FILEMETA_BY_ID"
	ErrGetFileMetas    = "FAILED_TO_GET_FILEMETAS"
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

func (fmr *FileMetaRepo) CreateFileMeta(ctx context.Context, tx pgx.Tx, fileMeta entity.FileMeta) error {
	if tx != nil {
		return fmr.queries.WithTx(tx).CreateFileMeta(ctx, queries.CreateFileMetaParams{
			ContentType: fileMeta.ContentType,
			OwnerType:   string(fileMeta.OwnerType),
			OwnerID:     fileMeta.OwnerID,
			UserID:      fileMeta.UserID,
			Uri:         fileMeta.URI,
		})
	}
	return fmr.queries.CreateFileMeta(ctx, queries.CreateFileMetaParams{
		ContentType: fileMeta.ContentType,
		OwnerType:   string(fileMeta.OwnerType),
		OwnerID:     fileMeta.OwnerID,
		UserID:      fileMeta.UserID,
		Uri:         fileMeta.URI,
	})
}
func (fmr *FileMetaRepo) DeleteFileMetaByID(ctx context.Context, tx pgx.Tx, id int) error {
	if tx != nil {
		return fmr.queries.WithTx(tx).DeleteFileMetaByID(ctx, id)
	}
	return fmr.queries.DeleteFileMetaByID(ctx, id)
}

func (fmr *FileMetaRepo) DeleteFileMetaByNoteID(ctx context.Context, tx pgx.Tx, ownerType entity.OwnerType, ownerID int) error {
	if tx != nil {
		return fmr.queries.WithTx(tx).DeleteFileMetaByNoteID(ctx, queries.DeleteFileMetaByNoteIDParams{
			OwnerType: string(ownerType),
			OwnerID:   ownerID,
		})
	}
	return fmr.queries.DeleteFileMetaByNoteID(ctx, queries.DeleteFileMetaByNoteIDParams{
		OwnerType: string(ownerType),
		OwnerID:   ownerID,
	})
}

func (fmr *FileMetaRepo) FileMetasExistsByIDAndUserID(ctx context.Context, tx pgx.Tx, id int, userID int) (bool, error) {
	if tx != nil {
		return fmr.queries.WithTx(tx).FileMetasExistsByIDAndUserID(ctx, queries.FileMetasExistsByIDAndUserIDParams{
			ID:     id,
			UserID: userID,
		})
	}
	return fmr.queries.FileMetasExistsByIDAndUserID(ctx, queries.FileMetasExistsByIDAndUserIDParams{
		ID:     id,
		UserID: userID,
	})
}

func (fmr *FileMetaRepo) GetFileMetaByID(ctx context.Context, tx pgx.Tx, id int) (entity.FileMeta, error) {
	var fileMeta queries.GetFileMetaByIDRow
	var err error
	switch tx {
	case nil:
		fileMeta, err = fmr.queries.GetFileMetaByID(ctx, id)
		if err != nil {
			log.Printf("Failed to get file meta by ID: %v", err)
			return entity.FileMeta{}, errors.New(ErrGetFileMetaByID)
		}
	default:
		fileMeta, err = fmr.queries.WithTx(tx).GetFileMetaByID(ctx, id)
		if err != nil {
			log.Printf("Failed to get file meta by ID: %v", err)
			return entity.FileMeta{}, errors.New(ErrGetFileMetaByID)
		}
	}

	return entity.FileMeta{
		ContentType: fileMeta.ContentType,
		OwnerType:   entity.OwnerType(fileMeta.OwnerType),
		OwnerID:     fileMeta.OwnerID,
		UserID:      fileMeta.UserID,
		URI:         fileMeta.Uri,
	}, nil
}

func (fmr *FileMetaRepo) GetFileMetaURI(ctx context.Context, tx pgx.Tx, id int) (string, error) {
	if tx != nil {
		return fmr.queries.WithTx(tx).GetFileMetaURI(ctx, id)
	}
	return fmr.queries.GetFileMetaURI(ctx, id)
}

func (fmr *FileMetaRepo) GetFileMetaIDByID(ctx context.Context, tx pgx.Tx, ownerType entity.OwnerType, ownerID int) ([]int, error) {
	if tx != nil {
		return fmr.queries.WithTx(tx).GetFileMetaIDByID(ctx, queries.GetFileMetaIDByIDParams{
			OwnerType: string(ownerType),
			OwnerID:   ownerID,
		})
	}
	return fmr.queries.GetFileMetaIDByID(ctx, queries.GetFileMetaIDByIDParams{
		OwnerType: string(ownerType),
		OwnerID:   ownerID,
	})
}

func (fmr *FileMetaRepo) GetFileMetas(ctx context.Context, tx pgx.Tx) ([]entity.FileMeta, error) {
	var res []queries.Filemeta
	var err error
	switch tx {
	case nil:
		res, err = fmr.queries.GetFileMetas(ctx)
		if err != nil {
			log.Printf("Failed to get file metas: %v", err)
			return nil, errors.New(ErrGetFileMetas)
		}
	default:
		res, err = fmr.queries.WithTx(tx).GetFileMetas(ctx)
		if err != nil {
			log.Printf("Failed to get file metas: %v", err)
			return nil, errors.New(ErrGetFileMetas)
		}
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

func (fmr *FileMetaRepo) GetFileMetasIDByUserID(ctx context.Context, tx pgx.Tx, userID int) ([]int, error) {
	if tx != nil {
		return fmr.queries.WithTx(tx).GetFileMetasIDByUserID(ctx, userID)
	}
	return fmr.queries.GetFileMetasIDByUserID(ctx, userID)

}
