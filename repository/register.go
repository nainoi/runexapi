package repository

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/omise/omise-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"thinkdev.app/think/runex/runexapi/config"
	"thinkdev.app/think/runex/runexapi/model"

	"thinkdev.app/think/runex/runexapi/api/mail"
)

type RegisterRepository interface {
	GetRegisterAll() ([]model.Register, error)
	AddRegister(register model.Register) (model.Register, error)
	AddRaceRegister(register model.Register) (model.Register, error)
	AddMerChant(userID string, eventID string, charge omise.Charge) error
	EditRegister(registerID string, register model.Register) error
	GetRegisterByEvent(eventID string) ([]model.Register, error)
	NotifySlipRegister(registerID string, slip model.SlipTransfer) error
	AdminNotifySlipRegister(registerID string, slip model.SlipTransfer) error
	CountByEvent(eventID string) (int64, error)
	GetRegEventByID(regID string) (model.Register, error)
	CheckUserRegisterEvent(eventID string, userID string) (bool, error)
	GetRegisterByUserID(userID string) ([]model.Register, error)
	GetRegisterByUserAndEvent(userID string, id string) (model.Register, error)
	SendMailRegister(registerID string) error
	GetRegisterActivateEvent(userID string) ([]model.EventRegInfo, error)
	GetRegisterReport(formRequest model.DataRegisterRequest) (model.ReportRegister, error)
	GetRegisterReportAll(formRequest model.DataRegisterRequest) (model.ReportRegister, error)
	FindPersonRegEvent(formRequest model.DataRegisterRequest) (model.ReportRegister, error)
	UpdateStatusRegister(registerID string, status string, userID string) error
}

type RegisterRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

const (
	registerCollection = "register"
	merchantCollection = "merchant"
)

func (registerMongo RegisterRepositoryMongo) GetRegisterAll() ([]model.Register, error) {
	var register []model.Register
	cur, err := registerMongo.ConnectionDB.Collection(registerCollection).Find(context.TODO(), bson.D{{}})
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {
		var u model.Register
		// decode the document
		if err := cur.Decode(&u); err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("post: %+v\n", p)
		register = append(register, u)
	}

	return register, err
}

func (registerMongo RegisterRepositoryMongo) AddRegister(register model.Register) (model.Register, error) {
	register.CreatedAt = time.Now()
	register.UpdatedAt = time.Now()
	res, err := registerMongo.ConnectionDB.Collection(registerCollection).InsertOne(context.TODO(), register)
	if err != nil {
		log.Fatal(res)
	}
	fmt.Println("Inserted a single document: ", res.InsertedID)

	regModel, err2 := registerMongo.GetRegEventByID(res.InsertedID.(primitive.ObjectID).Hex())
	if err2 != nil {
		return regModel, err
	}

	err3 := registerMongo.SendMailRegister(res.InsertedID.(primitive.ObjectID).Hex())
	if err3 != nil {
		log.Fatal(err3)
	}
	return regModel, err
}

//AddRaceRegister repo insert run event
func (registerMongo RegisterRepositoryMongo) AddRaceRegister(register model.Register) (model.Register, error) {
	register.CreatedAt = time.Now()
	register.UpdatedAt = time.Now()
	filterBib := bson.D{bson.E{Key: "event_id", Value: register.EventID}}
	var ebibEvent model.EbibEvent
	err := registerMongo.ConnectionDB.Collection(ebibCollection).FindOne(context.TODO(), filterBib).Decode(&ebibEvent)
	if err != nil {
		log.Println("find ebib in AddRaceRegister", err)
		// var ebib model.EbibEvent
		// ebib.EventID = register.EventID
		// ebib.LastNo = 0
		// ebib.CreatedAt = time.Now()
		// ebib.UpdatedAt = time.Now()
		// registerMongo.ConnectionDB.Collection(ebibCollection).InsertOne(context.TODO(),ebib)
	}
	for index, _ := range register.TicketOptions {
		log.Println(fmt.Sprintf("%05d", (ebibEvent.LastNo + int64(index+1))))
		// element.RegisterNumber = fmt.Sprintf("%05d", (ebibEvent.LastNo + int64(index + 1)))
		register.TicketOptions[index].RegisterNumber = fmt.Sprintf("%05d", (ebibEvent.LastNo + int64(index+1)))
	}
	res, err := registerMongo.ConnectionDB.Collection(registerCollection).InsertOne(context.TODO(), register)
	if err != nil {
		log.Fatal(res)
	}
	fmt.Println("Inserted a single document: ", res.InsertedID)
	regModel, err2 := registerMongo.GetRegEventByID(res.InsertedID.(primitive.ObjectID).Hex())
	if err2 != nil {
		return regModel, err
	}
	ebibEvent.LastNo += int64(len(regModel.TicketOptions))
	ebibEvent.UpdatedAt = time.Now()
	updated := bson.M{"$set": ebibEvent}
	_, err = registerMongo.ConnectionDB.Collection(ebibCollection).UpdateOne(context.TODO(), filterBib, updated)
	if err != nil {
		log.Print(res)
	}
	go registerMongo.SendRaceMailRegister(res.InsertedID.(primitive.ObjectID).Hex())
	// if err3 != nil {
	// 	log.Fatal(err3)
	// }
	return regModel, err
}

