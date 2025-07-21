package models

import "time"

type District struct {
	ID         string     `json:"id" bson:"id"`                   // ไอดีของอำเภอ
	NameTH     string     `json:"name_th" bson:"name_th"`         // ชื่ออำเภอภาษาไทย
	NameEN     string     `json:"name_en" bson:"name_en"`         // ชื่ออำเภอภาษาอังกฤษ
	ProvinceID string     `json:"province_id" bson:"province_id"` // ไอดีจังหวัด
	CreatedAt  time.Time  `json:"created_at" bson:"created_at"`   // วันที่สร้าง
	UpdatedAt  time.Time  `json:"updated_at" bson:"updated_at"`   // วันที่แก้ไขล่าสุด
	DeletedAt  *time.Time `json:"deleted_at" bson:"deleted_at"`   // วันที่ลบ (ถ้ามี)
}
