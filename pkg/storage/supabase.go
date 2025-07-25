package storage

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type SupabaseStorage struct {
	Client     *s3.S3
	Bucket     string
	Region     string
	Endpoint   string
	PublicBase string
}

func NewSupabaseStorage(cfg config.SupabaseConfig) (*SupabaseStorage, error) {
	awsConfig := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, ""),
		Endpoint:         aws.String(cfg.Endpoint),
		Region:           aws.String(cfg.Region),
		S3ForcePathStyle: aws.Bool(true),
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	client := s3.New(sess)

	// Verify connection by listing buckets
	_, err = client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		log.Printf("failed to connect to Supabase storage: %v", err)
		return nil, fmt.Errorf("failed to connect to Supabase storage: %w", err)
	}

	return &SupabaseStorage{
		Client:     client,
		Bucket:     cfg.Bucket,
		Region:     cfg.Region,
		Endpoint:   cfg.Endpoint,
		PublicBase: cfg.PublicBase,
	}, nil
}

func (s *SupabaseStorage) UploadFile(filePath, key string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	_, err = s.Client.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(s.Bucket),
		Key:           aws.String(key),
		Body:          file,
		ContentLength: aws.Int64(fileInfo.Size()),
		ContentType:   aws.String(determineContentType(filePath)),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

func (s *SupabaseStorage) GetFile(key, destinationPath string) error {
	output, err := s.Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to get file: %w", err)
	}
	defer output.Body.Close()

	file, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, output.Body)
	if err != nil {
		return fmt.Errorf("failed to write file to destination: %w", err)
	}

	return nil
}

func (s *SupabaseStorage) ListFileMetas(prefix string) ([]dto.FileMeta, error) {
	var files []dto.FileMeta

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.Bucket),
		Prefix: aws.String(prefix),
	}

	result, err := s.Client.ListObjectsV2(input)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	for _, item := range result.Contents {
		key := *item.Key
		url := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.PublicBase, s.Bucket, key)

		files = append(files, dto.FileMeta{
			Name: key,
			URL:  url,
			Size: *item.Size,
		})
	}

	return files, nil
}

func determineContentType(filePath string) string {
	if strings.HasSuffix(filePath, ".jpg") || strings.HasSuffix(filePath, ".jpeg") {
		return "image/jpeg"
	} else if strings.HasSuffix(filePath, ".png") {
		return "image/png"
	} else if strings.HasSuffix(filePath, ".pdf") {
		return "application/pdf"
	}
	return "application/octet-stream"
}

func (s *SupabaseStorage) GetFileURLByName(key string) (string, error) {
	_, err := s.Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", fmt.Errorf("file not found: %w", err)
	}

	url := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.PublicBase, s.Bucket, key)
	return url, nil
}

func (s *SupabaseStorage) DeleteFile(folder, fileName string) error {
	// รวม folder และ fileName เพื่อสร้าง key ที่ถูกต้อง
	fullKey := fmt.Sprintf("%s/%s", folder, fileName)

	_, err := s.Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(fullKey),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	err = s.Client.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(fullKey),
	})
	if err != nil {
		return fmt.Errorf("error while waiting for file deletion: %w", err)
	}

	return nil
}

func (s *SupabaseStorage) DownloadFile(folder, fileName string) ([]byte, error) {
	fullKey := fmt.Sprintf("%s/%s", folder, fileName)

	output, err := s.Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(fullKey),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer output.Body.Close()

	fileContent, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	return fileContent, nil
}
