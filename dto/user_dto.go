package dto

import "time"

// Request

type RequestCreateUser struct {
	Email          string    `json:"email" validate:"required"`         // อีเมลของผู้ใช้
	Password       string    `json:"password" validate:"required"`      // รหัสผ่าน
	TitleTH        string    `json:"title_th" validate:"required"`      // คำนำหน้าชื่อ (ภาษาไทย)
	TitleEN        string    `json:"title_en"`                          // คำนำหน้าชื่อ (ภาษาอังกฤษ)
	FirstNameTH    string    `json:"first_name_th" validate:"required"` // ชื่อจริงของพนักงาน (ภาษาไทย)
	LastNameTH     string    `json:"last_name_th" validate:"required"`  // นามสกุลของพนักงาน (ภาษาไทย)
	FirstNameEN    string    `json:"first_name_en"`                     // ชื่อจริงของพนักงาน (ภาษาอังกฤษ)
	LastNameEN     string    `json:"last_name_en"`                      // นามสกุลของพนักงาน (ภาษาอังกฤษ)
	IDCard         string    `json:"id_card" validate:"required"`       // หมายเลขบัตรประชาชน (อาจเข้ารหัสก่อนจัดเก็บ)
	Avatar         string    `json:"avatar"`                            // ลิงก์หรือที่อยู่รูปประจำตัวผู้ใช้
	Phone          string    `json:"phone"`                             // เบอร์โทรศัพท์ของพนักงาน
	EmployeeCode   string    `json:"employee_code"`                     // รหัสพนักงาน (อาจใช้สำหรับอ้างอิงภายใน)
	Gender         string    `json:"gender" validate:"required"`        // เพศของพนักงาน
	BirthDate      time.Time `json:"birth_date"`                        // วันเดือนปีเกิดของพนักงาน (รูปแบบ string)
	PositionID     string    `json:"position_id"`                       // รหัสตำแหน่งงาน (FK ไปยัง Positions)
	DepartmentID   string    `json:"department_id"`                     // รหัสแผนก (FK ไปยัง Departments)
	HireDate       time.Time `json:"hire_date"`                         // วันที่เริ่มงาน
	EmploymentType string    `json:"employment_type"`                   // ประเภทการจ้างงาน (เช่น full-time, part-time)
	Address        Address   `json:"address"`                           // ที่อยู่ของพนักงาน
	BankInfo       BankInfo  `json:"bank_info"`                         // ข้อมูลบัญชีธนาคารของพนักงาน
}

type Address struct {
	AddressLine1 string `json:"address_line1"`           // ที่อยู่บรรทัด 1
	AddressLine2 string `json:"address_line2,omitempty"` // ที่อยู่บรรทัด 2
	Subdistrict  string `json:"subdistrict"`             // ตำบล
	District     string `json:"district"`                // อำเภอ
	Province     string `json:"province"`                // จังหวัด
	PostalCode   string `json:"postal_code"`             // รหัสไปรษณีย์
	Country      string `json:"country"`                 // ประเทศ
}

type BankInfo struct {
	BankName    string `json:"bank_name"`    // ชื่อธนาคาร
	AccountNo   string `json:"account_no"`   // เลขที่บัญชี
	AccountName string `json:"account_name"` // ชื่อบัญชี
}

type Document struct {
	Name       string     `json:"name"`                 // ชื่อเอกสาร
	FileURL    string     `json:"file_url"`             // ลิงก์ไฟล์เอกสาร
	Type       string     `json:"type"`                 // ประเภทเอกสาร เช่น "id_card", "degree"
	CreatedAt  time.Time  `json:"created_at"`           // วันที่สร้างเอกสาร
	UploadedAt time.Time  `json:"uploaded_at"`          // วันที่อัปโหลดเอกสาร
	DeletedAt  *time.Time `json:"deleted_at,omitempty"` // วันที่ลบเอกสาร (soft delete)
}

type RequestGetUserAll struct {
	Page      int    `query:"page"`       // หมายเลขหน้าที่ต้องการดึงข้อมูล
	Limit     int    `query:"limit"`      // จำนวนรายการต่อหน้า
	Search    string `query:"search"`     // คำค้นหาสำหรับกรองข้อมูล
	Status    string `query:"status"`     // สถานะของผู้ใช้ (เช่น pending, approved, rejected, cancelled)
	SortBy    string `query:"sort_by"`    // คอลัมน์ที่ต้องการเรียงลำดับ
	SortOrder string `query:"sort_order"` // ทิศทางการเรียงลำดับ (asc หรือ desc)
}

