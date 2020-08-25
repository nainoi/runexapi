package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"thinkdev.app/think/runex/runexapi/model"
)

type CouponRepository interface {
	CreateCoupon(coupon model.Coupon) (string, error)
	ExistByCode(code string) (bool, error)
	ExistByCodeForEdit(code string, couponID string) (bool, error)
	EditCoupon(couponID string, coupon model.EditCouponForm) error
	DeleteCouponByID(couponID string) error
	GetCouponAll() ([]model.Coupon, error)
	GetCouponByCode(code string) (model.Coupon, error)
}

type CouponRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

const (
	couponCollection = "coupon"
)

func (couponMongo CouponRepositoryMongo) CreateCoupon(coupon model.Coupon) (string, error) {

	coupon.CreatedAt = time.Now()
	coupon.UpdatedAt = time.Now()
	res, err := couponMongo.ConnectionDB.Collection(couponCollection).InsertOne(context.TODO(), coupon)
	if err != nil {
		log.Fatal(res)
	}
	fmt.Println("Inserted a single document: ", res.InsertedID)
	return res.InsertedID.(primitive.ObjectID).Hex(), err

}

func (couponMongo CouponRepositoryMongo) ExistByCode(code string) (bool, error) {

	filter := bson.D{{"coupon_code", code}}

	count, err := couponMongo.ConnectionDB.Collection(couponCollection).CountDocuments(context.TODO(), filter)
	log.Printf("[info] count %s", count)
	if err != nil {
		log.Fatal(err)
		log.Println(err)
		return true, err
	}
	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (couponMongo CouponRepositoryMongo) ExistByCodeForEdit(code string, couponID string) (bool, error) {

	var coupon model.Coupon
	id, err := primitive.ObjectIDFromHex(couponID)
	filter := bson.M{"_id": id}
	err2 := couponMongo.ConnectionDB.Collection(couponCollection).FindOne(context.TODO(), filter).Decode(&coupon)
	log.Printf("[info] coupon %s", err2)
	if err2 != nil {
		log.Fatal(err2)
		return true, err
		//return true, err2
	}
	if coupon.CouponCode == code {
		return false, nil
	}

	filter2 := bson.D{{"coupon_code", code}}
	count, err := couponMongo.ConnectionDB.Collection(couponCollection).CountDocuments(context.TODO(), filter2)
	log.Printf("[info] count %s", count)
	if err != nil {
		log.Fatal(err)
		log.Println(err)
		return true, err
	}
	if count > 0 {
		return true, nil
	}

	return false, nil
}
func (couponMongo CouponRepositoryMongo) EditCoupon(couponID string, coupon model.EditCouponForm) error {
	objectID, err := primitive.ObjectIDFromHex(couponID)
	filter := bson.D{{"_id", objectID}}
	coupon.UpdatedAt = time.Now()
	updated := bson.M{"$set": coupon}
	res, err := couponMongo.ConnectionDB.Collection(couponCollection).UpdateOne(context.TODO(), filter, updated)
	if err != nil {
		//log.Fatal(res)
		log.Printf("[info] err %s", res)
		return err
	}

	return nil
}

func (couponMongo CouponRepositoryMongo) DeleteCouponByID(couponID string) error {

	id, err := primitive.ObjectIDFromHex(couponID)
	if err != nil {
		return err
		log.Fatal(err)
	}
	filter := bson.M{"_id": id}
	deleteResult, err2 := couponMongo.ConnectionDB.Collection(couponCollection).DeleteOne(context.TODO(), filter)
	log.Printf("[info] cur %s", err2)
	if err2 != nil {
		return err2
	}
	fmt.Printf("Deleted %v documents in the coupon collection\n", deleteResult.DeletedCount)
	return nil
}

func (couponMongo CouponRepositoryMongo) GetCouponAll() ([]model.Coupon, error) {
	var coupons []model.Coupon
	cur, err := couponMongo.ConnectionDB.Collection(couponCollection).Find(context.TODO(), bson.D{{}})
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {
		var u model.Coupon
		// decode the document
		if err := cur.Decode(&u); err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("post: %+v\n", p)
		coupons = append(coupons, u)
	}

	return coupons, err
}

func (couponMongo CouponRepositoryMongo) GetCouponByCode(code string) (model.Coupon, error) {

	var coupon model.Coupon

	filter := bson.M{"coupon_code": code}
	err2 := couponMongo.ConnectionDB.Collection(couponCollection).FindOne(context.TODO(), filter).Decode(&coupon)
	log.Printf("[info Event] cur %s", err2)
	// if err2 != nil {
	// 	return coupon, err2
	// }

	return coupon, err2
}
