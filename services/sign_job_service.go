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

type signJobService struct {
	config       config.Config
	signJobRepo  ports.SignJobRepository
	dropDownRepo ports.DropDownRepository
	taskRepo     ports.TaskRepository
}

func NewSignJobService(cfg config.Config, signJobRepo ports.SignJobRepository, dropDownRepo ports.DropDownRepository, taskRepo ports.TaskRepository) ports.SignJobService {
	return &signJobService{config: cfg, signJobRepo: signJobRepo, dropDownRepo: dropDownRepo, taskRepo: taskRepo}
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

	model := models.SignJob{
		JobID:          uuid.NewString(),
		CompanyName:    signJob.CompanyName,
		ContactPerson:  signJob.ContactPerson,
		Phone:          signJob.Phone,
		Email:          signJob.Email,
		CustomerTypeID: signJob.CustomerTypeID,
		Address:        signJob.Address,

		ProjectID:         signJob.ProjectID,
		ProjectName:       signJob.ProjectName,
		JobName:           signJob.JobName,
		SignTypeID:        signJob.SignTypeID,
		Width:             signJob.Width,
		Height:            signJob.Height,
		Quantity:          signJob.Quantity,
		PriceTHB:          signJob.PriceTHB,
		DepositAmount:     signJob.DepositAmount,
		OutstandingAmount: signJob.OutstandingAmount,
		Content:           signJob.Content,
		MainColor:         signJob.MainColor,

		PaymentMethod:  signJob.PaymentMethod,
		IsDeposit:      signJob.IsDeposit,
		ProductionTime: signJob.ProductionTime,
		DueDate:        due,

		DesignOption:  signJob.DesignOption,
		InstallOption: signJob.InstallOption,
		Notes:         signJob.Notes,

		Status:    "DPT001",
		CreatedBy: claims.UserID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.signJobRepo.CreateSignJob(ctx, model); err != nil {
		return err
	}
	return nil
}

func (s *signJobService) ListSignJobs(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, status string, sortBy string, sortOrder string) (dto.Pagination, error) {
	skip := int64((page - 1) * size)
	limit := int64(size)

	filter := bson.M{
		"deleted_at": nil,
	}

	status = strings.TrimSpace(status)
	if status != "" {
		filter["status"] = status
	}

	search = strings.TrimSpace(search)
	if search != "" {
		safe := regexp.QuoteMeta(search)
		re := primitive.Regex{Pattern: safe, Options: "i"}
		filter["$or"] = []bson.M{
			{"project_name": re},
			{"job_name": re},
			{"company_name": re},
			{"contact_person": re},
		}
	}

	projection := bson.M{}

	// sort
	allowedSortFields := map[string]string{
		"created_at":   "created_at",
		"updated_at":   "updated_at",
		"due_date":     "due_date",
		"job_name":     "job_name",
		"project_name": "project_name",
		"company_name": "company_name",
		"status":       "status",
		"price_thb":    "price_thb",
		"quantity":     "quantity",
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

	items, total, err := s.signJobRepo.GetListSignJobsByFilter(ctx, filter, projection, sort, skip, limit)
	if err != nil {
		return dto.Pagination{}, fmt.Errorf("list sign jobs: %w", err)
	}

	list := make([]interface{}, 0, len(items))
	for _, m := range items {

		SignTypeName := ""
		filter := bson.M{"type_id": m.SignTypeID, "deleted_at": nil}
		projection := bson.M{}

		signTypes, errOnGetSignTypes := s.dropDownRepo.GetSignTypes(ctx, filter, projection)
		if errOnGetSignTypes != nil {
			return dto.Pagination{}, errOnGetSignTypes
		}

		// avoid panic when slice is empty (len==0 but slice not nil)
		if len(signTypes) > 0 {
			SignTypeName = signTypes[0].NameTH
		}

		list = append(list, dto.SignJobDTO{
			JobID:             m.JobID,
			CompanyName:       m.CompanyName,
			ContactPerson:     m.ContactPerson,
			Phone:             m.Phone,
			Email:             m.Email,
			CustomerTypeID:    m.CustomerTypeID,
			Address:           m.Address,
			ProjectID:         m.ProjectID,
			ProjectName:       m.ProjectName,
			JobName:           m.JobName,
			SignTypeName:      SignTypeName,
			SignTypeID:        m.SignTypeID,
			Width:             m.Width,
			Height:            m.Height,
			Quantity:          m.Quantity,
			PriceTHB:          m.PriceTHB,
			DepositAmount:     m.DepositAmount,
			OutstandingAmount: m.OutstandingAmount,
			Content:           m.Content,
			MainColor:         m.MainColor,
			PaymentMethod:     m.PaymentMethod,
			IsDeposit:         m.IsDeposit,
			ProductionTime:    m.ProductionTime,
			DueDate:           m.DueDate,
			DesignOption:      m.DesignOption,
			InstallOption:     m.InstallOption,
			Notes:             m.Notes,
			Status:            m.Status,
			CreatedBy:         m.CreatedBy,
			CreatedAt:         m.CreatedAt,
			UpdatedAt:         m.UpdatedAt,
			DeletedAt:         m.DeletedAt,
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

func (s *signJobService) GetSignJobByJobID(ctx context.Context, jobID string, claims *dto.JWTClaims) (*dto.SignJobDTO, error) {

	filter := bson.M{"job_id": jobID, "deleted_at": nil}
	projection := bson.M{}

	m, err := s.signJobRepo.GetOneSignJobByFilter(ctx, filter, projection)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}

	SignTypeName := ""
	filterSignType := bson.M{"type_id": m.SignTypeID, "deleted_at": nil}
	projectionSignType := bson.M{}

	signTypes, errOnGetSignTypes := s.dropDownRepo.GetSignTypes(ctx, filterSignType, projectionSignType)
	if errOnGetSignTypes != nil {
		return nil, errOnGetSignTypes
	}

	if len(signTypes) > 0 { // prevent potential panic when empty slice returned
		SignTypeName = signTypes[0].NameTH
	}

	dtoObj := &dto.SignJobDTO{
		// ---------- ลูกค้า ----------
		JobID:          m.JobID,
		CompanyName:    m.CompanyName,
		ContactPerson:  m.ContactPerson,
		Phone:          m.Phone,
		Email:          m.Email,
		CustomerTypeID: m.CustomerTypeID,
		Address:        m.Address,
		// ---------- รายละเอียดงานป้าย ----------
		ProjectID:         m.ProjectID,
		ProjectName:       m.ProjectName,
		JobName:           m.JobName,
		SignTypeName:      SignTypeName,
		SignTypeID:        m.SignTypeID,
		Width:             m.Width,
		Height:            m.Height,
		Quantity:          m.Quantity,
		PriceTHB:          m.PriceTHB,
		DepositAmount:     m.DepositAmount,
		OutstandingAmount: m.OutstandingAmount,
		Content:           m.Content,
		MainColor:         m.MainColor,
		// ---------- การชำระเงิน ----------
		PaymentMethod: m.PaymentMethod,
		IsDeposit:     m.IsDeposit,
		// ---------- การผลิต / ไทม์ไลน์ ----------
		ProductionTime: m.ProductionTime,
		DueDate:        m.DueDate,
		// ---------- งานออกแบบ / การติดตั้ง ----------
		DesignOption:  m.DesignOption,
		InstallOption: m.InstallOption,
		// ---------- หมายเหตุ ----------
		Notes: m.Notes,
		// ---------- เมต้า ----------
		Status:    m.Status,
		CreatedBy: m.CreatedBy,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
	return dtoObj, nil
}

func (s *signJobService) UpdateSignJobByJobID(ctx context.Context, jobID string, update dto.UpdateSignJobDTO, claims *dto.JWTClaims) error {
	// ดึงข้อมูลเดิม
	IsEditJobName := false

	filter := bson.M{"job_id": jobID, "deleted_at": nil}
	existing, err := s.signJobRepo.GetOneSignJobByFilter(ctx, filter, bson.M{})
	if err != nil {
		return err
	}
	if existing == nil {
		return mongo.ErrNoDocuments
	}

	if update.CompanyName != "" {
		existing.CompanyName = update.CompanyName
	}
	if update.ContactPerson != "" {
		existing.ContactPerson = update.ContactPerson
	}
	if update.Phone != "" {
		existing.Phone = update.Phone
	}
	if update.Email != "" {
		existing.Email = update.Email
	}
	if update.CustomerTypeID != "" {
		existing.CustomerTypeID = update.CustomerTypeID
	}
	if update.Address != "" {
		existing.Address = update.Address
	}

	if update.ProjectID != "" {
		existing.ProjectID = update.ProjectID
	}
	if update.ProjectName != "" {
		existing.ProjectName = update.ProjectName
	}
	if update.JobName != "" {
		existing.JobName = update.JobName
		IsEditJobName = true
	}
	if update.SignTypeID != "" {
		existing.SignTypeID = update.SignTypeID
	}
	if update.Width > 0 {
		existing.Width = update.Width
	}
	if update.Height > 0 {
		existing.Height = update.Height
	}
	if update.Quantity > 0 {
		existing.Quantity = update.Quantity
	}
	if update.PriceTHB > 0 {
		existing.PriceTHB = update.PriceTHB
	}
	if update.Content != "" {
		existing.Content = update.Content
	}
	if update.MainColor != "" {
		existing.MainColor = update.MainColor
	}

	if update.PaymentMethod != "" {
		existing.PaymentMethod = update.PaymentMethod
	}
	if update.ProductionTime != "" {
		existing.ProductionTime = update.ProductionTime
	}
	if update.DueDate != "" {
		parsedDate, err := time.Parse("2006-01-02", update.DueDate)
		if err != nil {
			return err
		}
		existing.DueDate = parsedDate
	}

	if update.DesignOption != "" {
		existing.DesignOption = update.DesignOption
	}
	if update.InstallOption != "" {
		existing.InstallOption = update.InstallOption
	}
	if update.Notes != "" {
		existing.Notes = update.Notes
	}

	if update.DepositAmount > 0 {
		existing.DepositAmount = update.DepositAmount
	}
	if update.OutstandingAmount > 0 {
		existing.OutstandingAmount = update.OutstandingAmount
	}

	existing.IsDeposit = update.IsDeposit

	// update status only when a new status is provided (was previously checking existing.Status)
	if update.Status != "" {
		existing.Status = update.Status
	}

	existing.UpdatedAt = time.Now()

	updated, err := s.signJobRepo.UpdateSignJobByJobID(ctx, jobID, *existing)
	if err != nil {
		return err
	}
	if updated == nil {
		return mongo.ErrNoDocuments
	}

	if IsEditJobName {
		filterTask := bson.M{"job_id": existing.JobID}
		partialTaskUpdate := bson.M{"job_name": existing.JobName}

		_, errOnUpdateTask := s.taskRepo.UpdateManyTaskFields(ctx, filterTask, partialTaskUpdate)
		if errOnUpdateTask != nil {
			return errOnUpdateTask
		}
	}

	return nil
}

func (s *signJobService) DeleteSignJobByJobID(ctx context.Context, jobID string, claims *dto.JWTClaims) error {
	err := s.signJobRepo.SoftDeleteSignJobByJobID(ctx, jobID)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}
