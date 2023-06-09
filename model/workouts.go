package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Workouts struct model
type Workouts struct {
	ID                  primitive.ObjectID    `json:"id" bson:"_id,omitempty"`
	UserID              primitive.ObjectID    `json:"user_id" bson:"user_id"`
	WorkoutActivityInfo []WorkoutActivityInfo `json:"activity_info" bson:"activity_info,omitempty"`
	TotalDistance       float64               `json:"total_distance" bson:"total_distance"`
}

// WorkoutActivityInfo struct model
type WorkoutActivityInfo struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ActivityType     string             `json:"activity_type" json:"activity_type"`
	APP              string             `json:"app" bson:"app"`
	RefID            string             `json:"ref_id" bson:"ref_id"`
	Calory           float64            `json:"calory" bson:"calory"`
	Caption          string             `json:"caption" json:"caption"`
	Distance         float64            `json:"distance" bson:"distance"`
	Pace             float64            `json:"pace" bson:"pace"`
	Duration         int64              `json:"duration" bson:"duration"`
	TimeString       string             `json:"time_string" json:"time_string"`
	StartDate        time.Time          `json:"start_date" bson:"start_date"`
	EndDate          time.Time          `json:"end_date" bson:"end_date"`
	WorkoutDate      time.Time          `json:"workout_date" bson:"workout_date"`
	NetElevationGain float64            `json:"net_elevation_gain" bson:"net_elevation_gain"`
	IsSync           bool               `json:"is_sync" bson:"is_sync"`
	Locations        []Location         `json:"locations" bson:"locations"`
	LocURL           string             `json:"loc_url" bson:"loc_url"`
}

// AddWorkout struct model for request
type AddWorkout struct {
	UserID              primitive.ObjectID  `json:"user_id" bson:"user_id"`
	WorkoutActivityInfo WorkoutActivityInfo `form:"workout_info" json:"workout_info"`
}

// HookWorkout struct model for request
type HookWorkout struct {
	Provider            string                `form:"provider" json:"provider" binding:"required"`
	ProviderID          string                `form:"provider_id" json:"provider_id" binding:"required"`
	WorkoutActivityInfo []WorkoutActivityInfo `form:"workout_info" json:"workout_info"`
}

// AddMultiWorkout struct model for request
type AddMultiWorkout struct {
	WorkoutActivityInfos []WorkoutActivityInfo `form:"workouts" json:"workouts"`
}

// AddWorkoutForm struct model for request
type AddWorkoutForm struct {
	UserID           string     `form:"user_id" json:"user_id" bson:"user_id"`
	ActivityType     string     `json:"activity_type" json:"activity_type"`
	Calory           float64    `json:"calory" bson:"calory"`
	APP              string     `json:"app" bson:"app"`
	Caption          string     `json:"caption" json:"caption"`
	Distance         float64    `json:"distance" bson:"distance"`
	Pace             float64    `json:"pace" bson:"pace"`
	Duration         int64      `json:"duration" bson:"duration"`
	TimeString       string     `json:"time_string" json:"time_string"`
	StartDate        string     `json:"start_date" bson:"start_date"`
	EndDate          string     `json:"end_date" bson:"end_date"`
	WorkoutDate      string     `json:"workout_date" bson:"workout_date"`
	NetElevationGain float64    `json:"net_elevation_gain" bson:"net_elevation_gain"`
	IsSync           bool       `json:"is_sync" bson:"is_sync"`
	Locations        []Location `json:"locations" bson:"locations"`
}

// WorkoutHookForm struct model for request
type WorkoutHookForm struct {
	ActivityType     string     `json:"activity_type" json:"activity_type"`
	Calory           float64    `json:"calory" bson:"calory"`
	APP              string     `json:"app" bson:"app"`
	Caption          string     `json:"caption" json:"caption"`
	Distance         float64    `json:"distance" bson:"distance"`
	Pace             float64    `json:"pace" bson:"pace"`
	Duration         int64      `json:"duration" bson:"duration"`
	TimeString       string     `json:"time_string" json:"time_string"`
	StartDate        string     `json:"start_date" bson:"start_date"`
	EndDate          string     `json:"end_date" bson:"end_date"`
	WorkoutDate      string     `json:"workout_date" bson:"workout_date"`
	NetElevationGain float64    `json:"net_elevation_gain" bson:"net_elevation_gain"`
	IsSync           bool       `json:"is_sync" bson:"is_sync"`
	Locations        []Location `json:"locations" bson:"locations"`
}

// Location struct model
type Location struct {
	Timestamp     time.Time `form:"timestamp" json:"timestamp" bson:"timestamp"`
	Altitude      float64   `form:"altitude" json:"altitude" bson:"altitude"`
	Latitude      float64   `form:"latitude" json:"latitude" bson:"latitude"`
	Longitude     float64   `form:"longitude" json:"longitude" bson:"longitude"`
	Temp          float64   `form:"temp" json:"temp" bson:"temp"`
	HarthRate     float64   `form:"harth_rate" json:"harth_rate" bson:"harth_rate"`
	ElevationGain float64   `form:"elevation_gain" json:"elevation_gain" bson:"elevation_gain"`
}

// WorkoutHistoryDayFilter filter day
type WorkoutHistoryDayFilter struct {
	Year  int `json:"year" bson:"year"`
	Month int `json:"month" bson:"month"`
}

//WorkoutHistoryMonthFilter filter struct
type WorkoutHistoryMonthFilter struct {
	Year int `form:"year" json:"year"`
}

// WorkoutHistoryMonthInfo struct per month
type WorkoutHistoryMonthInfo struct {
	Month         int                   `form:"month" json:"month"`
	MonthName     string                `form:"month_name" json:"month_name"`
	TotalDistance float64               `form:"total_distance" json:"total_distance"`
	TimeString    string                `json:"time_string" json:"time_string"`
	Calory        float64               `json:"calory" bson:"calory"`
	HistoryDay    []WorkoutActivityInfo `json:"workout_day" bson:"workout_day"`
}

// WorkoutHistoryDayInfo struct per day
type WorkoutHistoryDayInfo struct {
	Year            int                   `form:"year" json:"year"`
	Month           int                   `form:"month" json:"month"`
	DistanceDayInfo []WorkoutActivityInfo `form:"distance_info" json:"distance_info"`
}

//WorkoutDistanceDayInfo struct
type WorkoutDistanceDayInfo struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Distance    float64            `json:"distance" bson:"distance"`
	WorkoutDate string             `json:"workout_date" bson:"workout_date"`
	WorkoutTime string             `json:"workout_time" json:"workout_time"`
}

//WorkoutHistoryAllInfo all
type WorkoutHistoryAllInfo struct {
	Year          int                   `form:"year" json:"year" bson:"year"`
	Month         int                   `form:"month" json:"month" bson:"month"`
	MonthName     string                `form:"month_name" json:"month_name" bson:"month_name"`
	TotalDistance float64               `form:"total_distance" json:"total_distance" bson:"total_distance"`
	TimeString    string                `json:"time_string" json:"time_string" bson:"time_string"`
	Calory        float64               `json:"calory" bson:"calory" bson:"calory"`
	HistoryDay    []WorkoutActivityInfo `json:"workout_day" bson:"workout_day" bson:"workout_day"`
}

type RemoveWorkoutReq struct {
	// ID                  primitive.ObjectID  `json:"id" bson:"_id,omitempty" binding:"required"`
	WorkoutActivityInfo WorkoutActivityInfo `json:"activity_info" bson:"activity_info" binding:"required"`
}
