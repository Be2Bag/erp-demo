package dto

import "time"

// ---------- Request DTO ----------
type CreateSignJobDTO struct {
	ProjectName   string `json:"project_name"`             // ชื่อโปรเจกต์
	JobName       string `json:"job_name"`                 // ชื่องาน
	CustomerName  string `json:"customer_name"`            // ชื่อลูกค้าหรือบริษัท
	ContactPerson string `json:"contact_person,omitempty"` // ผู้ติดต่อ
	Phone         string `json:"phone,omitempty"`          // เบอร์โทรศัพท์
	Email         string `json:"email,omitempty"`          // อีเมล

	CustomerTypeID string `json:"customer_type_id,omitempty"` // ID ประเภทลูกค้า (dropdown)
	Address        string `json:"address,omitempty"`          // ที่อยู่ติดตั้ง / จัดส่ง

	SignTypeID string `json:"sign_type_id,omitempty"` // ID ประเภทงานป้าย (dropdown)
	Size       string `json:"size,omitempty"`         // ขนาดป้าย (กว้างxสูง) หน่วยซม.
	Quantity   int    `json:"quantity,omitempty"`     // จำนวนป้าย
	Content    string `json:"content,omitempty"`      // รายละเอียดข้อความบนป้าย
	MainColor  string `json:"main_color,omitempty"`   // สีหลัก / โทนสี

	DesignOption   string `json:"design_option,omitempty"`   // ตัวเลือกการออกแบบ ("have_design", "need_design")
	ProductionTime string `json:"production_time,omitempty"` // ระยะเวลาในการผลิต
	DueDate        string `json:"due_date,omitempty"`        // วันที่ต้องการรับงาน
	InstallOption  string `json:"install_option,omitempty"`  // ตัวเลือกการติดตั้ง ("none", "self", "shop")

	Notes string `json:"notes,omitempty"` // หมายเหตุเพิ่มเติม
}

type UpdateSignJobDTO = CreateSignJobDTO

// ---------- Response DTO ----------
type SignJobDTO struct {
	JobID         string `json:"job_id"`                   // UUID ของงาน
	ProjectName   string `json:"project_name"`             // ชื่อโปรเจกต์
	JobName       string `json:"job_name"`                 // ชื่องาน
	CustomerName  string `json:"customer_name"`            // ชื่อลูกค้า/บริษัท
	ContactPerson string `json:"contact_person,omitempty"` // ผู้ติดต่อ
	Phone         string `json:"phone,omitempty"`          // เบอร์โทรศัพท์
	Email         string `json:"email,omitempty"`          // อีเมล

	CustomerTypeID string `json:"customer_type_id,omitempty"` // ID ประเภทลูกค้า
	Address        string `json:"address,omitempty"`          // ที่อยู่ติดตั้ง/จัดส่ง

	SignTypeID string `json:"sign_type_id,omitempty"` // ID ประเภทงานป้าย
	Size       string `json:"size,omitempty"`         // ขนาดป้าย
	Quantity   int    `json:"quantity,omitempty"`     // จำนวนป้าย
	Content    string `json:"content,omitempty"`      // รายละเอียดข้อความ
	MainColor  string `json:"main_color,omitempty"`   // สีหลัก / โทนสี

	DesignOption   string    `json:"design_option,omitempty"`   // การออกแบบ
	ProductionTime string    `json:"production_time,omitempty"` // ระยะเวลาในการผลิต
	DueDate        time.Time `json:"due_date,omitempty"`        // วันที่ต้องการรับงาน
	InstallOption  string    `json:"install_option,omitempty"`  // ตัวเลือกการติดตั้ง

	Notes string `json:"notes,omitempty"` // หมายเหตุเพิ่มเติม

	Status    string    `json:"status"`     // สถานะงาน
	CreatedBy string    `json:"created_by"` // ผู้สร้าง
	CreatedAt time.Time `json:"created_at"` // วันที่สร้าง
	UpdatedAt time.Time `json:"updated_at"` // วันที่แก้ไขล่าสุด
}
