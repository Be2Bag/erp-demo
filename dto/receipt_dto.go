package dto // ประกาศแพ็กเกจ dto สำหรับชุดโครงสร้างข้อมูล (DTO)

import "time" // นำเข้าแพ็กเกจ time สำหรับชนิดข้อมูลเวลา

// ---------- Request DTOs ---------- // ส่วนของโครงสร้างสำหรับรับคำขอจากผู้ใช้ (Request)
type CreateReceiptDTO struct { // โครงสร้างข้อมูลสำหรับสร้างใบเสร็จ
	ReceiptDate   string           `json:"receipt_date,omitempty"`            // YYYY-MM-DD รูปแบบวันที่ปี-เดือน-วัน (ไม่ส่งมาก็ได้)
	Customer      CustomerInfoDTO  `json:"customer" binding:"required"`       // ข้อมูลลูกค้า (จำเป็นต้องส่ง)
	Issuer        IssuerInfoDTO    `json:"issuer,omitempty"`                  // ข้อมูลผู้ออกเอกสาร (ไม่ส่งมาก็ได้)
	Items         []ReceiptItemDTO `json:"items" binding:"required"`          // รายการสินค้า/บริการ (จำเป็นต้องส่ง)
	Discount      float64          `json:"discount"`                          // ส่วนลดรวม (บาท) - ไม่ส่งมาก็ได้
	Remark        string           `json:"remark,omitempty"`                  // หมายเหตุ (ไม่ส่งมาก็ได้)
	PaymentDetail PaymentInfoDTO   `json:"payment_detail" binding:"required"` // รายละเอียดการชำระเงิน (จำเป็นต้องส่ง)
	Status        string           `json:"status,omitempty"`                  // สถานะใบเสร็จ เช่น paid, pending, credit (ไม่ส่งมาก็ได้)
	BillType      string           `json:"bill_type" binding:"required"`      // ประเภทบิล: quotation, delivery_note, receipt (จำเป็นต้องส่ง)
	ApprovedBy    string           `json:"approved_by,omitempty"`             // ผู้อนุมัติ (ไม่ส่งมาก็ได้)
	ReceivedBy    string           `json:"received_by,omitempty"`             // ผู้รับเงิน/ผู้รับเอกสาร (ไม่ส่งมาก็ได้)
	TypeReceipt   string           `json:"type_receipt" binding:"required"`   // ประเภทใบเสร็จ "company" หรือ "shop" (จำเป็นต้องส่ง)
	TaxID         string           `json:"tax_id,omitempty"`                  // เลขประจำตัวผู้เสียภาษีอากร (ไม่ส่งมาก็ได้)
	ShopDetail    string           `json:"shop_detail,omitempty"`             // รายละเอียดร้านค้า (ถ้ามี) (ไม่ส่งมาก็ได้)
} // จบโครงสร้าง CreateReceiptDTO

type UpdateReceiptDTO struct { // โครงสร้างข้อมูลสำหรับอัปเดตใบเสร็จ (ส่งเฉพาะฟิลด์ที่ต้องการแก้)
	ReceiptDate   string           `json:"receipt_date,omitempty"`   // YYYY-MM-DD รูปแบบวันที่ปี-เดือน-วัน (ไม่ส่งมาก็ได้)
	Customer      *CustomerInfoDTO `json:"customer,omitempty"`       // ข้อมูลลูกค้า (อาจไม่ส่ง)
	Issuer        *IssuerInfoDTO   `json:"issuer,omitempty"`         // ข้อมูลผู้ออกเอกสาร (อาจไม่ส่ง)
	Items         []ReceiptItemDTO `json:"items,omitempty"`          // รายการสินค้า/บริการ (อาจไม่ส่ง)
	Discount      *float64         `json:"discount,omitempty"`       // ส่วนลดรวม (บาท) (อาจไม่ส่ง)
	Remark        *string          `json:"remark,omitempty"`         // หมายเหตุ (อาจไม่ส่ง)
	PaymentDetail *PaymentInfoDTO  `json:"payment_detail,omitempty"` // รายละเอียดการชำระเงิน (อาจไม่ส่ง)
	Status        *string          `json:"status,omitempty"`         // สถานะใบเสร็จ (อาจไม่ส่ง)
	ApprovedBy    *string          `json:"approved_by,omitempty"`    // ผู้อนุมัติ (อาจไม่ส่ง)
	ReceivedBy    *string          `json:"received_by,omitempty"`    // ผู้รับเงิน/ผู้รับเอกสาร (อาจไม่ส่ง)
} // จบโครงสร้าง UpdateReceiptDTO

