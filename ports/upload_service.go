package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
)

type UpLoadService interface {
	UploadFile(ctx context.Context, filePath, key string) error
	ListFiles(ctx context.Context, prefix string) ([]dto.FileMeta, error)
	GetFileURL(ctx context.Context, req dto.RequestGetFile) (string, error)
	DeleteFileByID(ctx context.Context, req dto.RequestDeleteFile) error
	GetDownloadFile(ctx context.Context, req dto.RequestDownloadFile) ([]byte, error)
}
