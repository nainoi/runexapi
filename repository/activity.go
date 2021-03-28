package repository

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"thinkdev.app/think/runex/runexapi/model"
)

type ActivityRepository interface {
	AddActivity(activity model.AddActivity) error
	GetActivityByEvent(event_user string) ([]model.ActivityInfo, error)
	GetActivityByEvent2(event_user string) (model.Activity, error)
	GetHistoryDayByEvent(event_user string, year int, month int) (model.HistoryDayInfo, error)
	HistoryMonthByEvent(event_user string, year int) ([]model.HistoryMonthInfo, error)
	DeleteActivity(event_user string, activity_id string) error
	GetActivityAllEvent(eventID string) ([]model.ActivityAllInfo, error)
}

type ActivityRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

const (
	activityCollection = "activity"
)

func (activityMongo ActivityRepositoryMongo) AddActivity(activity model.AddActivity) error {
	//model := activity.
	filter := bson.M{"event_user": activity.EventUser}
	count, err := activityMongo.ConnectionDB.Collection(activityCollection).CountDocuments(context.TODO(), filter)
	log.Printf("[info] count %s", count)
	if err != nil {
		log.Println(err)
		return err
	}
	if count > 0 {
		var activityModel model.Activity
		err := activityMongo.ConnectionDB.Collection(activityCollection).FindOne(context.TODO(), filter).Decode(&activityModel)
		var totalDistance = activityModel.ToTalDistance + activity.ActivityInfo.Distance

		updateDistance := bson.M{"$set": bson.M{"total_distance": totalDistance}}
		_, err2 := activityMongo.ConnectionDB.Collection(activityCollection).UpdateOne(context.TODO(), filter, updateDistance)
		if err2 != nil {
			log.Fatal(err2)
			log.Printf("[info] err %s", err2)
			return err2
		}
		dataInfo := activity.ActivityInfo
		dataInfo.ID = primitive.NewObjectID()
		update := bson.M{"$push": bson.M{"activity_info": dataInfo}}
		res, err := activityMongo.ConnectionDB.Collection(activityCollection).UpdateOne(context.TODO(), filter, update)
		if err != nil {
			//log.Fatal(res)
			log.Printf("[info] err %s", res)
			return err
		}

	} else {
		var arrActivityInfo []model.ActivityInfo

		dataInfo := activity.ActivityInfo
		dataInfo.ID = primitive.NewObjectID()
		arrActivityInfo = append(arrActivityInfo, dataInfo)
		activityModel := model.Activity{
			EventUser:     activity.EventUser,
			UserID:        activity.UserID,
			EventCode:     activity.EventCode,
			ActivityInfo:  arrActivityInfo,
			ToTalDistance: activity.ActivityInfo.Distance,
		}
		log.Println(activityModel)
		res, err := activityMongo.ConnectionDB.Collection(activityCollection).InsertOne(context.TODO(), activityModel)
		if err != nil {
			log.Fatal(res)
			return err
		}
	}

	return nil
}

func (activityMongo ActivityRepositoryMongo) GetActivityByEvent(event_user string) ([]model.ActivityInfo, error) {
	var activity model.Activity
	var activityInfo []model.ActivityInfo
	filter := bson.D{{"event_user", event_user}}
	err := activityMongo.ConnectionDB.Collection(activityCollection).FindOne(context.TODO(), filter).Decode(&activity)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	activityInfo = activity.ActivityInfo

	return activityInfo, err
}

// GetActivityByEvent2 repo reg event detail
func (activityMongo ActivityRepositoryMongo) GetActivityByEvent2(eventUser string) (model.Activity, error) {
	var activity model.Activity
	//var activityInfo []model.ActivityInfo
	filter := bson.D{primitive.E{Key: "event_user", Value: eventUser}}
	count, err := activityMongo.ConnectionDB.Collection(activityCollection).CountDocuments(context.TODO(), filter)
	if count == 0 {
		return model.Activity{}, err
	}
	err = activityMongo.ConnectionDB.Collection(activityCollection).FindOne(context.TODO(), filter).Decode(&activity)

	if err != nil {
		log.Println(err)
		return activity, err
	}
	//activityInfo = activity.ActivityInfo

	return activity, err
}

// GetHistoryDayByEvent repo reg history event detail
func (activityMongo ActivityRepositoryMongo) GetHistoryDayByEvent(eventUser string, year int, month int) (model.HistoryDayInfo, error) {
	var activity model.Activity
	var activityInfo []model.ActivityInfo
	var historyDayInfo model.HistoryDayInfo
	//t1 := "2019-10-01T00:00:00.000Z"
	//t2 := "2019-10-31T00:00:00.000Z"
	//filterDate := bson.D{{"$gte", t1}, {"$lt", t2}}
	//{"activity_info.distance", bson.D{{"$gt", 15}}}
	filter := bson.D{primitive.E{Key: "event_user", Value: eventUser}}
	err := activityMongo.ConnectionDB.Collection(activityCollection).FindOne(context.TODO(), filter).Decode(&activity)

	if err != nil {
		log.Println(err)
		return historyDayInfo, err
	}
	activityInfo = activity.ActivityInfo

	//var activityInfoNew []model.ActivityInfo

	historyDayInfo.Year = year
	historyDayInfo.Month = month

	var distanceDayInfo []model.DistanceDayInfo
	//now := time.Now()
	//currentYear, currentMonth, _ := now.Date()
	//currentLocation := now.Location()
	firstday := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	lastday := firstday.AddDate(0, 1, 0).Add(time.Nanosecond * -1)
	log.Printf("[info] firstday %s", firstday)
	log.Printf("[info] lastday %s", lastday)
	for _, item := range activityInfo {
		if item.ActivityDate.Before(lastday) && item.ActivityDate.After(firstday) {
			distanceDayInfoNew := model.DistanceDayInfo{
				Distance:     item.Distance,
				ActivityDate: item.ActivityDate.Format("2006-01-02"),
			}
			distanceDayInfo = append(distanceDayInfo, distanceDayInfoNew)
		}
	}

	historyDayInfo.DistanceDayInfo = distanceDayInfo

	return historyDayInfo, err
}

