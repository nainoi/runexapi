package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"thinkdev.app/think/runex/runexapi/model"
)

// KaoRepository struct interface
type KaoRepository interface {
	AddKao(kaoObject model.KaoObject, userID string) error
	// GetWorkouts(userID primitive.ObjectID) (bool, model.Workouts, error)
	// UpdateWorkout(workout model.WorkoutActivityInfo, userID primitive.ObjectID) error
}

// KaoRepositoryMongo struct mongo db
type KaoRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

const (
	kaoCollection = "kao"
)

// AddKao repository for insert workouts
func (kaoMongo KaoRepositoryMongo) AddKao(kaoObject model.KaoObject, userID string) error {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	filter := bson.M{"user_id": userObjectID}
	count, err := kaoMongo.ConnectionDB.Collection(kaoCollection).CountDocuments(context.TODO(), filter)
	//log.Printf("[info] count %s", count)
	if err != nil {
		log.Println(err)
		return err
	}
	if count > 0 {
		filter := bson.M{"user_id": userObjectID, "kao_info.$.id": kaoObject.ID}
		count, err = kaoMongo.ConnectionDB.Collection(workoutsCollection).CountDocuments(context.TODO(), filter)
		if count == 0 {
			update := bson.M{"$push": bson.M{"kao_info": kaoObject}}
			_, err = kaoMongo.ConnectionDB.Collection(workoutsCollection).UpdateOne(context.TODO(), filter, update)
			if err != nil {
				log.Printf("[info] err %s", err)
				return err
			}
		}
	} else {
		var arrKaoObject []model.KaoObject
		arrKaoObject = append(arrKaoObject, kaoObject)
		workoutsModel := model.Kao{
			UserID:     userObjectID,
			KaoObjects: arrKaoObject,
		}
		//log.Println(workoutsModel)
		_, err := kaoMongo.ConnectionDB.Collection(workoutsCollection).InsertOne(context.TODO(), workoutsModel)
		if err != nil {
			//log.Fatal(res)
			return err
		}
	}

	return nil
}
