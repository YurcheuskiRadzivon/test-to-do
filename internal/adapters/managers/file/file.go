package filemanage

import (
	"errors"
	"mime/multipart"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/storages"
	"github.com/YurcheuskiRadzivon/test-to-do/pkg/generator"
	"github.com/gofiber/fiber/v2"
)

const (
	mainPath    = "uploaded_files/"
	contentType = "Content-Type"

	errUnsupportedFormat = "UNSUPPORTED_FORMAT"
)

type FileManager interface {
	UploadFiles(ctx *fiber.Ctx, files []*multipart.FileHeader) ([]string, error)
	UploadFile(ctx *fiber.Ctx, file *multipart.FileHeader) (string, error)
	DownloadFile(ctx *fiber.Ctx, objectName string) (string, error)
	DeleteFile(ctx *fiber.Ctx, objectName string) error
}

type FileManage struct {
	g       generator.UUVGenerator
	storage storages.FileStorage
}

func NewFileManage(g generator.UUVGenerator, storage storages.FileStorage) *FileManage {
	return &FileManage{
		g:       g,
		storage: storage,
	}
}

func (fm *FileManage) UploadFiles(ctx *fiber.Ctx, files []*multipart.FileHeader) ([]string, error) {
	uriList := make([]string, 0)
	for _, file := range files {
		uri, err := fm.UploadFile(ctx, file)
		if err != nil {
			return nil, err
		}
		uriList = append(uriList, uri)
	}
	return uriList, nil
}

func (fm *FileManage) UploadFile(ctx *fiber.Ctx, file *multipart.FileHeader) (string, error) {
	format := ""
	switch file.Header.Get(contentType) {
	case "image/jpeg":
		format = ".jpg"
	case "application/pdf":
		format = ".pdf"
	case "image/png":
		format = ".png"
	default:
		return "", errors.New(errUnsupportedFormat)
	}

	objectName := fm.g.NewFileName() + format
	err := fm.storage.UploadFile(ctx.Context(), objectName, file)
	if err != nil {
		return "", err
	}

	return objectName, nil
}

func (fm *FileManage) DownloadFile(ctx *fiber.Ctx, objectName string) (string, error) {
	uri, err := fm.storage.DownloadFile(objectName)
	if err != nil {
		return "", err
	}
	return uri, nil
}

func (fm *FileManage) DeleteFile(ctx *fiber.Ctx, objectName string) error {
	err := fm.storage.DeleteFile(objectName)
	if err != nil {
		return err
	}
	return nil
}
