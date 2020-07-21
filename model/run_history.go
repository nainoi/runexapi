package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddHistoryForm struct {
	ActivityType string  `form:"activity_type" json:"activity_type"`
	Calory       float32 `form:"calory" json:"calory" bson:"calory"`
	Caption      string  `form:"caption" json:"caption" bson:"caption"`
	Distance     float64 `form:"distance" json:"distance" bson:"distance" binding:"required"`
	Pace         float32 `form:"pace" json:"pace" bson:"pace"`
	Time         float32 `form:"time" json:"time" bson:"time"`
	ImagePath    string  `form:"image_path" json:"image_path" bson:"image_path"`
}

type RunHistory struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID         primitive.ObjectID `json:"user_id" bson:"user_id"`
	RunHistoryInfo []RunHistoryInfo   `json:"activity_info" bson:"activity_info,omitempty"`
	ToTalDistance  float64            `json:"total_distance" bson:"total_distance"`
}

type RunHistoryInfo struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ActivityType string             `form:"activity_type" json:"activity_type"`
	Calory       float32            `form:"calory" json:"calory" bson:"calory"`
	Caption      string             `form:"caption" json:"caption" bson:"caption"`
	Distance     float64            `form:"distance" json:"distance" bson:"distance" binding:"required"`
	Pace         float32            `form:"pace" json:"pace" bson:"pace"`
	Time         float32            `form:"time" json:"time" bson:"time"`
	ActivityDate time.Time          `json:"activity_date" bson:"activity_date"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	ImagePath    string             `form:"image_path" json:"image_path" bson:"image_path"`
}
