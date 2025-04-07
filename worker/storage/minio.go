package storage

import (
	"bytes"
	"context"
	"errors"
	"log"
	"time"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

func init() {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "minio:9000"
	}
	
	var err error
	for i := 0; i < 10; i++ { // Increased retries to 10
		// Initialize MinIO client
		minioClient, err = minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
			Secure: false,
		})
		if err != nil {
			log.Printf("MinIO client init attempt %d failed: %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}

		// Check bucket existence
		ctx := context.Background()
		exists, err := minioClient.BucketExists(ctx, "processed")
		if err == nil && exists {
			log.Println("MinIO bucket 'processed' already exists")
			return // Success, exit init
		}
		if err != nil {
			log.Printf("Bucket check attempt %d failed: %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}

		// Create bucket if it doesn’t exist
		err = minioClient.MakeBucket(ctx, "processed", minio.MakeBucketOptions{})
		if err == nil {
			log.Println("Created MinIO bucket 'processed'")
			return // Success, exit init
		}
		log.Printf("Bucket creation attempt %d failed: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	// If we get here, all retries failed
	log.Fatal("MinIO initialization failed after retries: ", err)
}

func SaveToMinIO(bucket, objectName string, data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if len(data) == 0 {
		return errors.New("empty image data")
	}

	// Verify we're saving actual JPEG data
    if !bytes.HasPrefix(data, []byte{0xFF, 0xD8}) {
        return errors.New("invalid JPEG data")
    }

	reader := bytes.NewReader(data)
	_, err := minioClient.PutObject(
		ctx,
		bucket,
		objectName,
		reader,
		int64(len(data)),
		minio.PutObjectOptions{
			ContentType: "image/jpeg",
		},
	)
	return err
}
