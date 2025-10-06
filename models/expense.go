package models

import "time"

const CollectionExpense = "expenses"

type Expense struct {
	ExpenseID     string     `bson:"expense_id" json:"expense_id"`
	CategoryID    string     `bson:"category_id" json:"category_id"`                           // ผูกกับ transaction_categories
	Description   string     `bson:"description" json:"description"`                           // รายละเอียด
	Amount        float64    `bson:"amount" json:"amount"`                                     // จำนวนเงิน
	Currency      string     `bson:"currency" json:"currency"`                                 // เช่น "THB"
	TxnDate       time.Time  `bson:"txn_date" json:"txn_date"`                                 // วันที่เกิดรายการ
	PaymentMethod string     `bson:"payment_method,omitempty" json:"payment_method,omitempty"` //เช่น "cash", "transfer", "credit_card"
	ReferenceNo   string     `bson:"reference_no,omitempty" json:"reference_no,omitempty"`     // เช่น เลขใบเสร็จ / หมายเลขธุรกรรมธนาคาร
	Note          *string    `bson:"note,omitempty" json:"note,omitempty"`
	CreatedBy     string     `bson:"created_by" json:"created_by"`
	CreatedAt     time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt     *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}
