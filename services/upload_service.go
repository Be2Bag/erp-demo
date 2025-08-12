package services

import (
	"context"
	"fmt"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/pkg/storage"
	"github.com/Be2Bag/erp-demo/ports"
	"go.mongodb.org/mongo-driver/bson"
)

type UpLoadService struct {
	config                   config.Config
	authRepo                 ports.AuthRepository
	upLoadRepo               ports.UpLoadRepository
	storageService           *storage.SupabaseStorage
	storageCloudflareService *storage.CloudflareStorage
	userRepo                 ports.UserRepository
}

func NewUpLoadService(cfg config.Config, authRepo ports.AuthRepository, upLoadRepo ports.UpLoadRepository, storageService *storage.SupabaseStorage, userRepo ports.UserRepository, storageCloudflareService *storage.CloudflareStorage) ports.UpLoadService {
	return &UpLoadService{config: cfg, authRepo: authRepo, upLoadRepo: upLoadRepo, storageService: storageService, userRepo: userRepo, storageCloudflareService: storageCloudflareService}
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

func (s *UpLoadService) DeleteFileByID(ctx context.Context, req dto.RequestDeleteFile) error {

	err := s.storageCloudflareService.DeleteFile(req.Type, req.Name)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	filter := bson.M{
		"user_id":        req.UserID,
		"documents.type": req.Type,
		"deleted_at":     nil,
	}
	update := bson.M{
		"$set": bson.M{
			"documents.$.name":     "",
			"documents.$.file_url": "",
		},
	}

	_, errOnUpdate := s.userRepo.UpdateUserByFilter(ctx, filter, update)
	if errOnUpdate != nil {
		return fmt.Errorf("failed to update user: %w", errOnUpdate)
	}

	return nil
}

func (s *UpLoadService) GetDownloadFile(ctx context.Context, req dto.RequestDownloadFile) ([]byte, error) {
	fileContent, err := s.storageCloudflareService.DownloadFile(req.Type, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	return fileContent, nil
}

func (s *UpLoadService) UploadFileCloudflare(ctx context.Context, filePath, key string) error {

	err := s.storageCloudflareService.UploadFile(filePath, key)
	if err != nil {
		return fmt.Errorf("failed to upload file to storage: %w", err)
	}
	return nil
}

func (s *UpLoadService) GetFileURLCloudflare(ctx context.Context, req dto.RequestGetFile) (string, error) {
	url, err := s.storageCloudflareService.GetFileURLByName(req.Folder + "/" + req.File)
	if err != nil {
		return "", fmt.Errorf("failed to get file URL: %w", err)
	}
	return url, nil
}
