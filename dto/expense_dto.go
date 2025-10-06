package dto

import "time"

// ---------- Request DTO ----------
type CreateExpenseDTO struct {
	CategoryID    string  `json:"category_id" binding:"required"`
	Description   string  `json:"description" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	Currency      string  `json:"currency" binding:"required"`
	TxnDate       string  `json:"txn_date" binding:"required"`
	PaymentMethod string  `json:"payment_method,omitempty"`
	ReferenceNo   string  `json:"reference_no,omitempty"`
	Note          *string `json:"note,omitempty"`
}

type UpdateExpenseDTO struct {
	CategoryID    string  `json:"category_id" binding:"required"`
	Description   string  `json:"description" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	Currency      string  `json:"currency" binding:"required"`
	TxnDate       string  `json:"txn_date" binding:"required"`
	PaymentMethod string  `json:"payment_method,omitempty"`
	ReferenceNo   string  `json:"reference_no,omitempty"`
	Note          *string `json:"note,omitempty"`
}

type RequestListExpense struct {
	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
	Search    string `query:"search"`
	SortBy    string `query:"sort_by"`
	SortOrder string `query:"sort_order"`
}

// ---------- Response DTO ----------
type ExpenseDTO struct {
	ExpenseID     string     `json:"expense_id"`
	CategoryID    string     `json:"category_id"`
	Description   string     `json:"description"`
	Amount        float64    `json:"amount"`
	Currency      string     `json:"currency"`
	TxnDate       time.Time  `json:"txn_date"`
	PaymentMethod string     `json:"payment_method,omitempty"`
	ReferenceNo   string     `json:"reference_no,omitempty"`
	Note          *string    `json:"note,omitempty"`
	CreatedBy     string     `json:"created_by"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}
