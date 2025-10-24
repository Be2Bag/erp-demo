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
	config       config.Config
	payablesRepo ports.PayableRepository
}

func NewPayablesService(cfg config.Config, payablesRepo ports.PayableRepository) ports.PayableService {
	return &payablesService{config: cfg, payablesRepo: payablesRepo}
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

func (s *payablesService) ListPayables(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, sortBy string, sortOrder string, status string) (dto.Pagination, error) {
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

	list := make([]interface{}, 0, len(items))
	for _, m := range items {

		list = append(list, dto.PayableDTO{
			IDPayable: m.IDPayable,
			Supplier:  m.Supplier,
			InvoiceNo: m.InvoiceNo,
			IssueDate: m.IssueDate,
			DueDate:   m.DueDate,
			Amount:    m.Amount,
			Balance:   m.Balance,
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

	filter := bson.M{"_id": payableID, "deleted_at": nil}
	projection := bson.M{}

	m, err := s.payablesRepo.GetOnePayableByFilter(ctx, filter, projection)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}

	dtoObj := &dto.PayableDTO{
		// ---------- รายละเอียดเจ้าหนี้ ----------
		IDPayable: m.IDPayable,
		Supplier:  m.Supplier,
		InvoiceNo: m.InvoiceNo,
		IssueDate: m.IssueDate,
		DueDate:   m.DueDate,
		Amount:    m.Amount,
		Balance:   m.Balance,
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

func (s *payablesService) SummaryPayableByFilter(ctx context.Context, claims *dto.JWTClaims) (dto.PayableSummaryDTO, error) {
	now := time.Now()

	// Today
	startToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endToday := startToday.Add(24 * time.Hour)
	filterToday := bson.M{
		"deleted_at": nil,
		"issue_date": bson.M{
			"$gte": startToday,
			"$lt":  endToday,
		},
	}
	payableToday, err := s.payablesRepo.GetAllPayablesByFilter(ctx, filterToday, nil)
	if err != nil {
		return dto.PayableSummaryDTO{}, err
	}
	var totalToday float64
	for _, payable := range payableToday {
		totalToday += payable.Amount
	}

	// This Month
	startMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endMonth := startMonth.AddDate(0, 1, 0)
	filterMonth := bson.M{
		"deleted_at": nil,
		"issue_date": bson.M{
			"$gte": startMonth,
			"$lt":  endMonth,
		},
	}
	payableMonth, err := s.payablesRepo.GetAllPayablesByFilter(ctx, filterMonth, nil)
	if err != nil {
		return dto.PayableSummaryDTO{}, err
	}
	var totalThisMonth float64
	for _, payable := range payableMonth {
		totalThisMonth += payable.Amount
	}

	// All
	filterAll := bson.M{
		"deleted_at": nil,
	}
	payableAll, err := s.payablesRepo.GetAllPayablesByFilter(ctx, filterAll, nil)
	if err != nil {
		return dto.PayableSummaryDTO{}, err
	}
	var totalAll float64
	for _, payable := range payableAll {
		totalAll += payable.Amount
	}

	return dto.PayableSummaryDTO{
		TotalAmount: totalAll,
		TotalPaid:   totalThisMonth,
		TotalDue:    totalToday,
	}, nil
}
