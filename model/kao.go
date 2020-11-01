package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SendActivityRequest struct for sent activity to kao
type SendActivityRequest struct {
	Distance float64 `json:"distance" bson:"distance"`
	Time     int32   `json:"time" bson:"time"`
}

// GetKaoActivityRequest struct for get activity from kao
type GetKaoActivityRequest struct {
	Slug string `json:"slug" bson:"slug"`
	EBIB string `json:"ebib" bson:"ebib"`
}

// Kao struct model
type Kao struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	KaoObjects []KaoObject        `json:"kao_info" bson:"kao_info"`
}

// KaoObject struct
type KaoObject struct {
	ID             int `json:"id"`
	Position       int `json:"position"`
	EventProductID int `json:"eventProductId"`
	EventProduct   struct {
		Name          string    `json:"name"`
		ImageURL      string    `json:"imageUrl"`
		AddonTitle    string    `json:"addonTitle"`
		EditableUntil time.Time `json:"editableUntil"`
		HasETicket    bool      `json:"hasETicket"`
	} `json:"eventProduct"`
	Price                  int    `json:"price"`
	Code                   string `json:"code"`
	HolderName             string `json:"holderName"`
	HolderPhone            string `json:"holderPhone"`
	HolderEmail            string `json:"holderEmail"`
	ETicketURL             string `json:"eTicketUrl"`
	VirtualRaceSubmissions []struct {
		ID                    int         `json:"id"`
		OrderItemID           int         `json:"orderItemId"`
		Status                string      `json:"status"`
		CreatedAt             time.Time   `json:"createdAt"`
		SubmitterMobileNumber interface{} `json:"submitterMobileNumber"`
		ImageURL              string      `json:"imageUrl"`
		Distance              float64     `json:"distance"`
		Time                  int         `json:"time"`
		RejectionReason       string      `json:"rejectionReason"`
	} `json:"virtualRaceSubmissions"`
	VirtualRaceProfile struct {
		OrderItemID     int     `json:"orderItemId"`
		HolderName      string  `json:"holderName"`
		Code            string  `json:"code"`
		ETicketURL      string  `json:"eTicketUrl"`
		Rank            int     `json:"rank"`
		Distance        float64 `json:"distance"`
		Time            int     `json:"time"`
		SubmissionCount int     `json:"submissionCount"`
	} `json:"virtualRaceProfile"`
	DailyDistance float64 `json:"dailyDistance"`
	LastRank      int     `json:"lastRank"`
}
