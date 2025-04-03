package storage

import (
	"bytes"
	"context"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

func init() {
	var err error
	minioClient, err = minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal("MinIO init failed:", err)
	}
}

func SaveToMinIO(bucket, objectName string, data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reader := bytes.NewReader(data)
	_, err := minioClient.PutObject(
		ctx,
		bucket,
		objectName,
		reader,
		int64(len(data)),
		minio.PutObjectOptions{},
	)
	return err
}
