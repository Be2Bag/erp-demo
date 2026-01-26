package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/pkg/util"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type signJobService struct {
	signJobRepo    ports.SignJobRepository
	dropDownRepo   ports.DropDownRepository
	taskRepo       ports.TaskRepository
	incomeRepo     ports.InComeRepository
	receivableRepo ports.ReceivableRepository
	config         config.Config
}

func NewSignJobService(cfg config.Config, signJobRepo ports.SignJobRepository, dropDownRepo ports.DropDownRepository, taskRepo ports.TaskRepository, incomeRepo ports.InComeRepository, receivableRepo ports.ReceivableRepository) ports.SignJobService {
	return &signJobService{config: cfg, signJobRepo: signJobRepo, dropDownRepo: dropDownRepo, taskRepo: taskRepo, incomeRepo: incomeRepo, receivableRepo: receivableRepo}
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

	// ถ้า Waitprice = true ให้ตั้งค่าราคาเป็น 0
	// ปัดเศษค่าเงินเป็น 2 ตำแหน่งเพื่อป้องกัน floating-point precision error
	priceTHB := util.Round2(signJob.PriceTHB)
	depositAmount := util.Round2(signJob.DepositAmount)
	outstandingAmount := util.Round2(signJob.OutstandingAmount)
	if signJob.WaitPrice {
		priceTHB = 0
		depositAmount = 0
		outstandingAmount = 0
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
		PriceTHB:          priceTHB,
		DepositAmount:     depositAmount,
		OutstandingAmount: outstandingAmount,
		Content:           signJob.Content,
		MainColor:         signJob.MainColor,

		PaymentMethod:  signJob.PaymentMethod,
		IsDeposit:      signJob.IsDeposit,
		ProductionTime: signJob.ProductionTime,
		DueDate:        due,

		DesignOption:  signJob.DesignOption,
		InstallOption: signJob.InstallOption,
		Notes:         signJob.Notes,
		WaitPrice:     signJob.WaitPrice,
		WaitConfirm:   signJob.WaitConfirm,

		Status:    "in_progress",
		CreatedBy: claims.UserID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.signJobRepo.CreateSignJob(ctx, model); err != nil {
		return err
	}

	// ถ้า WaitPrice = true หรือ WaitConfirm = true ให้ข้ามการทำงานเกี่ยวกับระบบบัญชีทั้งหมด (Income, Receivable)
	if signJob.WaitPrice || signJob.WaitConfirm {
		return nil
	}

	jobName := signJob.JobName
	nowUTC := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	if !signJob.IsDeposit {

		// Truncate to start of day in UTC
		// สร้าง Income เฉพาะกรณียอด >= 0.01 (มากกว่า 1 สตางค์)
		if signJob.PriceTHB >= 0.01 {
			modelIncome := models.Income{
				IncomeID:              uuid.NewString(),
				BankID:                config.DefaultBankAccountIDs.CompanyBank,
				TransactionCategoryID: config.DefaultTransactionCategoryIDs.CompanyIncome,
				Description:           signJob.Content,
				Amount:                signJob.PriceTHB,
				Currency:              "THB",
				TxnDate:               nowUTC,
				PaymentMethod:         signJob.PaymentMethod,
				ReferenceNo:           "", // เพิ่มเลขใบเสร็จ / หมายเลขธุรกรรมธนาคาร
				Note:                  &jobName,
				CreatedBy:             claims.UserID,
				CreatedAt:             now,
				UpdatedAt:             now,
			}

			if err := s.incomeRepo.CreateInCome(ctx, modelIncome); err != nil {
				return err
			}
		}

	} else {
		// สร้าง Receivable และ Income เฉพาะกรณียอด >= 0.01
		if signJob.PriceTHB >= 0.01 {
			prefix := fmt.Sprintf("AR-%s-", now.Format("02-01-06"))
			maxInvoiceNo, err := s.receivableRepo.GetMaxInvoiceNumber(ctx, prefix)
			if err != nil {
				return fmt.Errorf("get max invoice number: %w", err)
			}

			counter := 1
			if maxInvoiceNo != "" {
				// Extract counter from last invoice (e.g., "AR-25-01-15-001" -> 1)
				var lastCounter int
				_, scanErr := fmt.Sscanf(maxInvoiceNo, prefix+"%d", &lastCounter)
				if scanErr == nil {
					counter = lastCounter + 1
				}
			}

			invoiceNo := fmt.Sprintf("%s%03d", prefix, counter)

			modelReceivable := models.Receivable{
				IDReceivable: uuid.NewString(),
				BankID:       config.DefaultBankAccountIDs.CompanyBank,
				Customer:     signJob.CompanyName,
				InvoiceNo:    invoiceNo,
				IssueDate:    nowUTC,
				DueDate:      nowUTC.AddDate(0, 0, 30),
				Amount:       signJob.PriceTHB,
				Balance:      signJob.OutstandingAmount,
				Status:       "pending",
				Phone:        signJob.Phone,
				Address:      signJob.Address,
				CreatedBy:    claims.UserID,
				CreatedAt:    now,
				UpdatedAt:    now,
				Note:         jobName,
				JobID:        model.JobID,
			}

			if err := s.receivableRepo.CreateReceivable(ctx, modelReceivable); err != nil {
				return err
			}

			// สร้าง Income สำหรับมัดจำ เฉพาะกรณียอดมัดจำ >= 0.01
			if signJob.DepositAmount >= 0.01 {
				modelIncome := models.Income{
					IncomeID:              uuid.NewString(),
					BankID:                config.DefaultBankAccountIDs.CompanyBank,
					TransactionCategoryID: config.DefaultTransactionCategoryIDs.CompanyIncome,
					Description:           signJob.Content,
					Amount:                signJob.DepositAmount,
					Currency:              "THB",
					TxnDate:               nowUTC,
					PaymentMethod:         signJob.PaymentMethod,
					ReferenceNo:           invoiceNo,
					Note:                  &jobName,
					CreatedBy:             claims.UserID,
					CreatedAt:             now,
					UpdatedAt:             now,
				}

				if err := s.incomeRepo.CreateInCome(ctx, modelIncome); err != nil {
					return err
				}
			}
		}

	}

	// สร้าง Receivable สำหรับ credit เฉพาะกรณียอด >= 0.01
	if signJob.PaymentMethod == "credit" && !signJob.IsDeposit && signJob.PriceTHB >= 0.01 {

		prefix := fmt.Sprintf("AR-%s-", now.Format("02-01-06"))
		maxInvoiceNo, err := s.receivableRepo.GetMaxInvoiceNumber(ctx, prefix)
		if err != nil {
			return fmt.Errorf("get max invoice number: %w", err)
		}

		counter := 1
		if maxInvoiceNo != "" {
			// Extract counter from last invoice (e.g., "AR-25-01-15-001" -> 1)
			var lastCounter int
			_, scanErr := fmt.Sscanf(maxInvoiceNo, prefix+"%d", &lastCounter)
			if scanErr == nil {
				counter = lastCounter + 1
			}
		}

		invoiceNo := fmt.Sprintf("%s%03d", prefix, counter)

		modelReceivable := models.Receivable{
			IDReceivable: uuid.NewString(),
			BankID:       config.DefaultBankAccountIDs.CompanyBank,
			Customer:     signJob.CompanyName,
			InvoiceNo:    invoiceNo,
			IssueDate:    nowUTC,
			DueDate:      nowUTC.AddDate(0, 0, 30),
			Amount:       signJob.PriceTHB,
			Balance:      signJob.PriceTHB,
			Status:       "pending",
			Phone:        signJob.Phone,
			Address:      signJob.Address,
			CreatedBy:    claims.UserID,
			CreatedAt:    now,
			UpdatedAt:    now,
			Note:         jobName,
			JobID:        model.JobID,
		}

		if err := s.receivableRepo.CreateReceivable(ctx, modelReceivable); err != nil {
			return err
		}

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
	} else {
		filter["status"] = "in_progress"
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
			WaitPrice:         m.WaitPrice,
			WaitConfirm:       m.WaitConfirm,
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
		Notes:       m.Notes,
		WaitPrice:   m.WaitPrice,
		WaitConfirm: m.WaitConfirm,
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
	oldJobName := ""

	filter := bson.M{"job_id": jobID, "deleted_at": nil}
	existing, err := s.signJobRepo.GetOneSignJobByFilter(ctx, filter, bson.M{})
	if err != nil {
		return err
	}
	if existing == nil {
		return mongo.ErrNoDocuments
	}

	// เก็บสถานะ WaitPrice เดิมไว้ เพื่อตรวจสอบว่าเปลี่ยนจากรอราคาเป็นมีราคาหรือไม่
	wasWaitingForPrice := existing.WaitPrice

	// เก็บสถานะ IsDeposit เดิมไว้ เพื่อตรวจสอบว่าเปลี่ยนจากไม่มีมัดจำเป็นมีมัดจำหรือไม่
	wasNotDeposit := !existing.IsDeposit

	// เก็บ job_name เดิมไว้สำหรับอัพเดท Income
	oldJobName = existing.JobName

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
	if update.PriceTHB >= 0 {
		existing.PriceTHB = util.Round2(update.PriceTHB) // ปัดเศษ 2 ตำแหน่ง
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

	if update.DepositAmount >= 0 {
		existing.DepositAmount = util.Round2(update.DepositAmount) // ปัดเศษ 2 ตำแหน่ง
	}
	if update.OutstandingAmount >= 0 {
		existing.OutstandingAmount = util.Round2(update.OutstandingAmount) // ปัดเศษ 2 ตำแหน่ง
	}

	existing.IsDeposit = update.IsDeposit

	// อัพเดท WaitPrice และ WaitConfirm
	existing.WaitPrice = update.WaitPrice
	existing.WaitConfirm = update.WaitConfirm

	// update status only when a new status is provided (was previously checking existing.Status)
	if update.Status != "" {
		existing.Status = update.Status
	}

	now := time.Now()
	existing.UpdatedAt = now

	updated, err := s.signJobRepo.UpdateSignJobByJobID(ctx, jobID, *existing)
	if err != nil {
		return err
	}
	if updated == nil {
		return mongo.ErrNoDocuments
	}

	// ถ้าเปลี่ยนจากรอราคา (WaitPrice=true) เป็นมีราคาแล้ว (WaitPrice=false) และมีราคา > 0
	// ให้สร้าง Income/Receivable ย้อนหลัง
	if wasWaitingForPrice && !update.WaitPrice && existing.PriceTHB >= 0.01 {
		jobName := existing.JobName
		nowUTC := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

		if !existing.IsDeposit {
			// กรณีจ่ายเต็มจำนวน - สร้าง Income (เฉพาะกรณียอด >= 0.01)
			if existing.PriceTHB >= 0.01 {
				modelIncome := models.Income{
					IncomeID:              uuid.NewString(),
					BankID:                config.DefaultBankAccountIDs.CompanyBank,
					TransactionCategoryID: config.DefaultTransactionCategoryIDs.CompanyIncome,
					Description:           existing.Content,
					Amount:                existing.PriceTHB,
					Currency:              "THB",
					TxnDate:               nowUTC,
					PaymentMethod:         existing.PaymentMethod,
					ReferenceNo:           "",
					Note:                  &jobName,
					CreatedBy:             claims.UserID,
					CreatedAt:             now,
					UpdatedAt:             now,
				}

				if err := s.incomeRepo.CreateInCome(ctx, modelIncome); err != nil {
					return err
				}
			}

		} else {
			// กรณีมีมัดจำ - สร้าง Receivable และ Income สำหรับมัดจำ (เฉพาะกรณียอด >= 0.01)
			if existing.PriceTHB >= 0.01 {
				prefix := fmt.Sprintf("AR-%s-", now.Format("02-01-06"))
				maxInvoiceNo, err := s.receivableRepo.GetMaxInvoiceNumber(ctx, prefix)
				if err != nil {
					return fmt.Errorf("get max invoice number: %w", err)
				}

				counter := 1
				if maxInvoiceNo != "" {
					var lastCounter int
					_, scanErr := fmt.Sscanf(maxInvoiceNo, prefix+"%d", &lastCounter)
					if scanErr == nil {
						counter = lastCounter + 1
					}
				}

				invoiceNo := fmt.Sprintf("%s%03d", prefix, counter)

				modelReceivable := models.Receivable{
					IDReceivable: uuid.NewString(),
					BankID:       config.DefaultBankAccountIDs.CompanyBank,
					Customer:     existing.CompanyName,
					InvoiceNo:    invoiceNo,
					IssueDate:    nowUTC,
					DueDate:      nowUTC.AddDate(0, 0, 30),
					Amount:       existing.PriceTHB,
					Balance:      existing.OutstandingAmount,
					Status:       "pending",
					Phone:        existing.Phone,
					Address:      existing.Address,
					CreatedBy:    claims.UserID,
					CreatedAt:    now,
					UpdatedAt:    now,
					Note:         jobName,
					JobID:        existing.JobID,
				}

				if err := s.receivableRepo.CreateReceivable(ctx, modelReceivable); err != nil {
					return err
				}

				// สร้าง Income สำหรับมัดจำ (เฉพาะกรณียอดมัดจำ >= 0.01)
				if existing.DepositAmount >= 0.01 {
					modelIncome := models.Income{
						IncomeID:              uuid.NewString(),
						BankID:                config.DefaultBankAccountIDs.CompanyBank,
						TransactionCategoryID: config.DefaultTransactionCategoryIDs.CompanyIncome,
						Description:           existing.Content,
						Amount:                existing.DepositAmount,
						Currency:              "THB",
						TxnDate:               nowUTC,
						PaymentMethod:         existing.PaymentMethod,
						ReferenceNo:           invoiceNo,
						Note:                  &jobName,
						CreatedBy:             claims.UserID,
						CreatedAt:             now,
						UpdatedAt:             now,
					}

					if err := s.incomeRepo.CreateInCome(ctx, modelIncome); err != nil {
						return err
					}
				}
			}
		}

		// กรณี credit และไม่มีมัดจำ - สร้าง Receivable เพิ่ม (เฉพาะกรณียอด >= 0.01)
		if existing.PaymentMethod == "credit" && !existing.IsDeposit && existing.PriceTHB >= 0.01 {
			prefix := fmt.Sprintf("AR-%s-", now.Format("02-01-06"))
			maxInvoiceNo, err := s.receivableRepo.GetMaxInvoiceNumber(ctx, prefix)
			if err != nil {
				return fmt.Errorf("get max invoice number: %w", err)
			}

			counter := 1
			if maxInvoiceNo != "" {
				var lastCounter int
				_, scanErr := fmt.Sscanf(maxInvoiceNo, prefix+"%d", &lastCounter)
				if scanErr == nil {
					counter = lastCounter + 1
				}
			}

			invoiceNo := fmt.Sprintf("%s%03d", prefix, counter)

			modelReceivable := models.Receivable{
				IDReceivable: uuid.NewString(),
				BankID:       config.DefaultBankAccountIDs.CompanyBank,
				Customer:     existing.CompanyName,
				InvoiceNo:    invoiceNo,
				IssueDate:    nowUTC,
				DueDate:      nowUTC.AddDate(0, 0, 30),
				Amount:       existing.PriceTHB,
				Balance:      existing.PriceTHB,
				Status:       "pending",
				Phone:        existing.Phone,
				Address:      existing.Address,
				CreatedBy:    claims.UserID,
				CreatedAt:    now,
				UpdatedAt:    now,
				Note:         jobName,
				JobID:        existing.JobID,
			}

			if err := s.receivableRepo.CreateReceivable(ctx, modelReceivable); err != nil {
				return err
			}
		}
	}

	// กรณีเปลี่ยนจากจ่ายเต็มจำนวน (IsDeposit=false) เป็นมีมัดจำ (IsDeposit=true)
	// และไม่ได้เปลี่ยนจากรอราคา (ไม่ต้องสร้างซ้ำ หากสร้างไว้แล้วจาก wasWaitingForPrice)
	// ให้สร้าง Receivable ใหม่สำหรับยอดค้างชำระ
	if wasNotDeposit && update.IsDeposit && !wasWaitingForPrice && existing.PriceTHB >= 0.01 {
		jobName := existing.JobName
		nowUTC := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

		// ตรวจสอบว่ามี Receivable อยู่แล้วหรือไม่
		filterReceivableCheck := bson.M{"job_id": existing.JobID, "deleted_at": nil}
		existingReceivables, errCheck := s.receivableRepo.GetAllReceivablesByFilter(ctx, filterReceivableCheck, nil)
		if errCheck != nil && !errors.Is(errCheck, mongo.ErrNoDocuments) {
			return errCheck
		}

		// ถ้ายังไม่มี Receivable ให้สร้างใหม่
		if len(existingReceivables) == 0 {
			prefix := fmt.Sprintf("AR-%s-", now.Format("02-01-06"))
			maxInvoiceNo, err := s.receivableRepo.GetMaxInvoiceNumber(ctx, prefix)
			if err != nil {
				return fmt.Errorf("get max invoice number: %w", err)
			}

			counter := 1
			if maxInvoiceNo != "" {
				var lastCounter int
				_, scanErr := fmt.Sscanf(maxInvoiceNo, prefix+"%d", &lastCounter)
				if scanErr == nil {
					counter = lastCounter + 1
				}
			}

			invoiceNo := fmt.Sprintf("%s%03d", prefix, counter)

			modelReceivable := models.Receivable{
				IDReceivable: uuid.NewString(),
				BankID:       config.DefaultBankAccountIDs.CompanyBank,
				Customer:     existing.CompanyName,
				InvoiceNo:    invoiceNo,
				IssueDate:    nowUTC,
				DueDate:      nowUTC.AddDate(0, 0, 30),
				Amount:       existing.PriceTHB,
				Balance:      existing.OutstandingAmount,
				Status:       "pending",
				Phone:        existing.Phone,
				Address:      existing.Address,
				CreatedBy:    claims.UserID,
				CreatedAt:    now,
				UpdatedAt:    now,
				Note:         jobName,
				JobID:        existing.JobID,
			}

			if err := s.receivableRepo.CreateReceivable(ctx, modelReceivable); err != nil {
				return err
			}
		}
	}

	// อัพเดท Task ที่เกี่ยวข้องกับ SignJob นี้ (โดยใช้ job_id)
	filterTask := bson.M{"job_id": existing.JobID, "deleted_at": nil}
	partialTaskUpdate := bson.M{
		"project_id":   existing.ProjectID,
		"project_name": existing.ProjectName,
		"job_name":     existing.JobName,
		"description":  existing.Content,
		"updated_at":   now,
	}

	_, errOnUpdateTask := s.taskRepo.UpdateManyTaskFields(ctx, filterTask, partialTaskUpdate)
	if errOnUpdateTask != nil && !errors.Is(errOnUpdateTask, mongo.ErrNoDocuments) {
		return errOnUpdateTask
	}

	// อัพเดท Receivable ที่เกี่ยวข้องกับ SignJob นี้ (โดยใช้ job_id)
	filterReceivable := bson.M{"job_id": existing.JobID, "deleted_at": nil}
	receivables, errOnGetReceivables := s.receivableRepo.GetAllReceivablesByFilter(ctx, filterReceivable, nil)
	if errOnGetReceivables != nil && !errors.Is(errOnGetReceivables, mongo.ErrNoDocuments) {
		return errOnGetReceivables
	}

	for _, rec := range receivables {
		// อัพเดทข้อมูลลูกค้า
		rec.Customer = existing.CompanyName
		rec.Phone = existing.Phone
		rec.Address = existing.Address
		rec.Amount = existing.PriceTHB
		rec.Balance = existing.OutstandingAmount
		rec.Note = existing.JobName
		rec.UpdatedAt = now

		if _, errOnUpdateReceivable := s.receivableRepo.UpdateReceivableByID(ctx, rec.IDReceivable, *rec); errOnUpdateReceivable != nil {
			return errOnUpdateReceivable
		}
	}

	// อัพเดท Income ที่เกี่ยวข้อง (โดยใช้ note เป็น job_name เดิม)
	filterIncome := bson.M{"note": oldJobName, "deleted_at": nil}
	incomes, errOnGetIncomes := s.incomeRepo.GetAllInComeByFilter(ctx, filterIncome, nil)
	if errOnGetIncomes != nil && !errors.Is(errOnGetIncomes, mongo.ErrNoDocuments) {
		return errOnGetIncomes
	}

	for _, inc := range incomes {
		// อัพเดทข้อมูล
		inc.Description = existing.Content
		inc.PaymentMethod = existing.PaymentMethod
		newJobName := existing.JobName
		inc.Note = &newJobName
		inc.UpdatedAt = now

		// ถ้าเป็นการจ่ายแบบมีมัดจำ ให้อัพเดทจำนวนเงินเป็น DepositAmount
		if existing.IsDeposit {
			inc.Amount = existing.DepositAmount
		} else {
			inc.Amount = existing.PriceTHB
		}

		if _, errOnUpdateIncome := s.incomeRepo.UpdateInComeByID(ctx, inc.IncomeID, *inc); errOnUpdateIncome != nil {
			return errOnUpdateIncome
		}
	}

	return nil
}

func (s *signJobService) DeleteSignJobByJobID(ctx context.Context, jobID string, claims *dto.JWTClaims) error {
	err := s.signJobRepo.SoftDeleteSignJobByJobID(ctx, jobID)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil
	}
	return err
}

func (s *signJobService) VerifySignJob(ctx context.Context, jobID string, claims *dto.JWTClaims) error {

	if claims.Role != "admin" && !config.HasPermission(config.PermissionConfig.VerifySignJob, claims.UserID) {
		return fmt.Errorf("permission denied: requires admin role or verify permission")
	}

	filter := bson.M{"job_id": jobID, "deleted_at": nil}
	project := bson.M{"_id": 0, "status": 1, "job_name": 1}
	task, errOnGetTask := s.taskRepo.GetAllTaskByFilter(ctx, filter, project)
	if errOnGetTask != nil && !errors.Is(errOnGetTask, mongo.ErrNoDocuments) {
		return errOnGetTask
	}

	if len(task) > 0 {
		for _, t := range task {
			if t.Status != "done" {
				return fmt.Errorf("can not verify")
			}
		}

		filter := bson.M{"job_id": jobID, "deleted_at": nil}
		existing, err := s.signJobRepo.GetOneSignJobByFilter(ctx, filter, bson.M{})
		if err != nil {
			return err
		}
		if existing == nil {
			return mongo.ErrNoDocuments
		}
		existing.Status = "done"
		updated, err := s.signJobRepo.UpdateSignJobByJobID(ctx, jobID, *existing)
		if err != nil {
			return err
		}
		if updated == nil {
			return mongo.ErrNoDocuments
		}
	} else {
		return fmt.Errorf("no tasks found for this job")
	}

	return nil

}

// ConfirmSignJob เปลี่ยน WaitConfirm จาก true เป็น false (ยืนยันงาน) และสร้าง Income/Receivable
func (s *signJobService) ConfirmSignJob(ctx context.Context, jobID string, claims *dto.JWTClaims) error {
	filter := bson.M{"job_id": jobID, "deleted_at": nil}
	existing, err := s.signJobRepo.GetOneSignJobByFilter(ctx, filter, bson.M{})
	if err != nil {
		return err
	}
	if existing == nil {
		return mongo.ErrNoDocuments
	}

	// ตรวจสอบว่างานนี้อยู่ในสถานะรอยืนยันหรือไม่
	if !existing.WaitConfirm {
		return fmt.Errorf("งานนี้ไม่ได้อยู่ในสถานะรอยืนยัน")
	}

	now := time.Now()

	// เปลี่ยน WaitConfirm เป็น false
	existing.WaitConfirm = false
	existing.UpdatedAt = now

	updated, err := s.signJobRepo.UpdateSignJobByJobID(ctx, jobID, *existing)
	if err != nil {
		return err
	}
	if updated == nil {
		return mongo.ErrNoDocuments
	}

	// ถ้ายังรอราคาอยู่ (WaitPrice = true) ไม่ต้องสร้างระบบบัญชี
	if existing.WaitPrice {
		return nil
	}

	// ถ้าไม่รอราคาแล้ว (มีราคาแล้ว) และมีราคา >= 0.01 ให้สร้าง Income/Receivable
	if util.IsPositiveAmount(existing.PriceTHB) {
		jobName := existing.JobName
		nowUTC := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

		if !existing.IsDeposit {
			// กรณีจ่ายเต็มจำนวน - สร้าง Income
			modelIncome := models.Income{
				IncomeID:              uuid.NewString(),
				BankID:                config.DefaultBankAccountIDs.CompanyBank,
				TransactionCategoryID: config.DefaultTransactionCategoryIDs.CompanyIncome,
				Description:           existing.Content,
				Amount:                existing.PriceTHB,
				Currency:              "THB",
				TxnDate:               nowUTC,
				PaymentMethod:         existing.PaymentMethod,
				ReferenceNo:           "",
				Note:                  &jobName,
				CreatedBy:             claims.UserID,
				CreatedAt:             now,
				UpdatedAt:             now,
			}

			if err := s.incomeRepo.CreateInCome(ctx, modelIncome); err != nil {
				return err
			}

		} else {
			// กรณีมีมัดจำ - สร้าง Receivable และ Income สำหรับมัดจำ
			prefix := fmt.Sprintf("AR-%s-", now.Format("02-01-06"))
			maxInvoiceNo, err := s.receivableRepo.GetMaxInvoiceNumber(ctx, prefix)
			if err != nil {
				return fmt.Errorf("get max invoice number: %w", err)
			}

			counter := 1
			if maxInvoiceNo != "" {
				var lastCounter int
				_, scanErr := fmt.Sscanf(maxInvoiceNo, prefix+"%d", &lastCounter)
				if scanErr == nil {
					counter = lastCounter + 1
				}
			}

			invoiceNo := fmt.Sprintf("%s%03d", prefix, counter)

			modelReceivable := models.Receivable{
				IDReceivable: uuid.NewString(),
				BankID:       config.DefaultBankAccountIDs.CompanyBank,
				Customer:     existing.CompanyName,
				InvoiceNo:    invoiceNo,
				IssueDate:    nowUTC,
				DueDate:      nowUTC.AddDate(0, 0, 30),
				Amount:       existing.PriceTHB,
				Balance:      existing.OutstandingAmount,
				Status:       "pending",
				Phone:        existing.Phone,
				Address:      existing.Address,
				CreatedBy:    claims.UserID,
				CreatedAt:    now,
				UpdatedAt:    now,
				Note:         jobName,
				JobID:        existing.JobID,
			}

			if err := s.receivableRepo.CreateReceivable(ctx, modelReceivable); err != nil {
				return err
			}

			modelIncome := models.Income{
				IncomeID:              uuid.NewString(),
				BankID:                config.DefaultBankAccountIDs.CompanyBank,
				TransactionCategoryID: config.DefaultTransactionCategoryIDs.CompanyIncome,
				Description:           existing.Content,
				Amount:                existing.DepositAmount,
				Currency:              "THB",
				TxnDate:               nowUTC,
				PaymentMethod:         existing.PaymentMethod,
				ReferenceNo:           invoiceNo,
				Note:                  &jobName,
				CreatedBy:             claims.UserID,
				CreatedAt:             now,
				UpdatedAt:             now,
			}

			if err := s.incomeRepo.CreateInCome(ctx, modelIncome); err != nil {
				return err
			}
		}

		// กรณี credit และไม่มีมัดจำ - สร้าง Receivable เพิ่ม
		if existing.PaymentMethod == "credit" && !existing.IsDeposit {
			prefix := fmt.Sprintf("AR-%s-", now.Format("02-01-06"))
			maxInvoiceNo, err := s.receivableRepo.GetMaxInvoiceNumber(ctx, prefix)
			if err != nil {
				return fmt.Errorf("get max invoice number: %w", err)
			}

			counter := 1
			if maxInvoiceNo != "" {
				var lastCounter int
				_, scanErr := fmt.Sscanf(maxInvoiceNo, prefix+"%d", &lastCounter)
				if scanErr == nil {
					counter = lastCounter + 1
				}
			}

			invoiceNo := fmt.Sprintf("%s%03d", prefix, counter)

			modelReceivable := models.Receivable{
				IDReceivable: uuid.NewString(),
				BankID:       config.DefaultBankAccountIDs.CompanyBank,
				Customer:     existing.CompanyName,
				InvoiceNo:    invoiceNo,
				IssueDate:    nowUTC,
				DueDate:      nowUTC.AddDate(0, 0, 30),
				Amount:       existing.PriceTHB,
				Balance:      existing.PriceTHB,
				Status:       "pending",
				Phone:        existing.Phone,
				Address:      existing.Address,
				CreatedBy:    claims.UserID,
				CreatedAt:    now,
				UpdatedAt:    now,
				Note:         jobName,
				JobID:        existing.JobID,
			}

			if err := s.receivableRepo.CreateReceivable(ctx, modelReceivable); err != nil {
				return err
			}
		}
	}

	return nil
}
