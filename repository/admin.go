package repository

import (
	"context"
	"log"
	"time"

	"bitbucket.org/suthisakch/runex/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminRepository interface {
	EditShppingAddress(registerID string, shippingAddress model.ShipingAddressUpdateForm) error
	UpdateSlip(slipUpdate model.SlipUpdateForm, updateBy string) error
	GetSlipByReg(regID string) (model.SlipHistory, error)
	GetRegEventByID(regID string) (model.Register, error) 
}

type AdminRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

const (
	slipHistoryCollection = "slip_history"
)

func (adminMongo AdminRepositoryMongo) EditShppingAddress(registerID string, shippingAddress model.ShipingAddressUpdateForm) error {
	objectID, err := primitive.ObjectIDFromHex(registerID)
	//productObjectID, err := primitive.ObjectIDFromHex(product.ProductID)
	filter := bson.D{{"_id", objectID}}

	shippingAddress.UpdatedAt = time.Now()

	update := bson.M{"$set": bson.M{"shiping_address": shippingAddress}}

	res, err := adminMongo.ConnectionDB.Collection("register").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		//log.Fatal(res)
		log.Printf("[info] err %s", res)
		return err
	}

	return nil
}

func (adminMongo AdminRepositoryMongo) UpdateSlip(slipUpdate model.SlipUpdateForm, updateBy string) error {
	userID, err := primitive.ObjectIDFromHex(updateBy)
	regObjectID, err := primitive.ObjectIDFromHex(slipUpdate.RegID)
	filter := bson.D{{"reg_id", regObjectID}}

	count, err := adminMongo.ConnectionDB.Collection(slipHistoryCollection).CountDocuments(context.TODO(), filter)
	log.Printf("[info] count %s", count)
	if err != nil {
		log.Println(err)
		return err
	}

	if count > 0 {
		dataInfo := model.SlipTransferUpdate{
			Amount:       slipUpdate.Slip.Amount,
			DateTransfer: slipUpdate.Slip.DateTransfer,
			TimeTransfer: slipUpdate.Slip.TimeTransfer,
			Image:        slipUpdate.Slip.Image,
			Remark:       slipUpdate.Slip.Remark,
			OrderID:      slipUpdate.Slip.OrderID,
			BankAccount:  slipUpdate.Slip.BankAccount,
			Comment:      slipUpdate.Comment,
			CreatDated:   time.Now(),
			UpdateBy:     userID,
		}
		update := bson.M{"$push": bson.M{"slips": dataInfo}}
		res, err := adminMongo.ConnectionDB.Collection(slipHistoryCollection).UpdateOne(context.TODO(), filter, update)
		if err != nil {
			//log.Fatal(res)
			log.Printf("[info] err %s", res)
			return err
		}
	} else {
		var arrHistoryInfo []model.SlipTransferUpdate
		dataInfo := model.SlipTransferUpdate{
			Amount:       slipUpdate.Slip.Amount,
			DateTransfer: slipUpdate.Slip.DateTransfer,
			TimeTransfer: slipUpdate.Slip.TimeTransfer,
			Image:        slipUpdate.Slip.Image,
			Remark:       slipUpdate.Slip.Remark,
			OrderID:      slipUpdate.Slip.OrderID,
			BankAccount:  slipUpdate.Slip.BankAccount,
			Comment:      slipUpdate.Comment,
			CreatDated:   time.Now(),
			UpdateBy:     userID,
		}
		objectID, err := primitive.ObjectIDFromHex(slipUpdate.RegID)

		arrHistoryInfo = append(arrHistoryInfo, dataInfo)
		historyModel := model.SlipHistory{
			RegID: objectID,
			Slips: arrHistoryInfo,
		}
		res, err := adminMongo.ConnectionDB.Collection(slipHistoryCollection).InsertOne(context.TODO(), historyModel)
		if err != nil {
			log.Fatal(res)
			return err
		}
	}

	return nil
}

func (adminMongo AdminRepositoryMongo) GetSlipByReg(regID string) (model.SlipHistory, error) {
	objectID, err := primitive.ObjectIDFromHex(regID)
	var slips model.SlipHistory
	filter := bson.D{{"reg_id", objectID}}

	err = adminMongo.ConnectionDB.Collection(slipHistoryCollection).FindOne(context.TODO(), filter).Decode(&slips)
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}

	return slips, err
}

func (adminMongo AdminRepositoryMongo) GetRegEventByID(regID string) (model.Register, error) {

	var regEvent model.Register
	id, err := primitive.ObjectIDFromHex(regID)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.M{"_id": id}
	err2 := adminMongo.ConnectionDB.Collection(registerCollection).FindOne(context.TODO(), filter).Decode(&regEvent)
	log.Printf("[info RegEvent] cur %s", err2)
	if err2 != nil {
		log.Fatal(err2)
	}

	return regEvent, err2
}