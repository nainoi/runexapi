package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Register struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID         primitive.ObjectID `json:"user_id" bson:"user_id"`
	EventID        primitive.ObjectID `json:"event_id" bson:"event_id"`
	Event          EventReg           `json:"event" bson:"event"`
	Product        []RegisterProduct  `json:"product" bson:"product"`
	Tickets        []RegisterTicket   `json:"tickets" bson:"tickets"`
	Status         string             `json:"status" bson:"status"`
	PaymentType    string             `json:"payment_type" bson:"payment_type"`
	TotalPrice     float64            `json:"total_price" bson:"total_price"`
	DiscountPrice  float64            `json:"discount_price" bson:"discount_price"`
	PromoCode      string             `json:"promo_code" bson:"promo_code"`
	OrderID        string             `json:"order_id" bson:"order_id"`
	Ref2           string             `json:"ref2" bson:"ref2"`
	RegDate        time.Time          `json:"reg_date" bson:"reg_date"`
	RegisterNumber string             `json:"register_number" bson:"register_number"`
	ShipingAddress ShipingAddress     `json:"shiping_address" bson:"shiping_address"`
	Coupon         Coupon             `json:"coupon" bson:"coupon"`
	TicketOptions  []TicketOption     `json:"ticket_options" bson:"ticket_options"`
	Slip           SlipTransfer       `json:"slip" bson:"slip"`
	Phone          string             `json:"phone" bson:"phone"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}

type ShipingAddress struct {
	ShipingAddressID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Address          string             `json:"address" bson:"address"`
	Province         string             `json:"province" bson:"province"`
	District         string             `json:"district" bson:"district"`
	City             string             `json:"city" bson:"city"`
	Zipcode          string             `json:"zipcode" bson:"zipcode"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
}

type RegisterProduct struct {
	ProductID primitive.ObjectID    `json:"id" bson:"_id,omitempty"`
	Type      string                `json:"type" bson:"type"`
	Price     float32               `json:"price" bson:"price"`
	Product   RegisterProductDetail `json:"product" bson:"product"`
}

