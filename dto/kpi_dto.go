package dto

type KPITemplateDTO struct {
	Name        string            `json:"name"`         // ชื่อ Template
	Department  string            `json:"department"`   // แผนก
	Templates   []KPITemplateList `json:"templates"`    // รายการ KPI
	TargetValue int               `json:"target_value"` // ค่าเป้าหมายรวม (100%)
	IsActive    bool              `json:"is_active"`
}

type KPITemplateList struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	MaxScore    int    `json:"max_score"` // คะแนนเต็ม
	Value       int    `json:"value"`     // น้ำหนัก %
}