func (registerMongo RegisterRepositoryMongo) AddMerChant(userID string, eventID string, charge omise.Charge) error {
	uid, err := primitive.ObjectIDFromHex(userID)
	eid, err := primitive.ObjectIDFromHex(eventID)
	typePay := ""
	if charge.Source != nil {
		if len(charge.Source.Type) > 0 {
			typePay = charge.Source.Type
		}
	}

	var merchant model.Merchant
	merchant.UserID = uid
	merchant.EventID = eid
	merchant.PaymentType = typePay
	merchant.Status = string(charge.Status)
	merchant.OrderID = charge.Object
	merchant.OmiseID = charge.ID
	merchant.TotalPrice = charge.Amount / 100
	merchant.CreatedAt = time.Now()
	merchant.UpdatedAt = time.Now()
	//fmt.Println("Inserted a single document: ", merchant)
	res, err := registerMongo.ConnectionDB.Collection(merchantCollection).InsertOne(context.TODO(), merchant)
	if err != nil {
		log.Println(res)
	}
	return err
}

func (registerMongo RegisterRepositoryMongo) EditRegister(registerID string, register model.Register) error {
	objectID, err := primitive.ObjectIDFromHex(registerID)
	filter := bson.D{{"_id", objectID}}
	register.UpdatedAt = time.Now()
	updated := bson.M{"$set": register}
	res, err := registerMongo.ConnectionDB.Collection(registerCollection).UpdateOne(context.TODO(), filter, updated)
	if err != nil {
		//log.Fatal(res)
		log.Printf("[info] err %s", res)
		return err

	}

	// if register.Event.Category.Name == "Run" {
	// 	go registerMongo.SendRaceMailRegister(registerID)
	// } else {
	// 	go registerMongo.SendMailRegister(registerID)
	// }

	return nil
}

func (registerMongo RegisterRepositoryMongo) GetRegisterByEvent(eventID string) ([]model.Register, error) {

	var register []model.Register
	objectID, err := primitive.ObjectIDFromHex(eventID)
	filter := bson.D{{"event_id", objectID}}
	cur, err := registerMongo.ConnectionDB.Collection(registerCollection).Find(context.TODO(), filter)
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {
		var u model.Register
		// decode the document
		if err := cur.Decode(&u); err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("post: %+v\n", p)
		register = append(register, u)
	}

	return register, err
}

func (registerMongo RegisterRepositoryMongo) GetRegisterByUserID(userID string) ([]model.Register, error) {

	var register []model.Register
	objectID, err := primitive.ObjectIDFromHex(userID)
	filter := bson.D{{"user_id", objectID}}
	cur, err := registerMongo.ConnectionDB.Collection(registerCollection).Find(context.TODO(), filter)
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {
		var u model.Register

		// decode the document
		if err := cur.Decode(&u); err != nil {
			log.Print(err)
		}
		var event model.EventReg
		filter = bson.D{{"_id", u.EventID}}
		err := registerMongo.ConnectionDB.Collection("event").FindOne(context.TODO(), filter).Decode(&event)
		if err != nil {
			log.Print(err)
		}
		u.Event = event
		// var event model.Event
		// registerMongo.ConnectionDB.Collection(eventCollection).FindOne(context.TODO(), bson.D{{"_id", u.EventID}}).Decode(&event)
		// //fmt.Printf("post: %+v\n", p)
		//u.Event = event
		register = append(register, u)
	}

	return register, err
}

//GetRegisterByUserAndEvent detail reg
func (registerMongo RegisterRepositoryMongo) GetRegisterByUserAndEvent(userID string, id string) (model.Register, error) {

	var register model.Register
	objectID, err := primitive.ObjectIDFromHex(userID)
	regID, err := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"user_id", objectID}, {"_id", regID}}
	err = registerMongo.ConnectionDB.Collection(registerCollection).FindOne(context.TODO(), filter).Decode(&register)
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}

	// for cur.Next(context.TODO()) {
	// 	var u model.Register
	// 	// decode the document
	// 	if err := cur.Decode(&u); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	// var event model.Event
	// 	// registerMongo.ConnectionDB.Collection(eventCollection).FindOne(context.TODO(),bson.D{{"_id",u.EventID}}).Decode(&event)
	// 	// //fmt.Printf("post: %+v\n", p)
	// 	// u.Event = event
	// 	register = append(register, u)
	// }

	return register, err
}

func (registerMongo RegisterRepositoryMongo) NotifySlipRegister(registerID string, slip model.SlipTransfer) error {
	objectID, err := primitive.ObjectIDFromHex(registerID)
	filter := bson.D{{"_id", objectID}}

	updated := bson.M{"$set": bson.M{"slip": slip, "status": config.PAYMENT_WAITING_APPROVE, "updated_at": time.Now()}}

	res, err := registerMongo.ConnectionDB.Collection(registerCollection).UpdateOne(context.TODO(), filter, updated)
	if err != nil {
		//log.Fatal(res)
		log.Printf("[info] err %s", res)
		return err
	}

	var register model.Register
	err = registerMongo.ConnectionDB.Collection(registerCollection).FindOne(context.TODO(), filter).Decode(&register)
	if err == nil {
		if register.Event.Category.Name == "Run" {
			go registerMongo.SendRaceMailRegister(registerID)
		} else {
			go registerMongo.SendMailRegister(registerID)
		}
	}

	return nil
}

