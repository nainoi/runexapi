package repository

import (
	"context"
	"log"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"thinkdev.app/think/runex/runexapi/config/db"
	"thinkdev.app/think/runex/runexapi/model"
)

const (
	provinceCollection = "province"
	tambonCollection   = "tambon"
)

//AllProvince repo
func AllProvince() []model.Province {
	filter := bson.D{{}}
	curr, _ := db.DB.Collection(provinceCollection).Find(context.TODO(), filter)
	var provinces = []model.Province{}
	for curr.Next(context.TODO()) {
		var p model.Province
		// decode the document
		if err := curr.Decode(&p); err != nil {
			log.Println(err)
		}
		//fmt.Printf("post: %+v\n", p)
		provinces = append(provinces, p)
	}
	return provinces
}

//TambonByZipcode repo
func TambonByZipcode(zipcode string) []model.Tambon {
	code, _ := strconv.Atoi(zipcode)
	filter := bson.D{primitive.E{Key: "zipcode", Value: code}}
	curr, _ := db.DB.Collection(tambonCollection).Find(context.TODO(), filter)
	var tambons = []model.Tambon{}
	for curr.Next(context.TODO()) {
		var p model.Tambon
		// decode the document
		if err := curr.Decode(&p); err != nil {
			log.Println(err)
		}
		//fmt.Printf("post: %+v\n", p)
		tambons = append(tambons, p)
	}
	return tambons
}

//Province repo
func Province(province string) []model.Tambon {
	filter := bson.D{primitive.E{Key: "province", Value: bson.M{"$regex": province}}}
	curr, _ := db.DB.Collection(tambonCollection).Find(context.TODO(), filter)
	var tambons = []model.Tambon{}
	for curr.Next(context.TODO()) {
		var p model.Tambon
		// decode the document
		if err := curr.Decode(&p); err != nil {
			log.Println(err)
		}
		//fmt.Printf("post: %+v\n", p)
		tambons = append(tambons, p)
	}
	return tambons
}

//Amphoe repo
func Amphoe(amphoe string) []model.Tambon {
	filter := bson.D{primitive.E{Key: "amphoe", Value: bson.M{"$regex": amphoe}}}
	curr, _ := db.DB.Collection(tambonCollection).Find(context.TODO(), filter)
	var tambons = []model.Tambon{}
	for curr.Next(context.TODO()) {
		var p model.Tambon
		// decode the document
		if err := curr.Decode(&p); err != nil {
			log.Println(err)
		}
		//fmt.Printf("post: %+v\n", p)
		tambons = append(tambons, p)
	}
	return tambons
}

//District repo
func District(district string) []model.Tambon {
	filter := bson.D{primitive.E{Key: "district", Value: bson.M{"$regex": district}}}
	curr, _ := db.DB.Collection(tambonCollection).Find(context.TODO(), filter)
	var tambons = []model.Tambon{}
	for curr.Next(context.TODO()) {
		var p model.Tambon
		// decode the document
		if err := curr.Decode(&p); err != nil {
			log.Println(err)
		}
		//fmt.Printf("post: %+v\n", p)
		tambons = append(tambons, p)
	}
	return tambons
}

//TambonAll repo
func TambonAll() []model.Tambon {
	filter := bson.D{{}}
	curr, _ := db.DB.Collection(tambonCollection).Find(context.TODO(), filter)
	var tambons = []model.Tambon{}
	for curr.Next(context.TODO()) {
		var p model.Tambon
		// decode the document
		if err := curr.Decode(&p); err != nil {
			log.Println(err)
		}
		//fmt.Printf("post: %+v\n", p)
		tambons = append(tambons, p)
	}
	return tambons
}
