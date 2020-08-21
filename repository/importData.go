package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"thinkdev.app/think/runex/runexapi/config"
	"thinkdev.app/think/runex/runexapi/model"
)

type ImportDataRepository interface {
	ExistUserByEmail(email string) (string, bool, error)
	AddUser(user model.ExcelUserForm) (string, error)
	AddRegister(register model.Register) error
	ExistRegisterByUserAndEvent(userID string, eventID string) (bool, error)
}
type ImportDataRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

func (importDataMongo ImportDataRepositoryMongo) ExistUserByEmail(email string) (string, bool, error) {

	filter := bson.D{{"email", email}}
	var user model.UserProvider
	count, err := importDataMongo.ConnectionDB.Collection("user").CountDocuments(context.TODO(), filter)
	log.Printf("[info] count %s", count)
	if err != nil {
		log.Println(err)
	}
	if count > 0 {
		err := importDataMongo.ConnectionDB.Collection("user").FindOne(context.TODO(), filter).Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		userID := user.UserID.Hex()
		return userID, true, nil
	}

	return "", false, nil
}

func (importDataMongo ImportDataRepositoryMongo) AddUser(user model.ExcelUserForm) (string, error) {

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	if user.Role == "" {
		user.Role = config.MEMBER
	}

	res, err := importDataMongo.ConnectionDB.Collection("user").InsertOne(context.TODO(), user)
	fmt.Println("Inserted a single document: ", res.InsertedID)
	return res.InsertedID.(primitive.ObjectID).Hex(), err
}

func (importDataMongo ImportDataRepositoryMongo) AddRegister(register model.Register) error {
	register.CreatedAt = time.Now()
	register.UpdatedAt = time.Now()
	res, err := importDataMongo.ConnectionDB.Collection("register").InsertOne(context.TODO(), register)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", res.InsertedID)
	return nil
}

func (importDataMongo ImportDataRepositoryMongo) ExistRegisterByUserAndEvent(userID string, eventID string) (bool, error) {

	userObjectID, _ := primitive.ObjectIDFromHex(userID)
	eventObjectID, _ := primitive.ObjectIDFromHex(eventID)
	filter := bson.D{{"user_id", userObjectID}, {"event_id", eventObjectID}}
	count, err := importDataMongo.ConnectionDB.Collection("register").CountDocuments(context.TODO(), filter)
	log.Printf("[info] count %s", count)
	if err != nil {
		log.Println(err)
	}
	if count > 0 {

		return true, nil
	}

	return false, nil
}
