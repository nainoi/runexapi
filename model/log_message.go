package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LogUpdateRegisterStatus struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	LogName   string             `json:"log_name" bson:"log_name"`
	Status    string             `json:"status" bson:"status" `
	RegID     primitive.ObjectID `json:"reg_id" bson:"reg_id" `
	UpdateBy  primitive.ObjectID `json:"update_by" bson:"update_by" `
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
