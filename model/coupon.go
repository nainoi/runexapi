package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Coupon struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SaleID       primitive.ObjectID `json:"sale_id" bson:"sale_id"`
	Discount    float64            `json:"discount" bson:"discount"`
	CouponCode  string             `json:"coupon_code" bson:"coupon_code"`
	Description string             `json:"description" bson:"description"`
	StartDate   time.Time          `json:"start_date" bson:"start_date"`
	EndDate     time.Time          `json:"end_date" bson:"end_date"`
	Active      bool               `json:"active" bson:"active"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type EditCouponForm struct {
	Discount    float64   `json:"discount" bson:"discount"`
	SaleID       primitive.ObjectID `json:"sale_id" bson:"sale_id"`
	CouponCode  string    `json:"coupon_code" bson:"coupon_code"`
	Description string    `json:"description" bson:"description"`
	StartDate   time.Time `json:"start_date" bson:"start_date"`
	EndDate     time.Time `json:"end_date" bson:"end_date"`
	Active      bool      `json:"active" bson:"active"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

type ValidateCoupon struct {
	CouponCode string `json:"coupon_code" bson:"coupon_code" binding:"required"`
}
