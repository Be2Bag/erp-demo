package dto

type KPITemplateDTO struct {
	Name        string            `json:"name"`         // ชื่อ Template
	Department  string            `json:"department"`   // แผนก (เลือกจาก dropdown)
	KPIs        []KPITemplateItem `json:"kpis"`         // รายการ KPI
	TotalWeight int               `json:"total_weight"` // น้ำหนักรวม (ต้อง = 100)
}

type KPITemplateItem struct {
	Name        string `json:"name"`        // ชื่อ KPI
	Category    string `json:"category"`    // หมวดหมู่ (เลือกจาก dropdown)
	Description string `json:"description"` // คำอธิบาย
	MaxScore    int    `json:"max_score"`   // คะแนนเต็ม
	Weight      int    `json:"weight"`      // น้ำหนัก %
}
