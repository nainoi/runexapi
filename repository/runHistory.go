package repository

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"thinkdev.app/think/runex/runexapi/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RunHistoryRepository interface {
	AddHistory(userID string, history model.AddHistoryForm) error
	GetHistoryByUser(userID string) ([]model.RunHistory, error)
	DeleteActivity(user_id string, activity_id string) error
}
type RunHistoryRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

const (
	historyCollection = "run_history"
)

func (runHistoryRepositoryMongo RunHistoryRepositoryMongo) AddHistory(userID string, history model.AddHistoryForm) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	filter := bson.D{{"user_id", objectID}}
	count, err := runHistoryRepositoryMongo.ConnectionDB.Collection(historyCollection).CountDocuments(context.TODO(), filter)
	log.Printf("[info] count %s", count)
	if err != nil {
		log.Println(err)
		return err
	}
	if count > 0 {
		var historyModel model.RunHistory
		err := runHistoryRepositoryMongo.ConnectionDB.Collection(historyCollection).FindOne(context.TODO(), filter).Decode(&historyModel)
		var totalDistance = historyModel.ToTalDistance + history.Distance
		updateDistance := bson.M{"$set": bson.M{"total_distance": totalDistance}}
		_, err2 := runHistoryRepositoryMongo.ConnectionDB.Collection(historyCollection).UpdateOne(context.TODO(), filter, updateDistance)
		if err2 != nil {
			log.Printf("[info] err %s", err2)
			log.Fatal(err2)
			return err2
		}
		dataInfo := model.RunHistoryInfo{
			ID:           primitive.NewObjectID(),
			ActivityType: history.ActivityType,
			Calory:       history.Calory,
			Caption:      history.Caption,
			Distance:     history.Distance,
			Pace:         history.Pace,
			Time:         history.Time,
			ImagePath:    history.ImagePath,
			ActivityDate: time.Now(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		log.Println(dataInfo)
		log.Println(time.Now())
		update := bson.M{"$push": bson.M{"activity_info": dataInfo}}
		res, err := runHistoryRepositoryMongo.ConnectionDB.Collection(historyCollection).UpdateOne(context.TODO(), filter, update)
		if err != nil {
			//log.Fatal(res)
			log.Printf("[info] err %s", res)
			return err
		}
	} else {
		var arrHistoryInfo []model.RunHistoryInfo
		dataInfo := model.RunHistoryInfo{
			ID:           primitive.NewObjectID(),
			ActivityType: history.ActivityType,
			Calory:       history.Calory,
			Caption:      history.Caption,
			Distance:     history.Distance,
			Pace:         history.Pace,
			Time:         history.Time,
			ImagePath:    history.ImagePath,
			ActivityDate: time.Now(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		arrHistoryInfo = append(arrHistoryInfo, dataInfo)
		historyModel := model.RunHistory{
			UserID:         objectID,
			RunHistoryInfo: arrHistoryInfo,
			ToTalDistance:  history.Distance,
		}
		log.Println(historyModel)
		res, err := runHistoryRepositoryMongo.ConnectionDB.Collection(historyCollection).InsertOne(context.TODO(), historyModel)
		if err != nil {
			log.Fatal(res)
			return err
		}
	}

	return nil
}

func (runHistoryRepositoryMongo RunHistoryRepositoryMongo) GetHistoryByUser(userID string) ([]model.RunHistory, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	var history []model.RunHistory
	filter := bson.D{{"user_id", objectID}}

	options := options.Find()
	options.SetSort(bson.D{{"activity_info.activity_date", -1}})
	//options.SetSort(map[string]int{"activity_info.created_at": -1})
	options.SetSkip(0)
	options.SetLimit(10)

	cur, err := runHistoryRepositoryMongo.ConnectionDB.Collection(historyCollection).Find(context.TODO(), filter, options)
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {
		var u model.RunHistory
		// decode the document
		if err := cur.Decode(&u); err != nil {
			log.Println(err)
			log.Fatal(err)
		}
		//fmt.Printf("post: %+v\n", p)
		history = append(history, u)
	}

	return history, err
}

func (runHistoryRepositoryMongo RunHistoryRepositoryMongo) DeleteActivity(user_id string, activity_id string) error {
	objectID, _ := primitive.ObjectIDFromHex(activity_id)
	userObjectID, _ := primitive.ObjectIDFromHex(user_id)
	log.Printf("[info] user_id %s", user_id)
	log.Printf("[info] activity_id %s", activity_id)

	var runHistory model.RunHistory
	var runHistoryInfo model.RunHistoryInfo
	filterActivityInfo := bson.D{{"user_id", userObjectID}, {"activity_info._id", objectID}}
	err2 := runHistoryRepositoryMongo.ConnectionDB.Collection(historyCollection).FindOne(context.TODO(), filterActivityInfo).Decode(&runHistory)
	//log.Println("[info] activity %s", activity)
	if err2 != nil {
		log.Println(err2)
		return err2
	}
	for _, item := range runHistory.RunHistoryInfo {
		if objectID == item.ID {
			runHistoryInfo = item
			break
		}
	}
	var distance = runHistoryInfo.Distance
	var toTalDistance = runHistory.ToTalDistance - distance

	filterUpdate := bson.D{{"user_id", userObjectID}}

	delete := bson.M{"$pull": bson.M{"activity_info": bson.M{"_id": objectID}}}
	updated := bson.M{"$set": bson.M{"total_distance": toTalDistance}}

	res, err := runHistoryRepositoryMongo.ConnectionDB.Collection(historyCollection).UpdateOne(context.TODO(), filterUpdate, delete)
	if err != nil {
		//log.Fatal(res)
		log.Printf("[info] err %s", res)
		return err
	}
	log.Printf("[info] res %s", res)
	res2, err2 := runHistoryRepositoryMongo.ConnectionDB.Collection(historyCollection).UpdateOne(context.TODO(), filterUpdate, updated)
	if err2 != nil {
		log.Fatal(res2)
		return err2
	}
	//activityInfo = activity.ActivityInfo

	return err
}
