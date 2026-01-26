package dto

import "time"

// ---------- Request DTO ----------
type CreateIncomeDTO struct {
	Note                  *string `json:"note,omitempty"`
	BankID                string  `json:"bank_id" binding:"required"`
	TransactionCategoryID string  `json:"transaction_category_id" binding:"required"`
	Description           string  `json:"description" binding:"required"`
	Currency              string  `json:"currency" binding:"required"`
	TxnDate               string  `json:"txn_date" binding:"required"`
	PaymentMethod         string  `json:"payment_method,omitempty"`
	ReferenceNo           string  `json:"reference_no,omitempty"`
	Amount                float64 `json:"amount" binding:"required"`
}

type UpdateIncomeDTO struct {
	Note                  *string `json:"note,omitempty"`
	BankID                string  `json:"bank_id" binding:"required"`
	TransactionCategoryID string  `json:"transaction_category_id" binding:"required"`
	Description           string  `json:"description" binding:"required"`
	Currency              string  `json:"currency" binding:"required"`
	TxnDate               string  `json:"txn_date" binding:"required"`
	PaymentMethod         string  `json:"payment_method,omitempty"`
	ReferenceNo           string  `json:"reference_no,omitempty"`
	Amount                float64 `json:"amount" binding:"required"`
}

type RequestListIncome struct {
	Search                string `query:"search"`
	SortBy                string `query:"sort_by"`
	SortOrder             string `query:"sort_order"`
	TransactionCategoryID string `query:"transaction_category_id"`
	StartDate             string `query:"start_date"`
	EndDate               string `query:"end_date"`
	BankID                string `query:"bank_id"`
	Page                  int    `query:"page"`
	Limit                 int    `query:"limit"`
}

type RequestIncomeSummary struct {
	BankID    string `query:"bank_id"` // รหัสบัญชีธนาคารที่เกี่ยวข้อง
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
}

// ---------- Response DTO ----------
type IncomeDTO struct {
	TxnDate                   time.Time  `json:"txn_date"`
	CreatedAt                 time.Time  `json:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at"`
	Note                      *string    `json:"note,omitempty"`
	DeletedAt                 *time.Time `json:"deleted_at,omitempty"`
	IncomeID                  string     `json:"income_id"`
	BankID                    string     `json:"bank_id"`
	TransactionCategoryID     string     `json:"transaction_category_id"`
	TransactionCategoryNameTH string     `json:"transaction_category_name_th"`
	Description               string     `json:"description"`
	Currency                  string     `json:"currency"`
	PaymentMethod             string     `json:"payment_method,omitempty"`
	ReferenceNo               string     `json:"reference_no,omitempty"`
	CreatedBy                 string     `json:"created_by"`
	Amount                    float64    `json:"amount"`
}

type IncomeSummaryDTO struct {
	TotalToday     float64 `json:"total_today"`
	TotalThisMonth float64 `json:"total_this_month"`
	TotalAll       float64 `json:"total_all"`
}
