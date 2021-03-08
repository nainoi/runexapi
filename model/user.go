package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User object for db
type User struct {
	UserID           primitive.ObjectID   `json:"-" bson:"_id,omitempty"`
	Email            string               `json:"email" bson:"email"`
	Provider         string               `json:"provider" bson:"provider"`
	ProviderID       string               `json:"provider_id" bson:"provider_id"`
	FullName         string               `json:"fullname" bson:"fullname"`
	FirstName        string               `json:"firstname" bson:"firstname"`
	FirstNameTH      string               `json:"firstname_th" bson:"firstname_th"`
	LastNameTH       string               `json:"lastname_th" bson:"lastname_th"`
	LastName         string               `json:"lastname" bson:"lastname"`
	Password         string               `json:"-" bson:"password,omitempty"`
	Phone            string               `json:"phone" bson:"phone"`
	Avatar           string               `json:"avatar" bson:"avatar"`
	Role             string               `json:"role" bson:"role"`
	BirthDate        time.Time            `json:"birthdate" bson:"birthdate"`
	Gender           string               `json:"gender" bson:"gender"`
	CreatedAt        time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at" bson:"updated_at"`
	Confirm          bool                 `json:"confirm" bson:"confirm"`
	Address          []Address            `json:"address" bson:"address"`
	EmergencyContact string               `json:"emergency_contact" bson:"emergency_contact"`
	EmergencyPhone   string               `json:"emergency_phone" bson:"emergency_phone"`
	Nationality      string               `json:"nationality" bson:"nationality"`
	Passport         string               `json:"passport" bson:"passport"`
	CitycenID        string               `json:"citycen_id" bson:"citycen_id"`
	BloodType        string               `json:"blood_type" bson:"blood_type"`
	PF               string               `json:"pf" bson:"-"`
	Events           []primitive.ObjectID `json:"events" bson:"events"`
	StravaID         string               `json:"strava_id" bson:"strava_id"`
	StravaAvatar     string               `json:"strava_avatar" bson:"strava_avatar"`
	StravaFirstname  string               `json:"strava_firstname" bson:"strava_firstname"`
	StravaLastname   string               `json:"strava_latname" bson:"strava_latname"`
}

// UserMail object for register from email password
type UserMail struct {
	UserID      primitive.ObjectID `json:"user_id" bson:"_id,omitempty"`
	Email       string             `json:"email" bson:"email" binding:"exists,email"`
	FullName    string             `json:"fullname" bson:"fullname"`
	FirstName   string             `json:"firstname" bson:"firstname" binding:"required"`
	LastName    string             `json:"lastname" bson:"lastname" binding:"required"`
	FirstNameTH string             `json:"firstname_th" bson:"firstname_th"`
	LastNameTH  string             `json:"lastname_th" bson:"lastname_th"`
	Password    string             `json:"password" bson:"password" binding:"exists,min=8,max=255"`
	Phone       string             `json:"phone" bson:"phone"`
	Avatar      string             `json:"avatar" bson:"avatar"`
	Role        string             `json:"role" bson:"role"`
	PF          string             `json:"pf" bson:"pf" binding:"required"`
	BirthDate   time.Time          `json:"birthdate" bson:"birthdate" time_format:"2006-01-02"`
	Gender      string             `json:"gender" bson:"gender"`
	Confirm     bool               `json:"confirm" bson:"confirm"`
	Address     []Address          `json:"address" bson:"address"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// UserProvider object for register from provider
type UserProvider struct {
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

// Address customer
type Address struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Address   string             `json:"address" bson:"address" binding:"required"`
	Province  string             `json:"province" bson:"province" binding:"required"`
	District  string             `json:"district" bson:"district" binding:"required"`
	City      string             `json:"city" bson:"city" binding:"required"`
	ZipCode   string             `json:"zipcode" bson:"zipcode" binding:"required"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// UserEvent owner event
type UserEvent struct {
	UserID    primitive.ObjectID `json:"user_id" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email" binding:"exists,email"`
	FullName  string             `json:"fullname" bson:"fullname"`
	FirstName string             `json:"firstname" bson:"firstname"`
	LastName  string             `json:"lastname" bson:"lastname"`
	Phone     string             `json:"phone" bson:"phone"`
	Avatar    string             `json:"avatar" bson:"avatar"`
}

// Provider with provider open id
// swagger:model
type Provider struct {
	// in: body
	ProviderID   string `json:"provider_id" bson:"provider_id"`
	ProviderName string `json:"provider_name" bson:"provider_name"`
}

// UserProviderRequest object for register from provider
// swagger:model
type UserProviderRequest struct {
	// in: body
	Email      string    `json:"email" bson:"email"`
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

// UserStravaSyncRequest object for register from provider
// swagger:model
type UserStravaSyncRequest struct {
	// in: body
	Email           string `json:"email" binding:"exists,email"`
	Provider        string `json:"provider" bson:"provider" binding:"required"`
	ProviderID      string `json:"provider_id" bson:"provider_id" binding:"required"`
	StravaID        string `json:"strava_id" bson:"strava_id" binding:"required"`
	StravaAvatar    string `json:"strava_avatar"`
	StravaFirstname string `json:"strava_firstname"`
	StravaLastname  string `json:"strava_latname"`
}

// FirebaseUser model
type FirebaseUser struct {
	UserID         primitive.ObjectID `json:"user_id" bson:"user_id"`
	FirebaseTokens []string           `json:"firebase_tokens" bson:"firebase_tokens"`
}
