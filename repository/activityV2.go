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

// ActivityV2Repository interface repo
type ActivityV2Repository interface {
	AddActivity(activity model.AddActivityV2) error
	GetActivityByEvent(event_id string, user_id string) ([]model.ActivityInfo, error)
	GetActivityByEvent2(event_id string, user_id string) (model.ActivityV2, error)
	GetHistoryDayByEvent(event_id string, user_id string, year int, month int) (model.HistoryDayInfo, error)
	HistoryMonthByEvent(event_id string, user_id string, year int) ([]model.HistoryMonthInfo, error)
	DeleteActivity(event_id string, user_id string, activity_id string) error
	UpdateWorkout(workout model.WorkoutActivityInfo, userID primitive.ObjectID) error
	AddKaoLogActivity(activity model.LogSendKaoActivity) error
}

// ActivityV2RepositoryMongo db ref
type ActivityV2RepositoryMongo struct {
	ConnectionDB *mongo.Database
}

const (
	activityV2Collection = "activityV2"
	activityLog          = "activity_log"
	activityKaoLog       = "activity_kao_log"
)

// AddActivity repo add activity
func (activityMongo ActivityV2RepositoryMongo) AddActivity(activity model.AddActivityV2) error {
	//model := activity.
	//filter := bson.D{"event_id": activity.EventID}
	filter := bson.D{primitive.E{Key: "event_id", Value: activity.EventID}, primitive.E{Key: "activities.user_id", Value: activity.UserID}}
	count, err := activityMongo.ConnectionDB.Collection(activityV2Collection).CountDocuments(context.TODO(), filter)
	log.Printf("[info] count %d", count)
	if err != nil {
		log.Println(err)
		return err
	}
	if count > 0 {
		var activityModel model.ActivityV2
		err := activityMongo.ConnectionDB.Collection(activityV2Collection).FindOne(context.TODO(), filter).Decode(&activityModel)
		var totalDistance = activityModel.Activities.ToTalDistance + activity.ActivityInfo.Distance

		updateDistance := bson.M{"$set": bson.M{"activities.total_distance": totalDistance}}
		_, err2 := activityMongo.ConnectionDB.Collection(activityV2Collection).UpdateOne(context.TODO(), filter, updateDistance)
		if err2 != nil {
			log.Fatal(err2)
			log.Printf("[info] err %s", err2)
			return err2
		}
		dataInfo := activity.ActivityInfo
		update := bson.M{"$push": bson.M{"activities.activity_info": dataInfo}}
		_, err = activityMongo.ConnectionDB.Collection(activityV2Collection).UpdateOne(context.TODO(), filter, update)
		if err != nil {
			//log.Fatal(res)
			//log.Printf("[info] err %s", res)
			return err
		}
		activityLogInfo := model.LogActivityInfo{
			UserID:         activity.UserID,
			ActivityInfoID: dataInfo.ID,
			Distance:       dataInfo.Distance,
			ImageURL:       dataInfo.ImageURL,
			Caption:        dataInfo.Caption,
			Time:           dataInfo.Time,
			APP:            dataInfo.APP,
			ActivityDate:   dataInfo.ActivityDate,
			CreatedAt:      dataInfo.CreatedAt,
			UpdatedAt:      dataInfo.UpdatedAt,
		}
		_, _ = activityMongo.ConnectionDB.Collection(activityLog).InsertOne(context.TODO(), activityLogInfo)

	} else {
		var arrActivityInfo []model.ActivityInfo

		dataInfo := activity.ActivityInfo
		arrActivityInfo = append(arrActivityInfo, dataInfo)

		activities := model.Activities{
			UserID:        activity.UserID,
			ToTalDistance: activity.ActivityInfo.Distance,
			ActivityInfo:  arrActivityInfo,
		}

		activityModel := model.ActivityV2{
			EventID:    activity.EventID,
			Activities: activities,
		}
		log.Println(activityModel)
		res, err := activityMongo.ConnectionDB.Collection(activityV2Collection).InsertOne(context.TODO(), activityModel)
		if err != nil {
			log.Fatal(res)
			return err
		}

		activityLogInfo := model.LogActivityInfo{
			UserID:         activity.UserID,
			EventID:        activity.EventID,
			ActivityInfoID: dataInfo.ID,
			Distance:       dataInfo.Distance,
			ImageURL:       dataInfo.ImageURL,
			Caption:        dataInfo.Caption,
			Time:           dataInfo.Time,
			APP:            dataInfo.APP,
			ActivityDate:   dataInfo.ActivityDate,
			CreatedAt:      dataInfo.CreatedAt,
			UpdatedAt:      dataInfo.UpdatedAt,
		}
		_, _ = activityMongo.ConnectionDB.Collection(activityLog).InsertOne(context.TODO(), activityLogInfo)
	}

	return nil
}

