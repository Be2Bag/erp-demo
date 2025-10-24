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

type receivableService struct {
	config         config.Config
	receivableRepo ports.ReceivableRepository
}

func NewReceivableService(cfg config.Config, receivableRepo ports.ReceivableRepository) ports.ReceivableService {
	return &receivableService{config: cfg, receivableRepo: receivableRepo}
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
		Customer:     receivable.Customer,
		InvoiceNo:    receivable.InvoiceNo,
		IssueDate:    issue,
		DueDate:      due,
		Amount:       receivable.Amount,
		Balance:      receivable.Balance,
		Status:       "pending",
		CreatedBy:    claims.UserID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.receivableRepo.CreateReceivable(ctx, model); err != nil {
		return err
	}
	return nil
}

func (s *receivableService) ListReceivables(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, sortBy string, sortOrder string, Status string) (dto.Pagination, error) {
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

	list := make([]interface{}, 0, len(items))
	for _, m := range items {

		list = append(list, dto.ReceivableDTO{
			IDReceivable: m.IDReceivable,
			Customer:     m.Customer,
			InvoiceNo:    m.InvoiceNo,
			IssueDate:    m.IssueDate,
			DueDate:      m.DueDate,
			Amount:       m.Amount,
			Balance:      m.Balance,
			Status:       m.Status,
			CreatedBy:    m.CreatedBy,
			CreatedAt:    m.CreatedAt,
			UpdatedAt:    m.UpdatedAt,
			DeletedAt:    m.DeletedAt,
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

	filter := bson.M{"_id": receivableID, "deleted_at": nil}
	projection := bson.M{}

	m, err := s.receivableRepo.GetOneReceivableByFilter(ctx, filter, projection)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}

	dtoObj := &dto.ReceivableDTO{
		// ---------- รายละเอียดรายได้ ----------
		IDReceivable: m.IDReceivable,
		Customer:     m.Customer,
		InvoiceNo:    m.InvoiceNo,
		IssueDate:    m.IssueDate,
		DueDate:      m.DueDate,
		Amount:       m.Amount,
		Balance:      m.Balance,
		Status:       m.Status,
		CreatedBy:    m.CreatedBy,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		DeletedAt:    m.DeletedAt,
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

func (s *receivableService) SummaryReceivableByFilter(ctx context.Context, claims *dto.JWTClaims) (dto.ReceivableSummaryDTO, error) {
	now := time.Now()

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
	receivableToday, err := s.receivableRepo.GetAllReceivablesByFilter(ctx, filterToday, nil)
	if err != nil {
		return dto.ReceivableSummaryDTO{}, err
	}
	var totalToday float64
	for _, receivable := range receivableToday {
		totalToday += receivable.Amount
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
	receivableMonth, err := s.receivableRepo.GetAllReceivablesByFilter(ctx, filterMonth, nil)
	if err != nil {
		return dto.ReceivableSummaryDTO{}, err
	}
	var totalThisMonth float64
	for _, receivable := range receivableMonth {
		totalThisMonth += receivable.Amount
	}

	// All
	filterAll := bson.M{
		"deleted_at": nil,
	}
	receivableAll, err := s.receivableRepo.GetAllReceivablesByFilter(ctx, filterAll, nil)
	if err != nil {
		return dto.ReceivableSummaryDTO{}, err
	}
	var totalAll float64
	for _, receivable := range receivableAll {
		totalAll += receivable.Amount
	}

	return dto.ReceivableSummaryDTO{
		TotalAmount: totalAll,
		TotalPaid:   totalThisMonth,
		TotalDue:    totalToday,
	}, nil
}
