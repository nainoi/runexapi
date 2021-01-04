package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EventLists for event data list
type EventLists struct {
	Events []EventList `json:"events"`
}

// EventList for event data list
type EventList struct {
	Code           string    `json:"code"`
	Content        string    `json:"content"`
	Cover          string    `json:"cover"`
	EventDate      time.Time `json:"event_date"`
	EventEndDate   time.Time `json:"event_end_date"`
	EventStartDate time.Time `json:"event_start_date"`
	Title          string    `json:"title"`
}

// EventData for event data list
type EventData struct {
	Event struct {
		Agreement           string    `json:"agreement"`
		Category            string    `json:"category"`
		Code                string    `json:"code"`
		Contact             string    `json:"contact"`
		ContactFacebook     string    `json:"contactFacebook"`
		ContactLine         string    `json:"contactLine"`
		Content             string    `json:"content"`
		Cover               string    `json:"cover"`
		CoverThumbnail      string    `json:"coverThumbnail"`
		EventDate           string    `json:"eventDate"`
		EventEndDate        time.Time `json:"eventEndDate"`
		EventEndDateText    string    `json:"eventEndDateText"`
		EventStartDate      time.Time `json:"eventStartDate"`
		EventStartDateText  string    `json:"eventStartDateText"`
		ID                  int       `json:"id"`
		IsFreeEvent         bool      `json:"isFreeEvent"`
		IsRunexOnly         bool      `json:"isRunexOnly"`
		IsSendShirtByPost   bool      `json:"isSendShirtByPost"`
		Organizer           string    `json:"organizer"`
		PhotoBib            string    `json:"photoBib"`
		PhotoBibThumbnail   string    `json:"photoBibThumbnail"`
		PhotoCert           string    `json:"photoCert"`
		PhotoCertThumbnail  string    `json:"photoCertThumbnail"`
		PhotoMedal          string    `json:"photoMedal"`
		PhotoMedalThumbnail string    `json:"photoMedalThumbnail"`
		PhotoShirt          string    `json:"photoShirt"`
		PhotoShirtThumbnail string    `json:"photoShirtThumbnail"`
		Place               string    `json:"place"`
		Prizes              []struct {
			Description string `json:"description"`
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Photo       string `json:"photo"`
		} `json:"prizes"`
		RegisterEndDate       time.Time `json:"registerEndDate"`
		RegisterEndDateText   string    `json:"registerEndDateText"`
		RegisterStartDate     time.Time `json:"registerStartDate"`
		RegisterStartDateText string    `json:"registerStartDateText"`
		Schedules             []struct {
			Description string `json:"description"`
			ID          int    `json:"id"`
			Name        string `json:"name"`
		} `json:"schedules"`
		Shirts []interface{} `json:"shirts"`
		Title  string        `json:"title"`
		UserID string        `json:"userId"`
	} `json:"event"`
	Tickets []struct {
		Category   string      `json:"category"`
		CreatedAt  time.Time   `json:"created_at"`
		Detail     interface{} `json:"detail"`
		Distance   int         `json:"distance"`
		EventID    string      `json:"event_id"`
		ID         string      `json:"id"`
		Items      interface{} `json:"items"`
		Limit      int         `json:"limit"`
		PhotoMap   string      `json:"photo_map"`
		PhotoMedal string      `json:"photo_medal"`
		PhotoShirt string      `json:"photo_shirt"`
		Price      int         `json:"price"`
		Title      string      `json:"title"`
		UpdatedAt  time.Time   `json:"updated_at"`
	} `json:"tickets"`
}

// EventV2 for event data
type EventV2 struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name            string             `json:"name" bson:"name"`
	Description     string             `json:"description" bson:"description"`
	Body            string             `json:"body" bson:"body"`
	Cover           string             `json:"cover" bson:"cover"`
	CoverThumb      []CoverThumb       `json:"cover_thumb" bson:"cover_thumb"`
	Category        string             `json:"category" bson:"category"`
	Slug            string             `json:"slug" bson:"slug"`
	Ticket          []TicketEventV2    `json:"ticket" bson:"ticket"`
	OwnerID         primitive.ObjectID `json:"owner_id" bson:"owner_id"`
	Status          string             `json:"status" bson:"status"`
	Location        string             `json:"location" bson:"location"`
	ReceiveLocation string             `json:"receive_location" bson:"receive_location"`
	IsActive        bool               `json:"is_active" bson:"is_active"`
	IsFree          bool               `json:"is_free" bson:"is_free"`
	StartReg        time.Time          `json:"start_reg" bson:"start_reg"`
	EndReg          time.Time          `json:"end_reg" bson:"end_reg"`
	StartEvent      time.Time          `json:"start_event" bson:"start_event"`
	EndEvent        time.Time          `json:"end_event" bson:"end_event"`
	Inapp           bool               `json:"inapp" bson:"inapp"`
	IsPost          bool               `json:"is_post" bson:"is_post"`
	PostEndDate     time.Time          `json:"post_end_date" bson:"post_end_date"`
	Partner         PartnerEvent       `json:"partner" bson:"partner"`
	CreatedTime     time.Time          `json:"created_time" bson:"created_time"`
	UpdatedTime     time.Time          `json:"updated_time" bson:"updated_time"`
}