// GetActivityByEvent event and activity detail
func (activityMongo ActivityV2RepositoryMongo) GetActivityByEvent(eventID string, userID string) ([]model.ActivityInfo, error) {
	var activity model.ActivityV2
	var activityInfo = []model.ActivityInfo{}
	userObjectID, _ := primitive.ObjectIDFromHex(userID)
	eventObjectID, _ := primitive.ObjectIDFromHex(eventID)
	filter := bson.D{primitive.E{Key: "event_id", Value: eventObjectID}, primitive.E{Key: "activities.user_id", Value: userObjectID}}
	count, err := activityMongo.ConnectionDB.Collection(activityV2Collection).CountDocuments(context.TODO(), filter)
	if count > 0 {
		err = activityMongo.ConnectionDB.Collection(activityV2Collection).FindOne(context.TODO(), filter).Decode(&activity)

		if err != nil {
			log.Println(err)
			return activityInfo, err
		}
		//activityInfo = activity.ActivityInfo

		return activityInfo, err
	}
	return activityInfo, err
}

// GetActivityByEvent2 event and activity detail
func (activityMongo ActivityV2RepositoryMongo) GetActivityByEvent2(eventID string, userID string) (model.ActivityV2, error) {
	var activity = model.ActivityV2{}
	userObjectID, _ := primitive.ObjectIDFromHex(userID)
	eventObjectID, _ := primitive.ObjectIDFromHex(eventID)
	filter := bson.D{primitive.E{Key: "event_id", Value: eventObjectID}, primitive.E{Key: "activities.user_id", Value: userObjectID}}
	count, err := activityMongo.ConnectionDB.Collection(activityV2Collection).CountDocuments(context.TODO(), filter)
	if count > 0 {
		err = activityMongo.ConnectionDB.Collection(activityV2Collection).FindOne(context.TODO(), filter).Decode(&activity)

		if err != nil {
			log.Println(err)
			return activity, err
		}
		//activityInfo = activity.ActivityInfo

		return activity, err
	}
	return activity, err
}

