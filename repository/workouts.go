package repository

import (
	"context"
	"log"

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
		
		updateDistance := bson.M{"$set": bson.M{"total_distance": totalDistance, "activity_info": workoutsModel.WorkoutActivityInfo }}
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
