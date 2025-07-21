package dto

type ResponseGetPositions struct {
	PositionID   string `json:"position_id"`   // รหัสตำแหน่งงาน (ไม่ซ้ำกัน)
	PositionName string `json:"position_name"` // ชื่อตำแหน่งงาน
}

type ResponseGetDepartments struct {
	DepartmentID   string `json:"department_id"`   // รหัสแผนก (ไม่ซ้ำกัน)
	DepartmentName string `json:"department_name"` // ชื่อแผนก
}

type ResponseGetProvinces struct {
	ProvinceID   string `json:"province_id"`   // รหัสจังหวัด (ไม่ซ้ำกัน)
	ProvinceName string `json:"province_name"` // ชื่อจังหวัด
}

type ResponseGetDistricts struct {
	DistrictID   string `json:"district_id"`   // รหัสอำเภอ (
	DistrictName string `json:"district_name"` // ชื่ออำเภอ
}

type ResponseGetSubDistricts struct {
	SubDistrictID   string `json:"sub_district_id"`   // รหัสตำบล (ไม่ซ้ำกัน)
	SubDistrictName string `json:"sub_district_name"` // ชื่อตำบล
}
