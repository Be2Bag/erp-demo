package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KPITemplate struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	TemplateID  string             `bson:"template_id" json:"template_id"`
	Name        string             `bson:"name" json:"name"`
	Department  string             `bson:"department" json:"department"`
	TotalWeight int                `bson:"total_weight" json:"total_weight"`
	Items       []KPITemplateItem  `bson:"items" json:"items"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	Version     int                `bson:"version" json:"version"`
	CreatedBy   string             `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type KPITemplateItem struct {
	ItemID      string    `bson:"item_id" json:"item_id"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	Category    string    `bson:"category" json:"category"`
	MaxScore    int       `bson:"max_score" json:"max_score"`
	Weight      int       `bson:"weight" json:"weight"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}
