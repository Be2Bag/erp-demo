package services

import (
	"context"
	"fmt"
	"math"
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

type receivableService struct {
	config           config.Config
	receivableRepo   ports.ReceivableRepository
	bankAccountsRepo ports.BankAccountsRepository
	signJobRepo      ports.SignJobRepository
	incomeRepo       ports.InComeRepository
}

func NewReceivableService(cfg config.Config, receivableRepo ports.ReceivableRepository, bankAccountsRepo ports.BankAccountsRepository, signJobRepo ports.SignJobRepository, incomeRepo ports.InComeRepository) ports.ReceivableService {
	return &receivableService{config: cfg, receivableRepo: receivableRepo, bankAccountsRepo: bankAccountsRepo, signJobRepo: signJobRepo, incomeRepo: incomeRepo}
}

func (s *receivableService) CreateReceivable(ctx context.Context, receivable dto.CreateReceivableDTO, claims *dto.JWTClaims) error {
	now := time.Now()
	var due time.Time
	var issue time.Time
	if receivable.DueDate != "" {
		parsedDate, err := time.Parse("2006-01-02", receivable.DueDate)
		if err != nil {
			return err
		}
		due = parsedDate
	}

	if receivable.IssueDate != "" {
		parsedDate, err := time.Parse("2006-01-02", receivable.IssueDate)
		if err != nil {
			return err
		}
		issue = parsedDate
	}

	model := models.Receivable{
		IDReceivable: uuid.NewString(),
		BankID:       receivable.BankID,
		Customer:     receivable.Customer,
		InvoiceNo:    receivable.InvoiceNo,
		IssueDate:    issue,
		DueDate:      due,
		Amount:       receivable.Amount,
		Balance:      receivable.Balance,
		Status:       "pending",
		Phone:        receivable.Phone,
		Address:      receivable.Address,
		CreatedBy:    claims.UserID,
		CreatedAt:    now,
		UpdatedAt:    now,
		Note:         receivable.Note,
	}

	if err := s.receivableRepo.CreateReceivable(ctx, model); err != nil {
		return err
	}
	return nil
}

func (s *receivableService) ListReceivables(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, sortBy string, sortOrder string, Status string, startDate string, endDate string, bankID string) (dto.Pagination, error) {
	skip := int64((page - 1) * size)
	limit := int64(size)

	filter := bson.M{
		"deleted_at": nil,
	}
	if bankID != "" {
		filter["bank_id"] = bankID
	}

	// กรอง startDate และ endDate
	if startDate != "" || endDate != "" {
		txnDateFilter := bson.M{}

		if startDate != "" {
			parsedStartDate, err := time.ParseInLocation("2006-01-02", startDate, time.UTC)
			if err != nil {
				return dto.Pagination{}, fmt.Errorf("invalid startDate format: %w", err)
			}
			txnDateFilter["$gte"] = parsedStartDate
		}

		if endDate != "" {
			parsedEndDate, err := time.ParseInLocation("2006-01-02", endDate, time.UTC)
			if err != nil {
				return dto.Pagination{}, fmt.Errorf("invalid endDate format: %w", err)
			}
			// เพิ่ม 1 วันเพื่อให้ครบ 23:59:59 ของวันที่ endDate
			txnDateFilter["$lt"] = parsedEndDate.Add(24 * time.Hour)
		}

		filter["issue_date"] = txnDateFilter
	}

	search = strings.TrimSpace(search)
	if search != "" {
		safe := regexp.QuoteMeta(search)
		re := primitive.Regex{Pattern: safe, Options: "i"}
		filter["$or"] = []bson.M{
			{"description": re},
			{"customer": re},
			{"note": re},
			{"invoice_no": re},
		}
	}

	if Status != "" {
		filter["status"] = strings.ToLower(Status)
	}

	projection := bson.M{}

	// sort
	allowedSortFields := map[string]string{
		"created_at": "created_at",
		"updated_at": "updated_at",
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

	items, total, err := s.receivableRepo.GetListReceivablesByFilter(ctx, filter, projection, sort, skip, limit)
	if err != nil {
		return dto.Pagination{}, fmt.Errorf("list receivables: %w", err)
	}

	filterBankAccounts := bson.M{"deleted_at": nil}
	bankAccounts, errOnGetBankAccounts := s.bankAccountsRepo.GetAllBankAccountsByFilter(ctx, filterBankAccounts, nil)
	if errOnGetBankAccounts != nil {
		return dto.Pagination{}, fmt.Errorf("list bank accounts: %w", errOnGetBankAccounts)
	}

	// สร้าง map สำหรับ mapping BankID กับ BankName
	bankMap := make(map[string]string)
	for _, bank := range bankAccounts {
		bankMap[bank.BankID] = bank.BankName
	}

	list := make([]interface{}, 0, len(items))
	for _, m := range items {

		list = append(list, dto.ReceivableDTO{
			IDReceivable: m.IDReceivable,
			BankID:       m.BankID,
			BankName:     bankMap[m.BankID],
			Customer:     m.Customer,
			InvoiceNo:    m.InvoiceNo,
			IssueDate:    m.IssueDate,
			DueDate:      m.DueDate,
			Amount:       m.Amount,
			Balance:      m.Balance,
			Status:       m.Status,
			Phone:        m.Phone,
			Address:      m.Address,
			CreatedBy:    m.CreatedBy,
			CreatedAt:    m.CreatedAt,
			UpdatedAt:    m.UpdatedAt,
			Note:         m.Note,
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

func (s *receivableService) GetReceivableByID(ctx context.Context, receivableID string, claims *dto.JWTClaims) (*dto.ReceivableDTO, error) {

	filter := bson.M{"id_receivable": receivableID, "deleted_at": nil}
	projection := bson.M{}

	m, err := s.receivableRepo.GetOneReceivableByFilter(ctx, filter, projection)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}

	filterPaymentTransaction := bson.M{"ref_invoice_no": m.InvoiceNo, "transaction_type": "receivable", "deleted_at": nil}

	paymentTransaction, errPaymentTransaction := s.receivableRepo.GetAllPaymentTransactionsByFilter(ctx, filterPaymentTransaction, bson.M{})
	if errPaymentTransaction != nil {
		return nil, errPaymentTransaction
	}

	filterBankAccounts := bson.M{"deleted_at": nil}
	bankAccounts, errOnGetBankAccounts := s.bankAccountsRepo.GetAllBankAccountsByFilter(ctx, filterBankAccounts, nil)
	if errOnGetBankAccounts != nil {
		return nil, fmt.Errorf("list bank accounts: %w", errOnGetBankAccounts)
	}

	// สร้าง map สำหรับ mapping BankID กับ BankName
	bankMap := make(map[string]string)
	for _, bank := range bankAccounts {
		bankMap[bank.BankID] = bank.BankName
	}

	var transactions = make([]dto.PaymentTransactionDTO, 0, len(paymentTransaction))
	for _, pt := range paymentTransaction {
		transactions = append(transactions, dto.PaymentTransactionDTO{
			IDTransaction:   pt.IDTransaction,
			BankID:          pt.BankID,
			RefInvoiceNo:    pt.RefInvoiceNo,
			TransactionType: pt.TransactionType,
			PaymentDate:     pt.PaymentDate,
			Amount:          pt.Amount,
			PaymentMethod:   pt.PaymentMethod,
			PaymentRef:      pt.PaymentRef,
			Note:            pt.Note,
			CreatedBy:       pt.CreatedBy,
			CreatedAt:       pt.CreatedAt,
			UpdatedAt:       pt.UpdatedAt,
		})
	}

	dtoObj := &dto.ReceivableDTO{
		// ---------- รายละเอียดรายได้ ----------
		IDReceivable: m.IDReceivable,
		BankID:       m.BankID,
		BankName:     bankMap[m.BankID],
		Customer:     m.Customer,
		InvoiceNo:    m.InvoiceNo,
		IssueDate:    m.IssueDate,
		DueDate:      m.DueDate,
		Amount:       m.Amount,
		Balance:      m.Balance,
		Status:       m.Status,
		Phone:        m.Phone,
		Address:      m.Address,
		CreatedBy:    m.CreatedBy,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		Transactions: transactions,
		Note:         m.Note,
	}
	return dtoObj, nil
}

func (s *receivableService) UpdateReceivableByID(ctx context.Context, receivableID string, update dto.UpdateReceivableDTO, claims *dto.JWTClaims) error {
	// fetch existing
	filter := bson.M{"id_receivable": receivableID, "deleted_at": nil}
	existing, err := s.receivableRepo.GetOneReceivableByFilter(ctx, filter, bson.M{})
	if err != nil {
		return err
	}
	if existing == nil {
		return mongo.ErrNoDocuments
	}

	if strings.TrimSpace(update.Customer) != "" {
		existing.Customer = update.Customer
	}
	if strings.TrimSpace(update.BankID) != "" {
		existing.BankID = update.BankID
	}
	if strings.TrimSpace(update.InvoiceNo) != "" {
		existing.InvoiceNo = update.InvoiceNo
	}
	if update.IssueDate != "" {
		t, err := time.Parse("2006-01-02", update.IssueDate)
		if err != nil {
			return err
		}
		existing.IssueDate = t
	}
	if update.DueDate != "" {
		t, err := time.Parse("2006-01-02", update.DueDate)
		if err != nil {
			return err
		}
		existing.DueDate = t
	}
	if update.Amount > 0 {
		existing.Amount = update.Amount
	}
	if update.Balance > 0 {
		existing.Balance = update.Balance
	}
	if strings.TrimSpace(update.Status) != "" {
		switch strings.ToLower(update.Status) {
		case "pending", "paid", "overdue", "partial":
			existing.Status = strings.ToLower(update.Status)
		default:
			return fmt.Errorf("invalid status: %s", update.Status)
		}
	}

	if strings.TrimSpace(update.Note) != "" {
		existing.Note = update.Note
	}

	if strings.TrimSpace(update.Phone) != "" {
		existing.Phone = update.Phone
	}

	if update.Address != "" {
		existing.Address = update.Address
	}

	existing.UpdatedAt = time.Now()

	if _, err := s.receivableRepo.UpdateReceivableByID(ctx, receivableID, *existing); err != nil {
		return err
	}
	return nil
}

func (s *receivableService) DeleteReceivableByID(ctx context.Context, receivableID string, claims *dto.JWTClaims) error {
	err := s.receivableRepo.SoftDeleteReceivableByID(ctx, receivableID)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}

func (s *receivableService) SummaryReceivableByFilter(ctx context.Context, claims *dto.JWTClaims, report dto.RequestSummaryReceivable) (dto.ReceivableSummaryDTO, error) {
	now := time.Now()

	// base filter
	filter := bson.M{
		"deleted_at": nil,
	}
	if strings.TrimSpace(report.BankID) != "" {
		filter["bank_id"] = strings.TrimSpace(report.BankID)
	}

	// ถ้ามี StartDate หรือ EndDate ให้ใช้ช่วงวันที่ที่กำหนด
	if strings.TrimSpace(report.StartDate) != "" || strings.TrimSpace(report.EndDate) != "" {
		issueDateFilter := bson.M{}
		if strings.TrimSpace(report.StartDate) != "" {
			parsedStartDate, err := time.Parse("2006-01-02", report.StartDate)
			if err != nil {
				return dto.ReceivableSummaryDTO{}, fmt.Errorf("invalid start_date format: %w", err)
			}
			issueDateFilter["$gte"] = parsedStartDate
		}
		if strings.TrimSpace(report.EndDate) != "" {
			parsedEndDate, err := time.Parse("2006-01-02", report.EndDate)
			if err != nil {
				return dto.ReceivableSummaryDTO{}, fmt.Errorf("invalid end_date format: %w", err)
			}
			// เพิ่ม 1 วันเพื่อให้ครบ 23:59:59 ของวันที่ endDate
			issueDateFilter["$lt"] = parsedEndDate.Add(24 * time.Hour)
		}
		filter["issue_date"] = issueDateFilter
	}
	// ถ้าไม่มี StartDate/EndDate จะดึงข้อมูลทั้งหมด

	receivables, err := s.receivableRepo.GetAllReceivablesByFilter(ctx, filter, nil)
	if err != nil {
		return dto.ReceivableSummaryDTO{}, err
	}

	var totalAmount float64
	var totalDue float64
	overdueCount := 0

	for _, p := range receivables {
		totalAmount += p.Amount

		// outstanding amount
		if p.Balance > 0 {
			totalDue += p.Balance

			// overdue: due date passed and still has balance
			if !p.DueDate.IsZero() && p.DueDate.Before(now) {
				overdueCount++
			} else if strings.EqualFold(p.Status, "overdue") {
				// fallback if due date not set but status flagged
				overdueCount++
			}
		}
	}

	return dto.ReceivableSummaryDTO{
		TotalAmount:  totalAmount,
		TotalDue:     totalDue,
		OverdueCount: overdueCount,
	}, nil
}

func (s *receivableService) RecordReceipt(ctx context.Context, input dto.RecordReceiptDTO, claims *dto.JWTClaims) error { // บันทึกรายการรับชำระของลูกหนี้
	// 1) fetch receivable
	filter := bson.M{"id_receivable": strings.TrimSpace(input.ReceivableID), "deleted_at": nil} // เงื่อนไขค้นหาลูกหนี้ที่ยังไม่ถูกลบตาม ID
	rec, err := s.receivableRepo.GetOneReceivableByFilter(ctx, filter, bson.M{})                // ดึงข้อมูลลูกหนี้จากแหล่งข้อมูล
	if err != nil {
		return fmt.Errorf("get receivable: %w", err) // หากดึงข้อมูลผิดพลาด ส่งต่อ error ออกไป
	}
	if rec == nil {
		return mongo.ErrNoDocuments // ไม่พบข้อมูลลูกหนี้
	}

	// 2) validate amount
	amt := input.Amount // จำนวนเงินที่รับชำระ
	if amt <= 0 {       // ต้องเป็นจำนวนที่มากกว่า 0
		return fmt.Errorf("amount must be greater than 0") // แจ้งเตือนจำนวนไม่ถูกต้อง
	}
	if rec.Balance <= 0 { // หากยอดคงเหลือเป็น 0 หรือ น้อยกว่า
		return fmt.Errorf("receivable is already fully paid") // ถือว่าชำระครบแล้ว ไม่สามารถรับเพิ่ม
	}
	if amt > rec.Balance { // ไม่ให้ชำระเกินยอดคงเหลือ
		return fmt.Errorf("amount %.2f exceeds remaining balance %.2f", amt, rec.Balance) // แจ้งว่าเกินยอด
	}

	// 3) parse date
	now := time.Now()                               // เวลา ณ ปัจจุบัน
	payDate := now                                  // กำหนดวันที่รับชำระเริ่มต้นเป็นปัจจุบัน
	if strings.TrimSpace(input.PaymentDate) != "" { // หากมีระบุวันที่รับชำระมา
		t, err := time.Parse("2006-01-02", input.PaymentDate) // แปลงรูปแบบวันที่เป็น YYYY-MM-DD
		if err != nil {
			return fmt.Errorf("invalid payment_date, want YYYY-MM-DD: %w", err) // รูปแบบวันที่ไม่ถูกต้อง
		}
		payDate = t // ใช้วันที่ที่ผู้ใช้ระบุ
	}

	// 4) build transaction (incoming for receivable)
	refInvoice := rec.InvoiceNo              // อ้างอิงเลขที่ใบแจ้งหนี้
	if strings.TrimSpace(refInvoice) == "" { // หากไม่มีเลขที่ใบแจ้งหนี้
		refInvoice = rec.IDReceivable // ใช้รหัสลูกหนี้แทน
	}
	tx := models.PaymentTransaction{
		IDTransaction:   uuid.NewString(),                       // รหัสธุรกรรมใหม่แบบ UUID
		BankID:          rec.BankID,                             // รหัสบัญชีธนาคารที่เกี่ยวข้อง
		RefInvoiceNo:    refInvoice,                             // อ้างอิงเอกสาร
		TransactionType: "receivable",                           // ประเภทธุรกรรมเป็นลูกหนี้
		PaymentDate:     payDate,                                // วันที่ชำระเงิน
		Amount:          amt,                                    // จำนวนเงินที่รับ
		PaymentMethod:   strings.TrimSpace(input.PaymentMethod), // วิธีการชำระ (เช่น โอน/เงินสด)
		PaymentRef:      strings.TrimSpace(input.PaymentRef),    // เลขอ้างอิงการชำระ (เช่น เลขที่ธุรกรรม)
		Note:            strings.TrimSpace(input.Note),          // หมายเหตุ
		CreatedBy:       claims.UserID,                          // ผู้ทำรายการ
		CreatedAt:       now,                                    // เวลาสร้าง
		UpdatedAt:       now,                                    // เวลาอัปเดตล่าสุด
	}

	// 5) insert transaction first
	if err := s.receivableRepo.CreatePaymentTransaction(ctx, tx); err != nil { // บันทึกธุรกรรมการรับชำระก่อนเพื่อเก็บประวัติ
		return fmt.Errorf("insert payment tx: %w", err) // หากผิดพลาดให้คืน error
	}

	// 6) update receivable balance and status
	rec.Balance -= amt // หักยอดคงเหลือด้วยจำนวนที่รับชำระ

	// ปัดเศษเป็นทศนิยม 2 ตำแหน่งเพื่อป้องกัน floating-point precision error
	rec.Balance = math.Round(rec.Balance*100) / 100

	if rec.Balance < 0.01 { // ถ้าค่าใกล้ 0 มาก (น้อยกว่า 1 สตางค์) ให้ถือว่าเป็น 0
		rec.Balance = 0 // เซ็ตเป็นศูนย์
	}

	if rec.Balance == 0 { // หากชำระครบ
		rec.Status = "paid" // สถานะเป็นจ่ายครบ
	} else if !rec.DueDate.IsZero() && rec.DueDate.Before(now) { // ยังเหลือและเกินกำหนด
		rec.Status = "overdue" // สถานะค้างชำระเกินกำหนด
	} else {
		rec.Status = "partial" // ยังเหลือยอด -> จ่ายบางส่วน
	}
	rec.UpdatedAt = now // อัปเดตเวลาแก้ไขล่าสุด

	if _, err := s.receivableRepo.UpdateReceivableByID(ctx, rec.IDReceivable, *rec); err != nil { // บันทึกอัปเดตข้อมูลลูกหนี้
		return fmt.Errorf("update receivable: %w", err) // หากบันทึกไม่สำเร็จ ส่ง error ออกไป
	}

	// ตรวจสอบว่า Receivable มี JobID หรือไม่ เพื่ออัพเดท SignJob
	var signJob *models.SignJob
	if rec.JobID != "" {
		signJob, _ = s.signJobRepo.GetOneSignJobByFilter(ctx, bson.M{"job_id": rec.JobID, "deleted_at": nil}, bson.M{})

		if signJob != nil {
			signJob.OutstandingAmount = rec.Balance
			signJob.UpdatedAt = now
			_, errUpdateSignJob := s.signJobRepo.UpdateSignJobByJobID(ctx, signJob.JobID, *signJob)
			if errUpdateSignJob != nil {
				return fmt.Errorf("update sign job: %w", errUpdateSignJob)
			}
		}
	}

	// สร้าง Income ทุกครั้งที่รับชำระ
	var modelIncome models.Income
	if signJob != nil {
		// กรณีมี SignJob - ใช้ข้อมูลจาก SignJob
		modelIncome = models.Income{
			IncomeID:              uuid.NewString(),
			BankID:                config.DefaultBankAccountIDs.CompanyBank,
			TransactionCategoryID: config.DefaultTransactionCategoryIDs.CompanyIncome,
			Description:           signJob.Content,
			Amount:                amt,
			Currency:              "THB",
			TxnDate:               payDate,
			PaymentMethod:         input.PaymentMethod,
			ReferenceNo:           rec.InvoiceNo,
			Note:                  &signJob.JobName,
			CreatedBy:             claims.UserID,
			CreatedAt:             now,
			UpdatedAt:             now,
		}
	} else {
		// กรณีไม่มี SignJob - ใช้ข้อมูลจาก Receivable
		modelIncome = models.Income{
			IncomeID:              uuid.NewString(),
			BankID:                config.DefaultBankAccountIDs.CompanyBank,
			TransactionCategoryID: config.DefaultTransactionCategoryIDs.CompanyIncome,
			Description:           "รับชำระลูกหนี้ " + rec.Customer,
			Amount:                amt,
			Currency:              "THB",
			TxnDate:               payDate,
			PaymentMethod:         input.PaymentMethod,
			ReferenceNo:           rec.InvoiceNo,
			Note:                  &rec.Note,
			CreatedBy:             claims.UserID,
			CreatedAt:             now,
			UpdatedAt:             now,
		}
	}

	if err := s.incomeRepo.CreateInCome(ctx, modelIncome); err != nil {
		return err
	}

	return nil // สำเร็จ
}