type RequestUpdateUser struct {
	Email             string              `json:"email"`              // อีเมลของผู้ใช้
	TitleTH           string              `json:"title_th"`           // คำนำหน้าชื่อ (ภาษาไทย)
	TitleEN           string              `json:"title_en"`           // คำนำหน้าชื่อ (ภาษาอังกฤษ)
	FirstNameTH       string              `json:"first_name_th"`      // ชื่อจริงของพนักงาน
	LastNameTH        string              `json:"last_name_th"`       // นามสกุลของพนักงาน
	FirstNameEN       string              `json:"first_name_en"`      // ชื่อจริงของพนักงาน (ภาษาอังกฤษ)
	LastNameEN        string              `json:"last_name_en"`       // นามสกุลของพนักงาน (ภาษาอังกฤษ)
	IDCard            string              `json:"id_card"`            // หมายเลขบัตรประชาชน (อาจเข้ารหัสก่อนจัดเก็บ)
	Phone             string              `json:"phone"`              // เบอร์โทรศัพท์ของพนักงาน
	EmployeeCode      string              `json:"employee_code"`      // รหัสพนักงาน (อาจใช้สำหรับอ้างอิงภายใน)
	Gender            string              `json:"gender"`             // เพศของพนักงาน
	BirthDate         string              `json:"birth_date"`         // วันเดือนปีเกิดของพนักงาน (รูปแบบ string)
	PositionID        string              `json:"position_id"`        // รหัสตำแหน่งงาน (FK ไปยัง Positions)
	DepartmentID      string              `json:"department_id"`      // รหัสแผนก (FK ไปยัง Departments)
	HireDate          string              `json:"hire_date"`          // วันที่เริ่มงาน
	EmploymentType    string              `json:"employment_type"`    // ประเภทการจ้างงาน (เช่น full-time, part-time)
	EmploymentHistory []EmploymentHistory `json:"employment_history"` // ประวัติการจ้างงาน (อาจมีหลายรายการ)
	Address           Address             `json:"address"`            // ที่อยู่ของพนักงาน
	BankInfo          BankInfo            `json:"bank_info"`          // ข้อมูลบัญชีธนาคารของพนักงาน
	Documents         []Document          `json:"documents"`          // รายการเอกสารที่เกี่ยวข้องกับพนักงาน
}

type RequestUpdateDocuments struct {
	UserID  string `json:"user_id"`  // รหัสประจำตัวผู้ใช้ (ไม่ซ้ำกัน)
	Type    string `json:"type"`     // ประเภทเอกสาร เช่น "id_card", "degree"
	Name    string `json:"name"`     // ชื่อเอกสาร
	FileURL string `json:"file_url"` // ลิงก์ไฟล์เอกสาร
}

// Response

type ResponseGetUserByID struct {
	UserID            string              `json:"user_id"`            // รหัสประจำตัวผู้ใช้ (ไม่ซ้ำกัน)
	Email             string              `json:"email"`              // อีเมลของผู้ใช้
	TitleTH           string              `json:"title_th"`           // คำนำหน้าชื่อ (ภาษาไทย)
	TitleEN           string              `json:"title_en"`           // คำนำหน้าชื่อ (ภาษาอังกฤษ)
	FirstNameTH       string              `json:"first_name_th"`      // ชื่อจริงของพนักงาน
	LastNameTH        string              `json:"last_name_th"`       // นามสกุลของพนักงาน
	FirstNameEN       string              `json:"first_name_en"`      // ชื่อจริงของพนักงาน (ภาษาอังกฤษ)
	LastNameEN        string              `json:"last_name_en"`       // นามสกุลของพนักงาน (ภาษาอังกฤษ)
	IDCard            string              `json:"id_card"`            // หมายเลขบัตรประชาชน (อาจเข้ารหัสก่อนจัดเก็บ)
	Role              string              `json:"role"`               // บทบาทหรือสิทธิ์ของผู้ใช้ในระบบ (เช่น admin, user)
	Avatar            string              `json:"avatar"`             // ลิงก์หรือที่อยู่รูปประจำตัวผู้ใช้
	Phone             string              `json:"phone"`              // เบอร์โทรศัพท์ของพนักงาน
	Status            string              `json:"status"`             // สถานะของผู้ใช้ (เช่น active, inactive)
	EmployeeCode      string              `json:"employee_code"`      // รหัสพนักงาน (อาจใช้สำหรับอ้างอิงภายใน)
	Gender            string              `json:"gender"`             // เพศของพนักงาน
	BirthDate         time.Time           `json:"birth_date"`         // วันเดือนปีเกิดของพนักงาน (รูปแบบ string)
	Position          string              `json:"position"`           // รหัสตำแหน่งงาน (FK ไปยัง Positions)
	Department        string              `json:"department"`         // รหัสแผนก (FK ไปยัง Departments)
	HireDate          time.Time           `json:"hire_date"`          // วันที่เริ่มงาน
	EmploymentType    string              `json:"employment_type"`    // ประเภทการจ้างงาน (เช่น full-time, part-time)
	EmploymentHistory []EmploymentHistory `json:"employment_history"` // ประวัติการจ้างงาน (อาจมีหลายรายการ)
	Address           Address             `json:"address"`            // ที่อยู่ของพนักงาน
	BankInfo          BankInfo            `json:"bank_info"`          // ข้อมูลบัญชีธนาคารของพนักงาน
	Documents         []Document          `json:"documents"`
	CreatedAt         time.Time           `json:"created_at"` // วันที่สร้างข้อมูลนี้
	UpdatedAt         time.Time           `json:"updated_at"` // วันที่แก้ไขข้อมูลล่าสุด
	DeletedAt         *time.Time          `json:"deleted_at"` // วันที่ลบข้อมูล (soft delete)
}

