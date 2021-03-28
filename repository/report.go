package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"thinkdev.app/think/runex/runexapi/config"
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/repository/v2"
)

//ReportRepository interface
type ReportRepository interface {
	GetDashboardByEvent(eventID string) (model.ReportDashboard, error)
	GetDashboardByEventCode(code model.Event) (model.ReportDashboard, error)
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
	// objectID, err := primitive.ObjectIDFromHex(eventID)
	// if err != nil {
	// 	return dashboard, err
	// }

	var event model.Event
	event,err := repository.DetailEventByCode(eventID)
	if err != nil {
		log.Println(err)
		return dashboard, err
	}
	//dashboard.EventID = event.Event.ID
	dashboard.Code = eventID
	var ticket []model.Tickets
	ticket = event.Tickets
	var ticketTemp model.TicketSummary
	var amountTemp model.AmountSummary
	for _, item := range ticket {

		// matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"event_id": objectID}}}
		// projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$and": bson.A{bson.M{"$eq": bson.A{"$$regs.status", config.PAYMENT_SUCCESS}}, bson.M{"$eq": bson.A{"$$regs.ticket_options.tickets.ticket_name", "10 KM"}}}}}}}}}
		// //unwindStage := bson.D{bson.E{Key: "$unwind", Value: "$regs.ticket_options.tickets"}}
		// //projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs.status", config.PAYMENT_SUCCESS}}}}}}}
		// curPaid, err := reportMongo.ConnectionDB.Collection("register_v2").Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage})
		// if err != nil {
		// 	log.Println(err.Error())
		// 	return dashboard, err
		// }
		// log.Print("[info] ticketSummary Id %s", item.TicketID)
		// log.Print("[info] ticketTitle %s", item.Title)
		// paidCount := 0
		// for curPaid.Next(context.TODO()) {
		// 	var p PaymentState

		// 	// decode the document
		// 	if err := curPaid.Decode(&p); err != nil {
		// 		log.Print(err)
		// 	}

		// 	paidCount = len(p.Regs)
		// 	log.Print("[info] PaymentState %s", p)
		// }

		ticketTemp = model.TicketSummary{
			TicketID:                item.ID,
			Title:                   item.Title,
			RegisterCount:           0,
			PaidCount:               0,
			PaidWaitingApproveCount: 0,
		}
		ticketSummary = append(ticketSummary, ticketTemp)

		amountTemp = model.AmountSummary{
			TicketID:           item.ID,
			Title:              item.Title,
			PaidSuccess:        0.0,
			PaidWaiting:        0.0,
			PaidWaitingApprove: 0.0,
		}
		amountSummary = append(amountSummary, amountTemp)
	}

	dashboard.TicketSummary = ticketSummary
	dashboard.AmountSummary = amountSummary

	// for _, item := range ticketSummary {

	// 	matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"event_id": objectID}}}
	// 	projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$and": bson.A{bson.M{"$eq": bson.A{"$$regs.status", config.PAYMENT_SUCCESS}}, bson.M{"$eq": bson.A{"$$regs.ticket_options.tickets.ticket_id", item.TicketID}}}}}}}}}
	// 	curPaid, err := reportMongo.ConnectionDB.Collection("register_v2").Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage})
	// 	if err != nil {
	// 		log.Println(err.Error())
	// 		return dashboard, err
	// 	}
	// 	var p PaymentState
	// 	log.Print("[info] ticketSummary Id %s", item.TicketID)
	// 	// decode the document
	// 	if err := curPaid.Decode(&p); err != nil {
	// 		log.Print(err)
	// 	}
	// }

	filter := bson.M{"event_code": event.Code}
	count, err := reportMongo.ConnectionDB.Collection("register_v2").CountDocuments(context.TODO(), filter)
	if count == 0 {
		return dashboard, err
	}
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
	//var userTicket []model.RegisterTicketV2
	for _, item := range register {

		userTicket := item.TicketOptions

		if item.Status == "PAYMENT_SUCCESS" {
			paid = paid + item.TotalPrice
			registerPaid = registerPaid + 1

			for _, item2 := range userTicket {

				indexArr := 0
				for _, item3 := range ticketSummary {

					if item2.Tickets.ID == item3.ID {
						dashboard.TicketSummary[indexArr].RegisterCount = dashboard.TicketSummary[indexArr].RegisterCount + 1
						dashboard.TicketSummary[indexArr].PaidCount = dashboard.TicketSummary[indexArr].PaidCount + 1
						dashboard.AmountSummary[indexArr].PaidSuccess = dashboard.AmountSummary[indexArr].PaidSuccess + item2.TotalPrice
						break
					}
					indexArr = indexArr + 1
				}
			}
		}
		if item.Status == "PAYMENT_WAITING" {
			waitToPay = waitToPay + item.TotalPrice

			for _, item2 := range userTicket {

				indexArr := 0
				for _, item3 := range ticketSummary {

					if item2.Tickets.ID == item3.ID {
						dashboard.TicketSummary[indexArr].RegisterCount = dashboard.TicketSummary[indexArr].RegisterCount + 1
						dashboard.AmountSummary[indexArr].PaidWaiting = dashboard.AmountSummary[indexArr].PaidWaiting + item2.TotalPrice
						break
					}
					indexArr = indexArr + 1
				}
			}
		}
		if item.Status == "PAYMENT_WAITING_APPROVE" {
			waitToApprove = waitToApprove + item.TotalPrice

			for _, item2 := range userTicket {

				indexArr := 0
				for _, item3 := range ticketSummary {

					if item2.Tickets.ID == item3.ID {
						dashboard.TicketSummary[indexArr].RegisterCount = dashboard.TicketSummary[indexArr].RegisterCount + 1
						dashboard.TicketSummary[indexArr].PaidWaitingApproveCount = dashboard.TicketSummary[indexArr].PaidWaitingApproveCount + 1
						dashboard.AmountSummary[indexArr].PaidWaitingApprove = dashboard.AmountSummary[indexArr].PaidWaitingApprove + item2.TotalPrice
						break
					}
					indexArr = indexArr + 1
				}
			}
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
	matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"event_code": event.Code}}}
	//unwindStage := bson.D{{"$unwind", "$regs"}}
	//matchSubStage := bson.D{{"$match", bson.M{"regs.user_id": bson.M{"$eq": objectID}}}}
	//groupStage := bson.D{{"_id", "$_id"}, {"event_id", "$event_id"}, {"regs", bson.M{"$push": "$regs"}}}
	//filterStage := bson.D{{"$project", bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs.user_id", objectID}}}}}}}
	projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs.status", config.PAYMENT_SUCCESS}}}}}}}
	curPaid, err := reportMongo.ConnectionDB.Collection("register_v2").Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage})
	if err != nil {
		log.Println(err.Error())
		return dashboard, err
	}
	//log.Printf("[info] curPaid %s", curPaid)
	for curPaid.Next(context.TODO()) {
		var p PaymentState

		// decode the document
		if err := curPaid.Decode(&p); err != nil {
			log.Print(err)
		}
		dashboard.RegisterPaid = len(p.Regs)
	}
	projectStage = bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs.status", config.PAYMENT_WAITING}}}}}}}
	curUnpaid, err := reportMongo.ConnectionDB.Collection("register_v2").Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage})
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
	cur, err := reportMongo.ConnectionDB.Collection("register_v2").Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage, unwindStage, groupStage})
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

