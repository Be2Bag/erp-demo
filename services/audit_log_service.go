package services

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/ports"
	"go.mongodb.org/mongo-driver/bson"
)

type auditLogService struct {
	repo     ports.AuditLogRepository
	logQueue chan *models.AuditLog
	config   config.Config
}

// NewAuditLogService creates a new audit log service with async logging support
func NewAuditLogService(cfg config.Config, repo ports.AuditLogRepository) ports.AuditLogService {
	svc := &auditLogService{
		config:   cfg,
		repo:     repo,
		logQueue: make(chan *models.AuditLog, 1000), // Buffer 1000 logs
	}

	// Start background worker to process logs
	go svc.processLogQueue()

	return svc
}

// processLogQueue handles async log insertion
func (s *auditLogService) processLogQueue() {
	for log := range s.logQueue {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := s.repo.Create(ctx, log); err != nil {
			// Log error but don't fail - audit logging should not affect main operations
			// In production, you might want to log to a backup location or retry
		}
		cancel()
	}
}

// Log creates a new audit log entry asynchronously
func (s *auditLogService) Log(ctx context.Context, log *models.AuditLog) error {
	// Send to queue for async processing
	select {
	case s.logQueue <- log:
		return nil
	default:
		// Queue is full, log synchronously as fallback
		return s.repo.Create(ctx, log)
	}
}

// GetLogs retrieves audit logs with pagination and filtering
func (s *auditLogService) GetLogs(ctx context.Context, req dto.RequestListAuditLog) (dto.Pagination, error) {
	// Build filter
	filter := bson.M{}

	if req.UserID != "" {
		filter["user_id"] = req.UserID
	}
	if req.Action != "" {
		filter["action"] = req.Action
	}
	if req.Resource != "" {
		filter["resource"] = req.Resource
	}
	if req.Method != "" {
		filter["method"] = strings.ToUpper(req.Method)
	}

	// Date range filter
	if req.StartDate != "" || req.EndDate != "" {
		dateFilter := bson.M{}
		if req.StartDate != "" {
			startDate, err := time.Parse("2006-01-02", req.StartDate)
			if err == nil {
				dateFilter["$gte"] = startDate
			}
		}
		if req.EndDate != "" {
			endDate, err := time.Parse("2006-01-02", req.EndDate)
			if err == nil {
				// Add 1 day to include the end date
				dateFilter["$lte"] = endDate.Add(24 * time.Hour)
			}
		}
		if len(dateFilter) > 0 {
			filter["created_at"] = dateFilter
		}
	}

	// Search filter
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"path": bson.M{"$regex": req.Search, "$options": "i"}},
			{"email": bson.M{"$regex": req.Search, "$options": "i"}},
			{"full_name": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}

	// Default pagination
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20
	}

	logs, total, err := s.repo.GetByFilter(ctx, filter, req.Page, req.Limit, req.SortBy, req.SortOrder)
	if err != nil {
		return dto.Pagination{}, err
	}

	// Convert to response format
	items := make([]interface{}, 0, len(logs))
	for _, log := range logs {
		items = append(items, dto.ResponseAuditLog{
			LogID:        log.LogID,
			UserID:       log.UserID,
			Email:        log.Email,
			EmployeeCode: log.EmployeeCode,
			Role:         log.Role,
			FullName:     log.FullName,
			Method:       log.Method,
			Path:         log.Path,
			QueryParams:  log.QueryParams,
			RequestBody:  log.RequestBody,
			IPAddress:    log.IPAddress,
			UserAgent:    log.UserAgent,
			StatusCode:   log.StatusCode,
			ResponseTime: log.ResponseTime,
			Action:       log.Action,
			Resource:     log.Resource,
			ResourceID:   log.ResourceID,
			CreatedAt:    log.CreatedAt,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	return dto.Pagination{
		List:       items,
		TotalCount: int(total),
		Page:       req.Page,
		Size:       req.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetLogByID retrieves a single audit log by its ID
func (s *auditLogService) GetLogByID(ctx context.Context, logID string) (*dto.ResponseAuditLog, error) {
	log, err := s.repo.GetByID(ctx, logID)
	if err != nil {
		return nil, err
	}
	if log == nil {
		return nil, nil
	}

	return &dto.ResponseAuditLog{
		LogID:        log.LogID,
		UserID:       log.UserID,
		Email:        log.Email,
		EmployeeCode: log.EmployeeCode,
		Role:         log.Role,
		FullName:     log.FullName,
		Method:       log.Method,
		Path:         log.Path,
		QueryParams:  log.QueryParams,
		RequestBody:  log.RequestBody,
		IPAddress:    log.IPAddress,
		UserAgent:    log.UserAgent,
		StatusCode:   log.StatusCode,
		ResponseTime: log.ResponseTime,
		Action:       log.Action,
		Resource:     log.Resource,
		ResourceID:   log.ResourceID,
		CreatedAt:    log.CreatedAt,
	}, nil
}
