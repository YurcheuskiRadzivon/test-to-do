package minioclient

import (
	"fmt"
	"log"

	"github.com/minio/minio-go"
)

func NewMinioClientAndDebug(internalEndpoint string, accsessKey string, secretKey string, bucketName string) (*minio.Client, error) {
	minioClient, err := minio.New(
		internalEndpoint,
		accsessKey,
		secretKey,
		false,
	)
	if err != nil {
		return nil, fmt.Errorf("minio: %v", err)
	}

	//Create bucket
	exists, err := minioClient.BucketExists(bucketName)
	if err != nil {
		return nil, fmt.Errorf("minio: create bucket: %v", err)
	}
	if !exists {
		err := minioClient.MakeBucket(bucketName, "")
		if err != nil {
			return nil, fmt.Errorf("minio: create bucket: %v", err)
		}
	}

	//Debug connection test
	buckets, err := minioClient.ListBuckets()
	if err != nil {
		return nil, fmt.Errorf("minio: cannot connect to s3: %v", err)

	}
	log.Println("Succesfully connections. Buckets:", buckets)
	return minioClient, nil
}