type CustomerInfoDTO struct { // โครงสร้างข้อมูลลูกค้า
	Name                string `json:"name" binding:"required"`                  // ชื่อลูกค้า (จำเป็น)
	Address             string `json:"address" binding:"required"`               // ที่อยู่ลูกค้า (จำเป็น)
	Contact             string `json:"contact" binding:"required"`               // ช่องทางติดต่อ (จำเป็น)
	TaxIDCustomer       string `json:"tax_id_customer" binding:"required"`       // เลขประจำตัวผู้เสียภาษีอากรลูกค้า (จำเป็น)
	TypeReceiptCustomer string `json:"type_receipt_customer" binding:"required"` // ประเภทใบเสร็จลูกค้า "company" หรือ "shop" (จำเป็น)
	ShopDetailCustomer  string `json:"shop_detail_customer,omitempty"`           // รายละเอียดร้านค้าลูกค้า (ถ้ามี)
	Fax                 string `json:"fax,omitempty"`                            // เบอร์แฟกซ์
} // จบโครงสร้าง CustomerInfoDTO

type IssuerInfoDTO struct { // โครงสร้างข้อมูลผู้ออกเอกสาร
	Name       string `json:"name,omitempty"`        // ชื่อผู้ออกเอกสาร (ไม่ส่งมาก็ได้)
	Address    string `json:"address,omitempty"`     // ที่อยู่ผู้ออกเอกสาร (ไม่ส่งมาก็ได้)
	Contact    string `json:"contact,omitempty"`     // ช่องทางติดต่อ (ไม่ส่งมาก็ได้)
	Email      string `json:"email,omitempty"`       // อีเมล (ไม่ส่งมาก็ได้)
	PreparedBy string `json:"prepared_by,omitempty"` // ผู้จัดเตรียมเอกสาร (ไม่ส่งมาก็ได้)
} // จบโครงสร้าง IssuerInfoDTO

type ReceiptItemDTO struct { // โครงสร้างรายการในใบเสร็จ
	Description string  `json:"description" binding:"required"`      // รายละเอียดรายการ (จำเป็น)
	Quantity    int     `json:"quantity" binding:"required,gt=0"`    // จำนวน (ต้องมากกว่า 0)
	UnitPrice   float64 `json:"unit_price" binding:"required,gte=0"` // ราคาต่อหน่วย (ต้องมากกว่าหรือเท่ากับ 0)
	Total       float64 `json:"total,omitempty"`                     // ถ้าเป็น 0 ระบบจะคำนวณให้: จำนวน x ราคาต่อหน่วย + ค่าอื่นๆ
} // จบโครงสร้าง ReceiptItemDTO

type PaymentInfoDTO struct { // โครงสร้างรายละเอียดการชำระเงิน (Request)
	Method        string  `json:"method" binding:"required"`            // วิธีชำระเงิน (จำเป็น)
	BankName      string  `json:"bank_name,omitempty"`                  // ชื่อธนาคาร (ไม่ส่งมาก็ได้)
	AccountName   string  `json:"account_name,omitempty"`               // ชื่อบัญชี (ไม่ส่งมาก็ได้)
	AccountNumber string  `json:"account_number,omitempty"`             // เลขที่บัญชี (ไม่ส่งมาก็ได้)
	AmountPaid    float64 `json:"amount_paid" binding:"required,gte=0"` // จำนวนเงินที่ชำระ (ต้องมากกว่าหรือเท่ากับ 0 และจำเป็น)
	PaidDate      string  `json:"paid_date,omitempty"`                  // YYYY-MM-DD วันที่ชำระเงิน (ไม่ส่งมาก็ได้)
	Note          string  `json:"note,omitempty"`                       // หมายเหตุการชำระเงิน (ไม่ส่งมาก็ได้)
} // จบโครงสร้าง PaymentInfoDTO

// ---------- Query DTO ---------- // ส่วนของโครงสร้างสำหรับพารามิเตอร์การค้นหา/แบ่งหน้า
type RequestListReceipt struct { // โครงสร้างคำขอรายการใบเสร็จ
	Page        int    `query:"page"`         // หน้าที่ต้องการ
	Limit       int    `query:"limit"`        // จำนวนรายการต่อหน้า
	Search      string `query:"search"`       // คำค้นหา (เช่น เลขที่ใบเสร็จ/ชื่อลูกค้า)
	SortBy      string `query:"sort_by"`      // ฟิลด์ที่ใช้เรียงลำดับ
	SortOrder   string `query:"sort_order"`   // asc | desc
	Status      string `query:"status"`       // สถานะใบเสร็จสำหรับกรอง
	StartDate   string `query:"start_date"`   // YYYY-MM-DD วันที่เริ่มต้นช่วงค้นหา
	EndDate     string `query:"end_date"`     // YYYY-MM-DD วันที่สิ้นสุดช่วงค้นหา
	BillType    string `query:"bill_type"`    // ประเภทบิล: quotation, delivery_note, receipt
	TypeReceipt string `query:"type_receipt"` // ประเภทใบเสร็จ "company" หรือ "shop"
} // จบโครงสร้าง RequestListReceipt

