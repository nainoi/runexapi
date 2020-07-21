package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryMaster struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Active      bool               `json:"active" bson:"active"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type CategoryUpdateForm struct {
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Active      bool               `json:"active" bson:"active"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}
