package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Banner struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EventID   primitive.ObjectID `json:"event_id" bson:"event_id"`
	Active    bool               `json:"active" bson:"active"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type BannerAddForm struct {
	EventID string `json:"event_id" bson:"event_id" binding:"required"`
	Active  bool   `json:"active" bson:"active"`
}
