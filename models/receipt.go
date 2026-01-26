package models // แพ็กเกจ models

import "time" // นำเข้าแพ็กเกจ time สำหรับการจัดการเวลา

const CollectionReceipts = "receipts" // ชื่อคอลเลกชันใบเสร็จในฐานข้อมูล

type Receipt struct { // โครงสร้างข้อมูลใบเสร็จ
	ReceiptDate   time.Time     `json:"receipt_date" bson:"receipt_date"`                         // วันที่ออกใบเสร็จ
	CreatedAt     time.Time     `json:"created_at" bson:"created_at"`                             // วันที่สร้างข้อมูล
	UpdatedAt     time.Time     `json:"updated_at" bson:"updated_at"`                             // วันที่อัปเดตข้อมูลล่าสุด
	Customer      CustomerInfo  `json:"customer" bson:"customer"`                                 // ข้อมูลลูกค้า
	PaymentDetail PaymentInfo   `json:"payment_detail" bson:"payment_detail"`                     // ข้อมูลการชำระเงิน
	Issuer        IssuerInfo    `json:"issuer" bson:"issuer"`                                     // ข้อมูลผู้ออกใบเสร็จ
	IDReceipt     string        `json:"id_receipt" bson:"id_receipt"`                             // รหัสใบเสร็จ
	ReceiptNumber string        `json:"receipt_number" bson:"receipt_number"`                     // เลขที่ใบเสร็จ
	Remark        string        `json:"remark,omitempty" bson:"remark,omitempty"`                 // หมายเหตุ
	Status        string        `json:"status" bson:"status"`                                     // สถานะใบเสร็จ เช่น paid, pending, credit
	BillType      string        `json:"bill_type" bson:"bill_type"`                               // ประเภทบิล: quotation, delivery_note, receipt
	TypeReceipt   string        `json:"type_receipt" bson:"type_receipt"`                         // ประเภทใบเสร็จ "company" หรือ "shop"
	ApprovedBy    string        `json:"approved_by,omitempty" bson:"approved_by"`                 // ผู้อนุมัติ
	ReceivedBy    string        `json:"received_by,omitempty" bson:"received_by"`                 // ผู้รับเงิน/ผู้รับเอกสาร
	TaxID         string        `json:"tax_id" bson:"tax_id"`                                     // เลขประจำตัวผู้เสียภาษีอากร
	ShopDetail    string        `json:"shop_detail,omitempty" bson:"shop_detail,omitempty"`       // รายละเอียดร้านค้า (ถ้ามี)
	CopiedFromID  string        `json:"copied_from_id,omitempty" bson:"copied_from_id,omitempty"` // ID ของ receipt ต้นฉบับที่ copy มา (ถ้ามี)
	Items         []ReceiptItem `json:"items" bson:"items"`                                       // รายการสินค้า/บริการ
	SubTotal      float64       `json:"sub_total" bson:"sub_total"`                               // ยอดรวมก่อน VAT (บาท)
	Discount      float64       `json:"discount" bson:"discount"`                                 // ส่วนลดรวม (บาท)
	TotalVAT      float64       `json:"total_vat" bson:"total_vat"`                               // ค่าภาษีมูลค่าเพิ่ม VAT 7% (บาท)
	TotalAmount   float64       `json:"total_amount" bson:"total_amount"`                         // ยอดรวมสุทธิรวม VAT แล้ว (บาท)
} // ปิดโครงสร้าง Receipt

type CustomerInfo struct { // โครงสร้างข้อมูลลูกค้า
	Name                string `json:"name" bson:"name"`                                                     // ชื่อลูกค้า
	Address             string `json:"address" bson:"address"`                                               // ที่อยู่
	Contact             string `json:"contact" bson:"contact"`                                               // เบอร์ติดต่อ
	TaxIDCustomer       string `json:"tax_id_customer" bson:"tax_id_customer"`                               // เลขประจำตัวผู้เสียภาษีอากรลูกค้า
	TypeReceiptCustomer string `json:"type_receipt_customer" bson:"type_receipt_customer"`                   // ประเภทใบเสร็จลูกค้า "company" หรือ "shop"
	ShopDetailCustomer  string `json:"shop_detail_customer,omitempty" bson:"shop_detail_customer,omitempty"` // รายละเอียดร้านค้าลูกค้า (ถ้ามี)
	Fax                 string `json:"fax,omitempty" bson:"fax,omitempty"`                                   // เบอร์แฟกซ์
} // ปิดโครงสร้าง CustomerInfo

type IssuerInfo struct { // โครงสร้างข้อมูลผู้ออกใบเสร็จ/บริษัท
	Name       string `json:"name" bson:"name"`                       // ชื่อผู้ออกใบเสร็จ
	Address    string `json:"address" bson:"address"`                 // ที่อยู่บริษัท
	Contact    string `json:"contact" bson:"contact"`                 // เบอร์โทร
	Email      string `json:"email,omitempty" bson:"email,omitempty"` // อีเมล (ถ้ามี)
	PreparedBy string `json:"prepared_by" bson:"prepared_by"`         // ผู้จัดทำ
} // ปิดโครงสร้าง IssuerInfo

type ReceiptItem struct { // โครงสร้างข้อมูลรายการสินค้า/บริการ
	Description string  `json:"description" bson:"description"` // รายละเอียดรายการ
	Quantity    int     `json:"quantity" bson:"quantity"`       // จำนวน
	UnitPrice   float64 `json:"unit_price" bson:"unit_price"`   // ราคาต่อหน่วย
	Total       float64 `json:"total" bson:"total"`             // รวมเป็นเงิน
} // ปิดโครงสร้าง ReceiptItem

type PaymentInfo struct { // โครงสร้างข้อมูลการชำระเงิน
	PaidDate      *time.Time `json:"paid_date,omitempty" bson:"paid_date,omitempty"`           // วันที่ชำระเงิน (nil = ยังไม่ได้ชำระ)
	Method        string     `json:"method" bson:"method"`                                     // วิธีชำระ (เงินสด, โอน, อื่นๆ)
	BankName      string     `json:"bank_name,omitempty" bson:"bank_name,omitempty"`           // ชื่อธนาคาร (ถ้ามี)
	AccountName   string     `json:"account_name,omitempty" bson:"account_name,omitempty"`     // ชื่อบัญชี (ถ้ามี)
	AccountNumber string     `json:"account_number,omitempty" bson:"account_number,omitempty"` // เลขที่บัญชี (ถ้ามี)
	Note          string     `json:"note,omitempty" bson:"note,omitempty"`                     // รายละเอียดเพิ่มเติม
	AmountPaid    float64    `json:"amount_paid" bson:"amount_paid"`                           // จำนวนเงินที่ชำระ
} // ปิดโครงสร้าง PaymentInfo
