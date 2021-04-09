package model

import (
	//"thinkdev.app/think/runex/runexapi/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Ranking model
type Ranking struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID        primitive.ObjectID `json:"user_id" bson:"user_id"`
	EventID       int                `json:"event_id" bson:"event_id"`
	EventCode     string             `json:"event_code" bson:"event_code"`
	EventUser     string             `json:"event_user" bson:"event_user"`
	ActivityInfo  []ActivityInfo     `json:"activity_info" bson:"activity_info,omitempty"`
	ToTalDistance float64            `json:"total_distance" bson:"total_distance"`
	RankNo        int                `json:"rank_no"`
	UserInfo      UserEvent          `json:"user_info"`
	// Distance     float32   `json:"distance" bson:"distance"`
	// ImageURL     string    `json:"img_url" bson:"img_url"`
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
