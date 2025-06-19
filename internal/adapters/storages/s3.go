package storages

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"github.com/minio/minio-go"
)

const (
	ErrFailedReadFile            = "FAILED_TO_READ_FILE"
	ErrFailedCreateObject        = "FAILED_TO_CREATE_OBJECT"
	ErrFailedGettingURLForObject = "FAILED_TO_GETTING_URL_FOR_OBJECT"
	ErrDeletingObject            = "FAILED_DELETE_OBJECT"
)

type S3Storage struct {
	client           *minio.Client
	bucketName       string
	externalEndpoint string
	internalEndpoint string
}

func NewS3Storage(client *minio.Client, bucketName string, externalEndpoint string, internalEndpoint string) *S3Storage {
	return &S3Storage{
		client:           client,
		bucketName:       bucketName,
		externalEndpoint: externalEndpoint,
		internalEndpoint: internalEndpoint,
	}
}

func (s3s *S3Storage) UploadFile(ctx context.Context, objectName string, file *multipart.FileHeader) error {
	fileContent, err := file.Open()
	if err != nil {
		return err
	}

	fileBytes, err := io.ReadAll(fileContent)
	if err != nil {
		log.Printf("Failed to read file")
		return errors.New(ErrFailedReadFile)
	}

	reader := bytes.NewReader(fileBytes)

	_, err = s3s.client.PutObjectWithContext(ctx, s3s.bucketName, objectName, reader, file.Size, minio.PutObjectOptions{})
	if err != nil {
		log.Printf("Failed to create object %s: %v", file.Filename, err)
		return errors.New(ErrFailedCreateObject)
	}
	return nil
}

func (s3s *S3Storage) DownloadFile(objectName string) (string, error) {
	url, err := s3s.client.PresignedGetObject(s3s.bucketName, objectName, time.Second*24*60*60, nil)
	if err != nil {
		log.Printf("failed to getting url for object %s: %v", objectName, err)
		return "", errors.New(ErrFailedGettingURLForObject)
	}
	publicURL := strings.Replace(url.String(), s3s.internalEndpoint, s3s.externalEndpoint, 1)
	log.Println("Generated URL:", url.String())
	log.Println("Replacing:", s3s.internalEndpoint, "→", s3s.externalEndpoint)

	log.Println(publicURL)

	return publicURL, nil
}

func (s3s *S3Storage) DeleteFile(objectName string) error {
	err := s3s.client.RemoveObject(s3s.bucketName, objectName)
	if err != nil {
		log.Printf("failed when delete object %s: %v", objectName, err)
		return errors.New(ErrDeletingObject)
	}
	return nil
}
