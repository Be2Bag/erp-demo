package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type projectService struct {
	config      config.Config
	projectRepo ports.ProjectRepository
	userRepo    ports.UserRepository
	signJobRepo ports.SignJobRepository
	taskRepo    ports.TaskRepository
}

func NewProjectService(cfg config.Config, projectRepo ports.ProjectRepository, userRepo ports.UserRepository, signJobRepo ports.SignJobRepository, taskRepo ports.TaskRepository) ports.ProjectService {
	return &projectService{config: cfg, projectRepo: projectRepo, userRepo: userRepo, signJobRepo: signJobRepo, taskRepo: taskRepo}
}

func (s *projectService) CreateProject(ctx context.Context, createProject dto.CreateProjectDTO, claims *dto.JWTClaims) error {

	now := time.Now()

	model := models.Project{
		ProjectID:   uuid.NewString(),
		ProjectName: createProject.ProjectName,
		CreatedBy:   claims.UserID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.projectRepo.CreateProject(ctx, model); err != nil {
		return err
	}
	return nil

}

func (s *projectService) ListProject(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, sortBy string, sortOrder string) (dto.Pagination, error) {
	skip := int64((page - 1) * size)
	limit := int64(size)

	filter := bson.M{
		"deleted_at": nil,
	}

	search = strings.TrimSpace(search)
	if search != "" {
		safe := regexp.QuoteMeta(search)
		re := primitive.Regex{Pattern: safe, Options: "i"}
		filter["$or"] = []bson.M{
			{"project_name": re},
		}
	}

	projection := bson.M{}

	// sort
	allowedSortFields := map[string]string{
		"created_at":   "created_at",
		"updated_at":   "updated_at",
		"project_name": "project_name",
	}

	field, ok := allowedSortFields[sortBy]
	if !ok || field == "" {
		field = "created_at"
	}
	order := int32(-1)
	if strings.EqualFold(sortOrder, "asc") {
		order = 1
	}

	sort := bson.D{
		{Key: field, Value: order},
		{Key: "_id", Value: -1},
	}

	items, total, err := s.projectRepo.GetListProjectByFilter(ctx, filter, projection, sort, skip, limit)
	if err != nil {
		return dto.Pagination{}, fmt.Errorf("list projects: %w", err)
	}

	list := make([]interface{}, 0, len(items))
	for _, m := range items {

		createdByName := "ไม่พบผู้สร้าง"
		createdBy, _ := s.userRepo.GetByID(ctx, m.CreatedBy)
		if createdBy != nil {
			createdByName = fmt.Sprintf("%s %s %s", createdBy.TitleTH, createdBy.FirstNameTH, createdBy.LastNameTH)
		}

		list = append(list, dto.ProjectDTO{
			ProjectID:   m.ProjectID,
			ProjectName: m.ProjectName,
			CreatedBy:   createdByName,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
			DeletedAt:   m.DeletedAt,
		})
	}

	totalPages := 0
	if total > 0 && size > 0 {
		totalPages = int((total + int64(size) - 1) / int64(size))
	}

	return dto.Pagination{
		Page:       page,
		Size:       size,
		TotalCount: int(total),
		TotalPages: totalPages,
		List:       list,
	}, nil
}

func (s *projectService) GetProjectByID(ctx context.Context, projectID string, claims *dto.JWTClaims) (*dto.ProjectDTO, error) {

	filter := bson.M{"project_id": projectID}
	projection := bson.M{}

	m, err := s.projectRepo.GetOneProjectByFilter(ctx, filter, projection)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}

	createdByName := "ไม่พบผู้สร้าง"
	createdBy, _ := s.userRepo.GetByID(ctx, m.CreatedBy)
	if createdBy != nil {
		createdByName = fmt.Sprintf("%s %s %s", createdBy.TitleTH, createdBy.FirstNameTH, createdBy.LastNameTH)
	}

	dtoObj := &dto.ProjectDTO{
		ProjectID:   m.ProjectID,
		ProjectName: m.ProjectName,
		CreatedBy:   createdByName,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
	}
	return dtoObj, nil
}

func (s *projectService) UpdateProjectByID(ctx context.Context, projectID string, update dto.UpdateProjectDTO, claims *dto.JWTClaims) error {

	filter := bson.M{"project_id": projectID, "deleted_at": nil}
	existing, err := s.projectRepo.GetOneProjectByFilter(ctx, filter, bson.M{})
	if err != nil {
		return err
	}
	if existing == nil {
		return mongo.ErrNoDocuments
	}

	if update.ProjectName != "" {
		existing.ProjectName = update.ProjectName
	}

	if update.Note != "" {
		existing.Note = &update.Note
	}

	existing.UpdatedAt = time.Now()

	updated, err := s.projectRepo.UpdateProjectByID(ctx, projectID, *existing)
	if err != nil {
		return err
	}
	if updated == nil {
		return mongo.ErrNoDocuments
	}

	filterSignJob := bson.M{"project_id": existing.ProjectID}
	partialSignJobUpdate := bson.M{"project_name": existing.ProjectName}
	_, errOnUpdateSignJob := s.signJobRepo.UpdateManySignJobFields(ctx, filterSignJob, partialSignJobUpdate)
	if errOnUpdateSignJob != nil {
		return errOnUpdateSignJob
	}

	filterTask := bson.M{"project_id": existing.ProjectID}
	partialTaskUpdate := bson.M{"project_name": existing.ProjectName}

	_, errOnUpdateTask := s.taskRepo.UpdateManyTaskFields(ctx, filterTask, partialTaskUpdate)
	if errOnUpdateTask != nil {
		return errOnUpdateTask
	}

	return nil
}

func (s *projectService) DeleteProjectByID(ctx context.Context, projectID string, claims *dto.JWTClaims) error {
	err := s.projectRepo.SoftDeleteProjectByID(ctx, projectID, claims)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}
