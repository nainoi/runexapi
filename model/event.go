package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name            string             `json:"name" bson:"name"`
	Description     string             `json:"description" bson:"description"`
	Body            string             `json:"body" bson:"body"`
	Captions        []Caption          `json:"captions" bson:"captions"`
	Cover           string             `json:"cover" bson:"cover"`
	CoverThumb      string             `json:"cover_thumb" bson:"cover_thumb"`
	Category        Category           `json:"category" bson:"category"`
	Slug            string             `json:"slug" bson:"slug"`
	Product         []ProduceEvent     `json:"product" bson:"product"`
	Ticket          []TicketEvent      `json:"ticket" bson:"ticket"`
	Tag             []Tag              `json:"tags" bson:"tags"`
	OwnerID         primitive.ObjectID `json:"owner_id" bson:"owner_id"`
	Status          string             `json:"status" bson:"status"`
	Location        string             `json:"location" bson:"location"`
	ReceiveLocation string             `json:"receive_location" bson:"receive_location"`
	IsActive        bool               `json:"is_active" bson:"is_active"`
	IsFree          bool               `json:"is_free" bson:"is_free"`
	StartReg        time.Time          `json:"start_reg" bson:"start_reg"`
	EndReg          time.Time          `json:"end_reg" bson:"end_reg"`
	StartEvent      time.Time          `json:"start_event" bson:"start_event"`
	EndEvent        time.Time          `json:"end_event" bson:"end_event"`
	Inapp           bool               `json:"inapp" bson:"inapp"`
	IsPost          bool               `json:"is_post" bson:"is_post"`
	PostEndDate     time.Time          `json:"post_end_date" bson:"post_end_date"`
	CreatedTime     time.Time          `json:"created_time" bson:"created_time"`
	UpdatedTime     time.Time          `json:"updated_time" bson:"updated_time"`
}

type Category struct {
	CategoryID primitive.ObjectID `json:"id" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	Active     bool               `json:"active" bson:"active"`
}

type Caption struct {
	Image   string `json:"image" bson:"image"`
	Caption string `json:"caption" bson:"caption"`
}

type Tag struct {
	Name string `json:"name" bson:"name"`
}

type ProduceEvent struct {
	ProductID primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name" binding:"required"`
	Image     []ProductImage     `json:"image" bson:"image"`
	Detail    string             `json:"detail" bson:"detail"`
	Type      []ProductType      `json:"type" bson:"type"`
	Unit      int                `json:"unit" bson:"unit"`
	Currency  string             `json:"currency" bson:"currency"`
	Status    string             `json:"status" bson:"status"`
	Reuse     bool               `json:"reuse" bson:"reuse"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type ProductImage struct {
	PathURL string `json:"path_url" bson:"path_url"`
}

type ProductType struct {
	Name   string  `json:"name" bson:"name"`
	Remark string  `json:"remark" bson:"remark"`
	Price  float64 `json:"price" bson:"price"`
}

// type ProductSubType struct {
// 	Name string `json:"name" bson:"name"`
// }

type TicketEvent struct {
	TicketID    primitive.ObjectID `json:"id" bson:"_id"`
	ProductID   []ProductTicket    `json:"product" bson:"product"`
	Price       float64            `json:"price" bson:"price"`
	Unit        string             `json:"unit" bson:"unit"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Image       string             `json:"image" bson:"image"`
	SubType     []string           `json:"subtype" bson:"subtype"`
	Distance    float64            `json:"distance" bson:"distance"`
	Currency    string             `json:"currency" bson:"currency"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type ProductTicket struct {
	ProductTicketID primitive.ObjectID `json:"id" bson:"_id"`
	Show            bool               `json:"show" bson:"show"`
	Reuse           bool               `json:"reuse" bson:"reuse"`
}

type SearchEvent struct {
	Term string `json:"term" bson:"term" binding:"required"`
}

type EventReg struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Cover       string             `json:"cover" bson:"cover"`
	CoverThumb  string             `json:"cover_thumb" bson:"cover_thumb"`
	Category    Category           `json:"category" bson:"category"`
	// Product     []ProduceEvent     `json:"product" bson:"product"`
	// Ticket      []TicketEvent      `json:"ticket" bson:"ticket"`
	Status     string    `json:"status" bson:"status"`
	Location   string    `json:"location" bson:"location"`
	IsActive   bool      `json:"is_active" bson:"is_active"`
	Inapp      bool      `json:"inapp" bson:"inapp"`
	StartReg   time.Time `json:"start_reg" bson:"start_reg"`
	EndReg     time.Time `json:"end_reg" bson:"end_reg"`
	StartEvent time.Time `json:"start_event" bson:"start_event"`
	EndEvent   time.Time `json:"end_event" bson:"end_event"`
}

type EventRegInfo struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Body        string             `json:"body" bson:"body"`
	IsActive    bool               `json:"is_active" bson:"is_active"`
	StartEvent  time.Time          `json:"start_event" bson:"start_event"`
	EndEvent    time.Time          `json:"end_event" bson:"end_event"`
}

type Slug struct {
	Slug string `json:"slug" bson:"slug" binding:"required"`
}

// type ProductUpdateForm struct {
// 	ProductID primitive.ObjectID `json:"id" bson:"_id"`
// 	Name      string             `json:"name" bson:"name" binding:"required"`
// 	Image     []ProductImage     `json:"image" bson:"image"`
// 	Detail    string             `json:"detail" bson:"detail"`
// 	Type      []ProductType      `json:"type" bson:"type"`
// 	Unit      int                `json:"unit" bson:"unit"`
// 	Status    string             `json:"status" bson:"status"`
// 	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
// 	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
// }