//AdminNotifySlipRegister admin upload slip and update status
func (registerMongo RegisterRepositoryMongo) AdminNotifySlipRegister(registerID string, slip model.SlipTransfer) error {
	objectID, err := primitive.ObjectIDFromHex(registerID)
	filter := bson.D{{"_id", objectID}}

	var register model.Register
	err = registerMongo.ConnectionDB.Collection(registerCollection).FindOne(context.TODO(), filter).Decode(&register)
	paymentType := config.PAYMENT_TRANSFER
	if err == nil {
		if register.PaymentType != "" {
			paymentType = register.PaymentType
		}
	}

	updated := bson.M{"$set": bson.M{"slip": slip, "status": config.PAYMENT_WAITING_APPROVE, "payment_type": paymentType, "updated_at": time.Now()}}

	res, err := registerMongo.ConnectionDB.Collection(registerCollection).UpdateOne(context.TODO(), filter, updated)
	if err != nil {
		//log.Fatal(res)
		log.Printf("[info] err %s", res)
		return err
	}

	err = registerMongo.ConnectionDB.Collection(registerCollection).FindOne(context.TODO(), filter).Decode(&register)

	if err == nil {
		if register.Event.Category.Name == "Run" {
			go registerMongo.SendRaceMailRegister(registerID)
		} else {
			go registerMongo.SendMailRegister(registerID)
		}
	}

	return nil
}

func (registerMongo RegisterRepositoryMongo) CountByEvent(eventID string) (int64, error) {

	objectID, err := primitive.ObjectIDFromHex(eventID)
	filter := bson.D{{"event_id", objectID}}

	count, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filter)
	log.Printf("[info] count %s", count)
	if err != nil {
		log.Println(err)
	}

	return count, err
}

func (registerMongo RegisterRepositoryMongo) GetRegEventByID(regID string) (model.Register, error) {

	var regEvent model.Register
	id, err := primitive.ObjectIDFromHex(regID)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.M{"_id": id}
	err2 := registerMongo.ConnectionDB.Collection(registerCollection).FindOne(context.TODO(), filter).Decode(&regEvent)
	log.Printf("[info RegEvent] cur %s", err2)
	if err2 != nil {
		log.Fatal(err2)
	}

	return regEvent, err2
}

func (registerMongo RegisterRepositoryMongo) CheckUserRegisterEvent(eventID string, userID string) (bool, error) {
	userObjectID, _ := primitive.ObjectIDFromHex(userID)
	eventObjectID, _ := primitive.ObjectIDFromHex(eventID)
	filter := bson.D{{"user_id", userObjectID}, {"event_id", eventObjectID}}
	count, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filter)
	log.Printf("[info] count %s", count)
	if err != nil {
		log.Println(err)
	}
	if count > 0 {

		return true, nil
	}

	return false, nil
}

func (registerMongo RegisterRepositoryMongo) SendMailRegister(registerID string) error {
	var register model.Register
	var user model.UserProvider
	var mailTemplate model.EmailTemplateData2
	objectID, _ := primitive.ObjectIDFromHex(registerID)
	filter := bson.D{{"_id", objectID}}
	err := registerMongo.ConnectionDB.Collection("register").FindOne(context.TODO(), filter).Decode(&register)
	if err != nil {
		log.Fatal(err)
	}

	filterUser := bson.D{{"_id", register.UserID}}
	err = registerMongo.ConnectionDB.Collection("user").FindOne(context.TODO(), filterUser).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	filterEvent := bson.D{{"event_id", register.EventID}}
	count, err := registerMongo.ConnectionDB.Collection("register").CountDocuments(context.TODO(), filterEvent)
	registerNumber := fmt.Sprintf("%05d", count)
	address := register.ShipingAddress.Address + " City " + register.ShipingAddress.City + " District " + register.ShipingAddress.District + " Province " + register.ShipingAddress.Province + " " + register.ShipingAddress.Zipcode

	updated := bson.M{"$set": bson.M{"register_number": registerNumber}}
	res, err := registerMongo.ConnectionDB.Collection("register").UpdateOne(context.TODO(), filter, updated)
	if err != nil {
		log.Fatal(err)
		log.Printf("[info] err %s", res)
		return err
	}

	mailTemplate.Name = user.FirstName + " " + user.LastName
	mailTemplate.RefID = registerID
	mailTemplate.Email = user.Email
	mailTemplate.Phone = user.Phone
	mailTemplate.IdentificationNumber = user.CitycenID
	mailTemplate.ContactPhone = user.Phone
	mailTemplate.CompetitionType = register.Tickets[0].TicketDetail.Title
	mailTemplate.Price = register.TotalPrice
	mailTemplate.RegisterNumber = registerNumber
	mailTemplate.TicketName = register.Tickets[0].TicketDetail.Title
	mailTemplate.Status = register.Status
	mailTemplate.PaymentType = register.PaymentType
	mailTemplate.ShipingAddress = address

	fmt.Println("Mail Object :  ", mailTemplate)

	if register.PaymentType == "PAYMENT_FREE" {
		mail.SendRegFreeEventMail2(mailTemplate)
	} else {
		mail.SendRegEventMail2(mailTemplate)
	}

	return err
}

