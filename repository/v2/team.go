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
		update := bson.M{"$set": bson.M{"icon_url": form.IconURL}}
		res := db.DB.Collection(teamCollection).FindOneAndUpdate(context.TODO(), filter, update)
		return res.Err()
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