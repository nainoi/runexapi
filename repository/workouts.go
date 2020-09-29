package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"thinkdev.app/think/runex/runexapi/model"
)

// WorkoutsRepository struct interface
type WorkoutsRepository interface {
	AddWorkout(workout model.AddWorkout) error
}

// WorkoutsRepositoryMongo struct mongo db
type WorkoutsRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

const (
	workoutsCollection = "workouts"
)

// AddWorkout repository for insert workouts
func (workoutsMongo WorkoutsRepositoryMongo) AddWorkout(workout model.AddWorkout) error {

	filter := bson.M{"user_id": workout.UserID}
	count, err := workoutsMongo.ConnectionDB.Collection(workoutsCollection).CountDocuments(context.TODO(), filter)
	//log.Printf("[info] count %s", count)
	if err != nil {
		log.Println(err)
		return err
	}
	if count > 0 {

		dataInfo := workout.WorkoutActivityInfo
		dataInfo.ID = primitive.NewObjectID()
		update := bson.M{"$push": bson.M{"activity_info": dataInfo}}
		_, err := workoutsMongo.ConnectionDB.Collection(workoutsCollection).UpdateOne(context.TODO(), filter, update)
		if err != nil {
			//log.Fatal(res)
			//log.Printf("[info] err %s", res)
			return err
		}

	} else {
		var arrActivityInfo []model.WorkoutActivityInfo

		dataInfo := workout.WorkoutActivityInfo
		dataInfo.ID = primitive.NewObjectID()
		arrActivityInfo = append(arrActivityInfo, dataInfo)
		workoutsModel := model.Workouts{
			UserID:              workout.UserID,
			WorkoutActivityInfo: arrActivityInfo,
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
