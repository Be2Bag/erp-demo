package dto

import "time"

type CreateKPITemplateDTO struct {
	KPIName     string                     `json:"kpi_name"`
	Department  string                     `json:"department_id"` // แผนก
	TotalWeight int                        `json:"total_weight"`  // น้ำหนักรวม (ต้อง = 100)
	Items       []CreateKPITemplateItemDTO `json:"items"`         // รายการ KPI
}

type CreateKPITemplateItemDTO struct {
	Name        string `json:"name"`        // ชื่อ KPI
	Description string `json:"description"` // คำอธิบาย
	Category    string `json:"category"`    // หมวดหมู่
	MaxScore    int    `json:"max_score"`   // คะแนนเต็ม
	Weight      int    `json:"weight"`      // น้ำหนัก %
}

type UpdateKPITemplateDTO struct {
	KPIName     string                      `json:"kpi_name"`
	Department  string                      `json:"department_id"` // แผนก
	TotalWeight int                         `json:"total_weight"`  // น้ำหนักรวม (ต้อง = 100)
	Items       *[]CreateKPITemplateItemDTO `json:"items"`         // รายการ KPI
}

// added for list query
type KPITemplateListQuery struct {
	Page       int    `query:"page"`          // หมายเลขหน้าที่ต้องการดึงข้อมูล
	Limit      int    `query:"limit"`         // จำนวนรายการต่อหน้า
	Search     string `query:"search"`        // คำค้นหาสำหรับกรองข้อมูล
	Department string `query:"department_id"` // แผนก
	SortBy     string `query:"sort_by"`       // คอลัมน์ที่ต้องการเรียงลำดับ
	SortOrder  string `query:"sort_order"`    // ทิศทางการเรียงลำดับ (asc หรือ desc)
}

type KPITemplateDTO struct {
	KPIID       string               `json:"kpi_id"`
	KPIName     string               `json:"kpi_name"`
	Department  string               `json:"department_id"`
	TotalWeight int                  `json:"total_weight"`
	Items       []KPITemplateItemDTO `json:"items"`
	Version     int                  `json:"version"`
	CreatedBy   string               `json:"created_by"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

type KPITemplateItemDTO struct {
	ItemID      string    `json:"item_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	MaxScore    int       `json:"max_score"`
	Weight      int       `json:"weight"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
