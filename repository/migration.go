package repository

import (
	"context"
	"fmt"
	"log"
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
	options.SetLimit(2)
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
			durTime := time.Duration(item2.Time) * time.Minute
			modTime := time.Now().Round(0).Add(-(durTime))
			since := time.Since(modTime)
			durStr := fmtDuration(since)

			dataInfo := model.WorkoutActivityInfo{
				ActivityType:     item2.ActivityType,
				Calory:           float64(item2.Calory),
				Caption:          item2.Caption,
				Distance:         item2.Distance,
				Pace:             float64(item2.Pace),
				Duration:         int64(item2.Time),
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
		}

		_, err := migrationMongo.ConnectionDB.Collection(newCollection).InsertOne(context.TODO(), workoutsModel)
		if err != nil {
			//log.Fatal(res)
			return err
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
