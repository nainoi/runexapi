package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"thinkdev.app/think/runex/runexapi/config/db"
	"thinkdev.app/think/runex/runexapi/firebase"
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/repository/v2"
)

// WorkoutsRepository struct interface
type WorkoutsRepository interface {
	AddWorkout(workout model.AddWorkout) (model.WorkoutActivityInfo, error)
	WorkoutHook(workout model.HookWorkout) error
	AddMultiWorkout(userID string, workouts []model.WorkoutActivityInfo) error
	GetWorkouts(userID primitive.ObjectID) (bool, model.Workouts, error)
	UpdateWorkout(workout model.WorkoutActivityInfo, userID primitive.ObjectID) error
	HistoryMonth(userID primitive.ObjectID, year int) (bool, []model.WorkoutHistoryMonthInfo, error)
	HistoryAll(userID primitive.ObjectID) (bool, []model.WorkoutHistoryAllInfo, error)
	WorkoutInfo(userID primitive.ObjectID, workoutID primitive.ObjectID) (bool, model.WorkoutActivityInfo, error)
}

// WorkoutsRepositoryMongo struct mongo db
type WorkoutsRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

const (
	workoutsCollection = "workouts"
)

//HistoryState struct
type HistoryState struct {
	ID                  primitive.ObjectID          `json:"id" bson:"_id"`
	WorkoutActivityInfo []model.WorkoutActivityInfo `json:"activity_info" bson:"activity_info"`
	TotalDistance       float64                     `json:"totalDistance" bson:"totalDistance"`
	TotalCalory         float64                     `json:"totalCalory" bson:"totalCalory"`
	TotalDuration       int64                       `json:"totalDuration" bson:"totalDuration"`
}

//HistoryAllState struct
type HistoryAllState struct {
	ID                  HistoryDate                 `json:"id" bson:"_id"`
	WorkoutActivityInfo []model.WorkoutActivityInfo `json:"activity_info" bson:"activity_info"`
	TotalDistance       float64                     `json:"totalDistance" bson:"totalDistance"`
	TotalCalory         float64                     `json:"totalCalory" bson:"totalCalory"`
	TotalDuration       int64                       `json:"totalDuration" bson:"totalDuration"`
}

//WorkoutInfoState struct
type WorkoutInfoState struct {
	ID                  primitive.ObjectID        `json:"id" bson:"_id"`
	WorkoutActivityInfo model.WorkoutActivityInfo `json:"activity_info" bson:"activity_info"`
}

