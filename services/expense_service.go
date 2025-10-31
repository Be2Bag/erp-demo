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

type expenseService struct {
	config                  config.Config
	expenseRepo             ports.ExpenseRepository
	transactionCategoryRepo ports.TransactionCategoryRepository
}

func NewExpenseService(cfg config.Config, expenseRepo ports.ExpenseRepository, transactionCategoryRepo ports.TransactionCategoryRepository) ports.ExpenseService {
	return &expenseService{config: cfg, expenseRepo: expenseRepo, transactionCategoryRepo: transactionCategoryRepo}
}

func (s *expenseService) CreateExpense(ctx context.Context, expense dto.CreateExpenseDTO, claims *dto.JWTClaims) error {
	now := time.Now()
	var due time.Time
	if expense.TxnDate != "" {
		parsedDate, err := time.Parse("2006-01-02", expense.TxnDate)
		if err != nil {
			return err
		}
		due = parsedDate
	}

	model := models.Expense{
		ExpenseID:             uuid.NewString(),
		TransactionCategoryID: expense.TransactionCategoryID,
		BankID:                expense.BankID,
		Description:           expense.Description,
		Amount:                expense.Amount,
		Currency:              expense.Currency,
		TxnDate:               due,
		PaymentMethod:         expense.PaymentMethod,
		ReferenceNo:           expense.ReferenceNo,
		Note:                  expense.Note,
		CreatedBy:             claims.UserID,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	if err := s.expenseRepo.CreateExpense(ctx, model); err != nil {
		return err
	}
	return nil
}

func (s *expenseService) ListExpenses(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, sortBy string, sortOrder string, transactionCategoryID string, startDate string, endDate string, bankID string) (dto.Pagination, error) {
	skip := int64((page - 1) * size)
	limit := int64(size)

	filter := bson.M{
		"deleted_at": nil,
	}
	if transactionCategoryID != "" {
		filter["transaction_category_id"] = transactionCategoryID
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

		filter["txn_date"] = txnDateFilter
	}

	search = strings.TrimSpace(search)
	if search != "" {
		safe := regexp.QuoteMeta(search)
		re := primitive.Regex{Pattern: safe, Options: "i"}
		filter["$or"] = []bson.M{
			{"description": re},
		}
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

	filterTransactionCategory := bson.M{
		"type":       "expense",
		"deleted_at": nil,
	}
	transactionCategory, errOnGetTransactionCategory := s.transactionCategoryRepo.GetAllTransactionCategoryByFilter(ctx, filterTransactionCategory, nil)
	if errOnGetTransactionCategory != nil {
		return dto.Pagination{}, fmt.Errorf("list expenses: %w", errOnGetTransactionCategory)
	}

	// สร้าง map สำหรับ mapping TransactionCategoryID กับ TransactionCategoryNameTH
	categoryMap := make(map[string]string)
	for _, cat := range transactionCategory {
		categoryMap[cat.TransactionCategoryID] = cat.TransactionCategoryNameTH
	}

	items, total, err := s.expenseRepo.GetListExpensesByFilter(ctx, filter, projection, sort, skip, limit)
	if err != nil {
		return dto.Pagination{}, fmt.Errorf("list expenses: %w", err)
	}

	list := make([]interface{}, 0, len(items))
	for _, m := range items {
		list = append(list, dto.ExpenseDTO{
			ExpenseID:                 m.ExpenseID,
			TransactionCategoryID:     m.TransactionCategoryID,
			TransactionCategoryNameTH: categoryMap[m.TransactionCategoryID],
			BankID:                    m.BankID,
			Description:               m.Description,
			Amount:                    m.Amount,
			Currency:                  m.Currency,
			TxnDate:                   m.TxnDate,
			PaymentMethod:             m.PaymentMethod,
			ReferenceNo:               m.ReferenceNo,
			Note:                      m.Note,
			CreatedBy:                 m.CreatedBy,
			CreatedAt:                 m.CreatedAt,
			UpdatedAt:                 m.UpdatedAt,
			DeletedAt:                 m.DeletedAt,
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

func (s *expenseService) GetExpenseByID(ctx context.Context, expenseID string, claims *dto.JWTClaims) (*dto.ExpenseDTO, error) {

	filter := bson.M{"expense_id": expenseID, "deleted_at": nil}
	projection := bson.M{}

	m, err := s.expenseRepo.GetOneExpenseByFilter(ctx, filter, projection)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}

	filterTransactionCategory := bson.M{
		"transaction_category_id": m.TransactionCategoryID,
		"type":                    "expense",
		"deleted_at":              nil,
	}
	transactionCategory, errOnGetTransactionCategory := s.transactionCategoryRepo.GetAllTransactionCategoryByFilter(ctx, filterTransactionCategory, nil)
	if errOnGetTransactionCategory != nil {
		return nil, fmt.Errorf("get expense by id: %w", errOnGetTransactionCategory)
	}

	transactionCategoryNameTH := ""
	if len(transactionCategory) > 0 {
		transactionCategoryNameTH = transactionCategory[0].TransactionCategoryNameTH
	}

	dtoObj := &dto.ExpenseDTO{
		// ---------- รายละเอียดรายจ่าย ----------
		ExpenseID:                 m.ExpenseID,
		TransactionCategoryID:     m.TransactionCategoryID,
		TransactionCategoryNameTH: transactionCategoryNameTH,
		BankID:                    m.BankID,
		Description:               m.Description,
		Amount:                    m.Amount,
		Currency:                  m.Currency,
		TxnDate:                   m.TxnDate,
		PaymentMethod:             m.PaymentMethod,
		ReferenceNo:               m.ReferenceNo,
		Note:                      m.Note,
		CreatedBy:                 m.CreatedBy,
		CreatedAt:                 m.CreatedAt,
		UpdatedAt:                 m.UpdatedAt,
		DeletedAt:                 m.DeletedAt,
	}
	return dtoObj, nil
}

func (s *expenseService) UpdateExpenseByID(ctx context.Context, expenseID string, update dto.UpdateExpenseDTO, claims *dto.JWTClaims) error {
	// ดึงข้อมูลเดิม
	filter := bson.M{"expense_id": expenseID, "deleted_at": nil}
	existing, err := s.expenseRepo.GetOneExpenseByFilter(ctx, filter, bson.M{})
	if err != nil {
		return err
	}
	if existing == nil {
		return mongo.ErrNoDocuments
	}

	if update.TransactionCategoryID != "" {
		existing.TransactionCategoryID = update.TransactionCategoryID
	}

	if update.BankID != "" {
		existing.BankID = update.BankID
	}

	if update.Description != "" {
		existing.Description = update.Description
	}
	if update.Amount > 0 {
		existing.Amount = update.Amount
	}
	if update.Currency != "" {
		existing.Currency = update.Currency
	}
	if update.TxnDate != "" {

		var due time.Time

		parsedDate, err := time.Parse("2006-01-02", update.TxnDate)
		if err != nil {
			return err
		}
		due = parsedDate

		existing.TxnDate = due
	}
	if update.PaymentMethod != "" {
		existing.PaymentMethod = update.PaymentMethod
	}
	if update.ReferenceNo != "" {
		existing.ReferenceNo = update.ReferenceNo
	}
	if update.Note != nil {
		existing.Note = update.Note
	}

	if _, err := s.expenseRepo.UpdateExpenseByID(ctx, expenseID, *existing); err != nil {
		return err
	}
	return nil
}

func (s *expenseService) DeleteExpenseByID(ctx context.Context, expenseID string, claims *dto.JWTClaims) error {
	err := s.expenseRepo.SoftDeleteExpenseByID(ctx, expenseID)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}

func (s *expenseService) SummaryExpenseByFilter(ctx context.Context, claims *dto.JWTClaims, report dto.RequestExpenseSummary) (dto.ExpenseSummaryDTO, error) {
	now := time.Now()

	filter := bson.M{
		"deleted_at": nil,
	}

	if report.BankID != "" {
		filter["bank_id"] = report.BankID
	}

	// Today
	startToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endToday := startToday.Add(24 * time.Hour)
	filterToday := bson.M{
		"deleted_at": nil,
		"txn_date": bson.M{
			"$gte": startToday,
			"$lt":  endToday,
		},
	}
	expensesToday, err := s.expenseRepo.GetAllExpenseByFilter(ctx, filterToday, nil)
	if err != nil {
		return dto.ExpenseSummaryDTO{}, err
	}
	var totalToday float64
	for _, expense := range expensesToday {
		totalToday += expense.Amount
	}

	// This Month
	startMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endMonth := startMonth.AddDate(0, 1, 0)
	filterMonth := bson.M{
		"deleted_at": nil,
		"txn_date": bson.M{
			"$gte": startMonth,
			"$lt":  endMonth,
		},
	}
	expensesMonth, err := s.expenseRepo.GetAllExpenseByFilter(ctx, filterMonth, nil)
	if err != nil {
		return dto.ExpenseSummaryDTO{}, err
	}
	var totalThisMonth float64
	for _, expense := range expensesMonth {
		totalThisMonth += expense.Amount
	}

	// All
	filterAll := bson.M{
		"deleted_at": nil,
	}
	expensesAll, err := s.expenseRepo.GetAllExpenseByFilter(ctx, filterAll, nil)
	if err != nil {
		return dto.ExpenseSummaryDTO{}, err
	}
	var totalAll float64
	for _, expense := range expensesAll {
		totalAll += expense.Amount
	}

	return dto.ExpenseSummaryDTO{
		TotalToday:     totalToday,
		TotalThisMonth: totalThisMonth,
		TotalAll:       totalAll,
	}, nil
}
