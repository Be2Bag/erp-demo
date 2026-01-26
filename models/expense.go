package models

import "time"

const CollectionExpense = "expenses"

type Expense struct {
	TxnDate               time.Time  `bson:"txn_date" json:"txn_date"` // วันที่เกิดรายการ
	CreatedAt             time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt             time.Time  `bson:"updated_at" json:"updated_at"`
	Note                  *string    `bson:"note,omitempty" json:"note,omitempty"`
	DeletedAt             *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
	ExpenseID             string     `bson:"expense_id" json:"expense_id"`
	BankID                string     `bson:"bank_id" json:"bank_id"`
	TransactionCategoryID string     `bson:"transaction_category_id"`                                  // ผูกกับ transaction_categories
	Description           string     `bson:"description" json:"description"`                           // รายละเอียด
	Currency              string     `bson:"currency" json:"currency"`                                 // เช่น "THB"
	PaymentMethod         string     `bson:"payment_method,omitempty" json:"payment_method,omitempty"` //เช่น "cash", "transfer", "credit_card"
	ReferenceNo           string     `bson:"reference_no,omitempty" json:"reference_no,omitempty"`     // เช่น เลขใบเสร็จ / หมายเลขธุรกรรมธนาคาร
	CreatedBy             string     `bson:"created_by" json:"created_by"`
	Amount                float64    `bson:"amount" json:"amount"` // จำนวนเงิน
}
