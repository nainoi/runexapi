package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterRequest struct {
	EventID primitive.ObjectID `json:"event_id" bson:"event_id"`
	Regs    Regs               `json:"regs" bson:"regs"`
}

type RegisterV2 struct {
	EventID primitive.ObjectID `json:"event_id" bson:"event_id"`
	Regs    []Regs             `json:"regs" bson:"regs"`
}

type Regs struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID         primitive.ObjectID `json:"user_id" bson:"user_id"`
	Status         string             `json:"status" bson:"status"`
	PaymentType    string             `json:"payment_type" bson:"payment_type"`
	TotalPrice     float64            `json:"total_price" bson:"total_price"`
	DiscountPrice  float64            `json:"discount_price" bson:"discount_price"`
	PromoCode      string             `json:"promo_code" bson:"promo_code"`
	OrderID        string             `json:"order_id" bson:"order_id"`
	RegDate        time.Time          `json:"reg_date" bson:"reg_date"`
	PaymentDate    time.Time          `json:"payment_date" bson:"payment_date"`
	RegisterNumber string             `json:"register_number" bson:"register_number"`
	Coupon         Coupon             `json:"coupon" bson:"coupon"`
	TicketOptions  TicketOptionV2     `json:"ticket_options" bson:"ticket_options"`
	Partner        PartnerEvent       `json:"partner" bson:"partner"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}

type TicketOptionV2 struct {
	UserOption     UserOption         `json:"user_option" bson:"user_option"`
	TotalPrice     float64            `json:"total_price" bson:"total_price"`
	RegisterNumber string             `json:"register_number" bson:"register_number"`
	RecieptType    string             `json:"reciept_type" bson:"reciept_type"`
	Tickets        []RegisterTicketV2 `json:"tickets" bson:"tickets"`
}

type RegisterTicketV2 struct {
	TicketID   primitive.ObjectID `json:"ticket_id" bson:"ticket_id"`
	TicketName string             `json:"ticket_name" bson:"ticket_name"`
	Distance   float64            `json:"distance" bson:"distance"`
	TotalPrice float64            `json:"total_price" bson:"total_price"`
	Type       string             `json:"type" bson:"type"`
	Remark     string             `json:"remark" bson:"remark"`
	Product    []ProduceEventV2   `json:"product" bson:"product"`
}
