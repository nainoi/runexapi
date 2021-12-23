package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"thinkdev.app/think/runex/runexapi/firebase"
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/repository/v2"
)

// StravaRepository interface
type StravaRepository interface {
	AddActivity(stravaReq model.StravaAddRequest) error
	GetActivities(userID string) ([]model.StravaActivity, error)
}

//RepoStravaDB db connection struct
type RepoStravaDB struct {
	ConnectionDB *mongo.Database
}

const (
	stravaConlection = "strava"
)

//AddActivity repository for add activity from strava
func (repo RepoStravaDB) AddActivity(stravaReq model.StravaAddRequest) error {
	var stravaData model.StravaData

	userRepo := repository.RepoUserDB{
		ConnectionDB: repo.ConnectionDB,
	}
	user, err := repository.GetUserWithProvider(stravaReq.Provider, stravaReq.ProviderID)
	if err != nil {
		return err
	}
	if user.StravaID == "" {
		user.StravaID = stravaReq.StravaID
		userRepo.UpdateUser(user, user.UserID.Hex())
	}

	filter := bson.D{
		primitive.E{Key: "user_id", Value: user.UserID},
		primitive.E{Key: "strava_id", Value: stravaReq.StravaID},
	}

	stravaReq.StravaActivity.CreatedAt = time.Now()
	stravaReq.StravaActivity.IsSync = false
	err = repo.ConnectionDB.Collection(stravaConlection).FindOne(context.TODO(), filter).Decode(&stravaData)
	if err != nil {
		stravaData.UserID = user.UserID
		stravaData.StravaID = stravaReq.StravaID
		stravaData.Activities = []model.StravaActivity{stravaReq.StravaActivity}

		_, err := repo.ConnectionDB.Collection(stravaConlection).InsertOne(context.TODO(), stravaData)
		return err
	}
	filter = bson.D{
		primitive.E{Key: "user_id", Value: stravaReq.ProviderID},
		primitive.E{Key: "strava_id", Value: stravaReq.StravaID},
		primitive.E{Key: "activities.$.id", Value: stravaReq.StravaActivity.ID},
	}
	res := repo.ConnectionDB.Collection(stravaConlection).FindOne(context.TODO(), filter)
	if res.Err() != nil {
		update := bson.M{"$push": bson.M{"activities": stravaReq.StravaActivity}}
		_, err = repo.ConnectionDB.Collection(stravaConlection).UpdateOne(context.TODO(), filter, update)
		if err == nil {
			token, err := repository.GetFirebaseToken(user.UserID)
			if err == nil {
				if len(token.FirebaseTokens) > 0 {
					fcm := firebase.InitializeServiceAccountID()
					client, err := fcm.Messaging(context.Background())
					if err == nil {
						body := fmt.Sprintf("run %f km.", stravaReq.StravaActivity.Distance)
						go firebase.SendMulticastAndHandleErrors(context.Background(), client, token.FirebaseTokens, "RUNEX sync Strava", body)
					}
				}
			}
		}
	}

	return err
}

//GetActivities repository for get user activities from strava
func (repo RepoStravaDB) GetActivities(userID string) ([]model.StravaActivity, error) {
	var data model.StravaData
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return data.Activities, err
	}
	filter := bson.D{primitive.E{Key: "user_id", Value: objectID}}
	err = repo.ConnectionDB.Collection(stravaConlection).FindOne(context.TODO(), filter).Decode(&data)
	return data.Activities, err
}
