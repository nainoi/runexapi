package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SlipHistory struct {
	ID    primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	RegID primitive.ObjectID   `json:"reg_id" bson:"reg_id"`
	Slips []SlipTransferUpdate `json:"slips" bson:"slips"`
}

type SlipUpdateForm struct {
	RegID   string       `json:"reg_id" bson:"reg_id"`
	Slip    SlipTransfer `json:"slip" bson:"slip"`
	Comment string       `json:"comment" bson:"comment"`
}

type SlipTransferUpdate struct {
	Amount       float32            `json:"amount" bson:"amount"`
	DateTransfer string             `json:"date_tranfer" bson:"date_tranfer"`
	TimeTransfer string             `json:"time_tranfer" bson:"time_tranfer"`
	Image        string             `json:"image" bson:"image"`
	Remark       string             `json:"remark" bson:"remark"`
	OrderID      string             `json:"order_id" bson:"order_id"`
	BankAccount  BankAccount        `json:"bank_account" bson:"bank_account"`
	Comment      string             `json:"comment" bson:"comment"`
	CreatDated   time.Time          `json:"creat_dated" bson:"creat_dated"`
	UpdateBy     primitive.ObjectID `json:"update_by" bson:"update_by"`
}