func (registerMongo RegisterRepositoryMongo) SendRaceMailRegister(registerID string) error {
	var register model.Register
	var user model.UserProvider
	var mailTemplate model.EmailTemplateData2
	objectID, _ := primitive.ObjectIDFromHex(registerID)
	filter := bson.D{{"_id", objectID}}
	err := registerMongo.ConnectionDB.Collection("register").FindOne(context.TODO(), filter).Decode(&register)
	if err != nil {
		log.Fatal(err)
	}

	filterUser := bson.D{{"_id", register.UserID}}
	err = registerMongo.ConnectionDB.Collection("user").FindOne(context.TODO(), filterUser).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")

	address := register.ShipingAddress.Address + " City " + register.ShipingAddress.City + " District " + register.ShipingAddress.District + " Province " + register.ShipingAddress.Province + " " + register.ShipingAddress.Zipcode

	mailTemplate.Name = user.FirstName + " " + user.LastName
	mailTemplate.RefID = registerID
	mailTemplate.Email = user.Email
	mailTemplate.Phone = user.Phone
	mailTemplate.IdentificationNumber = user.CitycenID
	mailTemplate.ContactPhone = user.Phone
	mailTemplate.CompetitionType = register.Tickets[0].TicketDetail.Title
	mailTemplate.Price = register.TotalPrice
	mailTemplate.RegisterNumber = register.TicketOptions[0].RegisterNumber
	mailTemplate.TicketName = register.Tickets[0].TicketDetail.Title
	mailTemplate.EventName = register.Event.Name
	mailTemplate.Status = register.Status
	mailTemplate.PaymentType = register.PaymentType
	mailTemplate.ShipingAddress = address
	mailTemplate.TicketOptions = register.TicketOptions
	mailTemplate.RegisterDate = register.CreatedAt.In(loc)

	//fmt.Println("Mail Object :  ", mailTemplate)

	if register.PaymentType == "PAYMENT_FREE" {
		mail.SendRegRaceRun(mailTemplate)
	} else {
		mail.SendRegRaceRun(mailTemplate)
	}

	return err
}

/*func (registerMongo RegisterRepositoryMongo) GetRegisterActivateEvent(userID string) ([]model.EventRegInfo, error) {
	var register []model.Register
	objectID, err := primitive.ObjectIDFromHex(userID)
	filter := bson.D{{"user_id", objectID}, {"status", "PAYMENT_SUCCESS"}}
	cur, err := registerMongo.ConnectionDB.Collection(registerCollection).Find(context.TODO(), filter)
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {
		var u model.Register
		// decode the document
		if err := cur.Decode(&u); err != nil {
			log.Fatal(err)
		}

		register = append(register, u)
	}
	now := time.Now()
	var events []model.EventRegInfo
	for _, item := range register {
		skip := false
		fmt.Println("registerID :  ", item.ID)
		var event model.Event
		eventID := item.EventID
		//id, _ := primitive.ObjectIDFromHex(eventID)
		filter := bson.M{"_id": eventID}
		fmt.Println("eventID :  ", eventID)
		count, err := registerMongo.ConnectionDB.Collection("event").CountDocuments(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		if count > 0 {
			err := registerMongo.ConnectionDB.Collection("event").FindOne(context.TODO(), filter).Decode(&event)
			if err != nil {
				log.Fatal(err)
			}
			if event.StartEvent.Before(now) && event.EndEvent.After(now) && event.IsActive == true {

				eventsNew := model.EventRegInfo{
					ID:          event.ID,
					Name:        event.Name,
					Description: event.Description,
					Body:        event.Body,
					IsActive:    event.IsActive,
					StartEvent:  event.StartEvent,
					EndEvent:    event.EndEvent,
				}
				for _, u := range events {
					if eventsNew == u {
						skip = true
						break
					}
				}
				if !skip {
					events = append(events, eventsNew)
				}
			}
		}
	}
	//uniqueEvents := unique(events)
	return events, nil
}*/

