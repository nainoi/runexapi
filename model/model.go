package model

import (
	// "time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Users for all user
type Users struct {
	User []User `json:"user" bson:"user"`
}

//Projects for all projects

//LoginProvider struct
type LoginProvider struct {
	Provider   string `json:"provider" bson:"provider" binding:"required"`
	ProviderID string `json:"provider_id" bson:"provider_id" binding:"required"`
	PF         string `json:"pf" bson:"pf" binding:"required"`
}

//LoginEmail struct
type LoginEmail struct {
	Email    string `json:"email" bson:"email" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
	PF       string `json:"pf" bson:"pf" binding:"required"`
}

// UserAuth for user auth
type UserAuth struct {
	UserID      primitive.ObjectID `json:"id" bson:"_id"`
	Email       string             `json:"email" bson:"email"`
	Role        string             `json:"role" bson:"role"`
	PF          string             `json:"pf" bson:"pf"`
	Password    string             `json:"password" bson:"password"`
	NewPassword string             `json:"new_password"`
}
// UserForgot for forgot password
type UserForgot struct {
	UserID   primitive.ObjectID `json:"id" bson:"_id"`
	Email    string             `json:"email" bson:"email"`
	Role     string             `json:"role" bson:"role"`
	PF       string             `json:"pf" bson:"pf"`
	Token    string             `json:"token"`
	Fullname string             `json:"fullname"`
	Firstname string             `json:"firstname"`
	Lastname string             `json:"lastname"`
}
