package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionWorkflowTemplates = "workflow_templates"

type WorkFlowTemplate struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	TemplateID  string             `bson:"template_id" json:"template_id"` // UUID
	Name        string             `bson:"name" json:"name"`
	Department  string             `bson:"department" json:"department"`
	Description string             `bson:"description" json:"description"`
	TotalHours  float64            `bson:"total_hours" json:"total_hours"` // cache ผลรวมชั่วโมง
	Steps       []WorkFlowStep     `bson:"steps" json:"steps"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	Version     int                `bson:"version" json:"version"`
	CreatedBy   string             `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type WorkFlowStep struct {
	StepID      string    `bson:"step_id" json:"step_id"` // UUID
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description,omitempty" json:"description,omitempty"`
	Hours       float64   `bson:"hours" json:"hours"` // รองรับ 0.5 ฯลฯ
	Order       int       `bson:"order" json:"order"` // 1..N
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}
