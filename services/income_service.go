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

type inComeService struct {
	config                  config.Config
	inComeRepo              ports.InComeRepository
	transactionCategoryRepo ports.TransactionCategoryRepository
}

func NewInComeService(cfg config.Config, inComeRepo ports.InComeRepository, transactionCategoryRepo ports.TransactionCategoryRepository) ports.InComeService {
	return &inComeService{config: cfg, inComeRepo: inComeRepo, transactionCategoryRepo: transactionCategoryRepo}
}

func (s *inComeService) CreateInCome(ctx context.Context, inCome dto.CreateIncomeDTO, claims *dto.JWTClaims) error {
	now := time.Now()
	var due time.Time
	if inCome.TxnDate != "" {
		parsedDate, err := time.Parse("2006-01-02", inCome.TxnDate)
		if err != nil {
			return err
		}
		due = parsedDate
	}

	model := models.Income{
		IncomeID:              uuid.NewString(),
		BankID:                inCome.BankID,
		TransactionCategoryID: inCome.TransactionCategoryID,
		Description:           inCome.Description,
		Amount:                inCome.Amount,
		Currency:              inCome.Currency,
		TxnDate:               due,
		PaymentMethod:         inCome.PaymentMethod,
		ReferenceNo:           inCome.ReferenceNo,
		Note:                  inCome.Note,
		CreatedBy:             claims.UserID,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	if err := s.inComeRepo.CreateInCome(ctx, model); err != nil {
		return err
	}
	return nil
}

func (s *inComeService) ListInComes(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, sortBy string, sortOrder string, transactionCategoryID string) (dto.Pagination, error) {
	skip := int64((page - 1) * size)
	limit := int64(size)

	filter := bson.M{
		"deleted_at": nil,
	}

	if transactionCategoryID != "" {
		filter["transaction_category_id"] = transactionCategoryID
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

	items, total, err := s.inComeRepo.GetListInComesByFilter(ctx, filter, projection, sort, skip, limit)
	if err != nil {
		return dto.Pagination{}, fmt.Errorf("list incomes: %w", err)
	}

	filterTransactionCategory := bson.M{
		"type":       "income",
		"deleted_at": nil,
	}
	transactionCategory, errOnGetTransactionCategory := s.transactionCategoryRepo.GetAllTransactionCategoryByFilter(ctx, filterTransactionCategory, nil)
	if errOnGetTransactionCategory != nil {
		return dto.Pagination{}, fmt.Errorf("list incomes: %w", errOnGetTransactionCategory)
	}

	// สร้าง map สำหรับ mapping TransactionCategoryID กับ TransactionCategoryNameTH
	categoryMap := make(map[string]string)
	for _, cat := range transactionCategory {
		categoryMap[cat.TransactionCategoryID] = cat.TransactionCategoryNameTH
	}

	list := make([]interface{}, 0, len(items))
	for _, m := range items {

		list = append(list, dto.IncomeDTO{
			IncomeID:                  m.IncomeID,
			BankID:                    m.BankID,
			TransactionCategoryID:     m.TransactionCategoryID,
			TransactionCategoryNameTH: categoryMap[m.TransactionCategoryID],
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

func (s *inComeService) GetIncomeByID(ctx context.Context, incomeID string, claims *dto.JWTClaims) (*dto.IncomeDTO, error) {

	filter := bson.M{"income_id": incomeID, "deleted_at": nil}
	projection := bson.M{}

	m, err := s.inComeRepo.GetOneInComeByFilter(ctx, filter, projection)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}

	filterTransactionCategory := bson.M{
		"transaction_category_id": m.TransactionCategoryID,
		"type":                    "income",
		"deleted_at":              nil,
	}
	transactionCategory, errOnGetTransactionCategory := s.transactionCategoryRepo.GetAllTransactionCategoryByFilter(ctx, filterTransactionCategory, nil)
	if errOnGetTransactionCategory != nil {
		return nil, fmt.Errorf("get income by id: %w", errOnGetTransactionCategory)
	}

	transactionCategoryNameTH := ""
	if len(transactionCategory) > 0 {
		transactionCategoryNameTH = transactionCategory[0].TransactionCategoryNameTH
	}

	dtoObj := &dto.IncomeDTO{
		// ---------- รายละเอียดรายได้ ----------
		IncomeID:                  m.IncomeID,
		BankID:                    m.BankID,
		TransactionCategoryID:     m.TransactionCategoryID,
		TransactionCategoryNameTH: transactionCategoryNameTH,
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

func (s *inComeService) UpdateInComeByID(ctx context.Context, incomeID string, update dto.UpdateIncomeDTO, claims *dto.JWTClaims) error {
	// ดึงข้อมูลเดิม
	filter := bson.M{"income_id": incomeID, "deleted_at": nil}
	existing, err := s.inComeRepo.GetOneInComeByFilter(ctx, filter, bson.M{})
	if err != nil {
		return err
	}
	if existing == nil {
		return mongo.ErrNoDocuments
	}

	if update.BankID != "" {
		existing.BankID = update.BankID
	}
	if update.TransactionCategoryID != "" {
		existing.TransactionCategoryID = update.TransactionCategoryID
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

	if _, err := s.inComeRepo.UpdateInComeByID(ctx, incomeID, *existing); err != nil {
		return err
	}
	return nil
}

func (s *inComeService) DeleteInComeByInComeID(ctx context.Context, incomeID string, claims *dto.JWTClaims) error {
	err := s.inComeRepo.SoftDeleteInComeByincomeID(ctx, incomeID)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}

func (s *inComeService) SummaryInComeByFilter(ctx context.Context, claims *dto.JWTClaims, report dto.RequestIncomeSummary) (dto.IncomeSummaryDTO, error) {
	now := time.Now()

	filter := bson.M{
		"deleted_at": nil,
	}
	if strings.TrimSpace(report.BankID) != "" {
		filter["bank_id"] = strings.TrimSpace(report.BankID)
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
	incomesToday, err := s.inComeRepo.GetAllInComeByFilter(ctx, filterToday, nil)
	if err != nil {
		return dto.IncomeSummaryDTO{}, err
	}
	var totalToday float64
	for _, income := range incomesToday {
		totalToday += income.Amount
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
	incomesMonth, err := s.inComeRepo.GetAllInComeByFilter(ctx, filterMonth, nil)
	if err != nil {
		return dto.IncomeSummaryDTO{}, err
	}
	var totalThisMonth float64
	for _, income := range incomesMonth {
		totalThisMonth += income.Amount
	}

	// All
	filterAll := bson.M{
		"deleted_at": nil,
	}
	incomesAll, err := s.inComeRepo.GetAllInComeByFilter(ctx, filterAll, nil)
	if err != nil {
		return dto.IncomeSummaryDTO{}, err
	}
	var totalAll float64
	for _, income := range incomesAll {
		totalAll += income.Amount
	}

	return dto.IncomeSummaryDTO{
		TotalToday:     totalToday,
		TotalThisMonth: totalThisMonth,
		TotalAll:       totalAll,
	}, nil
}
