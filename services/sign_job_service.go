package services

import (
	"context"
	"time"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	now := time.Now()
	var due time.Time
	if signJob.DueDate != "" {
		parsedDate, err := time.Parse("2006-01-02", signJob.DueDate)
		if err != nil {
			return err
		}
		due = parsedDate
	}

	status := signJob.Status
	if status == "" {
		status = "แผนกออกแบบกราฟิก"
	}

	model := models.SignJob{
		ID:             primitive.NewObjectID(),
		JobID:          uuid.NewString(),
		CompanyName:    signJob.CompanyName,
		ContactPerson:  signJob.ContactPerson,
		Phone:          signJob.Phone,
		Email:          signJob.Email,
		CustomerTypeID: signJob.CustomerTypeID,
		Address:        signJob.Address,

		ProjectName: signJob.ProjectName,
		JobName:     signJob.JobName,
		SignTypeID:  signJob.SignTypeID,
		Width:       signJob.Width,
		Height:      signJob.Height,
		Quantity:    signJob.Quantity,
		PriceTHB:    signJob.PriceTHB,
		Content:     signJob.Content,
		MainColor:   signJob.MainColor,

		PaymentMethod:  signJob.PaymentMethod,
		ProductionTime: signJob.ProductionTime,
		DueDate:        due,

		DesignOption:  signJob.DesignOption,
		InstallOption: signJob.InstallOption,
		Notes:         signJob.Notes,

		Status:    status,
		CreatedBy: claims.UserID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.signJobRepo.CreateSignJob(ctx, model); err != nil {
		return err
	}
	return nil
}

func (s *signJobService) ListSignJobs(ctx context.Context, claims *dto.JWTClaims, page, size int, search string) (dto.Pagination, error) {
	if size <= 0 {
		size = 10
	}
	if page <= 0 {
		page = 1
	}

	items, total, err := s.signJobRepo.ListSignJobs(ctx, page, size, search)
	if err != nil {
		return dto.Pagination{}, err
	}

	list := make([]interface{}, 0, len(items))
	for _, m := range items {
		var duePtr *time.Time
		if !m.DueDate.IsZero() {
			d := m.DueDate
			duePtr = &d
		}
		list = append(list, dto.SignJobDTO{
			ID:             m.ID.Hex(),
			JobID:          m.JobID,
			CompanyName:    m.CompanyName,
			ContactPerson:  m.ContactPerson,
			Phone:          m.Phone,
			Email:          m.Email,
			CustomerTypeID: m.CustomerTypeID,
			Address:        m.Address,
			ProjectName:    m.ProjectName,
			JobName:        m.JobName,
			SignTypeID:     m.SignTypeID,
			Width:          m.Width,
			Height:         m.Height,
			Quantity:       m.Quantity,
			PriceTHB:       m.PriceTHB,
			Content:        m.Content,
			MainColor:      m.MainColor,
			PaymentMethod:  m.PaymentMethod,
			ProductionTime: m.ProductionTime,
			DueDate:        duePtr,
			DesignOption:   m.DesignOption,
			InstallOption:  m.InstallOption,
			Notes:          m.Notes,
			Status:         m.Status,
			CreatedBy:      m.CreatedBy,
			CreatedAt:      m.CreatedAt,
			UpdatedAt:      m.UpdatedAt,
			DeletedAt:      m.DeletedAt,
		})
	}

	totalPages := 0
	if size > 0 && total > 0 {
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

func (s *signJobService) GetSignJobByJobID(ctx context.Context, jobID string, claims *dto.JWTClaims) (*dto.SignJobDTO, error) {
	m, err := s.signJobRepo.GetSignJobByJobID(ctx, jobID)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}
	var duePtr *time.Time
	if !m.DueDate.IsZero() {
		d := m.DueDate
		duePtr = &d
	}
	dtoObj := &dto.SignJobDTO{
		ID:             m.ID.Hex(),
		JobID:          m.JobID,
		CompanyName:    m.CompanyName,
		ContactPerson:  m.ContactPerson,
		Phone:          m.Phone,
		Email:          m.Email,
		CustomerTypeID: m.CustomerTypeID,
		Address:        m.Address,
		ProjectName:    m.ProjectName,
		JobName:        m.JobName,
		SignTypeID:     m.SignTypeID,
		Width:          m.Width,
		Height:         m.Height,
		Quantity:       m.Quantity,
		PriceTHB:       m.PriceTHB,
		Content:        m.Content,
		MainColor:      m.MainColor,
		PaymentMethod:  m.PaymentMethod,
		ProductionTime: m.ProductionTime,
		DueDate:        duePtr,
		DesignOption:   m.DesignOption,
		InstallOption:  m.InstallOption,
		Notes:          m.Notes,
		Status:         m.Status,
		CreatedBy:      m.CreatedBy,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		DeletedAt:      m.DeletedAt,
	}
	return dtoObj, nil
}

func (s *signJobService) UpdateSignJobByJobID(ctx context.Context, jobID string, update dto.UpdateSignJobDTO, claims *dto.JWTClaims) error {
	var due time.Time
	if update.DueDate != "" {
		parsedDate, err := time.Parse("2006-01-02", update.DueDate)
		if err != nil {
			return err
		}
		due = parsedDate
	}
	status := update.Status
	if status == "" {
		status = "แผนกออกแบบกราฟิก"
	}
	now := time.Now()
	model := models.SignJob{
		CompanyName:    update.CompanyName,
		ContactPerson:  update.ContactPerson,
		Phone:          update.Phone,
		Email:          update.Email,
		CustomerTypeID: update.CustomerTypeID,
		Address:        update.Address,

		ProjectName: update.ProjectName,
		JobName:     update.JobName,
		SignTypeID:  update.SignTypeID,
		Width:       update.Width,
		Height:      update.Height,
		Quantity:    update.Quantity,
		PriceTHB:    update.PriceTHB,
		Content:     update.Content,
		MainColor:   update.MainColor,

		PaymentMethod:  update.PaymentMethod,
		ProductionTime: update.ProductionTime,
		DueDate:        due,

		DesignOption:  update.DesignOption,
		InstallOption: update.InstallOption,
		Notes:         update.Notes,

		Status:    status,
		UpdatedAt: now,
	}
	updated, err := s.signJobRepo.UpdateSignJobByJobID(ctx, jobID, model)
	if err != nil {
		return err
	}
	if updated == nil {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (s *signJobService) DeleteSignJobByJobID(ctx context.Context, jobID string, claims *dto.JWTClaims) error {
	err := s.signJobRepo.DeleteSignJobByJobID(ctx, jobID)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}
