package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RegisterRequest struct for request register event
type RegisterRequest struct {
	EventID    int64                 `json:"event_id" bson:"event_id"`
	EventCode  string                `json:"event_code" bson:"event_code"`
	Regs       Regs                  `json:"regs" bson:"regs"`
	KoaRequest GetKaoActivityRequest `json:"kao_request" bson:"kao_request"`
	Event      Event                 `json:"event" bson:"event"`
}

// RegisterV2 struct for register v2 event data
type RegisterV2 struct {
	OwnerID   string `json:"owner_id" bson:"owner_id"`
	UserCode  string `json:"user_code" bson:"user_code"`
	EventCode string `json:"event_code" bson:"event_code"`
	Ref2      string `json:"ref2" bson:"ref2"`
	Regs      []Regs `json:"regs" bson:"regs"`
	Event     Event  `json:"event" bson:"event"`
}

// ReportRegisterV2 struct for register v2 event data
type ReportRegisterV2 struct {
	OwnerID   string       `json:"owner_id" bson:"owner_id"`
	UserCode  string       `json:"user_code" bson:"user_code"`
	EventCode string       `json:"event_code" bson:"event_code"`
	Ref2      string       `json:"ref2" bson:"ref2"`
	Regs      []RegsReport `json:"regs" bson:"regs"`
}

// Regs struct for register v2 event data
type Regs struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID        primitive.ObjectID `json:"user_id" bson:"user_id"`
	EventID       primitive.ObjectID `json:"event_id" bson:"event_id"`
	EventCode     string             `json:"event_code" bson:"event_code"`
	Status        string             `json:"status" bson:"status"`
	PaymentType   string             `json:"payment_type" bson:"payment_type"`
	TotalPrice    float64            `json:"total_price" bson:"total_price"`
	DiscountPrice float64            `json:"discount_price" bson:"discount_price"`
	PromoCode     string             `json:"promo_code" bson:"promo_code"`
	OrderID       string             `json:"order_id" bson:"order_id"`
	RegDate       time.Time          `json:"reg_date" bson:"reg_date"`
	PaymentDate   time.Time          `json:"payment_date" bson:"payment_date"`
	Coupon        Coupon             `json:"coupon" bson:"coupon"`
	TicketOptions []TicketOptionV2   `json:"ticket_options" bson:"ticket_options"`
	Partner       PartnerEvent       `json:"partner" bson:"partner"`
	ParentRegID   primitive.ObjectID `json:"parent_reg_id" bson:"parent_reg_id"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at"`
}

// RegsReport struct for register v2 event data
type RegsReport struct {
	ID            primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	UserID        primitive.ObjectID   `json:"user_id" bson:"user_id"`
	EventCode     string               `json:"event_code" bson:"event_code"`
	Status        string               `json:"status" bson:"status"`
	PaymentType   string               `json:"payment_type" bson:"payment_type"`
	TotalPrice    float64              `json:"total_price" bson:"total_price"`
	OrderID       string               `json:"order_id" bson:"order_id"`
	RegDate       time.Time            `json:"reg_date" bson:"reg_date"`
	PaymentDate   time.Time            `json:"payment_date" bson:"payment_date"`
	TicketOptions []TicketOptionReport `json:"ticket_options" bson:"ticket_options"`
	// Event         Event              `json:"event" bson:"event"`
}

//TicketOptionV2 struct
type TicketOptionV2 struct {
	UserOption     UserOption `json:"user_option" bson:"user_option"`
	TotalPrice     float64    `json:"total_price" bson:"total_price"`
	RegisterNumber string     `json:"register_number" bson:"register_number"`
	RecieptType    string     `json:"reciept_type" bson:"reciept_type"`
	Tickets        Tickets    `json:"tickets" bson:"tickets"`
	Shirts         Shirts     `json:"shirts" bson:"shirts"`
}

//TicketOptionReport struct
type TicketOptionReport struct {
	UserOption     UserOptionReport `json:"user_option" bson:"user_option"`
	TotalPrice     float64          `json:"total_price" bson:"total_price"`
	RegisterNumber string           `json:"register_number" bson:"register_number"`
	Tickets        TicketsReport    `json:"tickets" bson:"tickets"`
}

// Tickets new model
type TicketsReport struct {
	ID    string `json:"id" bson:"id"`
	Title string `json:"title" bson:"title"`
}

//RegisterTicketV2 struct
type RegisterTicketV2 struct {
	TicketID   string  `json:"ticket_id" bson:"ticket_id"`
	TicketName string  `json:"ticket_name" bson:"ticket_name"`
	Category   string  `json:"category" bson:"category"`
	Distance   float64 `json:"distance" bson:"distance"`
	TotalPrice float64 `json:"total_price" bson:"total_price"`
	Type       string  `json:"type" bson:"type"`
	Remark     string  `json:"remark" bson:"remark"`
	// Product    []ProduceEventV2 `json:"product" bson:"product"`
}

//RegisterChargeRequest payment charge request
type RegisterChargeRequest struct {
	TokenOmise string  `json:"token" form:"token"`
	RegID      string  `json:"reg_id" form:"reg_id"`
	EventCode  string  `json:"event_code" form:"event_code"`
	Price      float64 `json:"price" form:"price"`
	OrderID    string  `json:"order_id" form:"order_id"`
	Ref2       string  `json:"ref2" bson:"ref2"`
}