type EmploymentHistory struct {
	UserID         string     `json:"user_id"`         // รหัสผู้ใช้ที่เกี่ยวข้อง
	PositionID     string     `json:"position_id"`     // ตำแหน่งในช่วงเวลานั้น
	DepartmentID   string     `json:"department_id"`   // แผนกในช่วงเวลานั้น
	FromDate       time.Time  `json:"from_date"`       // วันที่เริ่มต้น
	ToDate         *time.Time `json:"to_date"`         // วันที่สิ้นสุด (nullable ถ้ายังทำอยู่)
	EmploymentType string     `json:"employment_type"` // ประเภทการจ้าง (เช่น full-time, intern)
	Note           string     `json:"note,omitempty"`  // หมายเหตุ (ถ้ามี)
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"` // soft delete
}

type ResponseGetUserAll struct {
	UserID         string     `json:"user_id" example:"1d5855c2-7d14-4f8d-8b5d-ef20cb5cb3cf"` // รหัสประจำตัวผู้ใช้ (ไม่ซ้ำกัน)
	TitleTH        string     `json:"title_th" example:"นางสาว"`                              // คำนำหน้าชื่อ (ภาษาไทย)
	FirstNameTH    string     `json:"first_name_th" example:"กิตติยา"`                        // ชื่อจริงของพนักงาน
	LastNameTH     string     `json:"last_name_th" example:"จันทรสกุล"`                       // นามสกุลของพนักงาน
	TitleEN        string     `json:"title_en" example:"Miss"`                                // คำนำหน้าชื่อ (ภาษาอังกฤษ)
	FirstNameEN    string     `json:"first_name_en" example:"Kittiya"`                        // ชื่อจริงของพนักงาน (ภาษาอังกฤษ)
	LastNameEN     string     `json:"last_name_en" example:"Chanthasakul"`                    // นามสกุลของพนักงาน (ภาษาอังกฤษ)
	Avatar         string     `json:"avatar" example:"https://example.com/avatar.jpg"`        // ลิงก์หรือที่อยู่รูปประจำตัวผู้ใช้
	Email          string     `json:"email" example:"จันทรสกุล"`                              // อีเมลของผู้ใช้
	Phone          string     `json:"phone" example:"094-222-7788"`                           // เบอร์โทรศัพท์ของพนักงาน
	Position       string     `json:"position" example:"POS126"`                              // รหัสตำแหน่งงาน (FK ไปยัง Positions)
	Status         string     `json:"status" example:"approved"`                              // สถานะของผู้ใช้ (เช่น approved , pending, rejected)
	KPIScore       string     `json:"kpi_score" example:"85"`                                 // คะแนน KPI ของพนักงาน
	TasksTotal     string     `json:"tasks_total" example:"10"`                               // จำนวนงานทั้งหมดที่ได้รับมอบหมาย
	TasksCompleted string     `json:"tasks_completed" example:"8"`                            // จำนวนงานที่เสร็จสมบูรณ์
	Department     string     `json:"department" example:"DEP001"`                            // รหัสแผนก (FK ไปยัง Departments)
	CreatedAt      time.Time  `json:"created_at" example:"2025-07-11T08:25:08.526Z"`          // วันที่สร้างข้อมูลนี้
	UpdatedAt      time.Time  `json:"updated_at" example:"2025-07-11T08:25:08.526Z"`          // วันที่แก้ไขข้อมูลล่าสุด
	DeletedAt      *time.Time `json:"deleted_at" example:"null"`                              // วันที่ลบข้อมูล (soft delete)
}
