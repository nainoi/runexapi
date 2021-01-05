package repository

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"thinkdev.app/think/runex/runexapi/model"
)

// WorkoutsRepository struct interface
type WorkoutsRepository interface {
	AddWorkout(workout model.AddWorkout) (model.WorkoutActivityInfo, error)
	AddMultiWorkout(userID string, workouts []model.WorkoutActivityInfo) error
	GetWorkouts(userID primitive.ObjectID) (bool, model.Workouts, error)
	UpdateWorkout(workout model.WorkoutActivityInfo, userID primitive.ObjectID) error
	HistoryMonth(userID primitive.ObjectID, year int) (bool, []model.WorkoutHistoryMonthInfo, error)
}

// WorkoutsRepositoryMongo struct mongo db
type WorkoutsRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

const (
	workoutsCollection = "workouts"
)

type HistoryState struct {
	ID                  primitive.ObjectID          `json:"id" bson:"_id"`
	WorkoutActivityInfo []model.WorkoutActivityInfo `json:"activity_info" bson:"activity_info"`
	TotalDistance       float64                     `json:"totalDistance" bson:"totalDistance"`
	TotalCalory         float64                     `json:"totalCalory" bson:"totalCalory"`
	TotalDuration       int64                       `json:"totalDuration" bson:"totalDuration"`
}

// AddWorkout repository for insert workouts
func (workoutsMongo WorkoutsRepositoryMongo) AddWorkout(workout model.AddWorkout) (model.WorkoutActivityInfo, error) {

	filter := bson.M{"user_id": workout.UserID}
	dataInfo := workout.WorkoutActivityInfo
	count, err := workoutsMongo.ConnectionDB.Collection(workoutsCollection).CountDocuments(context.TODO(), filter)
	//log.Printf("[info] count %s", count)
	if err != nil {
		log.Println(err)
		return dataInfo, err
	}
	if count > 0 {
		var workoutsModel model.Workouts
		err := workoutsMongo.ConnectionDB.Collection(workoutsCollection).FindOne(context.TODO(), filter).Decode(&workoutsModel)
		var totalDistance = workoutsModel.TotalDistance + workout.WorkoutActivityInfo.Distance
		updateDistance := bson.M{"$set": bson.M{"total_distance": totalDistance}}
		_, err = workoutsMongo.ConnectionDB.Collection(workoutsCollection).UpdateOne(context.TODO(), filter, updateDistance)
		if err != nil {
			log.Printf("[info] err %s", err)
			return dataInfo, err
		}
		dataInfo.ID = primitive.NewObjectID()
		update := bson.M{"$push": bson.M{"activity_info": dataInfo}}
		_, err = workoutsMongo.ConnectionDB.Collection(workoutsCollection).UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Printf("[info] err %s", err)
			return dataInfo, err
		}

	} else {
		var arrActivityInfo []model.WorkoutActivityInfo

		dataInfo := workout.WorkoutActivityInfo
		dataInfo.ID = primitive.NewObjectID()
		arrActivityInfo = append(arrActivityInfo, dataInfo)
		workoutsModel := model.Workouts{
			UserID:              workout.UserID,
			WorkoutActivityInfo: arrActivityInfo,
			TotalDistance:       dataInfo.Distance,
		}
		//log.Println(workoutsModel)
		_, err := workoutsMongo.ConnectionDB.Collection(workoutsCollection).InsertOne(context.TODO(), workoutsModel)
		if err != nil {
			//log.Fatal(res)
			return dataInfo, err
		}
	}

	return dataInfo, nil
}