type EventRegV2 struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name            string             `json:"name" bson:"name"`
	Cover           string             `json:"cover" bson:"cover"`
	CoverThumb      []CoverThumb       `json:"cover_thumb" bson:"cover_thumb"`
	Category        string             `json:"category" bson:"category"`
	Slug            string             `json:"slug" bson:"slug"`
	OwnerID         primitive.ObjectID `json:"owner_id" bson:"owner_id"`
	Status          string             `json:"status" bson:"status"`
	Location        string             `json:"location" bson:"location"`
	ReceiveLocation string             `json:"receive_location" bson:"receive_location"`
	IsActive        bool               `json:"is_active" bson:"is_active"`
	IsFree          bool               `json:"is_free" bson:"is_free"`
	StartReg        time.Time          `json:"start_reg" bson:"start_reg"`
	EndReg          time.Time          `json:"end_reg" bson:"end_reg"`
	StartEvent      time.Time          `json:"start_event" bson:"start_event"`
	EndEvent        time.Time          `json:"end_event" bson:"end_event"`
	Inapp           bool               `json:"inapp" bson:"inapp"`
	IsPost          bool               `json:"is_post" bson:"is_post"`
	PostEndDate     time.Time          `json:"post_end_date" bson:"post_end_date"`
	Partner         PartnerEvent       `json:"partner" bson:"partner"`
	CreatedTime     time.Time          `json:"created_time" bson:"created_time"`
	UpdatedTime     time.Time          `json:"updated_time" bson:"updated_time"`
}

type CoverThumb struct {
	Image string `json:"image" bson:"image"`
	Size  string `json:"size" bson:"size"`
}

type TicketEventV2 struct {
	TicketID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Price       float64            `json:"price" bson:"price"`
	Description string             `json:"description" bson:"description"`
	Currency    string             `json:"currency" bson:"currency"`
	TicketType  string             `json:"ticket_type" bson:"ticket_type"`
	Team        int                `json:"team" bson:"team"`
	Quantity    int                `json:"quantity" bson:"quantity"`
	Distance    float64            `json:"distance" bson:"distance"`
	Products    ProduceEventV2     `json:"products" bson:"products"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type ProduceEventV2 struct {
	ProductID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name" binding:"required"`
	Image     []ProductImage     `json:"image" bson:"image"`
	Detail    string             `json:"detail" bson:"detail"`
	Status    string             `json:"status" bson:"status"`
	Reuse     bool               `json:"reuse" bson:"reuse"`
	IsShow    bool               `json:"is_show" bson:"is_show"`
	Sizes     []ProductSizes     `json:"sizes" bson:"sizes"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type ProductSizes struct {
	Name   string `json:"name" bson:"name"`
	Remark string `json:"remark" bson:"remark"`
}

// PartnerEvent struct
type PartnerEvent struct {
	PartnerID        primitive.ObjectID `form:"partner_id" json:"partner_id" bson:"partner_id"`
	PartnerName      string             `form:"partner_name" json:"partner_name" bson:"partner_name"`
	Slug             string             `form:"slug" json:"slug" bson:"slug"`
	RefEventKey      string             `form:"ref_event_key" json:"ref_event_key" bson:"ref_event_key"`
	RefActivityKey   string             `form:"ref_activity_key" json:"ref_activity_key" bson:"ref_activity_key"`
	RefEventValue    string             `form:"ref_event_value" json:"ref_event_value" bson:"ref_event_value"`
	RefActivityValue string             `form:"ref_activity_value" json:"ref_activity_value" bson:"ref_activity_value"`
	RefPhoneValue    string             `form:"ref_phone_value" json:"ref_phone_value" bson:"ref_phone_value"`
}
