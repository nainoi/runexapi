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

type CategoryRepository interface {
	AddCategory(category model.CategoryMaster) (string, error)
	GetCategoryAll() ([]model.CategoryMaster, error)
	EditCategory(categoryID string, category model.CategoryUpdateForm) error
	DeleteCategoryByID(categoryID string) error
}
type CategoryRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

const (
	categoryCollection = "category"
)

func (categoryMongo CategoryRepositoryMongo) AddCategory(category model.CategoryMaster) (string, error) {

	category.CreatedAt = time.Now()
	res, err := categoryMongo.ConnectionDB.Collection(categoryCollection).InsertOne(context.TODO(), category)
	if err != nil {
		log.Fatal(res)
	}
	fmt.Println("Inserted a single document: ", res.InsertedID)
	return res.InsertedID.(primitive.ObjectID).Hex(), err
}

func (categoryMongo CategoryRepositoryMongo) EditCategory(categoryID string, category model.CategoryUpdateForm) error {

	objectID, err := primitive.ObjectIDFromHex(categoryID)
	filter := bson.D{{"_id", objectID}}
	category.UpdatedAt = time.Now()
	updated := bson.M{"$set": category}
	res, err := categoryMongo.ConnectionDB.Collection(categoryCollection).UpdateOne(context.TODO(), filter, updated)
	if err != nil {
		//log.Fatal(res)
		log.Printf("[info] err %s", res)
		return err
	}

	return nil
}

func (categoryMongo CategoryRepositoryMongo) GetCategoryAll() ([]model.CategoryMaster, error) {
	var category []model.CategoryMaster
	cur, err := categoryMongo.ConnectionDB.Collection(categoryCollection).Find(context.TODO(), bson.D{{}})
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {
		var u model.CategoryMaster
		// decode the document
		if err := cur.Decode(&u); err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("post: %+v\n", p)
		category = append(category, u)
	}

	return category, err
}

func (categoryMongo CategoryRepositoryMongo) DeleteCategoryByID(categoryID string) error {

	id, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		log.Fatal(err)
		return err

	}
	filter := bson.M{"_id": id}
	deleteResult, err2 := categoryMongo.ConnectionDB.Collection(categoryCollection).DeleteOne(context.TODO(), filter)
	log.Printf("[info] cur %s", err2)
	if err2 != nil {
		return err2
	}
	fmt.Printf("Deleted %v documents in the Category collection\n", deleteResult.DeletedCount)
	return nil
}