func (activityMongo ActivityV2RepositoryMongo) GetHistoryDayByEvent(event_id string, user_id string, year int, month int) (model.HistoryDayInfo, error) {
	var activity model.ActivityV2
	var activityInfo []model.ActivityInfo
	var historyDayInfo model.HistoryDayInfo
	//t1 := "2019-10-01T00:00:00.000Z"
	//t2 := "2019-10-31T00:00:00.000Z"
	//filterDate := bson.D{{"$gte", t1}, {"$lt", t2}}
	//{"activity_info.distance", bson.D{{"$gt", 15}}}
	userObjectID, _ := primitive.ObjectIDFromHex(user_id)
	eventObjectID, _ := primitive.ObjectIDFromHex(event_id)
	filter := bson.D{{"event_id", eventObjectID}, {"activities.user_id", userObjectID}}

	err := activityMongo.ConnectionDB.Collection(activityV2Collection).FindOne(context.TODO(), filter).Decode(&activity)

	if err != nil {
		log.Println(err)
		return historyDayInfo, err
	}
	activityInfo = activity.Activities.ActivityInfo

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

func (activityMongo ActivityV2RepositoryMongo) HistoryMonthByEvent(event_id string, user_id string, year int) ([]model.HistoryMonthInfo, error) {

	var activity model.ActivityV2
	var activityInfo []model.ActivityInfo

	userObjectID, _ := primitive.ObjectIDFromHex(user_id)
	eventObjectID, _ := primitive.ObjectIDFromHex(event_id)
	filter := bson.D{{"event_id", eventObjectID}, {"activities.user_id", userObjectID}}

	err := activityMongo.ConnectionDB.Collection(activityV2Collection).FindOne(context.TODO(), filter).Decode(&activity)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	activityInfo = activity.Activities.ActivityInfo

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

func (activityMongo ActivityV2RepositoryMongo) DeleteActivity(event_id string, user_id string, activity_id string) error {
	objectID, _ := primitive.ObjectIDFromHex(activity_id)

	var activity model.ActivityV2
	var activityInfo model.ActivityInfo

	userObjectID, _ := primitive.ObjectIDFromHex(user_id)
	eventObjectID, _ := primitive.ObjectIDFromHex(event_id)
	//filter := bson.D{{"event_id", eventObjectID}, {"activities.user_id", userObjectID}}

	filterActivityInfo := bson.D{{"event_id", eventObjectID}, {"activities.user_id", userObjectID}, {"activities.activity_info._id", objectID}}
	err2 := activityMongo.ConnectionDB.Collection(activityV2Collection).FindOne(context.TODO(), filterActivityInfo).Decode(&activity)
	//log.Println("[info] activity %s", activity)
	if err2 != nil {
		log.Println(err2)
		return err2
	}
	for _, item := range activity.Activities.ActivityInfo {
		if objectID == item.ID {
			activityInfo = item
			break
		}
	}
	var distance = activityInfo.Distance
	var toTalDistance = activity.Activities.ToTalDistance - distance

	filterUpdate := bson.D{{"event_id", eventObjectID}, {"activities.user_id", userObjectID}}

	delete := bson.M{"$pull": bson.M{"activities.activity_info": bson.M{"_id": objectID}}}
	updated := bson.M{"$set": bson.M{"activities.total_distance": toTalDistance}}

	res, err := activityMongo.ConnectionDB.Collection(activityV2Collection).UpdateOne(context.TODO(), filterUpdate, delete)
	if err != nil {
		//log.Fatal(res)
		log.Printf("[info] err %s", res)
		return err
	}
	log.Printf("[info] res %s", res)
	res2, err2 := activityMongo.ConnectionDB.Collection(activityV2Collection).UpdateOne(context.TODO(), filterUpdate, updated)
	if err2 != nil {
		log.Fatal(res2)
		return err2
	}
	//activityInfo = activity.ActivityInfo

	return err
}

// UpdateWorkout repository for insert workouts
func (activityMongo ActivityV2RepositoryMongo) UpdateWorkout(workout model.WorkoutActivityInfo, userID primitive.ObjectID) error {
	filter := bson.D{primitive.E{Key: "user_id", Value: userID}, primitive.E{Key: "activity_info._id", Value: workout.ID}}
	update := bson.M{"$set": bson.M{"activity_info.$": workout}}
	_, err := activityMongo.ConnectionDB.Collection(workoutsCollection).UpdateOne(context.TODO(), filter, update)
	if err == nil {
		log.Println("update workout success")
	}
	return err
}

// AddKaoLogActivity log store send activity to Kao
func (activityMongo ActivityV2RepositoryMongo) AddKaoLogActivity(activity model.LogSendKaoActivity) error {
	_, err := activityMongo.ConnectionDB.Collection(activityKaoLog).InsertOne(context.TODO(), activity)
	if err == nil {
		log.Println("insert send activity kao success")
	}
	return err
}
