package services

import (
	"context"
	"time"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type signJobService struct {
	config      config.Config
	signJobRepo ports.SignJobRepository
}

func NewSignJobService(cfg config.Config, signJobRepo ports.SignJobRepository) ports.SignJobService {
	return &signJobService{config: cfg, signJobRepo: signJobRepo}
}

func (s *signJobService) CreateSignJob(ctx context.Context, signJob dto.CreateSignJobDTO, claims *dto.JWTClaims) error {

	var dueDate time.Time
	if signJob.DueDate != "" {
		parsedDate, err := time.Parse("2006-01-02", signJob.DueDate)
		if err != nil {
			return err
		}
		dueDate = parsedDate
	}

	now := time.Now()
	model := models.SignJob{
		JobID:          uuid.NewString(),
		ProjectName:    signJob.ProjectName,
		JobName:        signJob.JobName,
		CustomerName:   signJob.CustomerName,
		ContactPerson:  signJob.ContactPerson,
		Phone:          signJob.Phone,
		Email:          signJob.Email,
		CustomerTypeID: signJob.CustomerTypeID,
		Address:        signJob.Address,
		SignTypeID:     signJob.SignTypeID,
		Size:           signJob.Size,
		Quantity:       signJob.Quantity,
		Content:        signJob.Content,
		MainColor:      signJob.MainColor,
		DesignOption:   signJob.DesignOption,
		ProductionTime: signJob.ProductionTime,
		DueDate:        dueDate,
		InstallOption:  signJob.InstallOption,
		Notes:          signJob.Notes,
		Status:         "draft",
		CreatedBy:      claims.UserID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.signJobRepo.CreateSignJob(ctx, model); err != nil {
		return err
	}

	return nil
}

func (s *signJobService) ListSignJobs(ctx context.Context, claims *dto.JWTClaims, page, size int, search string) (dto.Pagination, error) {
	items, total, err := s.signJobRepo.ListSignJobs(ctx, claims.UserID, page, size, search)
	if err != nil {
		return dto.Pagination{}, err
	}

	list := make([]interface{}, 0, len(items))
	for _, m := range items {
		list = append(list, dto.SignJobDTO{
			JobID:          m.JobID,
			ProjectName:    m.ProjectName,
			JobName:        m.JobName,
			CustomerName:   m.CustomerName,
			ContactPerson:  m.ContactPerson,
			Phone:          m.Phone,
			Email:          m.Email,
			CustomerTypeID: m.CustomerTypeID,
			Address:        m.Address,
			SignTypeID:     m.SignTypeID,
			Size:           m.Size,
			Quantity:       m.Quantity,
			Content:        m.Content,
			MainColor:      m.MainColor,
			DesignOption:   m.DesignOption,
			ProductionTime: m.ProductionTime,
			DueDate:        m.DueDate,
			InstallOption:  m.InstallOption,
			Notes:          m.Notes,
			Status:         m.Status,
			CreatedBy:      m.CreatedBy,
			CreatedAt:      m.CreatedAt,
			UpdatedAt:      m.UpdatedAt,
		})
	}

	if size <= 0 {
		size = 20
	}
	totalPages := (int(total) + size - 1) / size

	p := dto.Pagination{
		Page:       page,
		Size:       size,
		TotalCount: int(total),
		TotalPages: totalPages,
		List:       list,
	}
	return p, nil
}

func (s *signJobService) GetSignJobByJobID(ctx context.Context, jobID string, claims *dto.JWTClaims) (*dto.SignJobDTO, error) {
	m, err := s.signJobRepo.GetSignJobByJobID(ctx, jobID, claims.UserID)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}
	dto := &dto.SignJobDTO{
		JobID:          m.JobID,
		ProjectName:    m.ProjectName,
		JobName:        m.JobName,
		CustomerName:   m.CustomerName,
		ContactPerson:  m.ContactPerson,
		Phone:          m.Phone,
		Email:          m.Email,
		CustomerTypeID: m.CustomerTypeID,
		Address:        m.Address,
		SignTypeID:     m.SignTypeID,
		Size:           m.Size,
		Quantity:       m.Quantity,
		Content:        m.Content,
		MainColor:      m.MainColor,
		DesignOption:   m.DesignOption,
		ProductionTime: m.ProductionTime,
		DueDate:        m.DueDate,
		InstallOption:  m.InstallOption,
		Notes:          m.Notes,
		Status:         m.Status,
		CreatedBy:      m.CreatedBy,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
	return dto, nil
}

func (s *signJobService) UpdateSignJobByJobID(ctx context.Context, jobID string, update dto.UpdateSignJobDTO, claims *dto.JWTClaims) error {

	var dueDate time.Time
	if update.DueDate != "" {
		parsedDate, err := time.Parse("2006-01-02", update.DueDate)
		if err != nil {
			return err
		}
		dueDate = parsedDate
	}

	now := time.Now()
	model := models.SignJob{
		ProjectName:    update.ProjectName,
		JobName:        update.JobName,
		CustomerName:   update.CustomerName,
		ContactPerson:  update.ContactPerson,
		Phone:          update.Phone,
		Email:          update.Email,
		CustomerTypeID: update.CustomerTypeID,
		Address:        update.Address,
		SignTypeID:     update.SignTypeID,
		Size:           update.Size,
		Quantity:       update.Quantity,
		Content:        update.Content,
		MainColor:      update.MainColor,
		DesignOption:   update.DesignOption,
		ProductionTime: update.ProductionTime,
		DueDate:        dueDate,
		InstallOption:  update.InstallOption,
		Notes:          update.Notes,
		Status:         "draft",
		UpdatedAt:      now,
	}
	updated, err := s.signJobRepo.UpdateSignJobByJobID(ctx, jobID, claims.UserID, model)
	if err != nil {
		return err
	}
	if updated == nil {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (s *signJobService) DeleteSignJobByJobID(ctx context.Context, jobID string, claims *dto.JWTClaims) error {
	err := s.signJobRepo.DeleteSignJobByJobID(ctx, jobID, claims.UserID)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}
