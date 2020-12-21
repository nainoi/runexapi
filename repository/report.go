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

	var dashboard model.ReportDashboard
	objectID, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		return dashboard, err
	}
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
