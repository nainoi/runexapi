package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Merchant struct
type Merchant struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	EventID     primitive.ObjectID `json:"event_id" bson:"event_id"`
	EventCode   string             `json:"event_code" bson:"event_code"`
	RegID       primitive.ObjectID `json:"reg_id" bson:"reg_id"`
	Status      string             `json:"status" bson:"status"`
	PaymentType string             `json:"payment_type" bson:"payment_type"`
	TotalPrice  int64              `json:"total_price" bson:"total_price"`
	OrderID     string             `json:"order_id" bson:"order_id"`
	OmiseID     string             `json:"omise_id" bson:"omise_id"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}
