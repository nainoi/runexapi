package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type ReportDashboard struct {
	EventID       primitive.ObjectID `json:"event_id" bson:"event_id"`
	Paid          float64            `json:"paid" bson:"paid"`
	WaitToPay     float64            `json:"wait_to_pay" bson:"wait_to_pay"`
	WaitToApprove float64            `json:"wait_to_approve" bson:"wait_to_approve"`
	RegisterCount int                `json:"register_count" bson:"register_count"`
	RegisterPaid  int                `json:"register_paid" bson:"register_paid"`
	ProductCount  int                `json:"product_count" bson:"product_count"`
	TicketSummary []TicketSummary    `json:"ticket_summary" bson:"ticket_summary"`
	AmountSummary []AmountSummary    `json:"amount_summary" bson:"amount_summary"`
}

type TicketSummary struct {
	TicketID                primitive.ObjectID `json:"ticket_id" bson:"ticket_id,omitempty"`
	Title                   string             `json:"title" bson:"title"`
	RegisterCount           int                `json:"register_count" bson:"register_count"`
	PaidCount               int                `json:"paid_count" bson:"paid_count"`
	PaidWaitingApproveCount int                `json:"paid_waiting_approve_count" bson:"paid_waiting_approve_count"`
}

type AmountSummary struct {
	TicketID           primitive.ObjectID `json:"ticket_id" bson:"ticket_id,omitempty"`
	Title              string             `json:"title" bson:"title"`
	PaidSuccess        float64            `json:"paid_success" bson:"paid_success"`
	PaidWaiting        float64            `json:"paid_waiting" bson:"paid_waiting"`
	PaidWaitingApprove float64            `json:"paid_waiting_approve" bson:"paid_waiting_approve"`
}
