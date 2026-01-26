package models

import "time"

const CollectionProvinces = "provinces"

type Province struct {
	CreatedAt   time.Time  `json:"created_at" bson:"created_at"`     // วันที่สร้างข้อมูล
	UpdatedAt   time.Time  `json:"updated_at" bson:"updated_at"`     // วันที่แก้ไขข้อมูลล่าสุด
	DeletedAt   *time.Time `json:"deleted_at" bson:"deleted_at"`     // วันที่ลบข้อมูล (ถ้ามี)
	ID          string     `json:"id" bson:"id"`                     // รหัสจังหวัด
	NameTH      string     `json:"name_th" bson:"name_th"`           // ชื่อจังหวัดภาษาไทย
	NameEN      string     `json:"name_en" bson:"name_en"`           // ชื่อจังหวัดภาษาอังกฤษ
	GeographyID string     `json:"geography_id" bson:"geography_id"` // รหัสภูมิภาค
}
