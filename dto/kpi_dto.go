package dto

import "time"

type CreateKPITemplateDTO struct {
	KPIName     string                     `json:"kpi_name"`
	Department  string                     `json:"department_id"` // แผนก
	Items       []CreateKPITemplateItemDTO `json:"items"`         // รายการ KPI
	TotalWeight int                        `json:"total_weight"`  // น้ำหนักรวม (ต้อง = 100)
}

type CreateKPITemplateItemDTO struct {
	Name        string `json:"name"`        // ชื่อ KPI
	Description string `json:"description"` // คำอธิบาย
	Category    string `json:"category"`    // หมวดหมู่
	MaxScore    int    `json:"max_score"`   // คะแนนเต็ม
	Weight      int    `json:"weight"`      // น้ำหนัก %
}

type UpdateKPITemplateDTO struct {
	Items       *[]CreateKPITemplateItemDTO `json:"items"` // รายการ KPI
	KPIName     string                      `json:"kpi_name"`
	Department  string                      `json:"department_id"` // แผนก
	TotalWeight int                         `json:"total_weight"`  // น้ำหนักรวม (ต้อง = 100)
}

// added for list query
type KPITemplateListQuery struct {
	Search     string `query:"search"`        // คำค้นหาสำหรับกรองข้อมูล
	Department string `query:"department_id"` // แผนก
	SortBy     string `query:"sort_by"`       // คอลัมน์ที่ต้องการเรียงลำดับ
	SortOrder  string `query:"sort_order"`    // ทิศทางการเรียงลำดับ (asc หรือ desc)
	Page       int    `query:"page"`          // หมายเลขหน้าที่ต้องการดึงข้อมูล
	Limit      int    `query:"limit"`         // จำนวนรายการต่อหน้า
}

type KPITemplateDTO struct {
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	KPIID       string               `json:"kpi_id"`
	KPIName     string               `json:"kpi_name"`
	Department  string               `json:"department_id"`
	CreatedBy   string               `json:"created_by"`
	Items       []KPITemplateItemDTO `json:"items"`
	TotalWeight int                  `json:"total_weight"`
	Version     int                  `json:"version"`
}

type KPITemplateItemDTO struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ItemID      string    `json:"item_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	MaxScore    int       `json:"max_score"`
	Weight      int       `json:"weight"`
}
