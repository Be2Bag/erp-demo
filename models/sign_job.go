package models

import (
	"time"
)

const CollectionSignJobs = "sign_jobs"

type SignJob struct {
	DueDate time.Time `bson:"due_date" json:"due_date"` // ใช้ pointer หากอาจว่าง

	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at" json:"deleted_at"` // สำหรับ soft delete
	// ---------- Keys ----------
	JobID string `bson:"job_id" json:"job_id"` // UUID/รหัสงาน (unique)

	// ---------- Customer ----------
	CompanyName    string `bson:"company_name" json:"company_name"`         // ชื่อบริษัท
	ContactPerson  string `bson:"contact_person" json:"contact_person"`     // ชื่อผู้ติดต่อ
	Phone          string `bson:"phone" json:"phone"`                       // เบอร์โทร
	Email          string `bson:"email" json:"email"`                       // อีเมล
	CustomerTypeID string `bson:"customer_type_id" json:"customer_type_id"` // ประเภทลูกค้า
	Address        string `bson:"address" json:"address"`                   // ที่อยู่ติดตั้ง / จัดส่ง

	// ---------- Sign detail ----------
	ProjectID   string `bson:"project_id" json:"project_id"`
	ProjectName string `bson:"project_name" json:"project_name"` // ชื่อโปรเจกต์
	JobName     string `bson:"job_name" json:"job_name"`         // ชื่องาน
	SignTypeID  string `bson:"sign_type_id" json:"sign_type_id"` // ประเภทป้าย
	Content     string `bson:"content" json:"content"`           // รายละเอียด
	MainColor   string `bson:"main_color" json:"main_color"`     // สีหลัก

	// ---------- Payment ----------
	PaymentMethod string `bson:"payment_method" json:"payment_method"` // deposit|cash|transfer|credit

	// ---------- Production / Timeline ----------
	ProductionTime string `bson:"production_time" json:"production_time"` // เช่น "5 วัน"

	// ---------- Design / Install ----------
	DesignOption  string `bson:"design_option" json:"design_option"`   // have_design|need_design
	InstallOption string `bson:"install_option" json:"install_option"` // none|self|shop

	// ---------- Notes ----------
	Notes string `bson:"notes" json:"notes"` // หมายเหตุ
	// ---------- Meta ----------
	Status            string  `bson:"status" json:"status"`                         // อยู่ในขั้นตอนไหนแล้ว
	CreatedBy         string  `bson:"created_by" json:"created_by"`                 // ใครสร้างงานนี้
	Width             float64 `bson:"width" json:"width"`                           // ซม.
	Height            float64 `bson:"height" json:"height"`                         // ซม.
	Quantity          int     `bson:"quantity" json:"quantity"`                     // จำนวน
	PriceTHB          float64 `bson:"price_thb" json:"price_thb"`                   // ราคาเป็นสตางค์หรือบาททั้งจำนวน เลือกแนวทางเดียวให้คงที่
	DepositAmount     float64 `bson:"deposit_amount" json:"deposit_amount"`         // เงินมัดจำ
	OutstandingAmount float64 `bson:"outstanding_amount" json:"outstanding_amount"` // ยอดค้างชำระ
	IsDeposit         bool    `bson:"is_deposit" json:"is_deposit"`                 // ใช้เงินมัดจำหรือไม่

	WaitPrice   bool `bson:"waitprice" json:"waitprice"`     // รอราคา
	WaitConfirm bool `bson:"waitconfirm" json:"waitconfirm"` // รอยืนยัน
}