//GetDashboardByEventCode repo for dashboard event
func (reportMongo ReportRepositoryMongo) GetDashboardByEventCode(event model.Event) (model.ReportDashboard, error) {
	var dashboard model.ReportDashboard
	var registerEvent model.RegisterV2
	var register []model.Regs
	var ticketSummary []model.TicketSummary
	var amountSummary []model.AmountSummary

	//var ticket []model.EventData.
	ticket := event.Tickets
	var ticketTemp model.TicketSummary
	var amountTemp model.AmountSummary
	for _, item := range ticket {

		// matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"event_id": objectID}}}
		// projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$and": bson.A{bson.M{"$eq": bson.A{"$$regs.status", config.PAYMENT_SUCCESS}}, bson.M{"$eq": bson.A{"$$regs.ticket_options.tickets.ticket_name", "10 KM"}}}}}}}}}
		// //unwindStage := bson.D{bson.E{Key: "$unwind", Value: "$regs.ticket_options.tickets"}}
		// //projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs.status", config.PAYMENT_SUCCESS}}}}}}}
		// curPaid, err := reportMongo.ConnectionDB.Collection("register_v2").Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage})
		// if err != nil {
		// 	log.Println(err.Error())
		// 	return dashboard, err
		// }
		// log.Print("[info] ticketSummary Id %s", item.TicketID)
		// log.Print("[info] ticketTitle %s", item.Title)
		// paidCount := 0
		// for curPaid.Next(context.TODO()) {
		// 	var p PaymentState

		// 	// decode the document
		// 	if err := curPaid.Decode(&p); err != nil {
		// 		log.Print(err)
		// 	}

		// 	paidCount = len(p.Regs)
		// 	log.Print("[info] PaymentState %s", p)
		// }

		ticketTemp = model.TicketSummary{
			ID:                      item.ID,
			Title:                   item.Title,
			RegisterCount:           0,
			PaidCount:               0,
			PaidWaitingApproveCount: 0,
		}
		ticketSummary = append(ticketSummary, ticketTemp)

		amountTemp = model.AmountSummary{
			ID:                 item.ID,
			Title:              item.Title,
			PaidSuccess:        0.0,
			PaidWaiting:        0.0,
			PaidWaitingApprove: 0.0,
		}
		amountSummary = append(amountSummary, amountTemp)
	}

	dashboard.TicketSummary = ticketSummary
	dashboard.AmountSummary = amountSummary

	// for _, item := range ticketSummary {

	// 	matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"event_id": objectID}}}
	// 	projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$and": bson.A{bson.M{"$eq": bson.A{"$$regs.status", config.PAYMENT_SUCCESS}}, bson.M{"$eq": bson.A{"$$regs.ticket_options.tickets.ticket_id", item.TicketID}}}}}}}}}
	// 	curPaid, err := reportMongo.ConnectionDB.Collection("register_v2").Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage})
	// 	if err != nil {
	// 		log.Println(err.Error())
	// 		return dashboard, err
	// 	}
	// 	var p PaymentState
	// 	log.Print("[info] ticketSummary Id %s", item.TicketID)
	// 	// decode the document
	// 	if err := curPaid.Decode(&p); err != nil {
	// 		log.Print(err)
	// 	}
	// }

	filter := bson.M{"event_code": event.Code}
	err := reportMongo.ConnectionDB.Collection("register_v2").FindOne(context.TODO(), filter).Decode(&registerEvent)

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
	//var userTicket []model.RegisterTicketV2
	for _, item := range register {

		userTicket := item.TicketOptions

		if item.Status == "PAYMENT_SUCCESS" {
			paid = paid + item.TotalPrice
			registerPaid = registerPaid + 1

			for _, item2 := range userTicket {

				indexArr := 0
				for _, item3 := range ticketSummary {

					if item2.Tickets.ID == item3.TicketID {
						dashboard.TicketSummary[indexArr].RegisterCount = dashboard.TicketSummary[indexArr].RegisterCount + 1
						dashboard.TicketSummary[indexArr].PaidCount = dashboard.TicketSummary[indexArr].PaidCount + 1
						dashboard.AmountSummary[indexArr].PaidSuccess = dashboard.AmountSummary[indexArr].PaidSuccess + item2.TotalPrice
						break
					}
					indexArr = indexArr + 1
				}
			}
		}
		if item.Status == "PAYMENT_WAITING" {
			waitToPay = waitToPay + item.TotalPrice

			for _, item2 := range userTicket {

				indexArr := 0
				for _, item3 := range ticketSummary {

					if item2.Tickets.ID == item3.ID {
						dashboard.TicketSummary[indexArr].RegisterCount = dashboard.TicketSummary[indexArr].RegisterCount + 1
						dashboard.AmountSummary[indexArr].PaidWaiting = dashboard.AmountSummary[indexArr].PaidWaiting + item2.TotalPrice
						break
					}
					indexArr = indexArr + 1
				}
			}
		}
		if item.Status == "PAYMENT_WAITING_APPROVE" {
			waitToApprove = waitToApprove + item.TotalPrice

			for _, item2 := range userTicket {

				indexArr := 0
				for _, item3 := range ticketSummary {

					if item2.Tickets.ID == item3.ID {
						dashboard.TicketSummary[indexArr].RegisterCount = dashboard.TicketSummary[indexArr].RegisterCount + 1
						dashboard.TicketSummary[indexArr].PaidWaitingApproveCount = dashboard.TicketSummary[indexArr].PaidWaitingApproveCount + 1
						dashboard.AmountSummary[indexArr].PaidWaitingApprove = dashboard.AmountSummary[indexArr].PaidWaitingApprove + item2.TotalPrice
						break
					}
					indexArr = indexArr + 1
				}
			}
		}
		registerCount = registerCount + 1
	}
	dashboard.Paid = paid
	dashboard.WaitToPay = waitToPay
	dashboard.RegisterCount = registerCount
	dashboard.RegisterPaid = registerPaid
	dashboard.ProductCount = productCount

	return dashboard, nil

	matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"event_code": event.Code}}}
	//Old algorithm

	//unwindStage := bson.D{{"$unwind", "$regs"}}
	//matchSubStage := bson.D{{"$match", bson.M{"regs.user_id": bson.M{"$eq": objectID}}}}
	//groupStage := bson.D{{"_id", "$_id"}, {"event_id", "$event_id"}, {"regs", bson.M{"$push": "$regs"}}}
	//filterStage := bson.D{{"$project", bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs.user_id", objectID}}}}}}}
	projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs.status", config.PAYMENT_SUCCESS}}}}}}}
	curPaid, err := reportMongo.ConnectionDB.Collection("register_v2").Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage})
	if err != nil {
		log.Println(err.Error())
		return dashboard, err
	}
	//log.Printf("[info] curPaid %s", curPaid)
	for curPaid.Next(context.TODO()) {
		var p PaymentState

		// decode the document
		if err := curPaid.Decode(&p); err != nil {
			log.Print(err)
		}
		//log.Println("[info] p %s", err.Error())
		dashboard.RegisterPaid = len(p.Regs)
	}
	projectStage = bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs.status", config.PAYMENT_WAITING}}}}}}}
	curUnpaid, err := reportMongo.ConnectionDB.Collection("register_v2").Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage})
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
	cur, err := reportMongo.ConnectionDB.Collection("register_v2").Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage, unwindStage, groupStage})
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
