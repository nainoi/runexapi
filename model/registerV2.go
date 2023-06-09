package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RegisterRequest struct for request register event
type RegisterRequest struct {
	EventCode  string                `json:"event_code" bson:"event_code" binding:"required"`
	Regs       Regs                  `json:"regs" bson:"regs"`
	KoaRequest GetKaoActivityRequest `json:"kao_request" bson:"kao_request"`
	Event      Event                 `json:"event" bson:"event"`
}

// RegisterRequest struct for request register event
type CheckRegisterERRequest struct {
	EventCode string `json:"event_code" bson:"event_code" binding:"required"`
	TicketID  string `json:"ticket_id" bson:"ticket_id" binding:"required"`
	CitycenID string `json:"citycen_id" bson:"citycen_id" binding:"required"`
}

// RegisterRequest struct for request register event
type AddTeamRequest struct {
	EventID     int64              `json:"event_id" bson:"event_id" binding:"required"`
	TeamUserID  primitive.ObjectID `json:"team_user_id" bson:"team_user_id" binding:"required"`
	EventCode   string             `json:"event_code" bson:"event_code" binding:"required"`
	ParentRegID primitive.ObjectID `json:"parent_reg_id" bson:"parent_reg_id" binding:"required"`
	Regs        Regs               `json:"regs" bson:"regs"`
	Event       Event              `json:"event" bson:"event"`
}

// RegisterRequest struct for request register event
type RegEventDashboardRequest struct {
	EventCode   string             `json:"event_code" bson:"event_code" binding:"required"`
	ParentRegID primitive.ObjectID `json:"parent_reg_id" bson:"parent_reg_id"`
	RegID       primitive.ObjectID `json:"reg_id" bson:"reg_id" binding:"required"`
	TicketID    string             `json:"ticket_id" bson:"ticket_id" binding:"required"`
}

// RegisterRequest struct for request register event
type RegUpdateUserInfoRequest struct {
	EventCode    string             `json:"event_code" bson:"event_code" binding:"required"`
	RegID        primitive.ObjectID `json:"reg_id" bson:"reg_id" binding:"required"`
	TicketOption TicketOptionV2     `json:"ticket_options" bson:"ticket_options" binding:"required"`
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
	TicketID      string             `json:"ticket_id" bson:"ticket_id"`
	EventCode     string             `json:"event_code" bson:"event_code"`
	Status        string             `json:"status" bson:"status"`
	PaymentType   string             `json:"payment_type" bson:"payment_type"`
	TotalPrice    float64            `json:"total_price" bson:"total_price"`
	DiscountPrice float64            `json:"discount_price" bson:"discount_price"`
	PromoCode     string             `json:"promo_code" bson:"promo_code"`
	OrderID       string             `json:"order_id" bson:"order_id"`
	IsTeamLead    bool               `json:"is_team_lead" bson:"is_team_lead"`
	RegDate       time.Time          `json:"reg_date" bson:"reg_date"`
	PaymentDate   time.Time          `json:"payment_date" bson:"payment_date"`
	Slip          string             `json:"slip" bson:"slip"`
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
	Slip          string               `json:"slip" bson:"slip"`
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
	Shirts         Shirts           `json:"shirts" bson:"shirts"`
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

//RegisterChargeRequest payment charge request
type RegisterAttachSlipRequest struct {
	RegID       primitive.ObjectID `json:"reg_id" form:"reg_id"`
	EventCode   string             `json:"event_code" form:"event_code"`
	OrderID     string             `json:"order_id" form:"order_id"`
	Ref2        string             `json:"ref2" bson:"ref2"`
	PaymentType string             `json:"payment_type" bson:"payment_type"`
	Image       string             `json:"image" bson:"image"`
	Status      string             `json:"status" bson:"status"`
}

//AdminAttachSlipRequest payment charge request
type AdminAttachSlipRequest struct {
	RegID       primitive.ObjectID `json:"reg_id" form:"reg_id" bson:"reg_id"`
	UserID      primitive.ObjectID `json:"user_id" form:"user_id" bson:"user_id"`
	EventCode   string             `json:"event_code" form:"event_code" bson:"event_code"`
	OrderID     string             `json:"order_id" form:"order_id" bson:"order_id"`
	Ref2        string             `json:"ref2" bson:"ref2" form:"ref2"`
	PaymentType string             `json:"payment_type" bson:"payment_type" form:"payment_type"`
	Image       string             `json:"image" bson:"image" form:"image"`
	Status      string             `json:"status" bson:"status" form:"status"`
}

type RegisterActivityInfo struct {
	EventCode           string  `json:"event_code" form:"event_code" bson:"event_code"`
	EventTitle          string  `json:"event_title" form:"event_title" bson:"event_title"`
	EventCover          string  `json:"event_cover" form:"event_cover" bson:"event_cover"`
	EventCoverThumbnail string  `json:"event_cover_thumbnail" form:"event_cover_thumbnail" bson:"event_cover_thumbnail"`
	EventDate           string  `json:"event_date" bson:"event_date" form:"event_date"`
	TicketId            string  `json:"ticket_id" bson:"ticket_id" form:"ticket_id"`
	TicketCategory      string  `json:"ticket_category" bson:"ticket_category" form:"ticket_category"`
	TicketDistance      string  `json:"ticket_distance" bson:"ticket_distance" form:"ticket_distance"`
	TicketTitle         string  `json:"ticket_title" bson:"ticket_title" form:"ticket_title"`
	TicketPrice         string  `json:"ticket_price" bson:"ticket_price" form:"ticket_price"`
	ProviderName        string  `json:"provider_name" bson:"provider_name" form:"provider_name"`
	ProviderId          string  `json:"provider_id" bson:"provider_id" form:"provider_id"`
	ShirtSize           string  `json:"shirt_size" bson:"shirt_size" form:"shirt_size"`
	Firstname           string  `json:"firstname" bson:"firstname" form:"firstname"`
	Lastname            string  `json:"lastname" bson:"lastname" form:"lastname"`
	CardId              string  `json:"card_id" bson:"card_id" form:"card_id"`
	PhoneNumber         string  `json:"phone_number" bson:"phone_number" form:"phone_number"`
	Bib                 string  `json:"bib" bson:"bib" form:"bib"`
	Gender              string  `json:"gender" bson:"gender" form:"gender"`
	Birthday            string  `json:"birthday" bson:"birthday" form:"birthday"`
	Blood               string  `json:"blood" bson:"blood" form:"blood"`
	Address             string  `json:"address" bson:"address" form:"address"`
	Moo                 string  `json:"moo" bson:"moo" form:"moo"`
	ZipCode             string  `json:"zip_code" bson:"zip_code" form:"zip_code"`
	District            string  `json:"district" bson:"district" form:"district"`
	Amphoe              string  `json:"amphoe" bson:"amphoe" form:"amphoe"`
	Province            string  `json:"province" bson:"province" form:"province"`
	TeamName            string  `json:"team_name" bson:"team_name" form:"team_name"`
	Options             Options `json:"options" bson:"options" form:"options"`
	RegisterId          string  `json:"register_id" bson:"register_id" form:"register_id"`
}

type Options struct {
	Color      string `json:"color" bson:"color" form:"color"`
	EmployeeId string `json:"employee_id" bson:"employee_id" form:"employee_id"`
}
