package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Team struct {
	RegID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	OwnerID primitive.ObjectID `json:"user_id" bson:"user_id"`
	Name    string             `json:"name" bson:"name"`
	Color   string             `json:"color" bson:"color"`
	Zone    string             `json:"zone" bson:"zone"`
	IconURL string             `json:"icon_url" bson:"icon_url"`
}

type TeamIcon struct {
	RegID   primitive.ObjectID `json:"reg_id" form:"reg_id" bson:"reg_id"`
	IconURL string             `json:"icon_url" form:"icon_url" bson:"icon_url"`
}
