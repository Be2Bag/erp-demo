package models

import "time"

const CollectionSProject = "projects"

type Project struct {
	ProjectID   string     `bson:"project_id" json:"project_id"`
	ProjectName string     `bson:"project_name" json:"project_name"`
	CreatedBy   string     `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `bson:"deleted_at" json:"deleted_at"`
	Note        *string    `bson:"note" json:"note"`
}
