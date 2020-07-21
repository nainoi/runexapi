package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExcelUserForm struct {
	UserID      primitive.ObjectID `json:"user_id" bson:"_id,omitempty"`
	Email       string             `json:"email" bson:"email" binding:"exists,email"`
	Provider    string             `json:"provider" bson:"provider" binding:"required"`
	ProviderID  string             `json:"provider_id" bson:"provider_id" binding:"required"`
	FullName    string             `json:"fullname" bson:"fullname"`
	FirstName   string             `json:"firstname" bson:"firstname"`
	LastName    string             `json:"lastname" bson:"lastname"`
	FirstNameTH string             `json:"firstname_th" bson:"firstname_th"`
	LastNameTH  string             `json:"lastname_th" bson:"lastname_th"`
	CitycenID   string             `json:"citycen_id" bson:"citycen_id"`
	Phone       string             `json:"phone" bson:"phone"`
	Avatar      string             `json:"avatar" bson:"avatar"`
	Role        string             `json:"role" bson:"role"`
	PF          string             `json:"pf" bson:"pf" binding:"required"`
	BirthDate   time.Time          `json:"birthdate" bson:"birthdate"`
	Gender      string             `json:"gender" bson:"gender"`
	Confirm     bool               `json:"confirm" bson:"confirm"`
	Address     []Address          `json:"address" bson:"address"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}
