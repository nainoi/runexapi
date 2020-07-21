package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EbibEvent for count bib
type EbibEvent struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EventID   primitive.ObjectID `json:"event_id" bson:"event_id"`
	LastNo    int64              `json:"last_no" bson:"last_no"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
