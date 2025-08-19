package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"go.mongodb.org/mongo-driver/bson"
)

type ProjectService interface {
	CreateProject(ctx context.Context, createProject dto.CreateProjectDTO, claims *dto.JWTClaims) error
	ListProject(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, sortBy string, sortOrder string) (dto.Pagination, error)
	GetProjectByID(ctx context.Context, projectID string, claims *dto.JWTClaims) (*dto.ProjectDTO, error)
	UpdateProjectByID(ctx context.Context, projectID string, update dto.UpdateProjectDTO, claims *dto.JWTClaims) error
	DeleteProjectByID(ctx context.Context, projectID string, claims *dto.JWTClaims) error
}
type ProjectRepository interface {
	CreateProject(ctx context.Context, project models.Project) error
	GetListProjectByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.Project, int64, error)
	GetOneProjectByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.Project, error)
	UpdateProjectByID(ctx context.Context, projectID string, update models.Project) (*models.Project, error)
	SoftDeleteProjectByID(ctx context.Context, projectID string, claims *dto.JWTClaims) error
}
