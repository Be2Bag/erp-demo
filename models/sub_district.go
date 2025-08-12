package models

import "time"

const CollectionSubDistricts = "sub_districts"

type SubDistrict struct {
	ID         string     `json:"id" bson:"id"`                   // ไอดีของตำบล
	ZipCode    string     `json:"zip_code" bson:"zip_code"`       // รหัสไปรษณีย์
	NameTH     string     `json:"name_th" bson:"name_th"`         // ชื่อตำบล
	NameEN     string     `json:"name_en" bson:"name_en"`         // ชื่อตำบลภาษาอังกฤษ
	DistrictID string     `json:"district_id" bson:"district_id"` // ไอดีอำเภอ
	CreatedAt  time.Time  `json:"created_at" bson:"created_at"`   // วันที่สร้าง
	UpdatedAt  time.Time  `json:"updated_at" bson:"updated_at"`   // วันที่แก้ไขล่าสุด
	DeletedAt  *time.Time `json:"deleted_at" bson:"deleted_at"`   // วันที่ลบ (ถ้ามี)
}
