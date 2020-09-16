package model

//StravaRecieve model
//model for strava hook activity id
type StravaRecieve struct {
	StravaID         string               `json:"-"`
	ActivityID       string               `json:"-"`
}