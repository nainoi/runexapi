package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//ConfigRepository interface user repository
type ConfigRepository interface {
	GetConfig() map[string]interface{}
}

// ConfigRepositoryMongo mongo ref
type ConfigRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

const (
	configCollection = "config"
)

// GetConfig repo
func (configMongo ConfigRepositoryMongo) GetConfig() map[string]interface{} {
	var configRes map[string]interface{}
	err := configMongo.ConnectionDB.Collection(configCollection).FindOne(context.TODO(), bson.D{}).Decode(&configRes)
	if err != nil {
		log.Println(err.Error())
	}
	return configRes
}
