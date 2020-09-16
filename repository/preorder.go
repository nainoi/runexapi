package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"thinkdev.app/think/runex/runexapi/model"
)

const (
	preorderCollection = "preorder"
)

//PreorderRepositoryMongo mongo reference type
type PreorderRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

// FindPreorder search pre order data
func (preorderMongo PreorderRepositoryMongo) FindPreorder(formRequest model.FindPreOrderRequest) ([]model.PreOrder, error) {

	//filter := bson.M{"event_id": formRequest.Keyword}
	filter := bson.M{
		"$or": []interface{}{
			bson.M{"FirstName": primitive.Regex{Pattern: formRequest.Keyword, Options: ""}},
			bson.M{"LastName": primitive.Regex{Pattern: formRequest.Keyword, Options: ""}},
			bson.M{"E-BiB": primitive.Regex{Pattern: formRequest.Keyword, Options: ""}},
			bson.M{"Tel_No": primitive.Regex{Pattern: formRequest.Keyword, Options: ""}},
		},
	}
	options := options.Find()
	options.SetLimit(20)
	var preOrders = []model.PreOrder{}
	cur, err := preorderMongo.ConnectionDB.Collection(preorderCollection).Find(context.TODO(), filter, options)
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}
	for cur.Next(context.TODO()) {
		var u model.PreOrder
		// decode the document
		if err := cur.Decode(&u); err != nil {
			log.Println(err)
		}
		//fmt.Printf("post: %+v\n", p)
		preOrders = append(preOrders, u)
	}

	return preOrders, err
}