//HistoryDate struct
type HistoryDate struct {
	Year  int `form:"year" json:"year" bson:"year"`
	Month int `form:"month" json:"month" bson:"month"`
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
		dataInfo.LocURL = SaveLocation(dataInfo, dataInfo.ID.Hex()).URL
		dataInfo.Locations = []model.Location{}
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
		dataInfo.LocURL = SaveLocation(dataInfo, dataInfo.ID.Hex()).URL
		dataInfo.Locations = []model.Location{}
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

// WorkoutHook repository for insert workouts
func (workoutsMongo WorkoutsRepositoryMongo) WorkoutHook(workout model.HookWorkout) error {
	var user model.User
	filterUser := bson.M{"provider_id": workout.ProviderID, "provider": workout.Provider}
	err := workoutsMongo.ConnectionDB.Collection("user").FindOne(context.TODO(), filterUser).Decode(&user)
	if err != nil {
		log.Println(err)
		return err
	}
	filter := bson.M{"user_id": user.UserID}
	count, err := workoutsMongo.ConnectionDB.Collection(workoutsCollection).CountDocuments(context.TODO(), filter)
	//log.Printf("[info] count %s", count)
	if err != nil {
		log.Println(err)
		return err
	}
	total := 0.0
	for index, s := range workout.WorkoutActivityInfo {
		total += s.Distance
		workout.WorkoutActivityInfo[index].LocURL = SaveLocation(workout.WorkoutActivityInfo[index], workout.WorkoutActivityInfo[index].ID.Hex()).URL
		workout.WorkoutActivityInfo[index].Locations = []model.Location{}
	}
	if count > 0 {
		var workoutsModel model.Workouts
		err := workoutsMongo.ConnectionDB.Collection(workoutsCollection).FindOne(context.TODO(), filter).Decode(&workoutsModel)

		var totalDistance = workoutsModel.TotalDistance + total
		updateDistance := bson.M{"$set": bson.M{"total_distance": totalDistance}}
		_, err = workoutsMongo.ConnectionDB.Collection(workoutsCollection).UpdateOne(context.TODO(), filter, updateDistance)
		if err != nil {
			log.Printf("[info] err %s", err)
			return err
		}
		for _, s := range workout.WorkoutActivityInfo {
			updateWorkout := bson.M{"$push": bson.M{"activity_info": s}}
			_, err = workoutsMongo.ConnectionDB.Collection(workoutsCollection).UpdateOne(context.TODO(), filter, updateWorkout)
		}

		log.Printf("[info] err %s", err)
		return err
		// update := bson.M{"$push": bson.M{"activity_info": workouts}}
		// _, err = workoutsMongo.ConnectionDB.Collection(workoutsCollection).UpdateOne(context.TODO(), filter, update)
		// if err != nil {
		// 	log.Printf("[info] err %s", err)
		// 	return err
		// }

	} else {

		workoutsModel := model.Workouts{
			UserID:              user.UserID,
			WorkoutActivityInfo: workout.WorkoutActivityInfo,
			TotalDistance:       total,
		}
		//log.Println(workoutsModel)
		_, err := workoutsMongo.ConnectionDB.Collection(workoutsCollection).InsertOne(context.TODO(), workoutsModel)
		if err != nil {
			//log.Fatal(res)
			return err
		}
	}

	token, err := repository.GetFirebaseToken(user.UserID)
	if err == nil {
		if len(token.FirebaseTokens) > 0 {
			fcm := firebase.InitializeServiceAccountID()
			client, _ := fcm.Messaging(context.Background())
			if len(workout.WorkoutActivityInfo) > 0 {
				body := fmt.Sprintf("new %d activity from %s", len(workout.WorkoutActivityInfo), workout.WorkoutActivityInfo[0].APP)
				go firebase.SendMulticastAndHandleErrors(context.Background(), client, token.FirebaseTokens, "RUNEX Activity", body)
			}
			// if err == nil {

			// }
		}
	}

	return nil
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
	for index, s := range workouts {
		workouts[index].LocURL = SaveLocation(workouts[index], s.ID.Hex()).URL
		workouts[index].Locations = []model.Location{}
		total += s.Distance
	}
	if count > 0 {
		var workoutsModel model.Workouts
		err := workoutsMongo.ConnectionDB.Collection(workoutsCollection).FindOne(context.TODO(), filter).Decode(&workoutsModel)

		var totalDistance = workoutsModel.TotalDistance + total

		updateDistance := bson.M{"$set": bson.M{"total_distance": totalDistance}}
		_, err = workoutsMongo.ConnectionDB.Collection(workoutsCollection).UpdateOne(context.TODO(), filter, updateDistance)
		if err != nil {
			log.Printf("[info] err %s", err)
			return err
		}
		for _, s := range workouts {
			updateWorkout := bson.M{"$push": bson.M{"activity_info": s}}
			_, err = workoutsMongo.ConnectionDB.Collection(workoutsCollection).UpdateOne(context.TODO(), filter, updateWorkout)
		}

		log.Printf("[info] err %s", err)
		return err
		// for _, s := range workouts {

		// 	workoutsModel.WorkoutActivityInfo = append(workoutsModel.WorkoutActivityInfo, s)
		// }

		// updateDistance := bson.M{"$set": bson.M{"total_distance": totalDistance, "activity_info": workoutsModel.WorkoutActivityInfo}}
		// _, err = workoutsMongo.ConnectionDB.Collection(workoutsCollection).UpdateOne(context.TODO(), filter, updateDistance)
		// if err != nil {
		// 	log.Printf("[info] err %s", err)
		// 	return err
		// }

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

//HistoryMonth filter
func (workoutsMongo WorkoutsRepositoryMongo) HistoryMonth(userID primitive.ObjectID, year int) (bool, []model.WorkoutHistoryMonthInfo, error) {
	//var workout model.Workouts
	var workoutInfo []model.WorkoutActivityInfo
	var historyMonthInfo []model.WorkoutHistoryMonthInfo
	var historyDayInfo []model.WorkoutActivityInfo
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

	//log.Printf("[info] workoutInfo %d", workoutInfo)

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()

	countMonth := 0

	if year == currentYear {

		countMonth = int(currentMonth)
	} else {
		countMonth = 12
	}

	for m := 1; m <= countMonth; m++ {
		historyDayInfo = nil
		workoutInfo = nil
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
			//log.Println("[error] curHistory %d", err2.Error())
			return true, historyMonthInfo, err2
		}
		//log.Println("[info] curHistory %d", curHistory)
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

			// historyDayInfoNew := model.WorkoutDistanceDayInfo{
			// 	Distance:    item.Distance,
			// 	WorkoutDate: item.WorkoutDate.Format("2006-01-02"),
			// 	WorkoutTime: item.TimeString,
			// }
			historyDayInfo = append(historyDayInfo, item)
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

//HistoryAll all day
func (workoutsMongo WorkoutsRepositoryMongo) HistoryAll(userID primitive.ObjectID) (bool, []model.WorkoutHistoryAllInfo, error) {

	var workoutInfo []model.WorkoutActivityInfo
	var historyAllInfo []model.WorkoutHistoryAllInfo
	var historyDayInfo []model.WorkoutActivityInfo
	filter := bson.M{"user_id": userID}
	count, err := workoutsMongo.ConnectionDB.Collection(workoutsCollection).CountDocuments(context.TODO(), filter)
	log.Printf("[info] count %d", count)
	if err != nil {
		log.Println(err)
		return true, historyAllInfo, err
	}
	if count == 0 {

		return true, historyAllInfo, nil
	}

	matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"user_id": userID}}}
	unwindStage := bson.D{primitive.E{Key: "$unwind", Value: bson.M{"path": "$activity_info"}}}
	projectStage := bson.D{primitive.E{Key: "$project", Value: bson.M{"year": bson.M{"$year": "$activity_info.workout_date"}, "month": bson.M{"$month": "$activity_info.workout_date"}, "activity_info": 1}}}
	groupStage := bson.D{bson.E{Key: "$group", Value: bson.M{"_id": bson.M{"year": bson.M{"$year": "$activity_info.workout_date"}, "month": bson.M{"$month": "$activity_info.workout_date"}}, "activity_info": bson.M{"$push": "$activity_info"}, "totalDistance": bson.M{"$sum": "$activity_info.distance"}, "totalCalory": bson.M{"$sum": "$activity_info.calory"}, "totalDuration": bson.M{"$sum": "$activity_info.duration"}}}}
	sortStage := bson.D{primitive.E{Key: "$sort", Value: bson.M{"_id.year": -1, "_id.month": -1}}}
	curHistory, err2 := workoutsMongo.ConnectionDB.Collection(workoutsCollection).Aggregate(context.TODO(), mongo.Pipeline{matchStage, unwindStage, projectStage, groupStage, sortStage})
	if err2 != nil {
		log.Println("[error] curHistory %@", err2.Error())
		return true, historyAllInfo, err2
	}
	totalDistance := float64(0.00)
	calory := float64(0.00)
	duration := int64(0)
	year := int(0)
	month := int(0)
	for curHistory.Next(context.TODO()) {
		historyDayInfo = nil
		workoutInfo = nil
		var p HistoryAllState

		// decode the document
		if err := curHistory.Decode(&p); err != nil {
			log.Print(err)
		}
		//log.Println("[info] p %s", p)
		year = p.ID.Year
		month = p.ID.Month
		workoutInfo = p.WorkoutActivityInfo
		totalDistance = p.TotalDistance
		calory = p.TotalCalory
		duration = p.TotalDuration

		// log.Println("[info] totalDistance %s", p.TotalDistance)
		// log.Println("[info] calory %s", p.TotalCalory)
		// log.Println("[info] duration %s", p.TotalDuration)

		for _, item := range workoutInfo {

			// historyDayInfoNew := model.WorkoutDistanceDayInfo{
			// 	ID:          item.ID,
			// 	Distance:    item.Distance,
			// 	WorkoutDate: item.WorkoutDate.Format("2006-01-02 15:04:05"),
			// 	WorkoutTime: item.TimeString,
			// }
			historyDayInfo = append(historyDayInfo, item)
		}

		durTime := time.Duration(duration) * time.Second
		modTime := time.Now().Round(0).Add(-(durTime))
		since := time.Since(modTime)
		durStr := fmtDuration(since)

		historyAllInfoNew := model.WorkoutHistoryAllInfo{
			Year:          year,
			Month:         month,
			MonthName:     time.Month(month).String(),
			TotalDistance: totalDistance,
			Calory:        calory,
			TimeString:    durStr,
			HistoryDay:    historyDayInfo,
		}

		historyAllInfo = append(historyAllInfo, historyAllInfoNew)
	}

	return false, historyAllInfo, err
}

//WorkoutInfo workout info
func (workoutsMongo WorkoutsRepositoryMongo) WorkoutInfo(userID primitive.ObjectID, workoutID primitive.ObjectID) (bool, model.WorkoutActivityInfo, error) {

	var workoutInfo model.WorkoutActivityInfo
	filter := bson.M{"user_id": userID}
	count, err := workoutsMongo.ConnectionDB.Collection(workoutsCollection).CountDocuments(context.TODO(), filter)
	log.Printf("[info] count %d", count)
	if err != nil {
		log.Println(err)
		return true, workoutInfo, err
	}
	if count == 0 {

		return true, workoutInfo, nil
	}

	matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"user_id": userID}}}
	unwindStage := bson.D{primitive.E{Key: "$unwind", Value: bson.M{"path": "$activity_info"}}}
	match2Stage := bson.D{primitive.E{Key: "$match", Value: bson.M{"activity_info._id": workoutID}}}
	projectStage := bson.D{primitive.E{Key: "$project", Value: bson.M{"activity_info": 1}}}

	curHistory, err2 := workoutsMongo.ConnectionDB.Collection(workoutsCollection).Aggregate(context.TODO(), mongo.Pipeline{matchStage, unwindStage, match2Stage, projectStage})
	if err2 != nil {
		log.Println("[error] curHistory %@", err2.Error())
		return true, workoutInfo, err2
	}

	for curHistory.Next(context.TODO()) {
		var p WorkoutInfoState

		// decode the document
		if err := curHistory.Decode(&p); err != nil {
			log.Print(err)
		}

		workoutInfo = p.WorkoutActivityInfo

	}

	return false, workoutInfo, err
}

