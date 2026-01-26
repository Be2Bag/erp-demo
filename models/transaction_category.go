package models

import "time"

const CollectionTransactionCategory = "transaction_categories"

type TransactionCategory struct {
	CreatedAt                 time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt                 time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt                 *time.Time `bson:"deleted_at" json:"deleted_at"`
	Note                      *string    `bson:"note" json:"note"`
	TransactionCategoryID     string     `bson:"transaction_category_id" json:"transaction_category_id"`
	Type                      string     `bson:"type" json:"type"` // ประเภทหมวดหมู่ (เช่น รายรับ, รายจ่าย)
	TransactionCategoryNameTH string     `bson:"transaction_category_name_th" json:"transaction_category_name_th"`
	Description               string     `bson:"description" json:"description"`

	CreatedBy string `bson:"created_by" json:"created_by"`
}
