package dto

import "time"

// ---------- Request DTO ----------
type CreateSignJobDTO struct { // DTO สำหรับสร้างงานป้ายใหม่
	// ---------- ลูกค้า ----------
	CompanyName    string `json:"company_name"`     // ชื่อบริษัท (จำเป็น)
	ContactPerson  string `json:"contact_person"`   // ชื่อผู้ติดต่อ
	Phone          string `json:"phone"`            // เบอร์โทร
	Email          string `json:"email"`            // อีเมล
	CustomerTypeID string `json:"customer_type_id"` // รหัสประเภทลูกค้า
	Address        string `json:"address"`          // ที่อยู่ติดตั้ง / จัดส่ง

	// ---------- รายละเอียดงานป้าย ----------
	ProjectName string  `json:"project_name"` // ชื่อโปรเจกต์ (จำเป็น)
	JobName     string  `json:"job_name"`     // ชื่องาน (จำเป็น)
	SignTypeID  string  `json:"sign_type_id"` // รหัสประเภทป้าย
	Width       float64 `json:"width"`        // ความกว้าง (ซม.)
	Height      float64 `json:"height"`       // ความสูง (ซม.)
	Quantity    int     `json:"quantity"`     // จำนวน
	PriceTHB    int64   `json:"price_thb"`    // ราคา (หน่วย: สตางค์หรือบาท เลือกใช้ให้คงที่)
	Content     string  `json:"content"`      // รายละเอียด / ข้อความบนป้าย
	MainColor   string  `json:"main_color"`   // สีหลัก

	// ---------- การชำระเงิน ----------
	PaymentMethod string `json:"payment_method"` // วิธีชำระเงิน: deposit|cash|transfer|credit

	// ---------- การผลิต / ไทม์ไลน์ ----------
	ProductionTime string `json:"production_time"` // ระยะเวลาผลิต เช่น "5 วัน"
	DueDate        string `json:"due_date"`        // กำหนดส่ง (อาจว่าง)

	// ---------- งานออกแบบ / การติดตั้ง ----------
	DesignOption  string `json:"design_option"`  // ตัวเลือกออกแบบ: have_design|need_design
	InstallOption string `json:"install_option"` // ตัวเลือกติดตั้ง: none|self|shop

	// ---------- หมายเหตุ ----------
	Notes string `json:"notes"` // หมายเหตุเพิ่มเติม

	// ---------- เมต้า ----------
	Status    string `json:"status"`     // สถานะงาน (อาจให้ระบบตั้ง)
	CreatedBy string `json:"created_by"` // ผู้สร้าง (อาจให้ระบบตั้ง)
}

type UpdateSignJobDTO = struct { // DTO สำหรับสร้างงานป้ายใหม่
	// ---------- ลูกค้า ----------
	CompanyName    string `json:"company_name"`     // ชื่อบริษัท (จำเป็น)
	ContactPerson  string `json:"contact_person"`   // ชื่อผู้ติดต่อ
	Phone          string `json:"phone"`            // เบอร์โทร
	Email          string `json:"email"`            // อีเมล
	CustomerTypeID string `json:"customer_type_id"` // รหัสประเภทลูกค้า
	Address        string `json:"address"`          // ที่อยู่ติดตั้ง / จัดส่ง

	// ---------- รายละเอียดงานป้าย ----------
	ProjectName string  `json:"project_name"` // ชื่อโปรเจกต์ (จำเป็น)
	JobName     string  `json:"job_name"`     // ชื่องาน (จำเป็น)
	SignTypeID  string  `json:"sign_type_id"` // รหัสประเภทป้าย
	Width       float64 `json:"width"`        // ความกว้าง (ซม.)
	Height      float64 `json:"height"`       // ความสูง (ซม.)
	Quantity    int     `json:"quantity"`     // จำนวน
	PriceTHB    int64   `json:"price_thb"`    // ราคา (หน่วย: สตางค์หรือบาท เลือกใช้ให้คงที่)
	Content     string  `json:"content"`      // รายละเอียด / ข้อความบนป้าย
	MainColor   string  `json:"main_color"`   // สีหลัก

	// ---------- การชำระเงิน ----------
	PaymentMethod string `json:"payment_method"` // วิธีชำระเงิน: deposit|cash|transfer|credit

	// ---------- การผลิต / ไทม์ไลน์ ----------
	ProductionTime string `json:"production_time"` // ระยะเวลาผลิต เช่น "5 วัน"
	DueDate        string `json:"due_date"`        // กำหนดส่ง (อาจว่าง)

	// ---------- งานออกแบบ / การติดตั้ง ----------
	DesignOption  string `json:"design_option"`  // ตัวเลือกออกแบบ: have_design|need_design
	InstallOption string `json:"install_option"` // ตัวเลือกติดตั้ง: none|self|shop

	// ---------- หมายเหตุ ----------
	Notes string `json:"notes"` // หมายเหตุเพิ่มเติม

	// ---------- เมต้า ----------
	Status string `json:"status"` // สถานะงาน (อาจให้ระบบตั้ง)
}