func WorkoutLocation() {

}

func DeleteWorkout(userID string, req model.RemoveWorkoutReq) error {
	var workout model.Workouts

	userObjectID, _ := primitive.ObjectIDFromHex(userID)
	//filter := bson.D{{"event_id", eventObjectID}, {"activities.user_id", userObjectID}}

	filterWorkoutInfo := bson.D{
		primitive.E{Key: "user_id", Value: userObjectID},
	}
	err := db.DB.Collection(workoutsCollection).FindOne(context.TODO(), filterWorkoutInfo).Decode(&workout)
	if err != nil {
		return err
	}

	var toTalDistance = workout.TotalDistance - req.WorkoutActivityInfo.Distance
	updated := bson.M{"$set": bson.M{"total_distance": toTalDistance}}

	_, err = db.DB.Collection(workoutsCollection).UpdateOne(context.TODO(), filterWorkoutInfo, updated)
	if err != nil {
		return err
	}

	delete := bson.M{"$pull": bson.M{"activity_info": bson.M{"_id": req.WorkoutActivityInfo.ID}}}

	_, err = db.DB.Collection(workoutsCollection).UpdateOne(context.TODO(), filterWorkoutInfo, delete)
	if err != nil {
		//log.Fatal(res)
		return err
	}

	//activityInfo = activity.ActivityInfo

	return err
}

