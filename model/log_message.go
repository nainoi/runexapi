package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LogUpdateRegisterStatus log store data
type LogUpdateRegisterStatus struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	LogName   string             `json:"log_name" bson:"log_name"`
	Status    string             `json:"status" bson:"status" `
	RegID     primitive.ObjectID `json:"reg_id" bson:"reg_id" `
	UpdateBy  primitive.ObjectID `json:"update_by" bson:"update_by" `
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// LogActivityInfo log store data
type LogActivityInfo struct {
	UserID         primitive.ObjectID `json:"user_id" bson:"user_id"`
	ActivityInfoID primitive.ObjectID `json:"activity_info_id" bson:"activity_info_id"`
	EventID        primitive.ObjectID `json:"event_id" bson:"event_id"`
	Distance       float64            `json:"distance" bson:"distance"`
	ImageURL       string             `json:"img_url" bson:"img_url"`
	Caption        string             `form:"caption" json:"caption"`
	APP            string             `form:"app" json:"app"`
	Time           int64              `form:"time" json:"time"`
	ActivityDate   time.Time          `json:"activity_date" bson:"activity_date"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}

// LogRegister log store data
type LogRegister struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID         primitive.ObjectID `json:"user_id" bson:"user_id"`
	EventID        primitive.ObjectID `json:"event_id" bson:"event_id"`
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

// LogSendKaoActivity log store data
type LogSendKaoActivity struct {
	UserID         primitive.ObjectID `json:"user_id" bson:"user_id"`
	ActivityInfoID primitive.ObjectID `json:"activity_info_id" bson:"activity_info_id"`
	EventID        primitive.ObjectID `json:"event_id" bson:"event_id"`
	Distance       float64            `json:"distance" bson:"distance"`
	ImageURL       string             `json:"img_url" bson:"img_url"`
	APP            string             `form:"app" json:"app"`
	Slug           string             `form:"slug" json:"slug"`
	Ebib           string             `form:"ebib" json:"ebib"`
	Time           int64              `form:"time" json:"time"`
	ActivityDate   time.Time          `json:"activity_date" bson:"activity_date"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}