// AddMultiWorkout repository for insert workouts
func (workoutsMongo WorkoutsRepositoryMongo) AddMultiWorkout(userID string, workouts []model.WorkoutActivityInfo) error {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	filter := bson.M{"user_id": userObjectID}
	count, err := workoutsMongo.ConnectionDB.Collection(workoutsCollection).CountDocuments(context.TODO(), filter)
	//log.Printf("[info] count %s", count)
	if err != nil {
		log.Println(err)
		return err
	}
	total := 0.0
	for _, s := range workouts {
		total += s.Distance
	}
	if count > 0 {
		var workoutsModel model.Workouts
		err := workoutsMongo.ConnectionDB.Collection(workoutsCollection).FindOne(context.TODO(), filter).Decode(&workoutsModel)

		var totalDistance = workoutsModel.TotalDistance + total
		for _, s := range workouts {
			workoutsModel.WorkoutActivityInfo = append(workoutsModel.WorkoutActivityInfo, s)
		}

		updateDistance := bson.M{"$set": bson.M{"total_distance": totalDistance, "activity_info": workoutsModel.WorkoutActivityInfo}}
		_, err = workoutsMongo.ConnectionDB.Collection(workoutsCollection).UpdateOne(context.TODO(), filter, updateDistance)
		if err != nil {
			log.Printf("[info] err %s", err)
			return err
		}

		// update := bson.M{"$push": bson.M{"activity_info": workouts}}
		// _, err = workoutsMongo.ConnectionDB.Collection(workoutsCollection).UpdateOne(context.TODO(), filter, update)
		// if err != nil {
		// 	log.Printf("[info] err %s", err)
		// 	return err
		// }

	} else {
		workoutsModel := model.Workouts{
			UserID:              userObjectID,
			WorkoutActivityInfo: workouts,
			TotalDistance:       total,
		}
		//log.Println(workoutsModel)
		_, err := workoutsMongo.ConnectionDB.Collection(workoutsCollection).InsertOne(context.TODO(), workoutsModel)
		if err != nil {
			//log.Fatal(res)
			return err
		}
	}

	return nil
}

// GetWorkouts repository for get workouts data
func (workoutsMongo WorkoutsRepositoryMongo) GetWorkouts(userID primitive.ObjectID) (bool, model.Workouts, error) {
	var workout model.Workouts
	filter := bson.M{"user_id": userID}
	option := options.FindOne()
	//option.SetSort(bson.D{primitive.E{Key: "activity_info.workout_date", Value: -1}})
	//option.SetSort(bson.D{primitive.E{Key: "activity_info.workout_date", Value: -1}})
	count, err := workoutsMongo.ConnectionDB.Collection(workoutsCollection).CountDocuments(context.TODO(), filter)
	log.Printf("[info] count %d", count)
	if err != nil {
		log.Println(err)
		return true, workout, err
	}
	if count == 0 {
		workout.UserID = userID
		workout.TotalDistance = 0
		workout.WorkoutActivityInfo = []model.WorkoutActivityInfo{}
		return false, workout, nil
	}

	//limit := bson.D{primitive.E{Key: "$limit", Value: 1}}
	err = workoutsMongo.ConnectionDB.Collection(workoutsCollection).FindOne(context.TODO(), filter, option).Decode(&workout)
	if err != nil {
		log.Println(err.Error())
		workout.UserID = userID
		workout.TotalDistance = 0
		workout.WorkoutActivityInfo = []model.WorkoutActivityInfo{}
		return false, workout, nil
	}
	return false, workout, err
}

// UpdateWorkout repository for insert workouts
func (workoutsMongo WorkoutsRepositoryMongo) UpdateWorkout(workout model.WorkoutActivityInfo, userID primitive.ObjectID) error {
	filter := bson.M{"user_id": userID, "activty_info._id": workout.ID}
	update := bson.M{"$set": bson.M{"activity_info.is_syn": true}}
	_, err := workoutsMongo.ConnectionDB.Collection(workoutsCollection).UpdateOne(context.TODO(), filter, update)
	return err
}

