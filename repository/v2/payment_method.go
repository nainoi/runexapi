package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"thinkdev.app/think/runex/runexapi/config/db"
	"thinkdev.app/think/runex/runexapi/model"
)

const (
	paymentMethodCollection   = "payment_method"
)

//Get active payment method
func Get() []model.PaymentMethod {
	filter := bson.D{primitive.E{Key: "is_active", Value: true}}
	curr, _ := db.DB.Collection(paymentMethodCollection).Find(context.TODO(), filter)
	var payments = []model.PaymentMethod{}
	for curr.Next(context.TODO()) {
		var p model.PaymentMethod
		// decode the document
		if err := curr.Decode(&p); err != nil {
			log.Println(err)
		}
		//fmt.Printf("post: %+v\n", p)
		payments = append(payments, p)
	}
	return payments
}

//All active payment method
func All() []model.PaymentMethod {
	filter := bson.D{{}}
	curr, _ := db.DB.Collection(paymentMethodCollection).Find(context.TODO(), filter)
	var payments = []model.PaymentMethod{}
	for curr.Next(context.TODO()) {
		var p model.PaymentMethod
		// decode the document
		if err := curr.Decode(&p); err != nil {
			log.Println(err)
		}
		//fmt.Printf("post: %+v\n", p)
		payments = append(payments, p)
	}
	return payments
}