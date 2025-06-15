package filemanage

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/gofiber/fiber/v2"
)

type FileManager interface {
	UploadFiles(ctx *fiber.Ctx, files []*multipart.FileHeader) ([]string, error)
	UploadFile(ctx *fiber.Ctx, file *multipart.FileHeader) (string, error)
	DownloadFile(ctx *fiber.Ctx, path string) error
	DeleteFile(ctx *fiber.Ctx, path string) error
}

type FileManage struct{}

func NewFileManage() *FileManage {
	return &FileManage{}
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
	err := ctx.SaveFile(file, "uploaded_files/"+file.Filename)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Ошибка при сохранении файла %s: %v", file.Filename, err))
	}
	uri := file.Filename
	return uri, nil
}

func (fm *FileManage) DownloadFile(ctx *fiber.Ctx, path string) error {
	return ctx.SendFile("uploaded_files/" + path)
}

func (fm *FileManage) DeleteFile(ctx *fiber.Ctx, path string) error {
	err := os.Remove("uploaded_files/" + path)
	if err != nil {
		return err
	}
	return nil
}
