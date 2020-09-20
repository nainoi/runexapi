package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"thinkdev.app/think/runex/runexapi/model"
)

// StravaRepository interface
type StravaRepository interface {
	AddActivity(stravaReq model.StravaAddRequest) error
	GetActivities(userID string) ([]model.StravaActivity, error)
}

//RepoDB db connection struct
type RepoDB struct {
	ConnectionDB *mongo.Database
}

const (
	stravaConlection = "strava"
)

//AddActivity repository for add activity from strava
func (repo RepoDB) AddActivity(stravaReq model.StravaAddRequest) error {
	var stravaData model.StravaData
	filter := bson.D{primitive.E{Key: "user_id", Value: stravaReq.UserID}}
	stravaReq.StravaActivity.CreatedAt = time.Now()
	stravaReq.StravaActivity.IsSync = false
	err := repo.ConnectionDB.Collection(stravaConlection).FindOne(context.TODO(), filter).Decode(&stravaData)
	if err != nil {
		stravaData.UserID = stravaReq.UserID
		stravaData.Activities = []model.StravaActivity{stravaReq.StravaActivity}

		_, err := repo.ConnectionDB.Collection(stravaConlection).InsertOne(context.TODO(), stravaData)
		return err
	}
	update := bson.M{"$push": bson.M{"activities": stravaReq.StravaActivity}}
	_, err = repo.ConnectionDB.Collection(stravaConlection).UpdateOne(context.TODO(), filter, update)
	return err
}

//GetActivities repository for get user activities from strava
func (repo RepoDB) GetActivities(userID string) ([]model.StravaActivity, error) {
	var data model.StravaData
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return data.Activities, err
	}
	filter := bson.D{primitive.E{Key: "user_id", Value: objectID}}
	err = repo.ConnectionDB.Collection(stravaConlection).FindOne(context.TODO(), filter).Decode(&data)
	return data.Activities, err
}