type RequestListSignJobs struct {
	Page      int    `query:"page"`       // หมายเลขหน้าที่ต้องการดึงข้อมูล
	Limit     int    `query:"limit"`      // จำนวนรายการต่อหน้า
	Search    string `query:"search"`     // คำค้นหาสำหรับกรองข้อมูล
	Status    string `query:"status"`     // สถานะ
	SortBy    string `query:"sort_by"`    // คอลัมน์ที่ต้องการเรียงลำดับ
	SortOrder string `query:"sort_order"` // ทิศทางการเรียงลำดับ (asc หรือ desc)
}

// ---------- Response DTO ----------
type SignJobDTO struct { // DTO สำหรับส่งกลับให้ฝั่งไคลเอนต์
	// ---------- คีย์ ----------
	JobID string `json:"job_id"` // รหัสงาน (UUID / Unique)

	// ---------- ลูกค้า ----------
	CompanyName    string `json:"company_name"`     // ชื่อบริษัท
	ContactPerson  string `json:"contact_person"`   // ชื่อผู้ติดต่อ
	Phone          string `json:"phone"`            // เบอร์โทร
	Email          string `json:"email"`            // อีเมล
	CustomerTypeID string `json:"customer_type_id"` // รหัสประเภทลูกค้า
	Address        string `json:"address"`          // ที่อยู่ติดตั้ง / จัดส่ง

	// ---------- รายละเอียดงานป้าย ----------
	ProjectName string  `json:"project_name"` // ชื่อโปรเจกต์
	JobName     string  `json:"job_name"`     // ชื่องาน
	SignTypeID  string  `json:"sign_type_id"` // รหัสประเภทป้าย
	Width       float64 `json:"width"`        // ความกว้าง (ซม.)
	Height      float64 `json:"height"`       // ความสูง (ซม.)
	Quantity    int     `json:"quantity"`     // จำนวน
	PriceTHB    int64   `json:"price_thb"`    // ราคา
	Content     string  `json:"content"`      // รายละเอียด / ข้อความบนป้าย
	MainColor   string  `json:"main_color"`   // สีหลัก

	// ---------- การชำระเงิน ----------
	PaymentMethod string `json:"payment_method"` // วิธีชำระเงิน

	// ---------- การผลิต / ไทม์ไลน์ ----------
	ProductionTime string    `json:"production_time"` // ระยะเวลาผลิต
	DueDate        time.Time `json:"due_date"`        // กำหนดส่ง

	// ---------- งานออกแบบ / การติดตั้ง ----------
	DesignOption  string `json:"design_option"`  // ตัวเลือกออกแบบ
	InstallOption string `json:"install_option"` // ตัวเลือกติดตั้ง

	// ---------- หมายเหตุ ----------
	Notes string `json:"notes"` // หมายเหตุ

	// ---------- เมต้า ----------
	Status    string     `json:"status"`     // สถานะงาน
	CreatedBy string     `json:"created_by"` // ผู้สร้าง
	CreatedAt time.Time  `json:"created_at"` // เวลาสร้าง
	UpdatedAt time.Time  `json:"updated_at"` // เวลาอัปเดตล่าสุด
	DeletedAt *time.Time `json:"deleted_at"` // เวลาเมื่อถูกลบ (soft delete)
}