func (registerMongo RegisterRepositoryMongo) GetRegisterReport(formRequest model.DataRegisterRequest) (model.ReportRegister, error) {
	var register []model.Register
	var report model.ReportRegister

	objectEventID, err := primitive.ObjectIDFromHex(formRequest.EventID)
	if err != nil {
		log.Fatal(err)

	}

	filterSuccess := bson.D{{"status", config.PAYMENT_SUCCESS}, {"event_id", objectEventID}}
	filterWaiting := bson.D{{"status", config.PAYMENT_WAITING}, {"event_id", objectEventID}}
	filterWaitingApprove := bson.D{{"status", config.PAYMENT_WAITING_APPROVE}, {"event_id", objectEventID}}

	countAll, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), bson.D{{"event_id", objectEventID}})

	countSuccess, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filterSuccess)

	countWaiting, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filterWaiting)

	countWaitingApprove, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filterWaitingApprove)

	if err != nil {
		log.Fatal(err)

	}

	report.PaymentAll = countAll
	report.PaymentSuccess = countSuccess
	report.PaymentWaiting = countWaiting
	report.PaymentWaitingApprove = countWaitingApprove

	skip := (formRequest.PageNumber - 1) * formRequest.NPerPage
	limit := formRequest.NPerPage

	options := options.Find()
	options.SetSkip(skip)
	options.SetLimit(limit)
	options.SetSort(bson.D{{"created_at", -1}})

	filter := bson.M{"event_id": objectEventID}
	if formRequest.Status != "" {
		filter = bson.M{"event_id": objectEventID, "status": formRequest.Status}
	}

	cur, err := registerMongo.ConnectionDB.Collection(registerCollection).Find(context.TODO(), filter, options)
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {
		var u model.Register
		// decode the document
		if err := cur.Decode(&u); err != nil {
			log.Println(err)
		}
		//fmt.Printf("post: %+v\n", p)
		register = append(register, u)
	}

	var dataRegister []model.DataRegister

	for _, item := range register {

		var user model.User
		userID := item.UserID
		// log.Printf("[info] userID %s", userID)
		filterUser := bson.D{{"_id", userID}}
		err := registerMongo.ConnectionDB.Collection("user").FindOne(context.TODO(), filterUser).Decode(&user)

		if err != nil {
			log.Println(err)
		}
		eventUser := item.EventID.Hex() + "." + userID.Hex()
		var activity model.Activity
		filter := bson.D{{"event_user", eventUser}}
		// log.Println(eventUser)
		distanceTotal := 0.0
		if item.Event.Category.Name != "Run" {
			err = registerMongo.ConnectionDB.Collection(activityCollection).FindOne(context.TODO(), filter).Decode(&activity)
			// log.Println(activity)
			if err == nil {
				distanceTotal = math.Round(activity.ToTalDistance*100) / 100
			}
		}
		var dataRegisterNew model.DataRegister
		if item.Event.Category.Name != "Run" {

			if len(item.Tickets) > 0 {
				dataRegisterNew = model.DataRegister{
					RegisterID:    item.ID,
					Firstname:     user.FirstName,
					Lastname:      user.LastName,
					Email:         user.Email,
					Address:       item.ShipingAddress,
					PaymentType:   item.PaymentType,
					PaymentStatus: item.Status,
					RegDate:       item.CreatedAt,
					PaymentDate:   item.UpdatedAt,
					TicketName:    item.Tickets[0].TicketDetail.Title,
					Distance:      item.Tickets[0].TicketDetail.Distance,
					ShirtSize:     item.Tickets[0].Type,
					Price:         item.TotalPrice,
					DistanceTotal: distanceTotal,
					Slip:          item.Slip,
					TicketProduct: item.Tickets,
					TicketOptions: item.TicketOptions,
				}
			} else {
				dataRegisterNew = model.DataRegister{
					RegisterID:    item.ID,
					Firstname:     user.FirstName,
					Lastname:      user.LastName,
					Email:         user.Email,
					Address:       item.ShipingAddress,
					PaymentType:   item.PaymentType,
					PaymentStatus: item.Status,
					RegDate:       item.CreatedAt,
					PaymentDate:   item.UpdatedAt,
					TicketName:    "",
					Distance:      0,
					ShirtSize:     "",
					DistanceTotal: distanceTotal,
					Price:         item.TotalPrice,
					TicketOptions: item.TicketOptions,
				}
			}

		} else {
			dataRegisterNew = model.DataRegister{
				RegisterID:    item.ID,
				Firstname:     user.FirstName,
				Lastname:      user.LastName,
				Email:         user.Email,
				Address:       item.ShipingAddress,
				PaymentType:   item.PaymentType,
				PaymentStatus: item.Status,
				RegDate:       item.CreatedAt,
				PaymentDate:   item.UpdatedAt,
				TicketName:    item.Tickets[0].TicketDetail.Title,
				Distance:      item.Tickets[0].TicketDetail.Distance,
				ShirtSize:     item.Tickets[0].Type,
				Price:         item.TotalPrice,
				DistanceTotal: distanceTotal,
				Slip:          item.Slip,
				TicketOptions: item.TicketOptions,
			}
		}

		dataRegister = append(dataRegister, dataRegisterNew)
	}

	report.Datas = dataRegister
	return report, nil
}

