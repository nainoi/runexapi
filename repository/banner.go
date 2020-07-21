package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"bitbucket.org/suthisakch/runex/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BannerRepository interface {
	AddBanner(banner model.BannerAddForm) (string, error)
	DeleteBannerByID(bannerID string) error
	GetBannerAll() ([]model.Banner, error)
	ExistByEventID(eventID string) (bool, error)
}

type BannerRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

const (
	bannerCollection = "banner"
)

func (bannerMongo BannerRepositoryMongo) AddBanner(banner model.BannerAddForm) (string, error) {

	var bannerModel model.Banner
	objectID, err := primitive.ObjectIDFromHex(banner.EventID)
	bannerModel.EventID = objectID
	bannerModel.CreatedAt = time.Now()
	bannerModel.Active = banner.Active
	res, err := bannerMongo.ConnectionDB.Collection(bannerCollection).InsertOne(context.TODO(), bannerModel)
	if err != nil {
		log.Fatal(res)
	}
	fmt.Println("Inserted a single document: ", res.InsertedID)
	return res.InsertedID.(primitive.ObjectID).Hex(), err
}

func (bannerMongo BannerRepositoryMongo) DeleteBannerByID(bannerID string) error {

	id, err := primitive.ObjectIDFromHex(bannerID)
	if err != nil {
		log.Fatal(err)
		return err

	}
	filter := bson.M{"_id": id}
	deleteResult, err2 := bannerMongo.ConnectionDB.Collection(bannerCollection).DeleteOne(context.TODO(), filter)
	log.Printf("[info] cur %s", err2)
	if err2 != nil {
		return err2
	}
	fmt.Printf("Deleted %v documents in the Banner collection\n", deleteResult.DeletedCount)
	return nil
}

func (bannerMongo BannerRepositoryMongo) GetBannerAll() ([]model.Banner, error) {

	var banner []model.Banner
	cur, err := bannerMongo.ConnectionDB.Collection(bannerCollection).Find(context.TODO(), bson.D{{}})
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {
		var u model.Banner
		// decode the document
		if err := cur.Decode(&u); err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("post: %+v\n", p)
		banner = append(banner, u)
	}

	return banner, err
}

func (bannerMongo BannerRepositoryMongo) ExistByEventID(eventID string) (bool, error) {
	id, err := primitive.ObjectIDFromHex(eventID)
	filter := bson.D{{"event_id", id}}
	count, err := bannerMongo.ConnectionDB.Collection(bannerCollection).CountDocuments(context.TODO(), filter)
	log.Printf("[info] count %s", count)
	if err != nil {
		log.Println(err)
	}
	if count > 0 {
		return true, nil
	}

	return false, nil
}
