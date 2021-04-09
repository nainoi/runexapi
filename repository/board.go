package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"thinkdev.app/think/runex/runexapi/model"
	v2 "thinkdev.app/think/runex/runexapi/repository/v2"
)

// BoardRepository interface
type BoardRepository interface {
	GetBoardByEvent(req model.RankingRequest, userID string) (model.Event, int64, []model.Ranking, []model.Ranking, error)
}

// BoardRepositoryMongo db struct
type BoardRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

const (
	boardCollection = "activity"
)

//GetBoardByEvent board ranking
func (boardMongo BoardRepositoryMongo) GetBoardByEvent(req model.RankingRequest, userID string) (model.Event, int64, []model.Ranking, []model.Ranking, error) {
	//var activity model.Activity
	//var activityInfo []model.ActivityInfo
	var event = model.Event{}
	var activities = []model.Ranking{}
	myActivities := []model.Ranking{}

	// objectEventID, err := primitive.ObjectIDFromHex(eventID)
	// if err != nil {
	// 	log.Println(err)
	// 	return event, 0, activities, myActivities, err
	// }
	objectUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return event, 0, activities, myActivities, err
	}

	var temps = []model.Ranking{}

	//filter := bson.D{primitive.E{Key: "event_id", Value: objectEventID}, primitive.E{Key: "activities.user_id", Value: objectUserID}}
	filter := bson.D{primitive.E{Key: "event_code", Value: req.EventCode}, primitive.E{Key: "ticket.id", Value: req.TicketID}}

	option := options.Find()

	//filterEvent := bson.D{primitive.E{Key: "_id", Value: objectEventID}}

	event, err = v2.DetailEventByCode(req.EventCode)
	if err != nil {
		return event, 0, activities, myActivities, err
	}
	//option.SetLimit(10)
	count, err := boardMongo.ConnectionDB.Collection(activityV2Collection).CountDocuments(context.TODO(), filter)
	if count == 0 {
		return event, count, activities, myActivities, err
	}
	//filter = bson.D{primitive.E{Key: "event_id", Value: objectEventID}, primitive.E{Key: "activities.user_id", Value: objectUserID}}
	option.SetSort(bson.D{primitive.E{Key: "total_distance", Value: -1}})
	cur, err := boardMongo.ConnectionDB.Collection(activityV2Collection).Find(context.TODO(), filter, option)

	if err != nil {
		log.Println(err)
		return event, count, activities, myActivities, err
	}

	//objectUserID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		log.Println(err)
		return event, count, activities, myActivities, err
	}

	n := 0

	for cur.Next(context.TODO()) {
		// if n >= 10 {
		// 	break
		// }

		var activity model.ActivityV2
		var a model.Ranking
		// decode the document
		if err := cur.Decode(&activity); err != nil {
			log.Fatal(err)
		}
		if n < 10 {
			var user model.UserEvent
			// log.Printf("[info] userID %s", userID)
			filterUser := bson.D{primitive.E{Key: "_id", Value: activity.UserID}}
			err := boardMongo.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filterUser).Decode(&user)

			if err != nil {
				log.Println(err)
			}
			a = model.Ranking{
				UserID:        activity.UserID,
				ActivityInfo:  activity.ActivityInfo,
				ToTalDistance: activity.ToTalDistance,
				UserInfo:      user,
				EventCode:     req.EventCode,
				ID:            activity.ID,
			}
			a.RankNo = n + 1
			//a.UserInfo = user
			activities = append(activities, a)
		}

		temps = append(temps, a)

		n++
	}

	index := -1
	for n, s := range temps {
		if s.UserID == objectUserID {
			index = n
			break
		}
	}

	//index := SliceIndex(len(activities), func(i int) bool { return activities[i].UserID == objectUserID })
	if index != -1 {
		if (index - 2) >= 0 {
			var user model.UserEvent
			// log.Printf("[info] userID %s", userID)
			ranking := temps[index-2]
			filterUser := bson.D{primitive.E{Key: "_id", Value: ranking.UserID}}
			err := boardMongo.ConnectionDB.Collection("user").FindOne(context.TODO(), filterUser).Decode(&user)

			if err != nil {
				log.Println(err)
			}
			ranking.RankNo = index - 1
			ranking.UserInfo = user
			myActivities = append(myActivities, ranking)
		}
		if (index - 1) >= 0 {
			var user model.UserEvent
			ranking := temps[index-1]
			filterUser := bson.D{primitive.E{Key: "_id", Value: ranking.UserID}}
			err := boardMongo.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filterUser).Decode(&user)

			if err != nil {
				log.Println(err)
			}
			ranking.RankNo = index
			ranking.UserInfo = user
			myActivities = append(myActivities, ranking)
		}
		var user model.UserEvent
		ranking := temps[index]
		filterUser := bson.D{primitive.E{Key: "_id", Value: ranking.UserID}}
		err := boardMongo.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filterUser).Decode(&user)

		if err != nil {
			log.Println(err)
		}
		ranking.RankNo = index + 1
		ranking.UserInfo = user
		myActivities = append(myActivities, ranking)
		if (index + 1) < len(temps) {
			var user model.UserEvent
			ranking := temps[index+1]
			filterUser := bson.D{primitive.E{Key: "_id", Value: ranking.UserID}}
			err := boardMongo.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filterUser).Decode(&user)

			if err != nil {
				log.Println(err)
			}
			ranking.RankNo = index + 2
			ranking.UserInfo = user
			myActivities = append(myActivities, ranking)
		}
		if (index + 2) < len(temps) {
			var user model.UserEvent
			ranking := temps[index+2]
			filterUser := bson.D{primitive.E{Key: "_id", Value: ranking.UserID}}
			err := boardMongo.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filterUser).Decode(&user)

			if err != nil {
				log.Println(err)
			}
			ranking.RankNo = index + 3
			ranking.UserInfo = user
			myActivities = append(myActivities, ranking)
		}
	}

	// //activityInfo = activity.ActivityInfo

	// filter = bson.D{{"event_id", objectEventID}}
	// option := options.Find()
	// option.SetLimit(10)
	// option.SetSort(bson.D{{"total_distance", -1}})
	// cur, err := boardMongo.ConnectionDB.Collection(activityCollection).Find(context.TODO(), filter, option)

	return event, count, activities, myActivities, err
}

//SliceIndex get index array object
func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}