func (registerMongo RegisterRepositoryMongo) GetRegisterReportAll(formRequest model.DataRegisterRequest) (model.ReportRegister, error) {
	var register []model.Register
	var report model.ReportRegister

	objectEventID, err := primitive.ObjectIDFromHex(formRequest.EventID)
	if err != nil {
		log.Fatal(err)

	}

	filterSuccess := bson.D{{"status", config.PAYMENT_SUCCESS}, {"event_id", objectEventID}}
	filterWaiting := bson.D{{"status", config.PAYMENT_WAITING}, {"event_id", objectEventID}}
	filterWaitingApprove := bson.D{{"status", config.PAYMENT_WAITING_APPROVE}, {"event_id", objectEventID}}

	countAll, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), bson.D{{"event_id", objectEventID}})

	countSuccess, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filterSuccess)

	countWaiting, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filterWaiting)

	countWaitingApprove, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filterWaitingApprove)

	if err != nil {
		log.Fatal(err)

	}

	report.PaymentAll = countAll
	report.PaymentSuccess = countSuccess
	report.PaymentWaiting = countWaiting
	report.PaymentWaitingApprove = countWaitingApprove

	// skip := (formRequest.PageNumber - 1) * formRequest.NPerPage
	// limit := formRequest.NPerPage

	options := options.Find()
	// options.SetSkip(skip)
	// options.SetLimit(limit)
	options.SetSort(bson.D{{"created_at", -1}})

	filter := bson.M{"event_id": objectEventID}
	if formRequest.Status != "" {
		filter = bson.M{"event_id": objectEventID, "status": formRequest.Status}
	}

	cur, err := registerMongo.ConnectionDB.Collection(registerCollection).Find(context.TODO(), filter, options)
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {
		var u model.Register
		// decode the document
		if err := cur.Decode(&u); err != nil {
			log.Println(err)
		}
		//fmt.Printf("post: %+v\n", p)
		register = append(register, u)
	}

	var dataRegister []model.DataRegister

	for _, item := range register {

		var user model.User
		userID := item.UserID
		// log.Printf("[info] userID %s", userID)
		filterUser := bson.D{{"_id", userID}}
		err := registerMongo.ConnectionDB.Collection("user").FindOne(context.TODO(), filterUser).Decode(&user)

		if err != nil {
			log.Println(err)
		}
		eventUser := item.EventID.Hex() + "." + userID.Hex()
		var activity model.Activity
		filter := bson.D{{"event_user", eventUser}}
		// log.Println(eventUser)
		distanceTotal := 0.0
		if item.Event.Category.Name != "Run" {
			err = registerMongo.ConnectionDB.Collection(activityCollection).FindOne(context.TODO(), filter).Decode(&activity)
			// log.Println(activity)
			if err == nil {
				distanceTotal = math.Round(activity.ToTalDistance*100) / 100
			}
		}
		var dataRegisterNew model.DataRegister
		if item.Event.Category.Name != "Run" {

			if len(item.Tickets) > 0 {
				dataRegisterNew = model.DataRegister{
					RegisterID:    item.ID,
					Firstname:     user.FirstName,
					Lastname:      user.LastName,
					Email:         user.Email,
					Address:       item.ShipingAddress,
					PaymentType:   item.PaymentType,
					PaymentStatus: item.Status,
					RegDate:       item.CreatedAt,
					PaymentDate:   item.UpdatedAt,
					TicketName:    item.Tickets[0].TicketDetail.Title,
					Distance:      item.Tickets[0].TicketDetail.Distance,
					ShirtSize:     item.Tickets[0].Type,
					Price:         item.TotalPrice,
					DistanceTotal: distanceTotal,
					Slip:          item.Slip,
					TicketProduct: item.Tickets,
					TicketOptions: item.TicketOptions,
				}
			} else {
				dataRegisterNew = model.DataRegister{
					RegisterID:    item.ID,
					Firstname:     user.FirstName,
					Lastname:      user.LastName,
					Email:         user.Email,
					Address:       item.ShipingAddress,
					PaymentType:   item.PaymentType,
					PaymentStatus: item.Status,
					RegDate:       item.CreatedAt,
					PaymentDate:   item.UpdatedAt,
					TicketName:    "",
					Distance:      0,
					ShirtSize:     "",
					Price:         item.TotalPrice,
					DistanceTotal: distanceTotal,
					TicketOptions: item.TicketOptions,
				}
			}

		} else {

			dataRegisterNew = model.DataRegister{
				RegisterID:    item.ID,
				Firstname:     user.FirstName,
				Lastname:      user.LastName,
				Email:         user.Email,
				Address:       item.ShipingAddress,
				PaymentType:   item.PaymentType,
				PaymentStatus: item.Status,
				RegDate:       item.CreatedAt,
				PaymentDate:   item.UpdatedAt,
				TicketName:    item.Tickets[0].TicketDetail.Title,
				Distance:      item.Tickets[0].TicketDetail.Distance,
				ShirtSize:     item.Tickets[0].Type,
				Price:         item.TotalPrice,
				DistanceTotal: distanceTotal,
				Slip:          item.Slip,
				TicketOptions: item.TicketOptions,
			}
		}

		dataRegister = append(dataRegister, dataRegisterNew)
	}

	report.Datas = dataRegister
	return report, nil
}

