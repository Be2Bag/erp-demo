package dto

type ResponseGetPositions struct {
	PositionID   string `json:"position_id"`   // รหัสตำแหน่งงาน (ไม่ซ้ำกัน)
	PositionName string `json:"position_name"` // ชื่อตำแหน่งงาน
}

type ResponseGetDepartments struct {
	DepartmentID   string `json:"department_id"`   // รหัสแผนก (ไม่ซ้ำกัน)
	DepartmentName string `json:"department_name"` // ชื่อแผนก
}
