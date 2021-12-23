package model

import (
	//"thinkdev.app/think/runex/runexapi/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Ranking model
type Ranking struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID        primitive.ObjectID `json:"user_id" bson:"user_id"`
	RegID         primitive.ObjectID `json:"reg_id" bson:"reg_id"`
	ParentRegID   primitive.ObjectID `json:"parent_reg_id" bson:"parent_reg_id"`
	EventID       int                `json:"event_id" bson:"event_id"`
	EventCode     string             `json:"event_code" bson:"event_code"`
	TicketID      string             `json:"ticket_id" bson:"ticket_id"`
	ActivityInfo  []ActivityInfo     `json:"activity_info" bson:"activity_info,omitempty"`
	ToTalDistance float64            `json:"total_distance" bson:"total_distance"`
	RankNo        int                `json:"rank_no"`
	UserInfo      UserOption         `json:"user_info"`
	ImageURL      string             `json:"img_url" bson:"img_url"`
	BibNo         string             `json:"bib_no" bson:"bib_no"`
	Teams         []UserOption       `json:"teams"`
	// ActivityDate time.Time `json:"activity_date" bson:"activity_date"`
	// CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	// UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
	// ActivityType string    `json:"activity_type" bson:"activity_type"`
}

//Ranking model
type RankingFinish struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID        primitive.ObjectID `json:"user_id" bson:"user_id"`
	RegID         primitive.ObjectID `json:"reg_id" bson:"reg_id"`
	ParentRegID   primitive.ObjectID `json:"parent_reg_id" bson:"parent_reg_id"`
	EventID       int                `json:"event_id" bson:"event_id"`
	EventCode     string             `json:"event_code" bson:"event_code"`
	TicketID      string             `json:"ticket_id" bson:"ticket_id"`
	ActivityInfo  []ActivityInfo     `json:"activity_info" bson:"activity_info,omitempty"`
	ToTalDistance float64            `json:"total_distance" bson:"total_distance"`
	RankNo        int                `json:"rank_no"`
	UserInfo      UserOption         `json:"user_info"`
	ImageURL      string             `json:"img_url" bson:"img_url"`
	BibNo         string             `json:"bib_no" bson:"bib_no"`
	Teams         []UserOption       `json:"teams"`
	DateFinished  time.Time          `json:"date_finished" bson:"date_finished"`
	// ActivityDate time.Time `json:"activity_date" bson:"activity_date"`
	// CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	// UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
	// ActivityType string    `json:"activity_type" bson:"activity_type"`
}

type RankingRequest struct {
	EventCode   string             `json:"event_code" bson:"event_code" binding:"required"`
	TicketID    string             `json:"ticket_id" bson:"ticket_id" binding:"required"`
	ParentRegID primitive.ObjectID `form:"parent_reg_id" json:"parent_reg_id"`
}

type AllRankingRequest struct {
	EventCode string `json:"event_code" bson:"event_code" binding:"required"`
	TicketID  string `json:"ticket_id" bson:"ticket_id" binding:"required"`
}
