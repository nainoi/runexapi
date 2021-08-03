package repository

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"thinkdev.app/think/runex/runexapi/config"
	"thinkdev.app/think/runex/runexapi/config/db"
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/repository/v2"
	"thinkdev.app/think/runex/runexapi/utils"
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
	//GetActivityWaitApprove(eventCode string) ([]model.ActivityInfo, error)
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
	filter := bson.D{primitive.E{Key: "event_code", Value: activity.EventCode}, primitive.E{Key: "user_id", Value: activity.UserID}, primitive.E{Key: "reg_id", Value: activity.RegID}}
	count, err := activityMongo.ConnectionDB.Collection(activityV2Collection).CountDocuments(context.TODO(), filter)
	log.Printf("[info] count %d", count)
	log.Printf("[info] count %s", activity.EventCode)
	if err != nil {
		log.Println(err)
		return err
	}
	if count > 0 {
		var activityModel model.ActivityV2
		err := activityMongo.ConnectionDB.Collection(activityV2Collection).FindOne(context.TODO(), filter).Decode(&activityModel)
		if activity.ActivityInfo.APP == "" {
			dataInfo := activity.ActivityInfo
			dataInfo.ID = primitive.NewObjectID()
			dataInfo.IsApprove = false
			dataInfo.Status = config.ACTIVITY_STATUS_WAITING
			update := bson.M{"$push": bson.M{"activity_info": dataInfo}}
			_, err = activityMongo.ConnectionDB.Collection(activityV2Collection).UpdateOne(context.TODO(), filter, update)
			if err != nil {
				return err
			}
			activityLogInfo := model.LogActivityInfo{
				UserID:         activity.UserID,
				ActivityInfoID: dataInfo.ID,
				Distance:       utils.ToFixed(dataInfo.Distance, 2),
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
			var totalDistance = activityModel.ToTalDistance + activity.ActivityInfo.Distance

			updateDistance := bson.M{"$set": bson.M{"total_distance": totalDistance}}
			_, err2 := activityMongo.ConnectionDB.Collection(activityV2Collection).UpdateOne(context.TODO(), filter, updateDistance)
			if err2 != nil {
				log.Printf("[info] err %s", err2)
				return err2
			}
			dataInfo := activity.ActivityInfo
			dataInfo.ID = primitive.NewObjectID()
			dataInfo.IsApprove = true
			dataInfo.Status = config.ACTIVITY_STATUS_APPROVE
			update := bson.M{"$push": bson.M{"activity_info": dataInfo}}
			_, err = activityMongo.ConnectionDB.Collection(activityV2Collection).UpdateOne(context.TODO(), filter, update)
			if err != nil {
				//log.Fatal(res)
				//log.Printf("[info] err %s", res)
				return err
			}
			activityLogInfo := model.LogActivityInfo{
				UserID:         activity.UserID,
				ActivityInfoID: dataInfo.ID,
				Distance:       utils.ToFixed(dataInfo.Distance, 2),
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
	} else {
		var arrActivityInfo []model.ActivityInfo

		dataInfo := activity.ActivityInfo
		if activity.ActivityInfo.APP == "" {
			dataInfo.ID = primitive.NewObjectID()
			dataInfo.IsApprove = false
			dataInfo.Status = config.ACTIVITY_STATUS_WAITING
			arrActivityInfo = append(arrActivityInfo, dataInfo)

			activityModel := model.ActivityV2{
				UserID:        activity.UserID,
				EventCode:     activity.EventCode,
				ActivityInfo:  arrActivityInfo,
				Ticket:        activity.Ticket,
				OrderID:       activity.OrderID,
				ToTalDistance: 0,
				RegID:         activity.RegID,
				ParentRegID:   activity.ParentRegID,
			}

			u, err := repository.GetUserInfo(activityModel)
			if err == nil {
				activityModel.UserInfo = u
			}

			_, err = activityMongo.ConnectionDB.Collection(activityV2Collection).InsertOne(context.TODO(), activityModel)

			activityLogInfo := model.LogActivityInfo{
				UserID:         activity.UserID,
				EventCode:      activity.EventCode,
				ActivityInfoID: dataInfo.ID,
				Distance:       utils.ToFixed(dataInfo.Distance, 2),
				ImageURL:       dataInfo.ImageURL,
				Caption:        dataInfo.Caption,
				Time:           dataInfo.Time,
				APP:            dataInfo.APP,
				ActivityDate:   dataInfo.ActivityDate,
				IsApprove:      false,
				CreatedAt:      dataInfo.CreatedAt,
				UpdatedAt:      dataInfo.UpdatedAt,
			}
			_, _ = activityMongo.ConnectionDB.Collection(activityLog).InsertOne(context.TODO(), activityLogInfo)
			if err != nil {
				return err
			}

		} else {
			dataInfo.ID = primitive.NewObjectID()
			dataInfo.IsApprove = true
			dataInfo.Status = config.ACTIVITY_STATUS_APPROVE
			arrActivityInfo = append(arrActivityInfo, dataInfo)

			// activities := model.Activities{
			// 	UserID:        activity.UserID,
			// 	ToTalDistance: activity.ActivityInfo.Distance,
			// 	ActivityInfo:  arrActivityInfo,
			// }

			activityModel := model.ActivityV2{
				UserID:        activity.UserID,
				EventCode:     activity.EventCode,
				ActivityInfo:  arrActivityInfo,
				Ticket:        activity.Ticket,
				OrderID:       activity.OrderID,
				ToTalDistance: activity.ActivityInfo.Distance,
				RegID:         activity.RegID,
				ParentRegID:   activity.ParentRegID,
			}
			u, err := repository.GetUserInfo(activityModel)
			if err == nil {
				activityModel.UserInfo = u
			}
			_, err = activityMongo.ConnectionDB.Collection(activityV2Collection).InsertOne(context.TODO(), activityModel)
			if err != nil {
				return err
			}

			activityLogInfo := model.LogActivityInfo{
				UserID:         activity.UserID,
				EventCode:      activity.EventCode,
				ActivityInfoID: dataInfo.ID,
				Distance:       utils.ToFixed(dataInfo.Distance, 2),
				ImageURL:       dataInfo.ImageURL,
				Caption:        dataInfo.Caption,
				Time:           dataInfo.Time,
				APP:            dataInfo.APP,
				ActivityDate:   dataInfo.ActivityDate,
				IsApprove:      true,
				CreatedAt:      dataInfo.CreatedAt,
				UpdatedAt:      dataInfo.UpdatedAt,
			}
			_, _ = activityMongo.ConnectionDB.Collection(activityLog).InsertOne(context.TODO(), activityLogInfo)
		}
	}

	return nil
}

// GetActivityByEvent event and activity detail
func (activityMongo ActivityV2RepositoryMongo) GetActivityByEvent(eventCode string, userID string) ([]model.ActivityInfo, error) {
	var activity model.ActivityV2
	var activityInfo = []model.ActivityInfo{}
	userObjectID, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.D{primitive.E{Key: "event_code", Value: eventCode}, primitive.E{Key: "activities.user_id", Value: userObjectID}}
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
	// eventObjectID, _ := primitive.ObjectIDFromHex(eventID)
	filter := bson.D{primitive.E{Key: "event_code", Value: eventID}, primitive.E{Key: "user_id", Value: userObjectID}}
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

// GetActivityByEvent2 event and activity detail
func GetActivityEventDashboard(req model.EventActivityDashboardReq, userID string) ([]model.ActivityV2, error) {
	var activity = []model.ActivityV2{}
	// userObjectID, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.D{primitive.E{Key: "event_code", Value: req.EventCode}, primitive.E{Key: "reg_id", Value: req.RegID}}
	if !req.ParentRegID.IsZero() {
		filter = bson.D{primitive.E{Key: "event_code", Value: req.EventCode}, primitive.E{Key: "parent_reg_id", Value: req.ParentRegID}}
	}

	count, err := db.DB.Collection(activityV2Collection).CountDocuments(context.TODO(), filter)
	if count > 0 {
		curr, err := db.DB.Collection(activityV2Collection).Find(context.TODO(), filter)

		if err != nil {
			log.Println(err)
			return activity, err
		}

		for curr.Next(context.TODO()) {
			var u model.ActivityV2
			// decode the document
			if err := curr.Decode(&u); err != nil {
				log.Println(err)
			}
			//fmt.Printf("post: %+v\n", p)
			activity = append(activity, u)
		}
		//activityInfo = activity.ActivityInfo

		return activity, err
	}
	return activity, err
}

// GetActivity event and activity detail
func GetActivity(regID primitive.ObjectID) (model.ActivityV2, error) {
	var activity = model.ActivityV2{}
	// eventObjectID, _ := primitive.ObjectIDFromHex(eventID)
	filter := bson.D{primitive.E{Key: "reg_id", Value: regID}}
	count, err := db.DB.Collection(activityV2Collection).CountDocuments(context.TODO(), filter)
	if count > 0 {
		err = db.DB.Collection(activityV2Collection).FindOne(context.TODO(), filter).Decode(&activity)
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
	// eventObjectID, _ := primitive.ObjectIDFromHex(event_id)
	filter := bson.D{primitive.E{Key: "event_code", Value: event_id}, primitive.E{Key: "activities.user_id", Value: userObjectID}}

	err := activityMongo.ConnectionDB.Collection(activityV2Collection).FindOne(context.TODO(), filter).Decode(&activity)

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

func (activityMongo ActivityV2RepositoryMongo) HistoryMonthByEvent(event_id string, user_id string, year int) ([]model.HistoryMonthInfo, error) {

	var activity model.ActivityV2
	var activityInfo []model.ActivityInfo

	userObjectID, _ := primitive.ObjectIDFromHex(user_id)
	// eventObjectID, _ := primitive.ObjectIDFromHex(event_id)
	filter := bson.D{primitive.E{Key: "event_code", Value: event_id}, primitive.E{Key: "user_id", Value: userObjectID}}

	err := activityMongo.ConnectionDB.Collection(activityV2Collection).FindOne(context.TODO(), filter).Decode(&activity)

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

func (activityMongo ActivityV2RepositoryMongo) DeleteActivity(eventCode string, userID string, activityID string) error {
	objectID, _ := primitive.ObjectIDFromHex(activityID)

	var activity model.ActivityV2
	var activityInfo model.ActivityInfo

	userObjectID, _ := primitive.ObjectIDFromHex(userID)
	//filter := bson.D{{"event_id", eventObjectID}, {"activities.user_id", userObjectID}}

	filterActivityInfo := bson.D{primitive.E{Key: "event_code", Value: eventCode}, primitive.E{Key: "activities.user_id", Value: userObjectID}, primitive.E{Key: "activities.activity_info._id", Value: objectID}}
	err2 := activityMongo.ConnectionDB.Collection(activityV2Collection).FindOne(context.TODO(), filterActivityInfo).Decode(&activity)
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

	filterUpdate := bson.D{primitive.E{Key: "event_code", Value: eventCode}, primitive.E{Key: "activities.user_id", Value: userObjectID}}

	delete := bson.M{"$pull": bson.M{"activity_info": bson.M{"_id": objectID}}}
	updated := bson.M{"$set": bson.M{"total_distance": toTalDistance}}

	_, err := activityMongo.ConnectionDB.Collection(activityV2Collection).UpdateOne(context.TODO(), filterUpdate, delete)
	if err != nil {
		//log.Fatal(res)
		return err
	}

	_, err = activityMongo.ConnectionDB.Collection(activityV2Collection).UpdateOne(context.TODO(), filterUpdate, updated)
	if err != nil {
		return err
	}
	//activityInfo = activity.ActivityInfo

	return err
}

func RemoveActivity(req model.EventActivityRemoveReq, userID string) error {

	var activity model.ActivityV2

	userObjectID, _ := primitive.ObjectIDFromHex(userID)
	//filter := bson.D{{"event_id", eventObjectID}, {"activities.user_id", userObjectID}}

	filterActivityInfo := bson.D{
		primitive.E{Key: "event_code", Value: req.EventCode},
		primitive.E{Key: "reg_id", Value: req.RegID},
		primitive.E{Key: "user_id", Value: userObjectID},
	}
	err := db.DB.Collection(activityV2Collection).FindOne(context.TODO(), filterActivityInfo).Decode(&activity)

	if req.ActivityInfo.Status == config.ACTIVITY_STATUS_APPROVE {
		var toTalDistance = activity.ToTalDistance - req.ActivityInfo.Distance

		updated := bson.M{"$set": bson.M{"total_distance": toTalDistance}}

		_, err = db.DB.Collection(activityV2Collection).UpdateOne(context.TODO(), filterActivityInfo, updated)
		if err != nil {
			return err
		}
	}

	delete := bson.M{"$pull": bson.M{"activity_info": bson.M{"_id": req.ActivityInfo.ID}}}

	_, err = db.DB.Collection(activityV2Collection).UpdateOne(context.TODO(), filterActivityInfo, delete)
	if err != nil {
		//log.Fatal(res)
		return err
	}

	//activityInfo = activity.ActivityInfo

	return err
}

func RemoveActivityByAdmin(req model.EventActivityRemoveReq) error {

	var activity model.ActivityV2
	//filter := bson.D{{"event_id", eventObjectID}, {"activities.user_id", userObjectID}}

	filterActivityInfo := bson.D{
		primitive.E{Key: "event_code", Value: req.EventCode},
		primitive.E{Key: "reg_id", Value: req.RegID},
	}
	err := db.DB.Collection(activityV2Collection).FindOne(context.TODO(), filterActivityInfo).Decode(&activity)
	
	if req.ActivityInfo.Status == config.ACTIVITY_STATUS_APPROVE {
		var toTalDistance = activity.ToTalDistance - req.ActivityInfo.Distance

		updated := bson.M{"$set": bson.M{"total_distance": toTalDistance}}

		_, err = db.DB.Collection(activityV2Collection).UpdateOne(context.TODO(), filterActivityInfo, updated)
		if err != nil {
			return err
		}
	}

	delete := bson.M{"$pull": bson.M{"activity_info": bson.M{"_id": req.ActivityInfo.ID}}}

	_, err = db.DB.Collection(activityV2Collection).UpdateOne(context.TODO(), filterActivityInfo, delete)
	if err != nil {
		//log.Fatal(res)
		return err
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

func GetActivityWaitApprove(eventCode string) ([]model.ActivityV2, error) {
	var activityInfos = []model.ActivityV2{}

	matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"event_code": eventCode}}}
	projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"activity_info": bson.M{"$filter": bson.M{"input": "$activity_info", "as": "activity_info", "cond": bson.M{"$eq": bson.A{"$$activity_info.status", config.ACTIVITY_STATUS_WAITING}}}}, "event_code": 1, "order_id": 1, "reg_id": 1, "ticket": 1, "parent_reg_id": 1, "user_id": 1, "total_distance": 1, "id": 1, "user_info": 1}}}
	curr, err := db.DB.Collection(activityV2Collection).Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage})

	// filterActivityInfo := bson.D{primitive.E{Key: "event_code", Value: eventCode}, primitive.E{Key: "activity_info.status", Value: config.ACTIVITY_STATUS_WAITING}}
	// curr, err := db.DB.Collection(activityV2Collection).Find(context.TODO(), filterActivityInfo)
	// //log.Println("[info] activity %s", activity)
	if err != nil {
		return activityInfos, err
	}

	for curr.Next(context.TODO()) {
		var u model.ActivityV2
		// decode the document
		if err := curr.Decode(&u); err != nil {
			log.Println(err)
		}
		//fmt.Printf("post: %+v\n", p)
		activityInfos = append(activityInfos, u)
	}
	return activityInfos, nil
}

