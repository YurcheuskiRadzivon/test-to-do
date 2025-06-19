package storages

import (
	"context"
	"errors"
	"log"
	"mime/multipart"
	"os"

	"github.com/valyala/fasthttp"
)

const (
	ErrSavingFile   = "FALED_TO_SAVE_FILE"
	ErrDeletingFile = "FALED_TO_DELETE_FILE"
)

type FSStorage struct {
	path               string
	fsEndpointExternal string
	appDomen           string
}

func NewFSStorage(path string, fsEndpointExternal string, appDomen string) *FSStorage {
	return &FSStorage{
		path:               path,
		fsEndpointExternal: fsEndpointExternal,
		appDomen:           appDomen,
	}
}

func (fss *FSStorage) UploadFile(ctx context.Context, objectName string, file *multipart.FileHeader) error {
	err := fasthttp.SaveMultipartFile(file, fss.path+objectName)
	if err != nil {
		log.Printf("Faled to save file in file system: %v", err)
		return errors.New(ErrSavingFile)
	}
	return nil
}

func (fss *FSStorage) DownloadFile(objectName string) (string, error) {
	return fss.appDomen + fss.fsEndpointExternal + objectName, nil
}

func (fss *FSStorage) DeleteFile(objectName string) error {
	err := os.Remove(fss.path + objectName)
	if err != nil {
		log.Printf("Faled to delete file in file system: %v", err)
		return errors.New(ErrDeletingFile)
	}
	return nil
}
