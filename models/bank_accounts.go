package models

import "time"

const CollectionSBankAccounts = "bank_accounts"

type BankAccount struct {
	BankID      string     `bson:"bank_id" json:"bank_id"`
	BankName    string     `bson:"bank_name" json:"bank_name"`
	AccountNo   string     `bson:"account_no" json:"account_no"`
	AccountName string     `bson:"account_name" json:"account_name"`
	CreatedBy   string     `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `bson:"deleted_at" json:"deleted_at"`
}
