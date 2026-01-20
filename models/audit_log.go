package models

import "time"

// CollectionAuditLogs is the MongoDB collection name for audit logs
const CollectionAuditLogs = "audit_logs"

// AuditLog represents a record of user actions in the system
type AuditLog struct {
	ID    string `bson:"_id,omitempty" json:"id"`
	LogID string `bson:"log_id" json:"log_id"`

	// User Information
	UserID       string `bson:"user_id" json:"user_id"`
	Email        string `bson:"email" json:"email"`
	EmployeeCode string `bson:"employee_code" json:"employee_code"`
	Role         string `bson:"role" json:"role"`
	FullName     string `bson:"full_name" json:"full_name"`

	// Request Information
	Method      string `bson:"method" json:"method"`             // GET, POST, PUT, DELETE
	Path        string `bson:"path" json:"path"`                 // /v1/user/create
	QueryParams string `bson:"query_params" json:"query_params"` // ?page=1&limit=10
	RequestBody string `bson:"request_body" json:"request_body"` // JSON body (sanitized)
	IPAddress   string `bson:"ip_address" json:"ip_address"`
	UserAgent   string `bson:"user_agent" json:"user_agent"`

	// Response Information
	StatusCode   int   `bson:"status_code" json:"status_code"`
	ResponseTime int64 `bson:"response_time_ms" json:"response_time_ms"` // milliseconds

	// Semantic Information
	Action     string `bson:"action" json:"action"`           // CREATE, READ, UPDATE, DELETE
	Resource   string `bson:"resource" json:"resource"`       // user, task, receipt, etc.
	ResourceID string `bson:"resource_id" json:"resource_id"` // ID of the affected resource

	// Timestamps
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
