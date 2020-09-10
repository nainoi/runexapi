package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"thinkdev.app/think/runex/runexapi/model"
)

// User object for db
type User struct {
	UserID           primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Email            string             `json:"email" bson:"email"`
	Provider         []model.Provider   `json:"provider" bson:"provider"`
	FullName         string             `json:"fullname" bson:"fullname"`
	FirstName        string             `json:"firstname" bson:"firstname"`
	FirstNameTH      string             `json:"firstname_th" bson:"firstname_th"`
	LastNameTH       string             `json:"lastname_th" bson:"lastname_th"`
	LastName         string             `json:"lastname" bson:"lastname"`
	Password         string             `json:"-" bson:"password,omitempty"`
	Phone            string             `json:"phone" bson:"phone"`
	Avatar           string             `json:"avatar" bson:"avatar"`
	Role             string             `json:"role" bson:"role"`
	PF               string             `json:"pf" bson:"-"`
	BirthDate        time.Time          `json:"birthdate" bson:"birthdate"`
	Gender           string             `json:"gender" bson:"gender"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
	Confirm          bool               `json:"confirm" bson:"confirm"`
	Address          []model.Address    `json:"address" bson:"address"`
	EmergencyContact string             `json:"emergency_contact" bson:"emergency_contact"`
	EmergencyPhone   string             `json:"emergency_phone" bson:"emergency_phone"`
	Nationality      string             `json:"nationality" bson:"nationality"`
	Passport         string             `json:"passport" bson:"passport"`
	CitycenID        string             `json:"citycen_id" bson:"citycen_id"`
	BloodType        string             `json:"blood_type" bson:"blood_type"`
}

// Address customer
// swagger:model
// type Address struct {
// 	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
// 	Address   string             `json:"address" bson:"address" binding:"required"`
// 	Province  string             `json:"province" bson:"province" binding:"required"`
// 	District  string             `json:"district" bson:"district" binding:"required"`
// 	City      string             `json:"city" bson:"city" binding:"required"`
// 	ZipCode   string             `json:"zipcode" bson:"zipcode" binding:"required"`
// 	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
// 	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
// }

// UserProviderRequest object for register from provider
// swagger:model
type UserProviderRequest struct {
	// in: body
	Email      string    `json:"email" bson:"email" binding:"exists,email"`
	Provider   string    `json:"provider" bson:"provider" binding:"required"`
	ProviderID string    `json:"provider_id" bson:"provider_id" binding:"required"`
	FullName   string    `json:"fullname" bson:"fullname"`
	FirstName  string    `json:"firstname" bson:"firstname"`
	LastName   string    `json:"lastname" bson:"lastname"`
	Avatar     string    `json:"avatar" bson:"avatar"`
	PF         string    `json:"pf" bson:"pf" binding:"required"`
	BirthDate  time.Time `json:"birthdate" bson:"birthdate"`
	Gender     string    `json:"gender" bson:"gender"`
}
