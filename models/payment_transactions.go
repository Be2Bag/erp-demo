package models

import "time"

const CollectionPaymentTransaction = "payment_transactions"

type PaymentTransaction struct {
	PaymentDate     time.Time  `json:"payment_date" bson:"payment_date"`                 // วันที่รับ/จ่ายเงิน
	CreatedAt       time.Time  `json:"created_at" bson:"created_at"`                     // วันที่บันทึก
	UpdatedAt       time.Time  `json:"updated_at" bson:"updated_at"`                     // วันที่แก้ไขล่าสุด
	DeletedAt       *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"` // วันที่ลบข้อมูล (ถ้ามี)
	IDTransaction   string     `json:"id_transaction" bson:"id_transaction"`             // รหัสเอกสาร (เช่น PAY-2024-001-001)
	BankID          string     `json:"bank_id" bson:"bank_id"`                           // รหัสบัญชีธนาคารที่เกี่ยวข้อง
	RefInvoiceNo    string     `json:"ref_invoice_no" bson:"ref_invoice_no"`             // อ้างอิงใบแจ้งหนี้ (InvoiceNo)
	TransactionType string     `json:"transaction_type" bson:"transaction_type"`         // ประเภท: receivable (ลูกหนี้) หรือ payable (เจ้าหนี้)
	PaymentMethod   string     `json:"payment_method" bson:"payment_method"`             // วิธีการชำระเงิน (cash, transfer, cheque)
	PaymentRef      string     `json:"payment_ref" bson:"payment_ref"`                   // หมายเลขอ้างอิง (เลขสลิป, เช็ค, ใบเสร็จ)
	Note            string     `json:"note" bson:"note"`                                 // หมายเหตุเพิ่มเติม
	CreatedBy       string     `json:"created_by" bson:"created_by"`                     // ผู้บันทึก
	Amount          float64    `json:"amount" bson:"amount"`                             // จำนวนเงินที่จ่าย/รับในครั้งนี้
}
