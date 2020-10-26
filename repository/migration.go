package repository

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"thinkdev.app/think/runex/runexapi/model"
)

type MigrationRepository interface {
	MigrateWorkout(newCollection string) error
}

type MigrationRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

func (migrationMongo MigrationRepositoryMongo) MigrateWorkout(newCollection string) error {
	var history []model.RunHistory

	options := options.Find()
	//options.SetLimit(2)
	cur, err := migrationMongo.ConnectionDB.Collection("run_history").Find(context.TODO(), bson.D{{}}, options)

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

	for _, item := range history {
		var arrActivityInfo []model.WorkoutActivityInfo
		var locationArr []model.Location
		// workoutsModel := model.Workouts{
		// 	UserID:              item.UserID,
		// 	WorkoutActivityInfo: arrActivityInfo,
		// }

		for _, item2 := range item.RunHistoryInfo {

			//duration_time := int(item2.Time)
			integer, fraction := math.Modf(float64(item2.Time))
			var modTimeMin int = int(integer * 60)
			var modTimeSec int = int(fraction * 100)
			var duration = modTimeMin + modTimeSec
			durTime := time.Duration(modTimeMin+modTimeSec) * time.Second
			modTime := time.Now().Round(0).Add(-(durTime))
			since := time.Since(modTime)
			durStr := fmtDuration(since)

			dataInfo := model.WorkoutActivityInfo{
				ActivityType:     item2.ActivityType,
				Calory:           float64(item2.Calory),
				Caption:          item2.Caption,
				Distance:         item2.Distance,
				Pace:             float64(item2.Pace),
				Duration:         int64(duration),
				TimeString:       durStr,
				StartDate:        item2.ActivityDate,
				EndDate:          item2.ActivityDate,
				WorkoutDate:      item2.ActivityDate,
				NetElevationGain: 0.0,
				IsSync:           false,
				Locations:        locationArr,
			}
			dataInfo.ID = primitive.NewObjectID()
			arrActivityInfo = append(arrActivityInfo, dataInfo)
		}
		workoutsModel := model.Workouts{
			UserID:              item.UserID,
			WorkoutActivityInfo: arrActivityInfo,
			TotalDistance:       item.ToTalDistance,
		}

		_, err := migrationMongo.ConnectionDB.Collection(newCollection).InsertOne(context.TODO(), workoutsModel)
		if err != nil {
			//log.Fatal(res)
			return err
		}

	}
	var workout []model.Workouts
	options.SetLimit(0)
	cur, err = migrationMongo.ConnectionDB.Collection("workouts").Find(context.TODO(), bson.D{{}}, options)
	for cur.Next(context.TODO()) {
		var u model.Workouts
		// decode the document
		if err := cur.Decode(&u); err != nil {
			log.Println(err)
			log.Fatal(err)
		}
		//fmt.Printf("post: %+v\n", p)
		workout = append(workout, u)
	}
	for _, item3 := range workout {
		filter := bson.D{{"user_id", item3.UserID}}
		count, err2 := migrationMongo.ConnectionDB.Collection(newCollection).CountDocuments(context.TODO(), filter)
		if err2 != nil {
			log.Println(err2)
			return err2
		}
		if count > 0 {
			var workoutsModel model.Workouts
			err := migrationMongo.ConnectionDB.Collection(newCollection).FindOne(context.TODO(), filter).Decode(&workoutsModel)
			var totalDistance = workoutsModel.TotalDistance
			for _, item4 := range item3.WorkoutActivityInfo {
				totalDistance = totalDistance + item4.Distance
				update := bson.M{"$push": bson.M{"activity_info": item4}}
				_, err = migrationMongo.ConnectionDB.Collection(newCollection).UpdateOne(context.TODO(), filter, update)
			}
			updateDistance := bson.M{"$set": bson.M{"total_distance": totalDistance}}
			_, err = migrationMongo.ConnectionDB.Collection(newCollection).UpdateOne(context.TODO(), filter, updateDistance)
			if err != nil {
				log.Printf("[info] err %s", err)
				log.Fatal(err)
				return err
			}
		} else {
			var workoutsModel model.Workouts
			workoutsModel = item3
			_, err := migrationMongo.ConnectionDB.Collection(newCollection).InsertOne(context.TODO(), workoutsModel)
			if err != nil {
				//log.Fatal(res)
				return err
			}
		}
	}

	return nil
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
