package dto

import "time"

// ---------- Request DTO ----------
type CreateExpenseDTO struct {
	BankID                string  `json:"bank_id" binding:"required"`
	TransactionCategoryID string  `json:"transaction_category_id" binding:"required"`
	Description           string  `json:"description" binding:"required"`
	Amount                float64 `json:"amount" binding:"required"`
	Currency              string  `json:"currency" binding:"required"`
	TxnDate               string  `json:"txn_date" binding:"required"`
	PaymentMethod         string  `json:"payment_method,omitempty"`
	ReferenceNo           string  `json:"reference_no,omitempty"`
	Note                  *string `json:"note,omitempty"`
}

type UpdateExpenseDTO struct {
	BankID                string  `json:"bank_id" binding:"required"`
	TransactionCategoryID string  `json:"transaction_category_id" binding:"required"`
	Description           string  `json:"description" binding:"required"`
	Amount                float64 `json:"amount" binding:"required"`
	Currency              string  `json:"currency" binding:"required"`
	TxnDate               string  `json:"txn_date" binding:"required"`
	PaymentMethod         string  `json:"payment_method,omitempty"`
	ReferenceNo           string  `json:"reference_no,omitempty"`
	Note                  *string `json:"note,omitempty"`
}

type RequestListExpense struct {
	Page                  int    `query:"page"`
	Limit                 int    `query:"limit"`
	Search                string `query:"search"`
	SortBy                string `query:"sort_by"`
	SortOrder             string `query:"sort_order"`
	TransactionCategoryID string `query:"transaction_category_id"`
	StartDate             string `query:"start_date"`
	EndDate               string `query:"end_date"`
}

type RequestExpenseSummary struct {
	BankID string `query:"bank_id"` // รหัสบัญชีธนาคารที่เกี่ยวข้อง
}

// ---------- Response DTO ----------
type ExpenseDTO struct {
	ExpenseID                 string     `json:"expense_id"`
	BankID                    string     `json:"bank_id"`
	TransactionCategoryID     string     `json:"transaction_category_id"`
	TransactionCategoryNameTH string     `json:"transaction_category_name_th"`
	Description               string     `json:"description"`
	Amount                    float64    `json:"amount"`
	Currency                  string     `json:"currency"`
	TxnDate                   time.Time  `json:"txn_date"`
	PaymentMethod             string     `json:"payment_method,omitempty"`
	ReferenceNo               string     `json:"reference_no,omitempty"`
	Note                      *string    `json:"note,omitempty"`
	CreatedBy                 string     `json:"created_by"`
	CreatedAt                 time.Time  `json:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at"`
	DeletedAt                 *time.Time `json:"deleted_at,omitempty"`
}

type ExpenseSummaryDTO struct {
	TotalToday     float64 `json:"total_today"`
	TotalThisMonth float64 `json:"total_this_month"`
	TotalAll       float64 `json:"total_all"`
}