func SaveLocation(activity model.WorkoutActivityInfo, filename string) model.UploadResponse {

	//pathDir := "." + config.UPLOAD_IMAGE
	// if _, err := os.Stat(pathDir); os.IsNotExist(err) {
	// 	os.MkdirAll(pathDir, os.ModePerm)
	// }
	location, _ := json.Marshal(activity)
	err := ioutil.WriteFile(filename+".json", location, 0644)
	if err != nil {
		log.Println(err.Error())
	}

	url := "https://storage.runex.co/upload-activities/"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	var resObject = model.UploadResponse{}

	writer.CreateFormField("file")
	part1, errFile1 := writer.CreateFormFile("file", fmt.Sprintf("%s", filename+".json"))
	// _, errFile1 = io.Copy(part1, file)
	part1.Write(location)

	// file, errFile1 := os.Open("/C:/Users/frogconn/Downloads/5f7324f2da1d9600135ed041vid1601548947586.mp4")
	//defer file.Close()
	//part1, errFile1 := writer.CreateFormField("file")

	if errFile1 != nil {
		fmt.Println(errFile1)
		log.Println("error create field")
	}
	err = writer.Close()
	if err != nil {
		log.Println("error close")
		fmt.Println(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println("error req")
		fmt.Println(err)
	}
	req.Header.Add("token", "5Dk2o03a4hVjQPglSueFEah577fCGQfM")
	// req.Header.Add("path", fmt.Sprintf("runex/workouts/"))
	req.Header.Add("Cookie", "__cfduid=dd42cd8b41a9c49d5b75f756dc64e01451604633984")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		log.Println("error do req")
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("error read upload")
		fmt.Println(err)
	}

	if res.StatusCode == 200 {
		err = json.Unmarshal(body, &resObject)
		if err != nil {
			log.Println(err)
		}
		log.Println(resObject)
	}
	return resObject
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
