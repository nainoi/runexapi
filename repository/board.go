package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"thinkdev.app/think/runex/runexapi/model"
)

type BoardRepository interface {
	GetBoardByEvent(eventID string, userID string) ([]model.Ranking, []model.Ranking, error)
}

type BoardRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

const (
	boardCollection = "activity"
)

//GetBoardByEvent board ranking
func (boardMongo BoardRepositoryMongo) GetBoardByEvent(eventID string, userID string) ([]model.Ranking, []model.Ranking, error) {
	//var activity model.Activity
	//var activityInfo []model.ActivityInfo
	objectEventID, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		log.Fatal(err)
	}
	var activities []model.Ranking
	var temps []model.Ranking
	filter := bson.D{{"event_id", objectEventID}}
	option := options.Find()
	//option.SetLimit(10)
	option.SetSort(bson.D{{"total_distance", -1}})
	cur, err := boardMongo.ConnectionDB.Collection(activityCollection).Find(context.TODO(), filter, option)

	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	objectUserID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		log.Fatal(err)
	}

	n := 0

	for cur.Next(context.TODO()) {
		// if n >= 10 {
		// 	break
		// }
		var a model.Ranking
		// decode the document
		if err := cur.Decode(&a); err != nil {
			log.Fatal(err)
		}
		if n < 10 {
			var user model.UserEvent
			// log.Printf("[info] userID %s", userID)
			filterUser := bson.D{{"_id", a.UserID}}
			err := boardMongo.ConnectionDB.Collection("user").FindOne(context.TODO(), filterUser).Decode(&user)

			if err != nil {
				log.Println(err)
			}
			a.UserInfo = user
			activities = append(activities, a)
		}

		temps = append(temps, a)

		n++
	}

	myActivities := []model.Ranking{}

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
			filterUser := bson.D{{"_id", ranking.UserID}}
			err := boardMongo.ConnectionDB.Collection("user").FindOne(context.TODO(), filterUser).Decode(&user)

			if err != nil {
				log.Println(err)
			}
			ranking.UserInfo = user
			myActivities = append(myActivities, ranking)
		}
		if (index - 1) >= 0 {
			var user model.UserEvent
			ranking := temps[index-1]
			filterUser := bson.D{{"_id", ranking.UserID}}
			err := boardMongo.ConnectionDB.Collection("user").FindOne(context.TODO(), filterUser).Decode(&user)

			if err != nil {
				log.Println(err)
			}
			ranking.UserInfo = user
			myActivities = append(myActivities, ranking)
		}
		var user model.UserEvent
		ranking := temps[index]
		filterUser := bson.D{{"_id", ranking.UserID}}
		err := boardMongo.ConnectionDB.Collection("user").FindOne(context.TODO(), filterUser).Decode(&user)

		if err != nil {
			log.Println(err)
		}
		ranking.UserInfo = user
		myActivities = append(myActivities, ranking)
		if (index + 1) < len(temps) {
			var user model.UserEvent
			ranking := temps[index+1]
			filterUser := bson.D{{"_id", ranking.UserID}}
			err := boardMongo.ConnectionDB.Collection("user").FindOne(context.TODO(), filterUser).Decode(&user)

			if err != nil {
				log.Println(err)
			}
			ranking.UserInfo = user
			myActivities = append(myActivities, ranking)
		}
		if (index + 2) < len(temps) {
			var user model.UserEvent
			ranking := temps[index+2]
			filterUser := bson.D{{"_id", ranking.UserID}}
			err := boardMongo.ConnectionDB.Collection("user").FindOne(context.TODO(), filterUser).Decode(&user)

			if err != nil {
				log.Println(err)
			}
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

	return activities, myActivities, err
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
