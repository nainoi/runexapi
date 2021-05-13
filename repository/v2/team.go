package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"thinkdev.app/think/runex/runexapi/config/db"
	"thinkdev.app/think/runex/runexapi/model"
)

const (
	teamCollection = "team"
)

func UpdateIcon(form model.TeamIcon) error {
	filter := bson.M{"reg_id": form.RegID}
	count, err := db.DB.Collection(teamCollection).CountDocuments(context.TODO(), filter)
	if count > 0 {
		_, err = db.DB.Collection(teamCollection).UpdateOne(context.TODO(), filter, form)
		return err
	}
	_, err = db.DB.Collection(teamCollection).InsertOne(context.TODO(), form)
	return err
}

func GetTeamIcon(form model.TeamIcon) (model.TeamIcon, error) {
	filter := bson.M{"reg_id": form.RegID}
	count, err := db.DB.Collection(teamCollection).CountDocuments(context.TODO(), filter)
	if count > 0 {
		err = db.DB.Collection(teamCollection).FindOne(context.TODO(), filter).Decode(&form)
		return form, err
	}
	return form, err
}