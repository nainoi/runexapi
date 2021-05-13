package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//AddActivityForm2 model request
type AddActivityForm2 struct {
	EventCode    string       `form:"event_code" json:"event_code" bson:"event_code" binding:"required"`
	ParentRegID  string       `form:"parent_reg_id" json:"parent_reg_id" bson:"parent_reg_id"`
	RegID        string       `form:"reg_id" json:"reg_id" bson:"reg_id" binding:"required"`
	OrderID      string       `form:"order_id" json:"order_id" bson:"order_id" binding:"required"`
	Ticket       Tickets      `form:"ticket" json:"ticket" bson:"ticket" binding:"required"`
	ActivityInfo ActivityInfo `form:"activity_info" json:"activity_info" bson:"activity_info"`
}

// ActivityV2 data per event
type ActivityV2 struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EventCode     string             `json:"event_code" bson:"event_code"`
	Ticket        Tickets            `json:"ticket" bson:"ticket"`
	OrderID       string             `json:"order_id" bson:"order_id"`
	RegID         primitive.ObjectID `json:"reg_id" bson:"reg_id"`
	ParentRegID   primitive.ObjectID `json:"parent_reg_id" bson:"parent_reg_id"`
	UserID        primitive.ObjectID `json:"user_id" bson:"user_id"`
	ToTalDistance float64            `json:"total_distance" bson:"total_distance"`
	UserInfo      UserOption         `json:"user_info" bson:"user_info"`
	ActivityInfo  []ActivityInfo     `json:"activity_info" bson:"activity_info,omitempty"`
	// Activities  Activities         `json:"activities" bson:"activities,omitempty"`
}

//ActivityDashboard struct
type ActivityDashboard struct {
	Activity     []ActivityV2 `json:"activities" bson:"activities"`
	RegisterData RegisterV2   `json:"register" bson:"register"`
}

// Activities data per user
type Activities struct {
	UserID        primitive.ObjectID `json:"user_id" bson:"user_id"`
	ToTalDistance float64            `json:"total_distance" bson:"total_distance"`
	ActivityInfo  []ActivityInfo     `json:"activity_info" bson:"activity_info,omitempty"`
}

//AddActivityV2 data
type AddActivityV2 struct {
	ID           primitive.ObjectID `form:"id" json:"id" bson:"_id,omitempty"`
	UserID       primitive.ObjectID `form:"user_id" json:"user_id" bson:"user_id"`
	EventCode    string             `form:"event_code" json:"event_code" bson:"event_code" binding:"required"`
	ParentRegID  primitive.ObjectID `form:"parent_reg_id" json:"parent_reg_id" bson:"parent_reg_id"`
	RegID        primitive.ObjectID `form:"reg_id" json:"reg_id" bson:"reg_id" binding:"required"`
	OrderID      string             `form:"order_id" json:"order_id" bson:"order_id" binding:"required"`
	Ticket       Tickets            `form:"ticket" json:"ticket" bson:"ticket" binding:"required"`
	ActivityInfo ActivityInfo       `form:"activity_info" json:"activity_info" bson:"activity_info"`
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
	EventID             primitive.ObjectID  `form:"event_id" json:"event_id"`
	EventCode           string              `json:"event_code" bson:"event_code"`
	// Partner             PartnerEvent        `json:"partner" bson:"partner"`
	ParentRegID string  `form:"parent_reg_id" json:"parent_reg_id" bson:"parent_reg_id"`
	RegID       string  `form:"reg_id" json:"reg_id" bson:"reg_id" binding:"required"`
	OrderID     string  `form:"order_id" json:"order_id" bson:"order_id" binding:"required"`
	Ticket      Tickets `form:"ticket" json:"ticket" bson:"ticket" binding:"required"`
}

//AddMultiActivityFormWorkout model request
type AddMultiActivityFormWorkout struct {
	WorkoutActivityInfo WorkoutActivityInfo `form:"workout_info" json:"workout_info"`
	EventActivity       []EventActivity     `form:"event_activity" json:"event_activity" binding:"required"`
}

// EventActivity for request add activity
type EventActivity struct {
	// EventID   string       `json:"event_id" bson:"event_id"`
	EventCode   string  `json:"event_code" bson:"event_code" binding:"required"`
	ParentRegID string  `form:"parent_reg_id" json:"parent_reg_id" bson:"parent_reg_id"`
	RegID       string  `form:"reg_id" json:"reg_id" bson:"reg_id" binding:"required"`
	OrderID     string  `form:"order_id" json:"order_id" bson:"order_id" binding:"required"`
	Ticket      Tickets `form:"ticket" json:"ticket" bson:"ticket" binding:"required"`
}

type EventActivityDashboardReq struct {
	// EventID   string       `json:"event_id" bson:"event_id"`
	EventCode   string             `form:"event_code" json:"event_code" bson:"event_code" binding:"required"`
	ParentRegID primitive.ObjectID `form:"parent_reg_id" json:"parent_reg_id" bson:"parent_reg_id"`
	RegID       primitive.ObjectID `form:"reg_id" json:"reg_id" bson:"reg_id" binding:"required"`
	OrderID     string             `form:"order_id" json:"order_id" bson:"order_id" binding:"required"`
}

type EventActivityRemoveReq struct {
	// EventID   string       `json:"event_id" bson:"event_id"`
	EventCode    string             `form:"event_code" json:"event_code" bson:"event_code" binding:"required"`
	TicketD      string             `form:"ticket_id" json:"ticket_id" bson:"ticket_id"`
	RegID        primitive.ObjectID `form:"reg_id" json:"reg_id" bson:"reg_id" binding:"required"`
	OrderID      string             `form:"order_id" json:"order_id" bson:"order_id" binding:"required"`
	ActivityInfo ActivityInfo       `form:"activity_info" json:"activity_info" bson:"activity_info" binding:"required"`
}

type UpdateActivityReq struct {
	// EventID   string       `json:"event_id" bson:"event_id"`
	EventCode  string             `form:"event_code" json:"event_code" bson:"event_code" binding:"required"`
	RegID      primitive.ObjectID `form:"reg_id" json:"reg_id" bson:"reg_id" binding:"required"`
	ActivityID primitive.ObjectID `form:"act_id" json:"act_id" bson:"act_id" binding:"required"`
	UserID     primitive.ObjectID `form:"user_id" json:"user_id" bson:"user_id" binding:"required"`
	OrderID    string             `form:"order_id" json:"order_id" bson:"order_id" binding:"required"`
	Status     string             `form:"status" json:"status" bson:"status" binding:"required"`
	Distance   float64            `form:"distance" json:"distance" bson:"distance" binding:"required"`
	Reason     string             `form:"reason" json:"reason" bson:"reason"`
}
