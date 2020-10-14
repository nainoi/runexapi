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
	AddWorkout(workout model.AddWorkout) (model.WorkoutActivityInfo, error)
	GetWorkouts(userID primitive.ObjectID) (bool, model.Workouts, error)
	UpdateWorkout(workout model.WorkoutActivityInfo, userID primitive.ObjectID) error
}

// WorkoutsRepositoryMongo struct mongo db
type WorkoutsRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

const (
	workoutsCollection = "workouts"
)

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
			log.Fatal(err)
			return dataInfo, err
		}
		dataInfo.ID = primitive.NewObjectID()
		update := bson.M{"$push": bson.M{"activity_info": dataInfo}}
		_, err = workoutsMongo.ConnectionDB.Collection(workoutsCollection).UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Fatal(err)
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

// GetWorkouts repository for get workouts data
func (workoutsMongo WorkoutsRepositoryMongo) GetWorkouts(userID primitive.ObjectID) (bool, model.Workouts, error) {
	var workout model.Workouts
	filter := bson.M{"user_id": userID}
	count, err := workoutsMongo.ConnectionDB.Collection(workoutsCollection).CountDocuments(context.TODO(), filter)
	//log.Printf("[info] count %s", count)
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
	err = workoutsMongo.ConnectionDB.Collection(workoutsCollection).FindOne(context.TODO(), filter).Decode(&workout)
	if err == nil {
		var total = 0.0
		for _, each := range workout.WorkoutActivityInfo {
			total += each.Distance
		}
		workout.TotalDistance = total
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