// ---------- Response DTOs ---------- // ส่วนของโครงสร้างสำหรับส่งข้อมูลตอบกลับ (Response)
type ReceiptDTO struct { // โครงสร้างข้อมูลใบเสร็จที่ส่งกลับ
	IDReceipt     string             `json:"id_receipt"`            // รหัสใบเสร็จ (ไอดีภายในระบบ)
	ReceiptNumber string             `json:"receipt_number"`        // เลขที่ใบเสร็จ
	ReceiptDate   time.Time          `json:"receipt_date"`          // วันที่ใบเสร็จ (ชนิดเวลา)
	Customer      CustomerInfoDTO    `json:"customer"`              // ข้อมูลลูกค้า
	Issuer        IssuerInfoDTO      `json:"issuer"`                // ข้อมูลผู้ออกเอกสาร
	Items         []ReceiptItemDTO   `json:"items"`                 // รายการสินค้า/บริการ
	SubTotal      float64            `json:"sub_total"`             // ยอดรวมก่อน VAT
	Discount      float64            `json:"discount"`              // ส่วนลดรวม (บาท)
	TotalVAT      float64            `json:"total_vat"`             // ค่าภาษีมูลค่าเพิ่ม VAT 7%
	TotalAmount   float64            `json:"total_amount"`          // ยอดรวมสุทธิรวม VAT แล้ว
	Remark        string             `json:"remark,omitempty"`      // หมายเหตุ (อาจว่าง)
	PaymentDetail PaymentInfoRespDTO `json:"payment_detail"`        // รายละเอียดการชำระเงิน (สำหรับตอบกลับ)
	Status        string             `json:"status"`                // สถานะใบเสร็จ
	BillType      string             `json:"bill_type"`             // ประเภทบิล: quotation, delivery_note, receipt
	TypeReceipt   string             `json:"type_receipt"`          // ประเภทใบเสร็จ "company" หรือ "shop"
	ApprovedBy    string             `json:"approved_by,omitempty"` // ผู้อนุมัติ (อาจว่าง)
	ReceivedBy    string             `json:"received_by,omitempty"` // ผู้รับเงิน/ผู้รับเอกสาร (อาจว่าง)
	CreatedAt     time.Time          `json:"created_at"`            // วันที่สร้างข้อมูล
	UpdatedAt     time.Time          `json:"updated_at"`            // วันที่แก้ไขล่าสุด
	TaxID         string             `json:"tax_id"`                // เลขประจำตัวผู้เสียภาษีอากร
	ShopDetail    string             `json:"shop_detail,omitempty"` // รายละเอียดร้านค้า (อาจว่าง)
} // จบโครงสร้าง ReceiptDTO

type PaymentInfoRespDTO struct { // โครงสร้างรายละเอียดการชำระเงิน (Response)
	Method        string     `json:"method"`                   // วิธีชำระเงิน
	BankName      string     `json:"bank_name,omitempty"`      // ชื่อธนาคาร (อาจว่าง)
	AccountName   string     `json:"account_name,omitempty"`   // ชื่อบัญชี (อาจว่าง)
	AccountNumber string     `json:"account_number,omitempty"` // เลขที่บัญชี (อาจว่าง)
	AmountPaid    float64    `json:"amount_paid"`              // จำนวนเงินที่ชำระ
	PaidDate      *time.Time `json:"paid_date,omitempty"`      // วันที่ชำระเงิน (nil = ยังไม่ได้ชำระ)
	Note          string     `json:"note,omitempty"`           // หมายเหตุการชำระเงิน (อาจว่าง)
} // จบโครงสร้าง PaymentInfoRespDTO

type ReceiptSummaryDTO struct { // โครงสร้างสรุปข้อมูลใบเสร็จ
	TotalAmount  float64 `json:"total_amount"`  // ยอดรวมทั้งหมดของใบเสร็จที่ค้นหา
	TotalPaid    float64 `json:"total_paid"`    // ยอดรวมที่ชำระแล้วของใบเสร็จที่ค้นหา
	PendingCount int     `json:"pending_count"` // จำนวนใบเสร็จที่ยังค้างชำระ
} // จบโครงสร้าง ReceiptSummaryDTO

type RequestSummaryReceipt struct {
	Report      string `query:"report"`       // รายงานประเภท day | month | all
	TypeReceipt string `query:"type_receipt"` // ประเภทใบเสร็จ "company" หรือ "shop"
	StartDate   string `query:"start_date"`
	EndDate     string `query:"end_date"`
}

type RequestCopyReceipt struct {
	IDReceipt string `json:"id_receipt" binding:"required"`
	BillType  string `json:"bill_type" binding:"required"`
}
