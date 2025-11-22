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

type receiptService struct {
	config           config.Config
	receiptRepo      ports.ReceiptRepository
	bankAccountsRepo ports.BankAccountsRepository
}

func NewReceiptService(cfg config.Config, receiptRepo ports.ReceiptRepository, bankAccountsRepo ports.BankAccountsRepository) ports.ReceiptService {
	return &receiptService{config: cfg, receiptRepo: receiptRepo, bankAccountsRepo: bankAccountsRepo}
}

func (s *receiptService) CreateReceipt(ctx context.Context, in dto.CreateReceiptDTO, claims *dto.JWTClaims) error {
	now := time.Now()

	// parse receipt date
	receiptDate := now
	if strings.TrimSpace(in.ReceiptDate) != "" {
		t, err := time.Parse("2006-01-02", in.ReceiptDate)
		if err != nil {
			return fmt.Errorf("invalid receipt_date (want YYYY-MM-DD): %w", err)
		}
		receiptDate = t
	}

	// items + total
	if len(in.Items) == 0 {
		return fmt.Errorf("items required")
	}
	items := make([]models.ReceiptItem, 0, len(in.Items))
	var total float64
	for _, it := range in.Items {
		qty := it.Quantity
		unit := it.UnitPrice
		other := it.Other
		itemTotal := it.Total
		if itemTotal <= 0 {
			itemTotal = float64(qty)*unit + other
		}
		items = append(items, models.ReceiptItem{
			Description: it.Description,
			Quantity:    qty,
			UnitPrice:   unit,
			Other:       other,
			Total:       itemTotal,
		})
		total += itemTotal
	}

	// Apply VAT 7% if TypeReceipt is "company"
	subTotal := total
	var totalVAT float64
	if strings.ToLower(strings.TrimSpace(in.TypeReceipt)) == "company" {
		totalVAT = subTotal * 0.07
		total = subTotal + totalVAT
	}

	// payment info
	paidDate := now
	if strings.TrimSpace(in.PaymentDetail.PaidDate) != "" {
		t, err := time.Parse("2006-01-02", in.PaymentDetail.PaidDate)
		if err != nil {
			return fmt.Errorf("invalid paid_date (want YYYY-MM-DD): %w", err)
		}
		paidDate = t
	}
	payment := models.PaymentInfo{
		Method:        in.PaymentDetail.Method,
		BankName:      in.PaymentDetail.BankName,
		AccountName:   in.PaymentDetail.AccountName,
		AccountNumber: in.PaymentDetail.AccountNumber,
		AmountPaid:    in.PaymentDetail.AmountPaid,
		PaidDate:      paidDate,
		Note:          in.PaymentDetail.Note,
	}

	// customer, issuer
	customer := models.CustomerInfo{
		Name:    in.Customer.Name,
		Address: in.Customer.Address,
		Contact: in.Customer.Contact,
	}
	preparedBy := strings.TrimSpace(in.Issuer.PreparedBy)
	if preparedBy == "" {
		preparedBy = claims.UserID
	}
	issuer := models.IssuerInfo{
		Name:       in.Issuer.Name,
		Address:    in.Issuer.Address,
		Contact:    in.Issuer.Contact,
		Email:      in.Issuer.Email,
		PreparedBy: preparedBy,
	}

	// status and received by
	status := strings.ToLower(strings.TrimSpace(in.Status))
	if status == "" {
		if payment.AmountPaid >= total && total > 0 {
			status = "paid"
		} else {
			status = "pending"
		}
	}
	receivedBy := strings.TrimSpace(in.ReceivedBy)
	if receivedBy == "" {
		receivedBy = claims.UserID
	}

	// IV014-DD-MM-YY-XXX
	year := now.Year() % 100
	yearPrefix := fmt.Sprintf("IV%03d", 14+year-25) // IV014 for 2025, IV015 for 2026, etc.
	datePrefix := fmt.Sprintf("%s-%02d-%02d-%02d", yearPrefix, now.Day(), int(now.Month()), year)

	maxNumber, err := s.receiptRepo.GetMaxReceiptNumber(ctx, datePrefix)
	if err != nil {
		return fmt.Errorf("failed to generate receipt number: %w", err)
	}

	sequence := 0
	if maxNumber != "" {
		parts := strings.Split(maxNumber, "-")
		if len(parts) == 5 {
			_, _ = fmt.Sscanf(parts[4], "%d", &sequence)
		}
	}

	receiptNumber := fmt.Sprintf("%s-%03d", datePrefix, sequence+1)

	model := models.Receipt{
		IDReceipt:     uuid.NewString(),
		ReceiptNumber: receiptNumber,
		ReceiptDate:   receiptDate,
		Customer:      customer,
		Issuer:        issuer,
		Items:         items,
		SubTotal:      subTotal,
		TotalVAT:      totalVAT,
		TotalAmount:   total,
		Remark:        in.Remark,
		PaymentDetail: payment,
		Status:        status,
		BillType:      in.BillType,
		TypeReceipt:   in.TypeReceipt,
		ApprovedBy:    in.ApprovedBy,
		ReceivedBy:    receivedBy,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.receiptRepo.CreateReceipt(ctx, model); err != nil {
		return err
	}
	return nil
}

func (s *receiptService) ListReceipts(ctx context.Context, claims *dto.JWTClaims, page, size int, search, sortBy, sortOrder, status, startDate, endDate, billType, typeReceipt string) (dto.Pagination, error) {
	skip := int64((page - 1) * size)
	limit := int64(size)

	filter := bson.M{
		"deleted_at": nil,
	}

	if strings.TrimSpace(typeReceipt) != "" {
		filter["type_receipt"] = strings.ToLower(strings.TrimSpace(typeReceipt))
	}

	// Bill type filter
	if strings.TrimSpace(billType) != "" {
		filter["bill_type"] = strings.ToLower(strings.TrimSpace(billType))
	}

	// Date range filter on receipt_date
	if strings.TrimSpace(startDate) != "" || strings.TrimSpace(endDate) != "" {
		dateFilter := bson.M{}

		if strings.TrimSpace(startDate) != "" {
			parsedStartDate, err := time.ParseInLocation("2006-01-02", startDate, time.UTC)
			if err != nil {
				return dto.Pagination{}, fmt.Errorf("invalid startDate format: %w", err)
			}
			dateFilter["$gte"] = parsedStartDate
		}
		if strings.TrimSpace(endDate) != "" {
			parsedEndDate, err := time.ParseInLocation("2006-01-02", endDate, time.UTC)
			if err != nil {
				return dto.Pagination{}, fmt.Errorf("invalid endDate format: %w", err)
			}
			// include entire endDate day
			dateFilter["$lt"] = parsedEndDate.Add(24 * time.Hour)
		}

		filter["receipt_date"] = dateFilter
	}

	// Search across common fields
	if sText := strings.TrimSpace(search); sText != "" {
		safe := regexp.QuoteMeta(sText)
		re := primitive.Regex{Pattern: safe, Options: "i"}
		filter["$or"] = []bson.M{
			{"receipt_number": re},
			{"customer.name": re},
			{"issuer.name": re},
			{"issuer.email": re},
			{"items.description": re},
			{"remark": re},
			{"payment_detail.method": re},
			{"payment_detail.bank_name": re},
			{"payment_detail.account_name": re},
			{"payment_detail.account_number": re},
		}
	}

	// Status filter
	if strings.TrimSpace(status) != "" {
		filter["status"] = strings.ToLower(strings.TrimSpace(status))
	}

	projection := bson.M{}

	// Sorting
	allowedSortFields := map[string]string{
		"created_at":   "created_at",
		"updated_at":   "updated_at",
		"receipt_date": "receipt_date",
		"total_amount": "total_amount",
	}
	field := allowedSortFields[strings.ToLower(strings.TrimSpace(sortBy))]
	if field == "" {
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

	items, total, err := s.receiptRepo.GetListReceiptsByFilter(ctx, filter, projection, sort, skip, limit)
	if err != nil {
		return dto.Pagination{}, fmt.Errorf("list receipts: %w", err)
	}

	list := make([]interface{}, 0, len(items))
	for _, m := range items {
		// map items
		dtoItems := make([]dto.ReceiptItemDTO, 0, len(m.Items))
		for _, it := range m.Items {
			dtoItems = append(dtoItems, dto.ReceiptItemDTO{
				Description: it.Description,
				Quantity:    it.Quantity,
				UnitPrice:   it.UnitPrice,
				Other:       it.Other,
				Total:       it.Total,
			})
		}

		// map DTO
		list = append(list, dto.ReceiptDTO{
			IDReceipt:     m.IDReceipt,
			ReceiptNumber: m.ReceiptNumber,
			ReceiptDate:   m.ReceiptDate,
			Customer: dto.CustomerInfoDTO{
				Name:    m.Customer.Name,
				Address: m.Customer.Address,
				Contact: m.Customer.Contact,
			},
			Issuer: dto.IssuerInfoDTO{
				Name:       m.Issuer.Name,
				Address:    m.Issuer.Address,
				Contact:    m.Issuer.Contact,
				Email:      m.Issuer.Email,
				PreparedBy: m.Issuer.PreparedBy,
			},
			Items:       dtoItems,
			SubTotal:    m.SubTotal,
			TotalVAT:    m.TotalVAT,
			TotalAmount: m.TotalAmount,
			Remark:      m.Remark,
			PaymentDetail: dto.PaymentInfoRespDTO{
				Method:        m.PaymentDetail.Method,
				BankName:      m.PaymentDetail.BankName,
				AccountName:   m.PaymentDetail.AccountName,
				AccountNumber: m.PaymentDetail.AccountNumber,
				AmountPaid:    m.PaymentDetail.AmountPaid,
				PaidDate:      m.PaymentDetail.PaidDate,
				Note:          m.PaymentDetail.Note,
			},
			Status:     m.Status,
			BillType:   m.BillType,
			ApprovedBy: m.ApprovedBy,
			ReceivedBy: m.ReceivedBy,
			CreatedAt:  m.CreatedAt,
			UpdatedAt:  m.UpdatedAt,
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

func (s *receiptService) GetReceiptByID(ctx context.Context, receiptID string, claims *dto.JWTClaims) (*dto.ReceiptDTO, error) {
	filter := bson.M{"id_receipt": strings.TrimSpace(receiptID), "deleted_at": nil}
	projection := bson.M{}

	m, err := s.receiptRepo.GetOneReceiptsByFilter(ctx, filter, projection)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}

	// map items
	dtoItems := make([]dto.ReceiptItemDTO, 0, len(m.Items))
	for _, it := range m.Items {
		dtoItems = append(dtoItems, dto.ReceiptItemDTO{
			Description: it.Description,
			Quantity:    it.Quantity,
			UnitPrice:   it.UnitPrice,
			Other:       it.Other,
			Total:       it.Total,
		})
	}

	// map to DTO
	dtoObj := &dto.ReceiptDTO{
		IDReceipt:     m.IDReceipt,
		ReceiptNumber: m.ReceiptNumber,
		ReceiptDate:   m.ReceiptDate,
		Customer: dto.CustomerInfoDTO{
			Name:    m.Customer.Name,
			Address: m.Customer.Address,
			Contact: m.Customer.Contact,
		},
		Issuer: dto.IssuerInfoDTO{
			Name:       m.Issuer.Name,
			Address:    m.Issuer.Address,
			Contact:    m.Issuer.Contact,
			Email:      m.Issuer.Email,
			PreparedBy: m.Issuer.PreparedBy,
		},
		Items:       dtoItems,
		SubTotal:    m.SubTotal,
		TotalVAT:    m.TotalVAT,
		TotalAmount: m.TotalAmount,
		Remark:      m.Remark,
		PaymentDetail: dto.PaymentInfoRespDTO{
			Method:        m.PaymentDetail.Method,
			BankName:      m.PaymentDetail.BankName,
			AccountName:   m.PaymentDetail.AccountName,
			AccountNumber: m.PaymentDetail.AccountNumber,
			AmountPaid:    m.PaymentDetail.AmountPaid,
			PaidDate:      m.PaymentDetail.PaidDate,
			Note:          m.PaymentDetail.Note,
		},
		Status:     m.Status,
		BillType:   m.BillType,
		ApprovedBy: m.ApprovedBy,
		ReceivedBy: m.ReceivedBy,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}

	return dtoObj, nil
}

func (s *receiptService) DeleteReceiptByID(ctx context.Context, receiptID string, claims *dto.JWTClaims) error {
	err := s.receiptRepo.SoftDeleteReceiptByID(ctx, receiptID)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}

func (s *receiptService) SummaryReceiptByFilter(ctx context.Context, claims *dto.JWTClaims, report dto.RequestSummaryReceipt) (dto.ReceiptSummaryDTO, error) {
	now := time.Now()

	// base filter
	filter := bson.M{
		"deleted_at": nil,
	}

	// date range by report type: day | month | all
	switch strings.ToLower(strings.TrimSpace(report.Report)) {
	case "day":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end := start.Add(24 * time.Hour)
		filter["receipt_date"] = bson.M{"$gte": start, "$lt": end}
	case "month":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 1, 0)
		filter["receipt_date"] = bson.M{"$gte": start, "$lt": end}
	case "all", "":
		// no date filter
	default:
		// unknown report type -> treat as all
	}

	receipts, err := s.receiptRepo.GetAllReceiptsByFilter(ctx, filter, nil)
	if err != nil {
		return dto.ReceiptSummaryDTO{}, err
	}

	var totalAmount float64
	var totalPaid float64
	var pendingCount int

	for _, r := range receipts {
		totalAmount += r.TotalAmount
		totalPaid += r.PaymentDetail.AmountPaid

		outstanding := r.TotalAmount - r.PaymentDetail.AmountPaid
		if strings.ToLower(strings.TrimSpace(r.Status)) != "paid" && outstanding > 0 {
			pendingCount++
		}
	}

	return dto.ReceiptSummaryDTO{
		TotalAmount:  totalAmount,
		TotalPaid:    totalPaid,
		PendingCount: pendingCount,
	}, nil
}
