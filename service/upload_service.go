package service

import (
	"context"
	"fmt"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/pkg/storage"
	"github.com/Be2Bag/erp-demo/ports"
)

type UpLoadService struct {
	config         config.Config
	authRepo       ports.AuthRepository
	upLoadRepo     ports.UpLoadRepository
	storageService *storage.SupabaseStorage
}

func NewUpLoadService(cfg config.Config, authRepo ports.AuthRepository, upLoadRepo ports.UpLoadRepository, storageService *storage.SupabaseStorage) ports.UpLoadService {
	return &UpLoadService{config: cfg, authRepo: authRepo, upLoadRepo: upLoadRepo, storageService: storageService}
}

func (s *UpLoadService) UploadFile(ctx context.Context, filePath, key string) error {

	err := s.storageService.UploadFile(filePath, key)
	if err != nil {
		return fmt.Errorf("failed to upload file to storage: %w", err)
	}
	return nil
}

func (s *UpLoadService) ListFiles(ctx context.Context, prefix string) ([]dto.FileMeta, error) {
	files, err := s.storageService.ListFileMetas(prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}
	return files, nil
}

func (s *UpLoadService) GetFileURL(ctx context.Context, req dto.RequestGetFile) (string, error) {
	url, err := s.storageService.GetFileURLByName(req.Folder + "/" + req.File)
	if err != nil {
		return "", fmt.Errorf("failed to get file URL: %w", err)
	}
	return url, nil
}
