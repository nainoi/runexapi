package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ActivityV2 data per event
type ActivityV2 struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EventID    primitive.ObjectID `json:"event_id" bson:"event_id"`
	Activities Activities         `json:"activities" bson:"activities,omitempty"`
}

// Activities data per user
type Activities struct {
	UserID        primitive.ObjectID `json:"user_id" bson:"user_id"`
	ToTalDistance float64            `json:"total_distance" bson:"total_distance"`
	ActivityInfo  []ActivityInfo     `json:"activity_info" bson:"activity_info,omitempty"`
}

//AddActivityV2 data
type AddActivityV2 struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	EventID      primitive.ObjectID `json:"event_id" bson:"event_id"`
	ActivityInfo ActivityInfo       `json:"activity_info" bson:"activity_info"`
}

// type ActivityInfoV2 struct {
// 	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
// 	Distance     float64            `json:"distance" bson:"distance"`
// 	ImageURL     string             `json:"img_url" bson:"img_url"`
// 	ActivityDate time.Time          `json:"activity_date" bson:"activity_date"`
// 	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
// 	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
// }

//AddActivityFormWorkout model request
type AddActivityFormWorkout struct {
	WorkoutActivityInfo WorkoutActivityInfo `form:"workout_info" json:"workout_info"`
	EventID             string              `form:"event_id" json:"event_id" binding:"required"`
}

//AddMultiActivityFormWorkout model request
type AddMultiActivityFormWorkout struct {
	WorkoutActivityInfo WorkoutActivityInfo `form:"workout_info" json:"workout_info"`
	EventActivity       []EventActivity     `form:"event_activity" json:"event_activity" binding:"required"`
}

// EventActivity for request add activity
type EventActivity struct {
	EventID string       `json:"event_id" bson:"event_id"`
	Partner PartnerEvent `json:"partner" bson:"partner"`
}
