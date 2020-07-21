package model

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Fcm struct {
	UserID  primitive.ObjectID  `json:"user_id" bson:"user_id" binding:"required"`
	TokenFCM  string `json:"token_fcm" bson:"token_fcm" binding:"required"`
	IsActive bool `json:"is_active" bson:"is_active" binding:"required"`
	PF string `json:"pf" bson:"pf" binding:"required"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}