func GetActivityWithStatus(req model.ActivityWithStatusReq) ([]model.ActivityV2, error) {
	var activityInfos = []model.ActivityV2{}

	matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"event_code": req.EventCode}}}

	projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"activity_info": bson.M{"$filter": bson.M{"input": "$activity_info", "as": "activity_info", "cond": bson.M{"$and": []interface{}{ bson.M{ "$eq": bson.A{"$$activity_info.status", req.Status}}, bson.M{"$gt": bson.A{"$$activity_info.activity_date", req.StartDate}}, bson.M{"$lt": bson.A{"$$activity_info.activity_date", req.EndDate}}}}}}, "event_code": 1, "order_id": 1, "reg_id": 1, "ticket": 1, "parent_reg_id": 1, "user_id": 1, "total_distance": 1, "id": 1, "user_info": 1}}}
	if req.Status == "" {
		projectStage = bson.D{bson.E{Key: "$project", Value: bson.M{"activity_info": bson.M{"$filter": bson.M{"input": "$activity_info", "as": "activity_info", "cond": bson.M{"$and": []interface{}{bson.M{"$gt": bson.A{"$$activity_info.activity_date", req.StartDate}}, bson.M{"$lt": bson.A{"$$activity_info.activity_date", req.EndDate}}}}}}, "event_code": 1, "order_id": 1, "reg_id": 1, "ticket": 1, "parent_reg_id": 1, "user_id": 1, "total_distance": 1, "id": 1, "user_info": 1}}}
	}
	curr, err := db.DB.Collection(activityV2Collection).Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage})

	// filterActivityInfo := bson.D{primitive.E{Key: "event_code", Value: eventCode}, primitive.E{Key: "activity_info.status", Value: config.ACTIVITY_STATUS_WAITING}}
	// curr, err := db.DB.Collection(activityV2Collection).Find(context.TODO(), filterActivityInfo)
	// //log.Println("[info] activity %s", activity)
	if err != nil {
		return activityInfos, err
	}

	for curr.Next(context.TODO()) {
		var u model.ActivityV2
		// decode the document
		if err := curr.Decode(&u); err != nil {
			log.Println(err)
		}
		//fmt.Printf("post: %+v\n", p)
		activityInfos = append(activityInfos, u)
	}
	return activityInfos, nil
}

