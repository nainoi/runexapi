package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
