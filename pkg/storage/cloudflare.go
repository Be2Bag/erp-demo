package storage

import (
	"fmt"
	"io"
	"os"
	"path" // ใช้กับ key ของ S3/R2
	"strings"
	"time"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// CloudflareStorage provides methods to interact with Cloudflare R2
type CloudflareStorage struct {
	Client        *s3.S3
	Bucket        string
	Region        string
	Endpoint      string
	PublicBaseURL string
}

// NewCloudflareStorage creates a new CloudflareStorage using the given
// configuration. It verifies the connection by listing available buckets.
func NewCloudflareStorage(cfg config.CloudflareConfig) (*CloudflareStorage, error) {

	awsCfg := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, ""),
		Endpoint:         aws.String(cfg.Endpoint),
		Region:           aws.String(cfg.Region),
		S3ForcePathStyle: aws.Bool(true),
	}

	sess, err := session.NewSession(awsCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	client := s3.New(sess)

	// verify connection
	if _, err := client.ListBuckets(&s3.ListBucketsInput{}); err != nil {
		return nil, fmt.Errorf("failed to connect to Cloudflare storage: %w", err)
	}

	return &CloudflareStorage{
		Client:        client,
		Bucket:        cfg.Bucket,
		Region:        cfg.Region,
		Endpoint:      cfg.Endpoint,
		PublicBaseURL: cfg.PublicBaseURL, // ตั้งค่าจาก config
	}, nil
}

// UploadFile uploads the local file at filePath to the bucket under the given key (folder/filename).
func (c *CloudflareStorage) UploadFile(filePath, key string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	_, err = c.Client.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(c.Bucket),
		Key:           aws.String(key),
		Body:          f,
		ContentLength: aws.Int64(info.Size()),
		ContentType:   aws.String(determineContentType(filePath)),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

// ListFileMetas lists metadata of files under the given prefix.
func (c *CloudflareStorage) ListFileMetas(prefix string) ([]dto.FileMeta, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(c.Bucket),
		Prefix: aws.String(prefix),
	}

	result, err := c.Client.ListObjectsV2(input)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	files := make([]dto.FileMeta, 0, len(result.Contents))

	// ถ้ามี PublicBaseURL ให้ใช้เป็นลิงก์สาธารณะ: {PublicBaseURL}/{bucket}/{key}
	// ถ้าไม่มี จะ fallback ไปใช้ endpoint เดิม (อาจเปิดไม่ได้ในเบราว์เซอร์ถ้าเป็น private)
	publicBase := ""
	if c.PublicBaseURL != "" {
		publicBase = fmt.Sprintf("%s/%s", strings.TrimRight(c.PublicBaseURL, "/"), c.Bucket)
	} else {
		publicBase = fmt.Sprintf("%s/%s", strings.TrimRight(c.Endpoint, "/"), c.Bucket)
	}

	for _, item := range result.Contents {
		key := *item.Key
		url := fmt.Sprintf("%s/%s", publicBase, key)
		files = append(files, dto.FileMeta{
			Name: key,
			URL:  url,
			Size: *item.Size,
		})
	}

	return files, nil
}

// DownloadFile retrieves a file from the specified folder and file name.
func (c *CloudflareStorage) DownloadFile(folder, fileName string) ([]byte, error) {
	// ใช้ path.Join เพื่อให้เป็น "/"
	key := path.Join(folder, fileName)

	out, err := c.Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer out.Body.Close()

	data, err := io.ReadAll(out.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	return data, nil
}

// DeleteFile removes the specified file from the bucket.
func (c *CloudflareStorage) DeleteFile(folder, fileName string) error {
	// ใช้ path.Join เพื่อให้เป็น "/"
	key := path.Join(folder, fileName)

	_, err := c.Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	err = c.Client.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("error while waiting for file deletion: %w", err)
	}

	return nil
}

func (c *CloudflareStorage) GetFileURLByName(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("key is required")
	}

	// Normalize to S3-style key separators
	cleanKey := strings.ReplaceAll(key, "\\", "/")
	cleanKey = strings.TrimPrefix(cleanKey, "/")

	var lastErr error
	for i := 0; i < 3; i++ {
		_, err := c.Client.HeadObject(&s3.HeadObjectInput{
			Bucket: aws.String(c.Bucket),
			Key:    aws.String(cleanKey),
		})
		if err == nil {
			base := ""
			if c.PublicBaseURL != "" {
				base = strings.TrimRight(c.PublicBaseURL, "/")
			} else {
				base = strings.TrimRight(c.Endpoint, "/")
			}
			return fmt.Sprintf("%s/%s", base, cleanKey), nil
		}
		lastErr = err
		time.Sleep(time.Duration(150*(i+1)) * time.Millisecond)
	}

	return "", fmt.Errorf("file not found: %w", lastErr)
}
