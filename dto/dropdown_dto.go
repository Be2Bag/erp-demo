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
	ZipCode         string `json:"zip_code"`          // รหัสไปรษณีย์
}

type ResponseGetSignTypes struct {
	TypeID string `json:"type_id"` // รหัสประเภทงาน
	NameTH string `json:"name_th"` // ชื่อประเภทงาน (ภาษาไทย)
	NameEN string `json:"name_en"` // ชื่อประเภทงาน (ภาษาอังกฤษ)
}

type ResponseGetCustomerTypes struct {
	TypeID string `json:"type_id"` // รหัสประเภทลูกค้า (ไม่ซ้ำกัน)
	NameTH string `json:"name_th"` // ชื่อประเภทลูกค้า (ภาษาไทย)
	NameEN string `json:"name_en"` // ชื่อประเภทลูกค้า (ภาษาอังกฤษ)
}

type ResponseGetSignList struct {
	JobID       string `json:"job_id"`       // รหัสงาน
	ProjectName string `json:"project_name"` // ชื่อโปรเจกต์
	JobName     string `json:"job_name"`     // ชื่องาน
	Content     string `json:"content"`      // รายละเอียด
}

type ResponseGetProjects struct {
	ProjectID   string `json:"project_id"`   // รหัสโครงการ (ไม่ซ้ำกัน)
	ProjectName string `json:"project_name"` // ชื่อโครงการ
}

type ResponseGetUsers struct {
	UserID     string `json:"user_id"`      // รหัสผู้ใช้ (ไม่ซ้ำกัน)
	FullNameTH string `json:"full_name_th"` // ชื่อเต็ม (ภาษาไทย)
}

type ResponseGetKPI struct {
	KPIID   string `json:"kpi_id"`   // รหัส KPI (ไม่ซ้ำกัน)
	KPIName string `json:"kpi_name"` // ชื่อ KPI
}

type ResponseGetWorkflows struct {
	WorkflowID   string `json:"workflow_id"`   // รหัส Workflow (ไม่ซ้ำกัน)
	WorkflowName string `json:"workflow_name"` // ชื่อ Workflow
}

type ResponseGetCategorys struct {
	CategoryID     string `json:"category_id"`      // รหัสหมวดหมู่ (ไม่ซ้ำกัน)
	CategoryNameTH string `json:"category_name_th"` // ชื่อหมวดหมู่
	CategoryNameEN string `json:"category_name_en"` // ชื่อหมวดหมู่ (ภาษาอังกฤษ)
}
