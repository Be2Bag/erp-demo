package models // แพ็กเกจ models สำหรับโครงสร้างข้อมูล workflow

import (
	"time" // ใช้จัดการเวลา สร้าง/แก้ไข
)

const CollectionWorkflowTemplates = "workflow_templates" // ชื่อคอลเลกชันในฐานข้อมูล

type WorkFlowTemplate struct {
	CreatedAt    time.Time      `bson:"created_at" json:"created_at"`       // วันเวลาสร้าง
	UpdatedAt    time.Time      `bson:"updated_at" json:"updated_at"`       // วันเวลาอัปเดตล่าสุด
	DeletedAt    *time.Time     `bson:"deleted_at" json:"deleted_at"`       // วันที่ลบ (soft delete)
	WorkFlowID   string         `bson:"workflow_id" json:"workflow_id"`     // รหัส Workflow (UUID)
	WorkFlowName string         `bson:"workflow_name" json:"workflow_name"` // ชื่อ Workflow
	Department   string         `bson:"department_id" json:"department_id"` // แผนกที่เกี่ยวข้อง
	Description  string         `bson:"description" json:"description"`     // รายละเอียดเพิ่มเติม
	CreatedBy    string         `bson:"created_by" json:"created_by"`       // ผู้สร้าง
	Steps        []WorkFlowStep `bson:"steps" json:"steps"`                 // ลำดับขั้นตอนทั้งหมด
	TotalHours   float64        `bson:"total_hours" json:"total_hours"`     // ชั่วโมงรวม (แคชจากผลรวม step)
	Version      int            `bson:"version" json:"version"`             // เวอร์ชันของ template
	IsActive     bool           `bson:"is_active" json:"is_active"`         // สถานะใช้งานหรือไม่
}

type WorkFlowStep struct {
	CreatedAt   time.Time  `bson:"created_at" json:"created_at"`                       // วันเวลาสร้าง
	UpdatedAt   time.Time  `bson:"updated_at" json:"updated_at"`                       // วันเวลาอัปเดตล่าสุด
	DeletedAt   *time.Time `bson:"deleted_at" json:"deleted_at"`                       // วันที่ลบ (soft delete)
	StepID      string     `bson:"step_id" json:"step_id"`                             // รหัส Step (UUID)
	StepName    string     `bson:"step_name" json:"step_name"`                         // ชื่อ Step
	Description string     `bson:"description,omitempty" json:"description,omitempty"` // รายละเอียด (ไม่บังคับ)
	Hours       float64    `bson:"hours" json:"hours"`                                 // ชั่วโมงที่ใช้ (รองรับทศนิยม เช่น 0.5)
	Order       int        `bson:"order" json:"order"`                                 // ลำดับขั้น (1..N)
}
