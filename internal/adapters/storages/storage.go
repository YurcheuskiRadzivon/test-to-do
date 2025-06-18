package storages

import (
	"context"
	"mime/multipart"
)

type FileStorage interface {
	UploadFile(ctx context.Context, objectName string, file *multipart.FileHeader) error
	DownloadFile(objectName string) (string, error)
	DeleteFile(objectName string) error
}
