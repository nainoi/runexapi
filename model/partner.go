package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Partner struct
type Partner struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name          string             `json:"name" bson:"name"`
	PartnerConfig PartnerConfig      `json:"config" bson:"config"`
}

// PartnerConfig partner struct
type PartnerConfig struct {
	IsConfirm    bool   `json:"is_confirm" bson:"is_confirm"`
	ConfirmTitle string `json:"confirm_title" bson:"confirm_title"`
	ConfirmMsg   string `json:"confirm_msg" bson:"confirm_msg"`
	PlaceHolder  string `json:"place_holder" bson:"place_holder"`
	RefID        string `json:"ref_id" bson:"ref_id"`
}

// Param partner struct
type Param struct {
	Name   string `json:"name" bson:"name"`
	URL    string `json:"url" bson:"url"`
	Method string `json:"method" bson:"method"`
}
