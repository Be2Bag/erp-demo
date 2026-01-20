package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
)

// AuditLogRepository defines the interface for audit log data access
type AuditLogRepository interface {
	Create(ctx context.Context, log *models.AuditLog) error
	GetByFilter(ctx context.Context, filter interface{}, page, limit int, sortBy, sortOrder string) ([]*models.AuditLog, int64, error)
	GetByID(ctx context.Context, logID string) (*models.AuditLog, error)
}

// AuditLogService defines the interface for audit log business logic
type AuditLogService interface {
	// Log creates a new audit log entry (should be async to not block requests)
	Log(ctx context.Context, log *models.AuditLog) error
	// GetLogs retrieves audit logs with pagination and filtering
	GetLogs(ctx context.Context, req dto.RequestListAuditLog) (dto.Pagination, error)
	// GetLogByID retrieves a single audit log by its ID
	GetLogByID(ctx context.Context, logID string) (*dto.ResponseAuditLog, error)
}