//GetRegisterReportDated report with dated register
func (registerMongo RegisterRepositoryMongo) GetRegisterReportDated(formRequest model.DataRegisterRequest) (model.ReportRegister, error) {
	var register []model.Register
	var report model.ReportRegister

	objectEventID, err := primitive.ObjectIDFromHex(formRequest.EventID)
	if err != nil {
		log.Fatal(err)

	}

	filterSuccess := bson.D{{"status", config.PAYMENT_SUCCESS}, {"event_id", objectEventID}}
	filterWaiting := bson.D{{"status", config.PAYMENT_WAITING}, {"event_id", objectEventID}}
	filterWaitingApprove := bson.D{{"status", config.PAYMENT_WAITING_APPROVE}, {"event_id", objectEventID}}

	countAll, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), bson.D{{"event_id", objectEventID}})

	countSuccess, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filterSuccess)

	countWaiting, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filterWaiting)

	countWaitingApprove, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filterWaitingApprove)

	if err != nil {
		log.Fatal(err)

	}

	report.PaymentAll = countAll
	report.PaymentSuccess = countSuccess
	report.PaymentWaiting = countWaiting
	report.PaymentWaitingApprove = countWaitingApprove

	skip := (formRequest.PageNumber - 1) * formRequest.NPerPage
	limit := formRequest.NPerPage

	options := options.Find()
	options.SetSkip(skip)
	options.SetLimit(limit)

	filter := bson.M{"event_id": objectEventID}

	cur, err := registerMongo.ConnectionDB.Collection(registerCollection).Find(context.TODO(), filter, options)
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {
		var u model.Register
		// decode the document
		if err := cur.Decode(&u); err != nil {
			log.Println(err)
		}
		//fmt.Printf("post: %+v\n", p)
		register = append(register, u)
	}

	var dataRegister []model.DataRegister

	for _, item := range register {

		var user model.User
		userID := item.UserID
		// log.Printf("[info] userID %s", userID)
		filterUser := bson.D{{"_id", userID}}
		err := registerMongo.ConnectionDB.Collection("user").FindOne(context.TODO(), filterUser).Decode(&user)

		if err != nil {
			log.Println(err)
		}
		eventUser := item.EventID.Hex() + "." + userID.Hex()
		var activity model.Activity
		filter := bson.D{{"event_user", eventUser}}
		// log.Println(eventUser)
		distanceTotal := 0.0
		if item.Event.Category.Name != "Run" {
			err = registerMongo.ConnectionDB.Collection(activityCollection).FindOne(context.TODO(), filter).Decode(&activity)
			// log.Println(activity)
			if err == nil {
				distanceTotal = math.Round(activity.ToTalDistance*100) / 100
			}
		}
		var dataRegisterNew model.DataRegister
		if item.Event.Category.Name != "Run" {
			if len(item.Tickets) > 0 {
				dataRegisterNew = model.DataRegister{
					RegisterID:    item.ID,
					Firstname:     user.FirstName,
					Lastname:      user.LastName,
					Email:         user.Email,
					Address:       item.ShipingAddress,
					PaymentType:   item.PaymentType,
					PaymentStatus: item.Status,
					RegDate:       item.CreatedAt,
					PaymentDate:   item.UpdatedAt,
					TicketName:    item.Tickets[0].TicketDetail.Title,
					Distance:      item.Tickets[0].TicketDetail.Distance,
					ShirtSize:     item.Tickets[0].Type,
					Price:         item.TotalPrice,
					DistanceTotal: distanceTotal,
					TicketOptions: item.TicketOptions,
				}
			} else {
				dataRegisterNew = model.DataRegister{
					RegisterID:    item.ID,
					Firstname:     user.FirstName,
					Lastname:      user.LastName,
					Email:         user.Email,
					Address:       item.ShipingAddress,
					PaymentType:   item.PaymentType,
					PaymentStatus: item.Status,
					RegDate:       item.CreatedAt,
					PaymentDate:   item.UpdatedAt,
					TicketName:    "",
					Distance:      0,
					ShirtSize:     "",
					Price:         item.TotalPrice,
					DistanceTotal: distanceTotal,
					TicketOptions: item.TicketOptions,
				}
			}

		} else {

			dataRegisterNew = model.DataRegister{
				RegisterID:    item.ID,
				Firstname:     user.FirstName,
				Lastname:      user.LastName,
				Email:         user.Email,
				Address:       item.ShipingAddress,
				PaymentType:   item.PaymentType,
				PaymentStatus: item.Status,
				RegDate:       item.CreatedAt,
				PaymentDate:   item.UpdatedAt,
				TicketName:    item.Tickets[0].TicketDetail.Title,
				Distance:      item.Tickets[0].TicketDetail.Distance,
				ShirtSize:     item.Tickets[0].Type,
				Price:         item.TotalPrice,
				DistanceTotal: distanceTotal,
				TicketOptions: item.TicketOptions,
			}
		}

		dataRegister = append(dataRegister, dataRegisterNew)
	}

	report.Datas = dataRegister
	return report, nil
}

