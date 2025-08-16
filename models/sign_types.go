package models

import "time"

const CollectionSignTypes = "sign_types"

type SignType struct {
	TypeID    string     `json:"type_id" bson:"type_id"`       // ไอดี (แนะนำใช้ UUID)
	NameTH    string     `json:"name_th" bson:"name_th"`       // ชื่อภาษาไทย
	NameEN    string     `json:"name_en" bson:"name_en"`       // ชื่อภาษาอังกฤษ
	CreatedAt time.Time  `json:"created_at" bson:"created_at"` // วันที่สร้าง
	UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"` // วันที่แก้ไขล่าสุด
	DeletedAt *time.Time `json:"deleted_at" bson:"deleted_at"` // วันที่ลบ (soft delete)
}
