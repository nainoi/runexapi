package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//AddActivityForm model request
type AddActivityForm struct {
	Caption      string  `form:"caption" json:"caption"`
	Distance     float64 `form:"distance" json:"distance" bson:"distance" binding:"required"`
	ActivityDate string  `form:"activity_date" json:"activity_date" bson:"distance"`
	EventCode    string  `form:"event_code" json:"event_code" bson:"event_code" binding:"required"`
	UserID       string  `form:"user_id" json:"user_id" bson:"user_id"`
	ImageURL     string  `form:"image_url" json:"image_url" bson:"image_url"`
}

//AddMultiActivityForm model request
type AddMultiActivityForm struct {
	Caption      string   `form:"caption" json:"caption"`
	Distance     float64  `form:"distance" json:"distance" bson:"distance" binding:"required"`
	ActivityDate string   `form:"activity_date" json:"activity_date" bson:"distance"`
	EventCodes   []string `form:"event_codes" json:"event_codes" bson:"event_codes" binding:"required"`
	UserID       string   `form:"user_id" json:"user_id" bson:"user_id"`
	ImageURL     string   `form:"image_url" json:"image_url" bson:"image_url"`
}

type AddActivity struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	EventCode    string             `json:"event_code" bson:"event_code"`
	EventUser    string             `json:"event_user" bson:"event_user"`
	ActivityInfo ActivityInfo       `json:"activity_info" bson:"activity_info"`
	// Distance     float32   `json:"distance" bson:"distance"`
	// ImageURL     string    `json:"img_url" bson:"img_url"`
	// ActivityDate time.Time `json:"activity_date" bson:"activity_date"`
	// CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	// UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
	// ActivityType string    `json:"activity_type" bson:"activity_type"`
}

type Activity struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID        primitive.ObjectID `json:"user_id" bson:"user_id"`
	EventCode     string             `json:"event_code" bson:"event_code"`
	EventUser     string             `json:"event_user" bson:"event_user"`
	ActivityInfo  []ActivityInfo     `json:"activity_info" bson:"activity_info,omitempty"`
	ToTalDistance float64            `json:"total_distance" bson:"total_distance"`
	// Distance     float32   `json:"distance" bson:"distance"`
	// ImageURL     string    `json:"img_url" bson:"img_url"`
	// ActivityDate time.Time `json:"activity_date" bson:"activity_date"`
	// CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	// UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
	// ActivityType string    `json:"activity_type" bson:"activity_type"`
}

// type EventActivity struct {
// 	EventID  string         `json:"event_id" bson:"event_id"`
// 	Activity []ActivityInfo `json:"activity" bson:"activity"`
// }

type ActivitySingel struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID        primitive.ObjectID `json:"user_id" bson:"user_id"`
	EventID       primitive.ObjectID `json:"event_id" bson:"event_id"`
	EventUser     string             `json:"event_user" bson:"event_user"`
	ActivityInfo  ActivityInfo       `json:"activity_info" bson:"activity_info,omitempty"`
	ToTalDistance float32            `json:"total_distance" bson:"total_distance"`
	// Distance     float32   `json:"distance" bson:"distance"`
	// ImageURL     string    `json:"img_url" bson:"img_url"`
	// ActivityDate time.Time `json:"activity_date" bson:"activity_date"`
	// CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	// UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
	// ActivityType string    `json:"activity_type" bson:"activity_type"`
}

// ActivityInfo info data
type ActivityInfo struct {
	ID           primitive.ObjectID `form:"id" json:"id" bson:"_id,omitempty"`
	Distance     float64            `form:"distance" json:"distance" bson:"distance"`
	ImageURL     string             `form:"img_url" json:"img_url" bson:"img_url"`
	Caption      string             `form:"caption" json:"caption" bson:"caption"`
	APP          string             `form:"app" json:"app" bson:"app"`
	Time         int64              `form:"time" json:"time" bson:"time"`
	IsApprove    bool               `form:"is_approve" json:"is_approve" bson:"is_approve"`
	Status       string             `form:"status" json:"status" bson:"status"`
	Reason       string             `form:"reason" json:"reason" bson:"reason"`
	ActivityDate time.Time          `form:"activity_date" json:"activity_date" bson:"activity_date"`
	CreatedAt    time.Time          `form:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `form:"updated_at" json:"updated_at" bson:"updated_at"`
}

type HistoryDayFilter struct {
	Year    int    `json:"year" bson:"year"`
	Month   int    `form:"month" json:"month"`
	EventID string `form:"event_id" json:"event_id"`
}

type HistoryMonthFilter struct {
	Year    int    `json:"year" bson:"year"`
	EventID string `form:"event_id" json:"event_id"`
}

type HistoryDayInfo struct {
	Year            int               `form:"year" json:"year"`
	Month           int               `form:"month" json:"month"`
	DistanceDayInfo []DistanceDayInfo `form:"distance_info" json:"distance_info"`
}

type HistoryMonthInfo struct {
	Month         int     `form:"month" json:"month"`
	MonthName     string  `form:"month_name" json:"month_name"`
	TotalDistance float64 `form:"total_distance" json:"total_distance"`
}

type DistanceDayInfo struct {
	Distance     float64 `json:"distance" bson:"distance"`
	ActivityDate string  `form:"activity_date" json:"activity_date" bson:"distance"`
}

type DeleteActivityForm struct {
	EventID    string `form:"event_id" json:"event_id" bson:"event_id"  binding:"required"`
	UserID     string `form:"user_id" json:"user_id" bson:"user_id"`
	ActivityID string `form:"activity_id" json:"activity_id" bson:"activity_id"`
}

// ActivityAllInfo for request get all activity
type ActivityAllInfo struct {
	UserInfo User     `json:"user_info"`
	Activity Activity `json:"activity"`
}