func UpdateActivity(req model.UpdateActivityReq) error {

	filter := bson.M{"$and": []interface{}{bson.M{"event_code": req.EventCode}, bson.M{"user_id": req.UserID}, bson.M{"reg_id": req.RegID}, bson.M{"activity_info._id": req.ActivityID}}}
	var act model.ActivityV2
	if req.Status == config.ACTIVITY_STATUS_APPROVE {
		err := db.DB.Collection(activityV2Collection).FindOne(context.TODO(), filter).Decode(&act)
		if err != nil {
			return err
		}
		distance := act.ToTalDistance + req.Distance
		update := bson.M{"$set": bson.M{"activity_info.$.status": req.Status, "activity_info.$.reason": req.Reason, "activity_info.$.updated_at": time.Now(), "activity_info.$.is_approve": true, "total_distance": utils.ToFixed(distance, 2), "activity_info.$.distance": utils.ToFixed(req.Distance, 2)}}
		result := db.DB.Collection(activityV2Collection).FindOneAndUpdate(context.TODO(), filter, update)
		return result.Err()
	}

	update := bson.M{"$set": bson.M{"activity_info.$.status": req.Status, "activity_info.$.updated_at": time.Now(), "activity_info.$.is_approve": true}}
	result := db.DB.Collection(activityV2Collection).FindOneAndUpdate(context.TODO(), filter, update)
	return result.Err()
}

func UpdateUserInfoActivity(userOption model.UserOption, eventCode string, userID primitive.ObjectID, regID primitive.ObjectID) error {

	filter := bson.M{"$and": []interface{}{bson.M{"event_code": eventCode}, bson.M{"user_id": userID}, bson.M{"reg_id": regID}}}
	count, err := db.DB.Collection(activityV2Collection).CountDocuments(context.TODO(), filter)
	if count > 0 {
		update := bson.M{"$set": bson.M{"user_info": userOption}}
		result := db.DB.Collection(activityV2Collection).FindOneAndUpdate(context.TODO(), filter, update)
		return result.Err()
	}
	return err
}