type RegisterProductDetail struct {
	ProductID primitive.ObjectID    `json:"id" bson:"_id,omitempty"`
	Name      string                `json:"name" bson:"name"`
	Price     float32               `json:"price" bson:"price"`
	Image     []ProductImage        `json:"image" bson:"image"`
	Detail    string                `json:"detail" bson:"detail"`
	Type      []RegisterProductType `json:"type" bson:"type"`
	Unit      float32               `json:"unit" bson:"unit"`
	Currency  string                `json:"currency" bson:"currency"`
	Status    string                `json:"status" bson:"status"`
	CreatedAt time.Time             `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time             `json:"updated_at" bson:"updated_at"`
}

type RegisterProductImage struct {
	PathURL string `json:"path_url" bson:"path_url"`
}

type RegisterProductType struct {
	Name   string  `json:"name" bson:"name"`
	Remark string  `json:"remark" bson:"remark"`
	Price  float32 `json:"price" bson:"price"`
}

type RegisterTicket struct {
	Product      RegisterProductDetail `json:"product" bson:"product"`
	Type         string                `json:"type" bson:"type"`
	Remark       string                `json:"remark" bson:"remark"`
	Distance     string                `json:"distance" bson:"distance"`
	TicketDetail TicketEvent           `json:"ticket" bson:"ticket"`
}

type RegisterTicketDetail struct {
	TicketID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Type     string             `json:"type" bson:"type"`
}

type ProductTicketDetail struct {
	ProductID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Show      bool               `json:"show" bson:"show"`
}

type SlipTransfer struct {
	Amount       float32     `json:"amount" bson:"amount"`
	DateTransfer string      `json:"date_tranfer" bson:"date_tranfer"`
	TimeTransfer string      `json:"time_tranfer" bson:"time_tranfer"`
	Image        string      `json:"image" bson:"image"`
	Remark       string      `json:"remark" bson:"remark"`
	OrderID      string      `json:"order_id" bson:"order_id"`
	BankAccount  BankAccount `json:"bank_account" bson:"bank_account"`
}

type BankAccount struct {
	Bank          string `json:"bank" bson:"bank"`
	AccountName   string `json:"account_name" bson:"account_name"`
	AccountNumber string `json:"account_number" bson:"account_number"`
}

type EmailTemplateData struct {
	Name                 string `json:"name" bson:"name"`
	Gender               string `json:"gender" bson:"gender"`
	Nationality          string `json:"nationality" bson:"nationality"`
	Birthdate            string `json:"birthdate" bson:"birthdate"`
	Email                string `json:"email" bson:"email" binding:"required"`
	Phone                string `json:"phone" bson:"phone"`
	IdentificationNumber string `json:"identification_number" bson:"identification_number"`
	ContactName          string `json:"contact_name" bson:"contact_name"`
	ContactRole          string `json:"contact_role" bson:"contact_role"`
	ContactPhone         string `json:"contact_phone" bson:"contact_phone"`
	CompetitionType      string `json:"competition_type" bson:"competition_type"`
	Price                string `json:"price" bson:"price"`
	Generation           string `json:"generation" bson:"generation"`
	Question1            string `json:"question1" bson:"question1"`
	Question2            string `json:"question2" bson:"question2"`
	Answer1              string `json:"answer1" bson:"answer1"`
	Answer2              string `json:"answer2" bson:"answer2"`
	RegisterNumber       string `json:"register_number" bson:"register_number"`
	ImageURL             string `json:"image_url" bson:"image_url"`
}

type EmailTemplateData2 struct {
	Name                 string         `json:"name" bson:"name"`
	Email                string         `json:"email" bson:"email" binding:"required"`
	Phone                string         `json:"phone" bson:"phone"`
	RefID                string         `json:"ref_id" bson:"ref_id"`
	IdentificationNumber string         `json:"identification_number" bson:"identification_number"`
	ContactPhone         string         `json:"contact_phone" bson:"contact_phone"`
	CompetitionType      string         `json:"competition_type" bson:"competition_type"`
	Price                float64        `json:"price" bson:"price"`
	Unit                 string         `json:"unit" bson:"unit"`
	RegisterNumber       string         `json:"register_number" bson:"register_number"`
	TicketName           string         `json:"ticket_name" bson:"ticket_name"`
	EventName            string         `json:"event_name" bson:"event_name"`
	Status               string         `json:"status" bson:"status"`
	RegisterDate         time.Time      `json:"register_date" bson:"register_date"`
	PaymentType          string         `json:"payment_type" bson:"payment_type"`
	ShipingAddress       string         `json:"shiping_address" bson:"shiping_address"`
	TicketOptions        []TicketOption `json:"ticket_options" bson:"ticket_options"`
}

type ReportRegister struct {
	PaymentAll            int64          `json:"payment_all" bson:"payment_all"`
	PaymentWaiting        int64          `json:"payment_waiting" bson:"payment_waiting"`
	PaymentWaitingApprove int64          `json:"payment_waiting_approve" bson:"payment_wiating_approve"`
	PaymentSuccess        int64          `json:"payment_success" bson:"payment_success"`
	Datas                 []DataRegister `json:"datas" bson:"datas"`
}

type DataRegister struct {
	RegisterID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Firstname     string             `json:"firstname" bson:"firstname"`
	Lastname      string             `json:"lastname" bson:"lastname"`
	Email         string             `json:"email" bson:"email"`
	Address       ShipingAddress     `json:"address" bson:"address"`
	TicketName    string             `json:"ticket_name" bson:"ticket_name"`
	Distance      float64            `json:"distance" bson:"distance"`
	PaymentType   string             `json:"payment_type" bson:"payment_type"`
	PaymentStatus string             `json:"payment_status" bson:"payment_status"`
	PaymentDate   time.Time          `json:"payment_date" bson:"payment_date"`
	RegDate       time.Time          `json:"reg_date" bson:"reg_date"`
	Price         float64            `json:"price" bson:"price"`
	ShirtSize     string             `json:"shirt_size" bson:"shirt_size"`
	DistanceTotal float64            `json:"distance_total" bson:"distance_total"`
	Slip          SlipTransfer       `json:"slip" bson:"slip"`
	TicketProduct []RegisterTicket   `json:"reg_tickets" bson:"reg_tickets"`
	TicketOptions []TicketOption     `json:"ticket_options" bson:"ticket_options"`
}

type DataRegisterRequest struct {
	EventID    string `json:"event_id" bson:"event_id"`
	PageNumber int64  `json:"pageNumber" bson:"pageNumber" binding:"required"`
	NPerPage   int64  `json:"nPerPage" bson:"nPerPage" binding:"required"`
	Status     string `json:"status" bson:"status"`
	KeyWord    string `json:"key_word" bson:"key_word"`
}

type OwnerRequest struct {
	EventCode string `json:"event_code" bson:"event_code" binding:"required"`
	OwnerID   string `json:"owner_id" bson:"owner_id" binding:"required"`
}

type UpdayeRegisterStatusRequest struct {
	RegisterID string `json:"register_id" bson:"register_id"`
	Status     string `json:"status" bson:"status" binding:"required"`
}

type TicketOption struct {
	UserOption     UserOption        `json:"user_option" bson:"user_option"`
	Product        []RegisterProduct `json:"product" bson:"product"`
	Tickets        []RegisterTicket  `json:"tickets" bson:"tickets"`
	TotalPrice     float64           `json:"total_price" bson:"total_price"`
	RecieptType    string            `json:"reciept_type" bson:"reciept_type"`
	RegisterNumber string            `json:"register_number" bson:"register_number"`
}

type UserOption struct {
	FirstName        string    `json:"firstname" bson:"firstname"`
	LastName         string    `json:"lastname" bson:"lastname"`
	FirstNameTH      string    `json:"firstname_th" bson:"firstname_th"`
	LastNameTH       string    `json:"lastname_th" bson:"lastname_th"`
	Phone            string    `json:"phone" bson:"phone"`
	BirthDate        time.Time `json:"birthdate" bson:"birthdate"`
	Gender           string    `json:"gender" bson:"gender"`
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`
	Confirm          bool      `json:"confirm" bson:"confirm"`
	EmergencyContact string    `json:"emergency_contact" bson:"emergency_contact"`
	EmergencyPhone   string    `json:"emergency_phone" bson:"emergency_phone"`
	Nationality      string    `json:"nationality" bson:"nationality"`
	Passport         string    `json:"passport" bson:"passport"`
	CitycenID        string    `json:"citycen_id" bson:"citycen_id"`
	BloodType        string    `json:"blood_type" bson:"blood_type"`
	Address          string    `json:"address" bson:"address"`
	HomeNo           string    `json:"home_no" bson:"home_no"`
	Moo              string    `json:"moo" bson:"moo"`
	Tambon           Tambon    `json:"tambon" bson:"tambon"`
	Team             string    `json:"team" bson:"team"`
	Color            string    `json:"color" bson:"color"`
	Zone             string    `json:"zone" bson:"zone"`
	EmpID            string    `json:"emp_id" bson:"emp_id"`
}

type UserOptionReport struct {
	FirstName string    `json:"firstname" bson:"firstname"`
	LastName  string    `json:"lastname" bson:"lastname"`
	Phone     string    `json:"phone" bson:"phone"`
	BirthDate time.Time `json:"birthdate" bson:"birthdate"`
	Gender    string    `json:"gender" bson:"gender"`
	Address   string    `json:"address" bson:"address"`
}

type ShipingAddressUpdateForm struct {
	Address   string    `json:"address" bson:"address"`
	Province  string    `json:"province" bson:"province"`
	District  string    `json:"district" bson:"district"`
	City      string    `json:"city" bson:"city"`
	Zipcode   string    `json:"zipcode" bson:"zipcode"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
