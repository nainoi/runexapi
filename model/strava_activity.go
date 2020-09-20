package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StravaActivity struct data activity from strava
type StravaActivity struct {
	IsSync                     bool              `json:"is_sync" bson:"is_sync"`
	CreatedAt                  time.Time         `json:"created_at" bson:"created_at"`
	ResourceState              int               `json:"resource_state" bson:"resource_state"`
	Athlete                    Athlete           `json:"athlete" bson:"athlete"`
	Name                       string            `json:"name" bson:"name"`
	Distance                   int               `json:"distance" bson:"distance"`
	MovingTime                 int               `json:"moving_time" bson:"moving_time"`
	ElapsedTime                int               `json:"elapsed_time" bson:"elapsed_time"`
	TotalElevationGain         int               `json:"total_elevation_gain" bson:"total_elevation_gain"`
	Type                       string            `json:"type" bson:"type"`
	WorkoutType                int               `json:"workout_type" bson:"workout_type"`
	ID                         int64             `json:"id" bson:"id"`
	ExternalID                 string            `json:"external_id" bson:"external_id"`
	UploadID                   int64             `json:"upload_id" bson:"upload_id"`
	StartDate                  time.Time         `json:"start_date" bson:"start_date"`
	StartDateLocal             time.Time         `json:"start_date_local" bson:"start_date_local"`
	Timezone                   string            `json:"timezone" bson:"timezone"`
	UtcOffset                  int               `json:"utc_offset" bson:"utc_offset"`
	StartLatlng                []float64         `json:"start_latlng" bson:"start_latlng"`
	EndLatlng                  []float64         `json:"end_latlng" bson:"end_latlng"`
	LocationCity               string            `json:"location_city" bson:"location_city"`
	LocationState              string            `json:"location_state" bson:"location_state"`
	LocationCountry            string            `json:"location_country" bson:"location_country"`
	StartLatitude              float64           `json:"start_latitude" bson:"start_latitude"`
	StartLongitude             float64           `json:"start_longitude" bson:"start_longitude"`
	AchievementCount           int               `json:"achievement_count" bson:"achievement_count"`
	KudosCount                 int               `json:"kudos_count" bson:"kudos_count"`
	CommentCount               int               `json:"comment_count" bson:"comment_count"`
	AthleteCount               int               `json:"athlete_count" bson:"athlete_count"`
	PhotoCount                 int               `json:"photo_count" bson:"photo_count"`
	Map                        Map               `json:"map" bson:"map"`
	Trainer                    bool              `json:"trainer" bson:"trainer"`
	Commute                    bool              `json:"commute" bson:"commute"`
	Manual                     bool              `json:"manual" bson:"manual"`
	Private                    bool              `json:"private" bson:"private"`
	Visibility                 string            `json:"visibility" bson:"visibility"`
	Flagged                    bool              `json:"flagged" bson:"flagged"`
	GearID                     string            `json:"gear_id" bson:"gear_id"`
	FromAcceptedTag            bool              `json:"from_accepted_tag" bson:"from_accepted_tag"`
	UploadIDStr                string            `json:"upload_id_str" bson:"upload_id_str"`
	AverageSpeed               int               `json:"average_speed" bson:"average_speed"`
	MaxSpeed                   int               `json:"max_speed" bson:"max_speed"`
	HasHeartrate               bool              `json:"has_heartrate" bson:"has_heartrate"`
	HeartrateOptOut            bool              `json:"heartrate_opt_out" bson:"heartrate_opt_out"`
	DisplayHideHeartrateOption bool              `json:"display_hide_heartrate_option" bson:"display_hide_heartrate_option"`
	ElevHigh                   int               `json:"elev_high" bson:"elev_high"`
	ElevLow                    int               `json:"elev_low" bson:"elev_low"`
	PrCount                    int               `json:"pr_count" bson:"pr_count"`
	TotalPhotoCount            int               `json:"total_photo_count" bson:"total_photo_count"`
	HasKudoed                  bool              `json:"has_kudoed" bson:"has_kudoed"`
	Description                string            `json:"description" bson:"description"`
	Calories                   int               `json:"calories" bson:"calories"`
	PerceivedExertion          string            `json:"perceived_exertion" bson:"perceived_exertion"`
	PreferPerceivedExertion    bool              `json:"prefer_perceived_exertion" bson:"prefer_perceived_exertion"`
	SegmentEfforts             []string          `json:"segment_efforts" bson:"segment_efforts"`
	SplitsMetric               []SplitsMetric    `json:"splits_metric" bson:"splits_metric"`
	SplitsStandard             []SplitsStandard  `json:"splits_standard" bson:"splits_standard"`
	Laps                       []string          `json:"laps" bson:"laps"`
	BestEfforts                []string          `json:"best_efforts" bson:"best_efforts"`
	Photos                     Photos            `json:"photos" bson:"photos"`
	DeviceName                 string            `json:"device_name" bson:"device_name"`
	EmbedToken                 string            `json:"embed_token" bson:"embed_token"`
	SimilarActivities          SimilarActivities `json:"similar_activities" bson:"similar_activities"`
	AvailableZones             []string          `json:"available_zones" bson:"available_zones"`
}

