package models

import "time"

const CollectionCategory = "categorys"

type Category struct {
	CategoryID     string     `bson:"category_id" json:"category_id"`
	DepartmentID   string     `bson:"department_id" json:"department_id"` // รหัสแผนก (FK ไปยัง Department)
	CategoryNameTH string     `bson:"category_name_th" json:"category_name_th"`
	CategoryNameEN string     `bson:"category_name_en" json:"category_name_en"`
	Description    string     `bson:"description" json:"description"`
	CreatedBy      string     `bson:"created_by" json:"created_by"`
	CreatedAt      time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time `bson:"deleted_at" json:"deleted_at"`
	Note           *string    `bson:"note" json:"note"`
}
