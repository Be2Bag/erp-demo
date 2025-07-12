package models

import (
	"time"
)

// User ใช้เก็บข้อมูลพนักงานในระบบ HR พร้อมเชื่อมโยงตำแหน่งงาน (PositionID) และแผนก (DepartmentID) ผ่าน FK
type User struct {
	UserID            string              `bson:"user_id" json:"user_id"`                       // รหัสประจำตัวผู้ใช้ (ไม่ซ้ำกัน)
	Username          string              `bson:"username" json:"username"`                     // ชื่อผู้ใช้สำหรับเข้าสู่ระบบ
	Email             string              `bson:"email" json:"email"`                           // อีเมลของผู้ใช้
	Password          string              `bson:"password" json:"password"`                     // รหัสผ่าน (ควรเข้ารหัสก่อนจัดเก็บ)
	TitleTH           string              `bson:"title_th" json:"title_th"`                     // คำนำหน้าชื่อ (ภาษาไทย)
	TitleEN           string              `bson:"title_en" json:"title_en"`                     // คำนำหน้าชื่อ (ภาษาอังกฤษ)
	FirstNameTH       string              `bson:"first_name_th" json:"first_name_th"`           // ชื่อจริงของพนักงาน
	LastNameTH        string              `bson:"last_name_th" json:"last_name_th"`             // นามสกุลของพนักงาน
	FirstNameEN       string              `bson:"first_name_en" json:"first_name_en"`           // ชื่อจริงของพนักงาน (ภาษาอังกฤษ)
	LastNameEN        string              `bson:"last_name_en" json:"last_name_en"`             // นามสกุลของพนักงาน (ภาษาอังกฤษ)
	IDCard            string              `bson:"id_card" json:"id_card"`                       // หมายเลขบัตรประชาชน (อาจเข้ารหัสก่อนจัดเก็บ)
	Role              string              `bson:"role" json:"role"`                             // บทบาทหรือสิทธิ์ของผู้ใช้ในระบบ (เช่น admin, user)
	Avatar            string              `bson:"avatar" json:"avatar"`                         // ลิงก์หรือที่อยู่รูปประจำตัวผู้ใช้
	Phone             string              `bson:"phone" json:"phone"`                           // เบอร์โทรศัพท์ของพนักงาน
	Status            string              `bson:"status" json:"status"`                         // สถานะของผู้ใช้ (เช่น active, inactive)
	EmployeeCode      string              `bson:"employee_code" json:"employee_code"`           // รหัสพนักงาน (อาจใช้สำหรับอ้างอิงภายใน)
	Gender            string              `bson:"gender" json:"gender"`                         // เพศของพนักงาน
	BirthDate         time.Time           `bson:"birth_date" json:"birth_date"`                 // วันเดือนปีเกิดของพนักงาน (รูปแบบ string)
	PositionID        string              `bson:"position_id" json:"position_id"`               // รหัสตำแหน่งงาน (FK ไปยัง Positions)
	DepartmentID      string              `bson:"department_id" json:"department_id"`           // รหัสแผนก (FK ไปยัง Departments)
	HireDate          time.Time           `bson:"hire_date" json:"hire_date"`                   // วันที่เริ่มงาน
	EmploymentType    string              `bson:"employment_type" json:"employment_type"`       // ประเภทการจ้างงาน (เช่น full-time, part-time)
	EmploymentHistory []EmploymentHistory `bson:"employment_history" json:"employment_history"` // ประวัติการจ้างงาน (อาจมีหลายรายการ)
	Address           Address             `bson:"address" json:"address"`                       // ที่อยู่ของพนักงาน
	BankInfo          BankInfo            `bson:"bank_info" json:"bank_info"`                   // ข้อมูลบัญชีธนาคารของพนักงาน
	Documents         []Document          `bson:"documents" json:"documents"`
	CreatedAt         time.Time           `bson:"created_at" json:"created_at"` // วันที่สร้างข้อมูลนี้
	UpdatedAt         time.Time           `bson:"updated_at" json:"updated_at"` // วันที่แก้ไขข้อมูลล่าสุด
	DeletedAt         *time.Time          `bson:"deleted_at" json:"deleted_at"` // วันที่ลบข้อมูล (soft delete)
}

type Address struct {
	AddressLine1 string `bson:"address_line1" json:"address_line1"`                     // ที่อยู่บรรทัด 1
	AddressLine2 string `bson:"address_line2,omitempty" json:"address_line2,omitempty"` // ที่อยู่บรรทัด 2
	Subdistrict  string `bson:"subdistrict" json:"subdistrict"`                         // ตำบล
	District     string `bson:"district" json:"district"`                               // อำเภอ
	Province     string `bson:"province" json:"province"`                               // จังหวัด
	PostalCode   string `bson:"postal_code" json:"postal_code"`                         // รหัสไปรษณีย์
	Country      string `bson:"country" json:"country"`                                 // ประเทศ
}

type BankInfo struct {
	BankName    string `bson:"bank_name"`    // ชื่อธนาคาร
	AccountNo   string `bson:"account_no"`   // เลขที่บัญชี
	AccountName string `bson:"account_name"` // ชื่อบัญชี
}

type Document struct {
	Name       string     `bson:"name"`                 // ชื่อเอกสาร
	FileURL    string     `bson:"file_url"`             // ลิงก์ไฟล์เอกสาร
	Type       string     `bson:"type"`                 // ประเภทเอกสาร เช่น "id_card", "degree"
	CreatedAt  time.Time  `bson:"created_at"`           // วันที่สร้างเอกสาร
	UploadedAt time.Time  `bson:"uploaded_at"`          // วันที่อัปโหลดเอกสาร
	DeletedAt  *time.Time `bson:"deleted_at,omitempty"` // วันที่ลบเอกสาร (soft delete)
}
