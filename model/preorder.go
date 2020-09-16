package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// PreOrder model order product by event
type PreOrder struct {
	UserID          primitive.ObjectID `json:"ID" bson:"_id,omitempty"`
	FirstName       string             `json:"FirstName" bson:"FirstName"`
	LastName        string             `json:"LastName" bson:"LastName"`
	EBib            string             `json:"EBib" bson:"E-BiB"`
	TelNo           string             `json:"Tel_No" bson:"Tel_No"`
	ShirtType       string             `json:"Shirt_Type" bson:"Shirt_Type"`
	ShirtSize       string             `json:"Shirt_Size" bson:"Shirt_Size"`
	ShippingAddress string             `json:"Shipping_Address" bson:"Shipping_Address"`
}

// FindPreOrderRequest model for request
type FindPreOrderRequest struct {
	Keyword string `json:"keyword"`
}