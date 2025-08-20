package dto

import "time"

type CreateKPITemplateDTO struct {
	KPIName     string                     `json:"kpi_name"`
	Department  string                     `json:"department"`   // แผนก
	TotalWeight int                        `json:"total_weight"` // น้ำหนักรวม (ต้อง = 100)
	Items       []CreateKPITemplateItemDTO `json:"items"`        // รายการ KPI
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
	Department  string                      `json:"department"`   // แผนก
	TotalWeight int                         `json:"total_weight"` // น้ำหนักรวม (ต้อง = 100)
	Items       *[]CreateKPITemplateItemDTO `json:"items"`        // รายการ KPI
}

// added for list query
type KPITemplateListQuery struct {
	Page       int    `query:"page"`       // หมายเลขหน้าที่ต้องการดึงข้อมูล
	Limit      int    `query:"limit"`      // จำนวนรายการต่อหน้า
	Search     string `query:"search"`     // คำค้นหาสำหรับกรองข้อมูล
	Department string `query:"department"` // แผนก
	SortBy     string `query:"sort_by"`    // คอลัมน์ที่ต้องการเรียงลำดับ
	SortOrder  string `query:"sort_order"` // ทิศทางการเรียงลำดับ (asc หรือ desc)
}

type KPITemplateDTO struct {
	KPIID       string               `bson:"kpi_id" json:"kpi_id"`
	KPIName     string               `bson:"kpi_name" json:"kpi_name"`
	Department  string               `bson:"department" json:"department"`
	TotalWeight int                  `bson:"total_weight" json:"total_weight"`
	Items       []KPITemplateItemDTO `bson:"items" json:"items"`
	Version     int                  `bson:"version" json:"version"`
	CreatedBy   string               `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time            `bson:"updated_at" json:"updated_at"`
}

type KPITemplateItemDTO struct {
	ItemID      string    `bson:"item_id" json:"item_id"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	Category    string    `bson:"category" json:"category"`
	MaxScore    int       `bson:"max_score" json:"max_score"`
	Weight      int       `bson:"weight" json:"weight"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}
