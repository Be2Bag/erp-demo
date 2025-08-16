package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionSignJobs = "sign_jobs"

type SignJob struct {
	ID    primitive.ObjectID `bson:"_id " json:"_id"`      // ObjectId ใน MongoDB
	JobID string             `bson:"job_id" json:"job_id"` // UUID ของงาน (ใช้ค้นหาแบบ unique)

	ProjectName string `bson:"project_name" json:"project_name"` // ชื่อโปรเจกต์
	JobName     string `bson:"job_name" json:"job_name"`         // ชื่องาน

	CustomerName  string `bson:"customer_name" json:"customer_name"`     // ชื่อบริษัทหรือลูกค้า
	ContactPerson string `bson:"contact_person " json:"contact_person "` // ชื่อผู้ติดต่อ
	Phone         string `bson:"phone " json:"phone "`                   // เบอร์โทรศัพท์ผู้ติดต่อ
	Email         string `bson:"email " json:"email "`                   // อีเมลผู้ติดต่อ

	CustomerTypeID string `bson:"customer_type_id " json:"customer_type_id "` // ID ประเภทลูกค้า (dropdown อ้างอิง master data)
	Address        string `bson:"address " json:"address "`                   // ที่อยู่ติดตั้งหรือจัดส่ง

	SignTypeID string `bson:"sign_type_id " json:"sign_type_id "` // ID ประเภทงานป้าย (dropdown)
	Size       string `bson:"size " json:"size "`                 // ขนาดป้าย (กว้างxสูง) หน่วยซม.
	Quantity   int    `bson:"quantity " json:"quantity "`         // จำนวนป้าย
	Content    string `bson:"content " json:"content "`           // รายละเอียดข้อความหรือเนื้อหาบนป้าย
	MainColor  string `bson:"main_color " json:"main_color "`     // สีหลัก / โทนสี

	DesignOption   string    `bson:"design_option " json:"design_option "`     // ตัวเลือกการออกแบบ ("have_design" = มีแบบแล้ว, "need_design" = ให้ช่วยออกแบบ)
	ProductionTime string    `bson:"production_time " json:"production_time "` // ระยะเวลาในการผลิต
	DueDate        time.Time `bson:"due_date" json:"due_date"`                 // วันที่ต้องการรับงาน (ISO-8601 หรือ string)
	InstallOption  string    `bson:"install_option " json:"install_option "`   // ตัวเลือกการติดตั้ง ("none", "self", "shop")

	Notes string `bson:"notes " json:"notes "` // หมายเหตุเพิ่มเติม

	Status    string    `bson:"status" json:"status"`         // สถานะงาน ("draft", "active", "cancelled")
	CreatedBy string    `bson:"created_by" json:"created_by"` // UUID ผู้สร้าง
	CreatedAt time.Time `bson:"created_at" json:"created_at"` // วันที่สร้าง
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"` // วันที่แก้ไขล่าสุด
}
