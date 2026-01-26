package dto

import "time"

// RequestListAuditLog represents query parameters for listing audit logs.
type RequestListAuditLog struct {
	Search    string `query:"search"`     // Search in path, email, full_name
	SortBy    string `query:"sort_by"`    // created_at (default)
	SortOrder string `query:"sort_order"` // asc, desc (default: desc)

	// Filters
	UserID    string `query:"user_id"`
	Action    string `query:"action"`     // CREATE, READ, UPDATE, DELETE
	Resource  string `query:"resource"`   // user, task, receipt, etc.
	Method    string `query:"method"`     // GET, POST, PUT, DELETE
	StartDate string `query:"start_date"` // YYYY-MM-DD
	EndDate   string `query:"end_date"`   // YYYY-MM-DD
	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
}

// ResponseAuditLog represents the response format for a single audit log.
type ResponseAuditLog struct {
	CreatedAt    time.Time `json:"created_at"`
	LogID        string    `json:"log_id"`
	UserID       string    `json:"user_id"`
	Email        string    `json:"email"`
	EmployeeCode string    `json:"employee_code"`
	Role         string    `json:"role"`
	FullName     string    `json:"full_name"`

	Method      string `json:"method"`
	Path        string `json:"path"`
	QueryParams string `json:"query_params,omitempty"`
	RequestBody string `json:"request_body,omitempty"`
	IPAddress   string `json:"ip_address"`
	UserAgent   string `json:"user_agent,omitempty"`

	Action     string `json:"action"`
	Resource   string `json:"resource"`
	ResourceID string `json:"resource_id,omitempty"`

	StatusCode   int   `json:"status_code"`
	ResponseTime int64 `json:"response_time_ms"`
}

// ResponseAuditLogList represents the paginated response for audit logs.
type ResponseAuditLogList struct {
	Items      []ResponseAuditLog `json:"items"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}