func (activityMongo ActivityRepositoryMongo) HistoryMonthByEvent(event_user string, year int) ([]model.HistoryMonthInfo, error) {

	var activity model.Activity
	var activityInfo []model.ActivityInfo

	filter := bson.D{{"event_user", event_user}}
	err := activityMongo.ConnectionDB.Collection(activityCollection).FindOne(context.TODO(), filter).Decode(&activity)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	activityInfo = activity.ActivityInfo

	var historyMonthInfo []model.HistoryMonthInfo

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()

	count_month := 0

	if year == currentYear {

		count_month = int(currentMonth)
	} else {
		count_month = 12
	}

	for m := 1; m <= count_month; m++ {

		firstday := time.Date(year, time.Month(m), 1, 0, 0, 0, 0, time.Local)
		lastday := firstday.AddDate(0, 1, 0).Add(time.Nanosecond * -1)

		total_step := float64(0.00)
		for _, item := range activityInfo {
			if item.ActivityDate.Before(lastday) && item.ActivityDate.After(firstday) {

				total_step = total_step + item.Distance
			}
		}
		historyMonthInfoNew := model.HistoryMonthInfo{
			Month:         m,
			MonthName:     time.Month(m).String(),
			TotalDistance: total_step,
		}

		historyMonthInfo = append(historyMonthInfo, historyMonthInfoNew)

	}

	//log.Printf("[info] firstday %s", firstday)
	//log.Printf("[info] lastday %s", lastday)

	return historyMonthInfo, err
}

//DeleteActivity delete activity
func (activityMongo ActivityRepositoryMongo) DeleteActivity(event_user string, activity_id string) error {
	objectID, _ := primitive.ObjectIDFromHex(activity_id)
	log.Printf("[info] event_user %s", event_user)
	log.Printf("[info] activity_id %s", activity_id)
	// var activity model.Activity

	// //var activityInfo []model.ActivityInfo
	// log.Printf("[info] event_user %s", event_user)

	// filter := bson.D{{"event_user", event_user}}

	// err := activityMongo.ConnectionDB.Collection(activityCollection).FindOne(context.TODO(), filter).Decode(&activity)
	// log.Println("[info] activity %s", activity)
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }

	var activity model.Activity
	var activityInfo model.ActivityInfo
	filterActivityInfo := bson.D{{"event_user", event_user}, {"activity_info._id", objectID}}
	err2 := activityMongo.ConnectionDB.Collection(activityCollection).FindOne(context.TODO(), filterActivityInfo).Decode(&activity)
	//log.Println("[info] activity %s", activity)
	if err2 != nil {
		log.Println(err2)
		return err2
	}
	for _, item := range activity.ActivityInfo {
		if objectID == item.ID {
			activityInfo = item
			break
		}
	}
	var distance = activityInfo.Distance
	var toTalDistance = activity.ToTalDistance - distance

	filterUpdate := bson.D{{"event_user", event_user}}

	delete := bson.M{"$pull": bson.M{"activity_info": bson.M{"_id": objectID}}}
	updated := bson.M{"$set": bson.M{"total_distance": toTalDistance}}

	res, err := activityMongo.ConnectionDB.Collection(activityCollection).UpdateOne(context.TODO(), filterUpdate, delete)
	if err != nil {
		//log.Fatal(res)
		log.Printf("[info] err %s", res)
		return err
	}
	log.Printf("[info] res %s", res)
	res2, err2 := activityMongo.ConnectionDB.Collection(activityCollection).UpdateOne(context.TODO(), filterUpdate, updated)
	if err2 != nil {
		log.Println(res2)
		return err2
	}
	//activityInfo = activity.ActivityInfo

	return err
}

//GetActivityAllEvent for eventer
func (activityMongo ActivityRepositoryMongo) GetActivityAllEvent(eventID string) ([]model.ActivityAllInfo, error) {
	var allInfos []model.ActivityAllInfo
	//var activityInfo []model.ActivityInfo
	objectID, _ := primitive.ObjectIDFromHex(eventID)
	filter := bson.D{primitive.E{Key: "event_id", Value: objectID}}
	curr, err := activityMongo.ConnectionDB.Collection(activityCollection).Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return allInfos, err
	}
	for curr.Next(context.TODO()) {
		var a model.Activity
		// decode the document
		if err := curr.Decode(&a); err != nil {
			log.Println(err)
		}
		//fmt.Printf("post: %+v\n", p)
		var user model.User
		filterUser := bson.D{primitive.E{Key: "_id", Value: a.UserID}}
		err := activityMongo.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filterUser).Decode(&user)

		if err != nil {
			log.Println(err)
		}
		var info = model.ActivityAllInfo{
			UserInfo: user,
			Activity: a,
		}
		allInfos = append(allInfos, info)
	}
	//activityInfo = activity.ActivityInfo

	return allInfos, err
}
