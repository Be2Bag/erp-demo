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

type payablesService struct {
	config           config.Config
	payablesRepo     ports.PayableRepository
	bankAccountsRepo ports.BankAccountsRepository
}

func NewPayablesService(cfg config.Config, payablesRepo ports.PayableRepository, bankAccountsRepo ports.BankAccountsRepository) ports.PayableService {
	return &payablesService{config: cfg, payablesRepo: payablesRepo, bankAccountsRepo: bankAccountsRepo}
}

func (s *payablesService) CreatePayable(ctx context.Context, payable dto.CreatePayableDTO, claims *dto.JWTClaims) error {
	now := time.Now()

	var issue, due time.Time
	if strings.TrimSpace(payable.IssueDate) != "" {
		t, err := time.Parse("2006-01-02", payable.IssueDate)
		if err != nil {
			return err
		}
		issue = t
	}
	if strings.TrimSpace(payable.DueDate) != "" {
		t, err := time.Parse("2006-01-02", payable.DueDate)
		if err != nil {
			return err
		}
		due = t
	}

	balance := payable.Balance
	if balance <= 0 {
		balance = payable.Amount
	}

	model := models.Payable{
		IDPayable:  uuid.NewString(),
		BankID:     strings.TrimSpace(payable.BankID),
		Supplier:   strings.TrimSpace(payable.Supplier),
		PurchaseNo: strings.TrimSpace(payable.PurchaseNo),
		InvoiceNo:  strings.TrimSpace(payable.InvoiceNo),
		IssueDate:  issue,
		DueDate:    due,
		Amount:     payable.Amount,
		Balance:    balance,
		Status:     "pending",
		PaymentRef: strings.TrimSpace(payable.PaymentRef),
		Note:       strings.TrimSpace(payable.Note),
		CreatedBy:  claims.UserID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.payablesRepo.CreatePayable(ctx, model); err != nil {
		return err
	}
	return nil
}

func (s *payablesService) ListPayables(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, sortBy string, sortOrder string, status string, startDate string, endDate string, bankID string) (dto.Pagination, error) {
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
		}
	}

	if status != "" {
		filter["status"] = strings.ToLower(status)
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

	items, total, err := s.payablesRepo.GetListPayablesByFilter(ctx, filter, projection, sort, skip, limit)
	if err != nil {
		return dto.Pagination{}, fmt.Errorf("list payables: %w", err)
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

		list = append(list, dto.PayableDTO{
			IDPayable:  m.IDPayable,
			BankID:     m.BankID,
			BankName:   bankMap[m.BankID],
			Supplier:   m.Supplier,
			InvoiceNo:  m.InvoiceNo,
			PurchaseNo: m.PurchaseNo,
			IssueDate:  m.IssueDate,
			DueDate:    m.DueDate,
			Amount:     m.Amount,
			Balance:    m.Balance,
			Status:     m.Status,
			PaymentRef: m.PaymentRef,
			Note:       m.Note,
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

func (s *payablesService) GetPayableByID(ctx context.Context, payableID string, claims *dto.JWTClaims) (*dto.PayableDTO, error) {

	filter := bson.M{"id_payable": payableID, "deleted_at": nil}
	projection := bson.M{}

	m, err := s.payablesRepo.GetOnePayableByFilter(ctx, filter, projection)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}

	filterPaymentTransaction := bson.M{"ref_invoice_no": m.InvoiceNo, "transaction_type": "payable", "deleted_at": nil}

	paymentTransaction, errPaymentTransaction := s.payablesRepo.GetAllPaymentTransactionByFilter(ctx, filterPaymentTransaction, bson.M{})
	if errPaymentTransaction != nil {
		return nil, errPaymentTransaction
	}

	filterBankAccounts := bson.M{"deleted_at": nil}
	bankAccounts, errOnGetBankAccounts := s.bankAccountsRepo.GetAllBankAccountsByFilter(ctx, filterBankAccounts, nil)
	if errOnGetBankAccounts != nil {
		return nil, errOnGetBankAccounts
	}

	// สร้าง map สำหรับ mapping BankID กับ BankName
	bankMap := make(map[string]string)
	for _, bank := range bankAccounts {
		bankMap[bank.BankID] = bank.BankName
	}

	PaymentTransactions := make([]dto.PaymentTransactionDTO, 0, len(paymentTransaction))
	if len(paymentTransaction) > 0 {

		for _, pt := range paymentTransaction {
			PaymentTransactions = append(PaymentTransactions, dto.PaymentTransactionDTO{
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
	}

	dtoObj := &dto.PayableDTO{
		// ---------- รายละเอียดเจ้าหนี้ ----------
		IDPayable:    m.IDPayable,
		BankID:       m.BankID,
		BankName:     bankMap[m.BankID],
		Supplier:     m.Supplier,
		InvoiceNo:    m.InvoiceNo,
		IssueDate:    m.IssueDate,
		DueDate:      m.DueDate,
		Amount:       m.Amount,
		Balance:      m.Balance,
		Status:       m.Status,
		PaymentRef:   m.PaymentRef,
		Note:         m.Note,
		Transactions: PaymentTransactions,
	}
	return dtoObj, nil
}

func (s *payablesService) UpdatePayableByID(ctx context.Context, payableID string, update dto.UpdatePayableDTO, claims *dto.JWTClaims) error {
	// fetch existing
	filter := bson.M{"id_payable": payableID, "deleted_at": nil}
	existing, err := s.payablesRepo.GetOnePayableByFilter(ctx, filter, bson.M{})
	if err != nil {
		return err
	}
	if existing == nil {
		return mongo.ErrNoDocuments
	}

	if strings.TrimSpace(update.Supplier) != "" {
		existing.Supplier = update.Supplier
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

	existing.UpdatedAt = time.Now()

	if _, err := s.payablesRepo.UpdatePayableByID(ctx, payableID, *existing); err != nil {
		return err
	}
	return nil
}

func (s *payablesService) DeletePayableByID(ctx context.Context, payableID string, claims *dto.JWTClaims) error {
	err := s.payablesRepo.SoftDeletePayableByID(ctx, payableID)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}

func (s *payablesService) SummaryPayableByFilter(ctx context.Context, claims *dto.JWTClaims, report dto.RequestSummaryPayable) (dto.PayableSummaryDTO, error) {
	now := time.Now()

	// base filter
	filter := bson.M{
		"deleted_at": nil,
	}
	if strings.TrimSpace(report.BankID) != "" {
		filter["bank_id"] = strings.TrimSpace(report.BankID)
	}

	// date range by report type: day | month | all
	switch strings.ToLower(strings.TrimSpace(report.Report)) {
	case "day":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end := start.Add(24 * time.Hour)
		filter["issue_date"] = bson.M{"$gte": start, "$lt": end}
	case "month":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 1, 0)
		filter["issue_date"] = bson.M{"$gte": start, "$lt": end}
	case "all", "":
		// no date filter
	default:
		// unknown report type -> treat as all
	}

	payables, err := s.payablesRepo.GetAllPayablesByFilter(ctx, filter, nil)
	if err != nil {
		return dto.PayableSummaryDTO{}, err
	}

	var totalAmount float64
	var totalDue float64
	overdueCount := 0

	for _, p := range payables {
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

	return dto.PayableSummaryDTO{
		TotalAmount:  totalAmount,
		TotalDue:     totalDue,
		OverdueCount: overdueCount,
	}, nil
}

func (s *payablesService) RecordPayment(ctx context.Context, input dto.RecordPaymentDTO, claims *dto.JWTClaims) error {
	// 1) fetch payable
	filter := bson.M{"id_payable": strings.TrimSpace(input.PayableID), "deleted_at": nil}
	payable, err := s.payablesRepo.GetOnePayableByFilter(ctx, filter, bson.M{})
	if err != nil {
		return fmt.Errorf("get payable: %w", err)
	}
	if payable == nil {
		return mongo.ErrNoDocuments
	}

	// 2) validate amount
	amt := input.Amount
	if amt <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	if payable.Balance <= 0 {
		return fmt.Errorf("payable is already fully paid")
	}
	if amt > payable.Balance {
		return fmt.Errorf("amount %.2f exceeds remaining balance %.2f", amt, payable.Balance)
	}

	// 3) parse date
	now := time.Now()
	payDate := now
	if strings.TrimSpace(input.PaymentDate) != "" {
		t, err := time.Parse("2006-01-02", input.PaymentDate)
		if err != nil {
			return fmt.Errorf("invalid payment_date, want YYYY-MM-DD: %w", err)
		}
		payDate = t
	}

	// 4) build transaction
	refInvoice := payable.InvoiceNo
	if strings.TrimSpace(refInvoice) == "" {
		refInvoice = payable.IDPayable
	}
	tx := models.PaymentTransaction{
		IDTransaction:   uuid.NewString(),
		BankID:          strings.TrimSpace(input.BankID),
		RefInvoiceNo:    refInvoice,
		TransactionType: "payable",
		PaymentDate:     payDate,
		Amount:          amt,
		PaymentMethod:   strings.TrimSpace(input.PaymentMethod),
		PaymentRef:      strings.TrimSpace(input.PaymentRef),
		Note:            strings.TrimSpace(input.Note),
		CreatedBy:       claims.UserID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// 5) insert transaction first
	if err := s.payablesRepo.CreatePaymentTransaction(ctx, tx); err != nil {
		return fmt.Errorf("insert payment tx: %w", err)
	}

	// 6) update payable balance and status
	newBalance := payable.Balance - amt
	if newBalance < 0 {
		newBalance = 0
	}
	payable.Balance = newBalance

	// compute status
	if payable.Balance == 0 {
		payable.Status = "paid"
	} else if !payable.DueDate.IsZero() && payable.DueDate.Before(now) {
		payable.Status = "overdue"
	} else {
		payable.Status = "partial"
	}
	payable.UpdatedAt = now

	if _, err := s.payablesRepo.UpdatePayableByID(ctx, payable.IDPayable, *payable); err != nil {
		return fmt.Errorf("update payable: %w", err)
	}

	return nil
}
