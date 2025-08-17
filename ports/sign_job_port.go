package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
)

type SignJobRepository interface {
	CreateSignJob(ctx context.Context, signJob models.SignJob) error
	ListSignJobs(ctx context.Context, createdBy string, page, size int, search string) ([]models.SignJob, int64, error)
	GetSignJobByJobID(ctx context.Context, jobID string, createdBy string) (*models.SignJob, error)
	UpdateSignJobByJobID(ctx context.Context, jobID string, createdBy string, update models.SignJob) (*models.SignJob, error)
	DeleteSignJobByJobID(ctx context.Context, jobID string, createdBy string) error
}

type SignJobService interface {
	CreateSignJob(ctx context.Context, signJob dto.CreateSignJobDTO, claims *dto.JWTClaims) error
	ListSignJobs(ctx context.Context, claims *dto.JWTClaims, page, size int, search string) (dto.Pagination, error)
	GetSignJobByJobID(ctx context.Context, jobID string, claims *dto.JWTClaims) (*dto.SignJobDTO, error)
	UpdateSignJobByJobID(ctx context.Context, jobID string, update dto.UpdateSignJobDTO, claims *dto.JWTClaims) error
	DeleteSignJobByJobID(ctx context.Context, jobID string, claims *dto.JWTClaims) error
}