// Athlete struct from strava
type Athlete struct {
	ID            int `json:"id" bson:"id"`
	ResourceState int `json:"resource_state" bson:"resource_state"`
}

// Map struct from strava
type Map struct {
	ID              string `json:"id" bson:"id"`
	Polyline        string `json:"polyline" bson:"polyline"`
	ResourceState   int    `json:"resource_state" bson:"resource_state"`
	SummaryPolyline string `json:"summary_polyline" bson:"summary_polyline"`
}

// SplitsMetric struct from strava
type SplitsMetric struct {
	Distance                  int `json:"distance" bson:"distance"`
	ElapsedTime               int `json:"elapsed_time" bson:"elapsed_time"`
	ElevationDifference       int `json:"elevation_difference" bson:"elevation_difference"`
	MovingTime                int `json:"moving_time" bson:"moving_time"`
	Split                     int `json:"split" bson:"split"`
	AverageSpeed              int `json:"average_speed" bson:"average_speed"`
	AverageGradeAdjustedSpeed int `json:"average_grade_adjusted_speed" bson:"average_grade_adjusted_speed"`
	PaceZone                  int `json:"pace_zone" bson:"pace_zone"`
}

// SplitsStandard struct from strava
type SplitsStandard struct {
	Distance                  int `json:"distance" bson:"distance"`
	ElapsedTime               int `json:"elapsed_time" bson:"elapsed_time"`
	ElevationDifference       int `json:"elevation_difference" bson:"elevation_difference"`
	MovingTime                int `json:"moving_time" bson:"moving_time"`
	Split                     int `json:"split" bson:"split"`
	AverageSpeed              int `json:"average_speed" bson:"average_speed"`
	AverageGradeAdjustedSpeed int `json:"average_grade_adjusted_speed" bson:"average_grade_adjusted_speed"`
	PaceZone                  int `json:"pace_zone" bson:"pace_zone"`
}

// Photos struct from strava
type Photos struct {
	Primary string `json:"primary" bson:"primary"`
	Count   int    `json:"count" bson:"count"`
}

// SimilarActivities struct from strava
type SimilarActivities struct {
	EffortCount        int    `json:"effort_count" bson:"effort_count"`
	AverageSpeed       int    `json:"average_speed" bson:"average_speed"`
	MinAverageSpeed    int    `json:"min_average_speed" bson:"min_average_speed"`
	MidAverageSpeed    int    `json:"mid_average_speed" bson:"mid_average_speed"`
	MaxAverageSpeed    int    `json:"max_average_speed" bson:"max_average_speed"`
	PrRank             string `json:"pr_rank" bson:"pr_rank"`
	FrequencyMilestone string `json:"frequency_milestone" bson:"frequency_milestone"`
	Trend              Trend  `json:"trend" bson:"trend"`
	ResourceState      int    `json:"resource_state" bson:"resource_state"`
}

// Trend struct from strava
type Trend struct {
	Speeds               []string `json:"speeds" bson:"speeds"`
	CurrentActivityIndex string   `json:"current_activity_index" bson:"current_activity_index"`
	MinSpeed             int      `json:"min_speed" bson:"min_speed"`
	MidSpeed             int      `json:"mid_speed" bson:"mid_speed"`
	MaxSpeed             int      `json:"max_speed" bson:"max_speed"`
	Direction            int      `json:"direction" bson:"direction"`
}

// StravaAddRequest struct for request add strava activity
type StravaAddRequest struct {
	UserID         primitive.ObjectID `json:"user_id"`
	StravaActivity StravaActivity     `json:"strava_activity"`
}

// StravaData struct for request add strava activity
type StravaData struct {
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	Activities []StravaActivity   `json:"activities" bson:"activities"`
}
