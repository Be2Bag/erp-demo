package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"go.mongodb.org/mongo-driver/bson"
)

type SignJobRepository interface {
	CreateSignJob(ctx context.Context, signJob models.SignJob) error
	UpdateSignJobByJobID(ctx context.Context, jobID string, update models.SignJob) (*models.SignJob, error)
	SoftDeleteSignJobByJobID(ctx context.Context, jobID string) error
	GetAllSignJobByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.SignJob, error)
	GetOneSignJobByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.SignJob, error)
	GetListSignJobsByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.SignJob, int64, error)
	UpdateManySignJobByFilter(ctx context.Context, filter interface{}, update models.SignJob) (int64, error)

	UpdateManySignJobFields(ctx context.Context, filter interface{}, update bson.M) (int64, error)
}

type SignJobService interface {
	CreateSignJob(ctx context.Context, signJob dto.CreateSignJobDTO, claims *dto.JWTClaims) error
	ListSignJobs(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, status string, sortBy string, sortOrder string) (dto.Pagination, error)
	GetSignJobByJobID(ctx context.Context, jobID string, claims *dto.JWTClaims) (*dto.SignJobDTO, error)
	UpdateSignJobByJobID(ctx context.Context, jobID string, update dto.UpdateSignJobDTO, claims *dto.JWTClaims) error
	DeleteSignJobByJobID(ctx context.Context, jobID string, claims *dto.JWTClaims) error
}
