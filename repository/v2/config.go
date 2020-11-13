package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"thinkdev.app/think/runex/runexapi/model"
)

//ConfigRepository interface user repository
type ConfigRepository interface {
	GetConfig() model.ConfigModel
}

// ConfigRepositoryMongo mongo ref
type ConfigRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

const (
	configCollection = "config"
)

// GetConfig repo
func (configMongo ConfigRepositoryMongo) GetConfig() model.ConfigModel {
	configRes := model.ConfigModel{}
	err := configMongo.ConnectionDB.Collection(configCollection).FindOne(context.TODO(), bson.D{}).Decode(&configRes)
	if err != nil {
		log.Println(err.Error())
	}
	return configRes
}