func (registerMongo RegisterRepositoryMongo) FindPersonRegEvent(formRequest model.DataRegisterRequest) (model.ReportRegister, error) {
	var register []model.Register
	var report model.ReportRegister

	objectEventID, err := primitive.ObjectIDFromHex(formRequest.EventID)
	if err != nil {
		log.Fatal(err)

	}

	filterSuccess := bson.D{{"status", config.PAYMENT_SUCCESS}, {"event_id", objectEventID}}
	filterWaiting := bson.D{{"status", config.PAYMENT_WAITING}, {"event_id", objectEventID}}
	filterWaitingApprove := bson.D{{"status", config.PAYMENT_WAITING_APPROVE}, {"event_id", objectEventID}}

	countAll, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), bson.D{{"event_id", objectEventID}})

	countSuccess, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filterSuccess)

	countWaiting, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filterWaiting)

	countWaitingApprove, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filterWaitingApprove)

	if err != nil {
		log.Fatal(err)

	}

	report.PaymentAll = countAll
	report.PaymentSuccess = countSuccess
	report.PaymentWaiting = countWaiting
	report.PaymentWaitingApprove = countWaitingApprove

	skip := (formRequest.PageNumber - 1) * formRequest.NPerPage
	limit := formRequest.NPerPage

	options := options.Find()
	options.SetSkip(skip)
	options.SetLimit(limit)

	filter := bson.M{"event_id": objectEventID}

	cur, err := registerMongo.ConnectionDB.Collection(registerCollection).Find(context.TODO(), filter, options)
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {
		var u model.Register
		// decode the document
		if err := cur.Decode(&u); err != nil {
			log.Println(err)
		}
		//fmt.Printf("post: %+v\n", p)
		register = append(register, u)
	}

	var dataRegister []model.DataRegister

	for _, item := range register {

		var user model.User
		userID := item.UserID
		// log.Printf("[info] userID %s", userID)
		// filterUser := bson.D{{"_id", userID}}
		filterUser := bson.M{
			"_id": userID,
			"$or": []interface{}{
				bson.M{"firstname": primitive.Regex{Pattern: formRequest.KeyWord, Options: ""}},
				bson.M{"lastname": primitive.Regex{Pattern: formRequest.KeyWord, Options: ""}},
				bson.M{"employeeid": primitive.Regex{Pattern: formRequest.KeyWord, Options: ""}},
			},
		}
		err := registerMongo.ConnectionDB.Collection("user").FindOne(context.TODO(), filterUser).Decode(&user)

		if err == nil {
			eventUser := item.EventID.Hex() + "." + userID.Hex()
			var activity model.Activity
			filter := bson.D{{"event_user", eventUser}}
			// log.Println(eventUser)
			distanceTotal := 0.0
			if item.Event.Category.Name != "Run" {
				err = registerMongo.ConnectionDB.Collection(activityCollection).FindOne(context.TODO(), filter).Decode(&activity)
				// log.Println(activity)
				if err == nil {
					distanceTotal = math.Round(activity.ToTalDistance*100) / 100
				}
			}
			var dataRegisterNew model.DataRegister
			if item.Event.Category.Name != "Run" {
				if len(item.Tickets) > 0 {
					dataRegisterNew = model.DataRegister{
						RegisterID:    item.ID,
						Firstname:     user.FirstName,
						Lastname:      user.LastName,
						Email:         user.Email,
						Address:       item.ShipingAddress,
						PaymentType:   item.PaymentType,
						PaymentStatus: item.Status,
						PaymentDate:   item.UpdatedAt,
						TicketName:    item.Tickets[0].TicketDetail.Title,
						Distance:      item.Tickets[0].TicketDetail.Distance,
						ShirtSize:     item.Tickets[0].Type,
						Price:         item.TotalPrice,
						TicketOptions: item.TicketOptions,
						DistanceTotal: distanceTotal,
					}
				} else {
					dataRegisterNew = model.DataRegister{
						RegisterID:    item.ID,
						Firstname:     user.FirstName,
						Lastname:      user.LastName,
						Email:         user.Email,
						Address:       item.ShipingAddress,
						PaymentType:   item.PaymentType,
						PaymentStatus: item.Status,
						PaymentDate:   item.UpdatedAt,
						TicketName:    "",
						Distance:      0,
						ShirtSize:     "",
						Price:         item.TotalPrice,
						DistanceTotal: distanceTotal,
					}
				}

			} else {

				dataRegisterNew = model.DataRegister{
					RegisterID:    item.ID,
					Firstname:     user.FirstName,
					Lastname:      user.LastName,
					Email:         user.Email,
					Address:       item.ShipingAddress,
					PaymentType:   item.PaymentType,
					PaymentStatus: item.Status,
					PaymentDate:   item.UpdatedAt,
					TicketName:    item.Tickets[0].TicketDetail.Title,
					Distance:      item.Tickets[0].TicketDetail.Distance,
					ShirtSize:     item.Tickets[0].Type,
					Price:         item.TotalPrice,
					DistanceTotal: distanceTotal,
					TicketOptions: item.TicketOptions,
				}
			}

			dataRegister = append(dataRegister, dataRegisterNew)
		}
	}

	report.Datas = dataRegister
	return report, nil
}

func (registerMongo RegisterRepositoryMongo) UpdateStatusRegister(registerID string, status string, userID string) error {

	objectID, err := primitive.ObjectIDFromHex(registerID)

	filter := bson.D{{"_id", objectID}}

	updated := bson.M{"$set": bson.M{"status": status}}

	res, err := registerMongo.ConnectionDB.Collection(registerCollection).UpdateOne(context.TODO(), filter, updated)
	if err != nil {
		//log.Fatal(res)
		log.Printf("[info] err %s", res)
		return err
	}

	objectUserID, err := primitive.ObjectIDFromHex(userID)

	logInsert := model.LogUpdateRegisterStatus{
		LogName:   "Update register status",
		Status:    status,
		RegID:     objectID,
		UpdateBy:  objectUserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err2 := registerMongo.ConnectionDB.Collection("log_message").InsertOne(context.TODO(), logInsert)
	if err2 != nil {
		log.Printf("[info] err %s", err2)
		return err2
	}

	var register model.Register

	err = registerMongo.ConnectionDB.Collection(registerCollection).FindOne(context.TODO(), filter).Decode(&register)

	if err == nil {
		if register.Event.Category.Name == "Run" {
			go registerMongo.SendRaceMailRegister(registerID)
		} else {
			go registerMongo.SendMailRegister(registerID)
		}
	}

	return nil

}
