package dto

import "time"

// ---------- Request DTO ----------

type CreateReceivableDTO struct {
	Customer  string  `json:"customer"`
	InvoiceNo string  `json:"invoice_no"`
	IssueDate string  `json:"issue_date"`
	DueDate   string  `json:"due_date"`
	Amount    float64 `json:"amount"`
	Balance   float64 `json:"balance"`
}

type UpdateReceivableDTO struct {
	Customer  string  `json:"customer,omitempty"`
	InvoiceNo string  `json:"invoice_no,omitempty"`
	IssueDate string  `json:"issue_date,omitempty"`
	DueDate   string  `json:"due_date,omitempty"`
	Amount    float64 `json:"amount,omitempty"`
	Balance   float64 `json:"balance,omitempty"`
	Status    string  `json:"status,omitempty"`
}

type RequestListReceivable struct {
	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
	Search    string `query:"search"`
	SortBy    string `query:"sort_by"`
	SortOrder string `query:"sort_order"`
	Status    string `query:"status"`
}

// ---------- Response DTO ----------

type ReceivableDTO struct {
	IDReceivable string     `json:"id_receivable"`
	Customer     string     `json:"customer"`
	InvoiceNo    string     `json:"invoice_no"`
	IssueDate    time.Time  `json:"issue_date"`
	DueDate      time.Time  `json:"due_date"`
	Amount       float64    `json:"amount"`
	Balance      float64    `json:"balance"`
	Status       string     `json:"status"`
	CreatedBy    string     `json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

type ReceivableSummaryDTO struct {
	TotalAmount float64 `json:"total_amount"`
	TotalPaid   float64 `json:"total_paid"`
	TotalDue    float64 `json:"total_due"`
}
