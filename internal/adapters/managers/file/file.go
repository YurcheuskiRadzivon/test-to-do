package filemanage

import (
	"errors"
	"mime/multipart"
	"os"

	"github.com/YurcheuskiRadzivon/test-to-do/pkg/generator"
	"github.com/gofiber/fiber/v2"
)

const (
	mainPath             = "uploaded_files/"
	contentType          = "Content-Type"
	errUnsupportedFormat = "UNSUPPORTED_FORMAT"
)

type FileManager interface {
	UploadFiles(ctx *fiber.Ctx, files []*multipart.FileHeader) ([]string, error)
	UploadFile(ctx *fiber.Ctx, file *multipart.FileHeader) (string, error)
	DownloadFile(ctx *fiber.Ctx, path string) error
	DeleteFile(ctx *fiber.Ctx, path string) error
}

type FileManage struct {
	g generator.UUVGenerator
}

func NewFileManage(g generator.UUVGenerator) *FileManage {
	return &FileManage{
		g: g,
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

	filename := fm.g.NewFileName()

	err := ctx.SaveFile(file, mainPath+filename+format)
	if err != nil {
		return "", err
	}
	uri := filename + format
	return uri, nil
}

func (fm *FileManage) DownloadFile(ctx *fiber.Ctx, path string) error {
	return ctx.SendFile(mainPath + path)
}

func (fm *FileManage) DeleteFile(ctx *fiber.Ctx, path string) error {
	err := os.Remove(mainPath + path)
	if err != nil {
		return err
	}
	return nil
}
