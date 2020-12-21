package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"thinkdev.app/think/runex/runexapi/config"
	"thinkdev.app/think/runex/runexapi/model"
)

//ReportRepository interface
type ReportRepository interface {
	GetDashboardByEvent(eventID string) (model.ReportDashboard, error)
}

//ReportRepositoryMongo mongo ref
type ReportRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

//PaymentState struct total payment with status
type PaymentState struct {
	ID   string       `bson:"_id"`
	Regs []model.Regs `bson:"regs"`
}

//PaymentSumState struct sum payment with status
type PaymentSumState struct {
	Status string `bson:"_id"`
	Total  int    `bson:"total"`
}

//GetDashboardByEvent repo for dashboard event
func (reportMongo ReportRepositoryMongo) GetDashboardByEvent(eventID string) (model.ReportDashboard, error) {
	log.Println(eventID)
	var dashboard model.ReportDashboard
	var registerEvent model.RegisterV2
	var register []model.Regs
	var ticketSummary []model.TicketSummary
	var amountSummary []model.AmountSummary
	objectID, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		return dashboard, err
	}
	dashboard.EventID = objectID

	var event model.EventV2
	filterEvent := bson.M{"_id": objectID}
	err = reportMongo.ConnectionDB.Collection("event_v2").FindOne(context.TODO(), filterEvent).Decode(&event)
	if err != nil {
		log.Println(err)
		return dashboard, err
	}
	var ticket []model.TicketEventV2
	ticket = event.Ticket
	var ticketTemp model.TicketSummary
	var amountTemp model.AmountSummary
	for _, item := range ticket {
		ticketTemp = model.TicketSummary{
			TicketID:                item.TicketID,
			Title:                   item.Title,
			RegisterCount:           0,
			PaidCount:               0,
			PaidWaitingApproveCount: 0,
		}
		ticketSummary = append(ticketSummary, ticketTemp)

		amountTemp = model.AmountSummary{
			TicketID:           item.TicketID,
			Title:              item.Title,
			PaidSuccess:        0.0,
			PaidWaiting:        0.0,
			PaidWaitingApprove: 0.0,
		}
		amountSummary = append(amountSummary, amountTemp)
	}

	dashboard.TicketSummary = ticketSummary
	dashboard.AmountSummary = amountSummary

	//filter := bson.D{bson.E{Key: "event_id", Value: objectID}}
	filter := bson.M{"event_id": objectID}
	err = reportMongo.ConnectionDB.Collection("register_v2").FindOne(context.TODO(), filter).Decode(&registerEvent)

	if err != nil {
		log.Println(err)
		return dashboard, err
	}
	register = registerEvent.Regs
	paid := 0.0
	waitToPay := 0.0
	waitToApprove := 0.0
	registerCount := 0
	registerPaid := 0
	productCount := 0
	for _, item := range register {
		if item.Status == "PAYMENT_SUCCESS" {
			paid = paid + item.TotalPrice
			registerPaid = registerPaid + 1
		}
		if item.Status == "PAYMENT_WAITING" {
			waitToPay = waitToPay + item.TotalPrice
		}
		if item.Status == "PAYMENT_WAITING_APPROVE" {
			waitToApprove = waitToApprove + item.TotalPrice
		}
		registerCount = registerCount + 1
	}
	dashboard.Paid = paid
	dashboard.WaitToPay = waitToPay
	dashboard.RegisterCount = registerCount
	dashboard.RegisterPaid = registerPaid
	dashboard.ProductCount = productCount

	return dashboard, nil

	//Old algorithm
	matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"event_id": objectID}}}
	//unwindStage := bson.D{{"$unwind", "$regs"}}
	//matchSubStage := bson.D{{"$match", bson.M{"regs.user_id": bson.M{"$eq": objectID}}}}
	//groupStage := bson.D{{"_id", "$_id"}, {"event_id", "$event_id"}, {"regs", bson.M{"$push": "$regs"}}}
	//filterStage := bson.D{{"$project", bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs.user_id", objectID}}}}}}}
	projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs.status", config.PAYMENT_SUCCESS}}}}}}}
	curPaid, err := reportMongo.ConnectionDB.Collection(registerCollection).Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage})
	if err != nil {
		log.Println(err.Error())
		return dashboard, err
	}
	for curPaid.Next(context.TODO()) {
		var p PaymentState

		// decode the document
		if err := curPaid.Decode(&p); err != nil {
			log.Print(err)
		}
		dashboard.RegisterPaid = len(p.Regs)
	}
	projectStage = bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs.status", config.PAYMENT_WAITING}}}}}}}
	curUnpaid, err := reportMongo.ConnectionDB.Collection(registerCollection).Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage})
	for curUnpaid.Next(context.TODO()) {
		var p PaymentState

		// decode the document
		if err := curPaid.Decode(&p); err != nil {
			log.Print(err)
		}
		dashboard.RegisterCount = len(p.Regs)
	}
	unwindStage := bson.D{bson.E{Key: "$unwind", Value: bson.M{"path": "$regs"}}}
	groupStage := bson.D{bson.E{Key: "$group", Value: bson.M{"_id": "$regs.status", "total": bson.M{"$sum": 1}}}}
	cur, err := reportMongo.ConnectionDB.Collection(registerCollection).Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage, unwindStage, groupStage})
	if err != nil {
		log.Println(err.Error())
		return dashboard, err
	}
	for cur.Next(context.TODO()) {
		var p PaymentSumState

		// decode the document
		if err := cur.Decode(&p); err != nil {
			log.Print(err)
		}
		dashboard.ProductCount = p.Total
	}
	return dashboard, nil
}