func (workoutsMongo WorkoutsRepositoryMongo) HistoryMonth(userID primitive.ObjectID, year int) (bool, []model.WorkoutHistoryMonthInfo, error) {
	//var workout model.Workouts
	var workoutInfo []model.WorkoutActivityInfo
	var historyMonthInfo []model.WorkoutHistoryMonthInfo
	var historyDayInfo []model.WorkoutDistanceDayInfo
	filter := bson.M{"user_id": userID}
	//option := options.FindOne()
	//option.SetSort(bson.D{primitive.E{Key: "activity_info.workout_date", Value: -1}})
	//option.SetSort(bson.D{primitive.E{Key: "activity_info.workout_date", Value: -1}})
	count, err := workoutsMongo.ConnectionDB.Collection(workoutsCollection).CountDocuments(context.TODO(), filter)
	log.Printf("[info] count %d", count)
	if err != nil {
		log.Println(err)
		return true, historyMonthInfo, err
	}
	if count == 0 {

		return true, historyMonthInfo, nil
	}

	//limit := bson.D{primitive.E{Key: "$limit", Value: 1}}
	// err = workoutsMongo.ConnectionDB.Collection(workoutsCollection).FindOne(context.TODO(), filter, option).Decode(&workout)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return false, nil, err
	// }
	// workoutInfo = workout.WorkoutActivityInfo

	log.Printf("[info] workoutInfo %d", workoutInfo)

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()

	count_month := 0

	if year == currentYear {

		count_month = int(currentMonth)
	} else {
		count_month = 12
	}

	for m := 1; m <= count_month; m++ {

		// firstday := time.Date(year, time.Month(m), 1, 0, 0, 0, 0, time.Local)
		// lastday := firstday.AddDate(0, 1, 0).Add(time.Nanosecond * -1)
		log.Printf("[info] year %d", year)
		log.Printf("[info] month %d", m)
		matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"user_id": userID}}}
		unwindStage := bson.D{primitive.E{Key: "$unwind", Value: bson.M{"path": "$activity_info"}}}
		projectStage := bson.D{primitive.E{Key: "$project", Value: bson.M{"year": bson.M{"$year": "$activity_info.workout_date"}, "month": bson.M{"$month": "$activity_info.workout_date"}, "activity_info": 1}}}
		matchStage2 := bson.D{primitive.E{Key: "$match", Value: bson.M{"year": year, "month": m}}}
		groupStage := bson.D{bson.E{Key: "$group", Value: bson.M{"_id": "$_id", "activity_info": bson.M{"$push": "$activity_info"}, "totalDistance": bson.M{"$sum": "$activity_info.distance"}, "totalCalory": bson.M{"$sum": "$activity_info.calory"}, "totalDuration": bson.M{"$sum": "$activity_info.duration"}}}}

		curHistory, err2 := workoutsMongo.ConnectionDB.Collection(workoutsCollection).Aggregate(context.TODO(), mongo.Pipeline{matchStage, unwindStage, projectStage, matchStage2, groupStage})
		if err2 != nil {
			log.Println("[error] curHistory %d", err2.Error())
			return true, historyMonthInfo, err2
		}
		log.Println("[info] curHistory %d", curHistory)
		totalDistance := float64(0.00)
		calory := float64(0.00)
		duration := int64(0)
		for curHistory.Next(context.TODO()) {
			var p HistoryState

			// decode the document
			if err := curHistory.Decode(&p); err != nil {
				log.Print(err)
			}
			//log.Println("[info] p %s", p)
			workoutInfo = p.WorkoutActivityInfo
			totalDistance = p.TotalDistance
			calory = p.TotalCalory
			duration = p.TotalDuration

			// log.Println("[info] totalDistance %s", p.TotalDistance)
			// log.Println("[info] calory %s", p.TotalCalory)
			// log.Println("[info] duration %s", p.TotalDuration)
		}

		for _, item := range workoutInfo {

			historyDayInfoNew := model.WorkoutDistanceDayInfo{
				Distance:    item.Distance,
				WorkoutDate: item.WorkoutDate.Format("2006-01-02"),
				WorkoutTime: item.TimeString,
			}
			historyDayInfo = append(historyDayInfo, historyDayInfoNew)
		}

		durTime := time.Duration(duration) * time.Second
		modTime := time.Now().Round(0).Add(-(durTime))
		since := time.Since(modTime)
		durStr := fmtDuration(since)

		historyMonthInfoNew := model.WorkoutHistoryMonthInfo{
			Month:         m,
			MonthName:     time.Month(m).String(),
			TotalDistance: totalDistance,
			Calory:        calory,
			TimeString:    durStr,
			HistoryDay:    historyDayInfo,
		}

		historyMonthInfo = append(historyMonthInfo, historyMonthInfoNew)
	}
	return false, historyMonthInfo, err
}

// func fmtDuration(d time.Duration) string {
// 	d = d.Round(time.Second)
// 	h := d / time.Hour
// 	d -= h * time.Hour
// 	m := d / time.Minute
// 	d -= m * time.Minute
// 	s := d / time.Second
// 	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
// }
