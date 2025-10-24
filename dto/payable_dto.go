package dto

import "time"

// ---------- Request DTO ----------

type CreatePayableDTO struct {
	Supplier   string  `json:"supplier"`
	PurchaseNo string  `json:"purchase_no"`
	InvoiceNo  string  `json:"invoice_no"`
	IssueDate  string  `json:"issue_date"`
	DueDate    string  `json:"due_date"`
	Amount     float64 `json:"amount"`
	Balance    float64 `json:"balance"`
	PaymentRef string  `json:"payment_ref,omitempty"`
	Note       string  `json:"note,omitempty"`
}

type UpdatePayableDTO struct {
	Supplier   string  `json:"supplier,omitempty"`
	PurchaseNo string  `json:"purchase_no,omitempty"`
	InvoiceNo  string  `json:"invoice_no,omitempty"`
	IssueDate  string  `json:"issue_date,omitempty"`
	DueDate    string  `json:"due_date,omitempty"`
	Amount     float64 `json:"amount,omitempty"`
	Balance    float64 `json:"balance,omitempty"`
	Status     string  `json:"status,omitempty"`
	PaymentRef string  `json:"payment_ref,omitempty"`
	Note       string  `json:"note,omitempty"`
}

type RequestListPayable struct {
	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
	Search    string `query:"search"`
	SortBy    string `query:"sort_by"`
	SortOrder string `query:"sort_order"`
	Status    string `query:"status"`
	Supplier  string `query:"supplier"`
}

// ---------- Response DTO ----------

type PayableDTO struct {
	IDPayable  string    `json:"id_payable"`
	Supplier   string    `json:"supplier"`
	PurchaseNo string    `json:"purchase_no"`
	InvoiceNo  string    `json:"invoice_no"`
	IssueDate  time.Time `json:"issue_date"`
	DueDate    time.Time `json:"due_date"`
	Amount     float64   `json:"amount"`
	Balance    float64   `json:"balance"`
}

type PayableSummaryDTO struct {
	TotalAmount float64 `json:"total_amount"`
	TotalPaid   float64 `json:"total_paid"`
	TotalDue    float64 `json:"total_due"`
}
