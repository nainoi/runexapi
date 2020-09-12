package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Workouts struct {
	ID                  primitive.ObjectID    `json:"id" bson:"_id,omitempty"`
	UserID              primitive.ObjectID    `json:"user_id" bson:"user_id"`
	WorkoutActivityInfo []WorkoutActivityInfo `json:"activity_info" bson:"activity_info,omitempty"`
}

type WorkoutActivityInfo struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ActivityType string             `form:"activity_type" json:"activity_type"`
	Calory       float64            `json:"calory" bson:"calory"`
	Caption      string             `form:"caption" json:"caption"`
	Distance     float64            `json:"distance" bson:"distance"`
	Pace         float64            `json:"pace" bson:"pace"`
	Time         float64            `json:"time" bson:"time"`
	ActivityDate time.Time          `json:"activity_date" bson:"activity_date"`
	ImagePath    string             `json:"image_path" bson:"image_path"`
	GpxData      string             `json:"gpx_data" bson:"gpx_data"`
}

type AddWorkout struct {
	UserID              primitive.ObjectID  `json:"user_id" bson:"user_id"`
	WorkoutActivityInfo WorkoutActivityInfo `form:"workout_info" json:"workout_info"`
}

type AddWorkoutForm struct {
	UserID       string  `form:"user_id" json:"user_id" bson:"user_id"`
	ActivityType string  `form:"activity_type" json:"activity_type"`
	Calory       float64 `form:"calory" json:"calory" bson:"calory"`
	Caption      string  `form:"caption" json:"caption"`
	Distance     float64 `form:"distance" json:"distance" bson:"distance"`
	Pace         float64 `form:"pace" json:"pace" bson:"pace"`
	Time         float64 `form:"time" json:"time" bson:"time"`
	ActivityDate string  `form:"activity_date" json:"activity_date" bson:"activity_date"`
	ImagePath    string  `form:"image_path" json:"image_path" bson:"image_path"`
	GpxData      string  `form:"gpx_data" json:"gpx_data" bson:"gpx_data"`
}
