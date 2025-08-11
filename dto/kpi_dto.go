package dto

type KPITemplateDTO struct {
	Name        string               `json:"name"`         // ชื่อ Template
	Department  string               `json:"department"`   // แผนก
	TotalWeight int                  `json:"total_weight"` // น้ำหนักรวม (ต้อง = 100)
	Items       []KPITemplateItemDTO `json:"items"`        // รายการ KPI
	Version     int                  `json:"version"`      // เวอร์ชันของ Template
}

type KPITemplateItemDTO struct {
	Name        string `json:"name"`        // ชื่อ KPI
	Description string `json:"description"` // คำอธิบาย
	Category    string `json:"category"`    // หมวดหมู่
	MaxScore    int    `json:"max_score"`   // คะแนนเต็ม
	Weight      int    `json:"weight"`      // น้ำหนัก %
}
