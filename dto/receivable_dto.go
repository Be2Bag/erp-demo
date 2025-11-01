package dto

import "time"

// ---------- Request DTO ----------

type CreateReceivableDTO struct {
	Customer  string  `json:"customer"`
	BankID    string  `json:"bank_id"`
	InvoiceNo string  `json:"invoice_no"`
	IssueDate string  `json:"issue_date"`
	DueDate   string  `json:"due_date"`
	Amount    float64 `json:"amount"`
	Balance   float64 `json:"balance"`
}

type UpdateReceivableDTO struct {
	Customer  string  `json:"customer,omitempty"`
	BankID    string  `json:"bank_id,omitempty"`
	InvoiceNo string  `json:"invoice_no,omitempty"`
	IssueDate string  `json:"issue_date,omitempty"`
	DueDate   string  `json:"due_date,omitempty"`
	Amount    float64 `json:"amount,omitempty"`
	Balance   float64 `json:"balance,omitempty"`
	Status    string  `json:"status,omitempty"`
	Note      string  `json:"note,omitempty"`
}

type RequestListReceivable struct {
	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
	Search    string `query:"search"`
	SortBy    string `query:"sort_by"`
	SortOrder string `query:"sort_order"`
	Status    string `query:"status"`
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
	BankID    string `query:"bank_id"`
}

// ---------- Response DTO ----------

type ReceivableDTO struct {
	IDReceivable string                  `json:"id_receivable"`
	BankID       string                  `json:"bank_id"`
	BankName     string                  `json:"bank_name"`
	Customer     string                  `json:"customer"`
	InvoiceNo    string                  `json:"invoice_no"`
	IssueDate    time.Time               `json:"issue_date"`
	DueDate      time.Time               `json:"due_date"`
	Amount       float64                 `json:"amount"`
	Balance      float64                 `json:"balance"`
	Status       string                  `json:"status"`
	CreatedBy    string                  `json:"created_by"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
	Transactions []PaymentTransactionDTO `json:"transactions"` // รายการชำระเงิน
	Note         string                  `json:"note"`
}

type ReceivableSummaryDTO struct {
	TotalAmount  float64 `json:"total_amount"`  // ยอดรวมทั้งหมด
	TotalDue     float64 `json:"total_due"`     // ยอดคงค้าง
	OverdueCount int     `json:"overdue_count"` // จำนวนรายการเกินกำหนด
}
