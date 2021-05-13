package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/omise/omise-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"thinkdev.app/think/runex/runexapi/config"
	"thinkdev.app/think/runex/runexapi/config/db"
	"thinkdev.app/think/runex/runexapi/firebase"
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/request"
	"thinkdev.app/think/runex/runexapi/utils"

	"thinkdev.app/think/runex/runexapi/api/mail"
)

// RegisterRepository interface
type RegisterRepository interface {
	GetRegisterAll() ([]model.RegisterV2, error)
	AddRegister(register model.RegisterRequest) (model.RegisterV2, error)
	AddRaceRegister(register model.RegisterRequest) (model.RegisterV2, error)
	AddMerChant(userID string, eventID string, regID string, charge omise.Charge, orderID string) error
	EditRegister(registerID string, register model.RegisterRequest) error
	GetRegisterByEvent(eventID string) ([]model.ReportRegisterV2, error)
	NotifySlipRegister(registerID string, slip model.SlipTransfer) error
	AdminNotifySlipRegister(registerID string, slip model.SlipTransfer) error
	CountByEvent(eventID string) (int64, error)
	GetRegEventByID(regID string) (model.RegisterV2, error)
	GetRegEventByIDNew(regID string, eventID string) (model.RegisterV2, error)
	CheckUserRegisterEvent(eventID string, userID string) (bool, error)
	CheckUserRegisterEventCode(code string, userID string) (bool, error)
	GetRegisterByUserID(userID string) ([]model.RegisterV2, error)
	GetRegisterByUserAndEvent(userID string, id string) (model.RegisterV2, error)
	GetRegisterByUserEventAndRegID(userID string, eventID string, regID string) (model.RegisterV2, error)
	SendMailRegister(registerID string) error
	SendMailRegisterNew(registerID string, eventID string) error
	GetRegisterActivateEvent(userID string) ([]model.RegisterV2, error)
	GetRegisterReport(formRequest model.DataRegisterRequest) (model.ReportRegister, error)
	GetRegisterReportAll(formRequest model.DataRegisterRequest) (model.ReportRegister, error)
	FindPersonRegEvent(formRequest model.DataRegisterRequest) (model.ReportRegister, error)
	UpdateStatusRegister(registerID string, status string, userID string) error
}

// RegisterRepositoryMongo mongo ref
type RegisterRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

type Reg struct {
	OwnerID   string      `json:"owner_id" bson:"owner_id"`
	UserCode  string      `json:"user_code" bson:"user_code"`
	EventCode string      `json:"event_code" bson:"event_code"`
	Ref2      string      `json:"ref2" bson:"ref2"`
	Regs      model.Regs  `json:"regs" bson:"regs"`
	Event     model.Event `json:"event" bson:"event"`
}

const (
	registerCollection    = "register_v2"
	merchantCollection    = "merchant"
	activityCollection    = "activityV2"
	registerAllCollection = "register_all"
	scbPaymentCollection  = "scbpayment"
)

type CountBiB struct {
	Count int `json:"count" bson:"count"`
}

// GetRegisterAll repo
func (registerMongo RegisterRepositoryMongo) GetRegisterAll() ([]model.RegisterV2, error) {
	var register []model.RegisterV2
	cur, err := registerMongo.ConnectionDB.Collection(registerCollection).Find(context.TODO(), bson.D{{}})
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {
		var u model.RegisterV2
		// decode the document
		if err := cur.Decode(&u); err != nil {
			log.Println(err)
		}
		//fmt.Printf("post: %+v\n", p)
		register = append(register, u)
	}

	return register, err
}

//AddRegister repo
func (registerMongo RegisterRepositoryMongo) AddRegister(register model.RegisterRequest) (model.RegisterV2, error) {
	filter := bson.M{"event_code": register.EventCode}
	dataInfo := register.Regs
	count, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filter)
	//log.Printf("[info] count %s", count)
	if err != nil {
		log.Println(err)
		return model.RegisterV2{}, err
	}

	if dataInfo.Partner.PartnerName == config.PartnerKao {
		// if register.KoaRequest.Slug == "" || register.KoaRequest.EBIB == "" {
		// 	return regModel, err
		// }
		urlS := fmt.Sprintf("https://kaokonlakao.com/api/%s/bib/%s", register.KoaRequest.Slug, register.KoaRequest.EBIB)
		var bearer = "Bearer olcgZVpqDXQikRDG"
		response, err := request.GetRequest(urlS, bearer, nil)
		if err != nil {
			log.Println(err)
			return model.RegisterV2{}, fmt.Errorf("ไม่พบ E-BIB นี้ในรายการ ก้าวเพื่อน้อง")
		}
		var koaObject model.KaoObject
		err = json.Unmarshal(response, &koaObject)
		if err != nil {
			log.Println(err)
			return model.RegisterV2{}, fmt.Errorf("ไม่พบ E-BIB นี้ในรายการ ก้าวเพื่อน้อง")
		}

		dataInfo.Partner.RefActivityValue = register.KoaRequest.EBIB
		dataInfo.Partner.RefEventValue = strconv.Itoa(koaObject.VirtualRaceProfile.OrderItemID)
		dataInfo.Partner.RefPhoneValue = koaObject.HolderPhone
	}

	if count > 0 {
		filter := bson.D{primitive.E{Key: "event_code", Value: register.EventCode}, primitive.E{Key: "regs.user_id", Value: dataInfo.UserID}}
		var regModel model.RegisterV2
		options := options.FindOne()
		options.SetProjection(bson.M{"regs.$": 1})
		err := registerMongo.ConnectionDB.Collection(registerCollection).FindOne(context.TODO(), filter, options).Decode(&regModel)
		if err == nil {
			return regModel, err
		}
		if err != nil {
			log.Println("not found")
			if err.Error() != "mongo: no documents in result" {
				return regModel, err
			}
		}
		dataInfo.ID = primitive.NewObjectID()
		dataInfo.CreatedAt = time.Now()
		dataInfo.UpdatedAt = time.Now()
		dataInfo.OrderID = utils.OrderIDGenerate()
		dataInfo.EventCode = register.EventCode

		matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"event_code": register.EventCode}}}
		unwindStage := bson.D{primitive.E{Key: "$unwind", Value: "$regs"}}
		countStage := bson.D{bson.E{Key:"$count",Value: "count"}}
		cur, err := db.DB.Collection(registerCollection).Aggregate(context.TODO(), mongo.Pipeline{matchStage, unwindStage, countStage})
		if err != nil {
			log.Println(err)
			return regModel,err
		}
		var c CountBiB
		for cur.Next(context.TODO()) {
			// decode the document
			if err := cur.Decode(&c); err != nil {
				return regModel, err
			}
			//fmt.Printf("post: %+v\n", p)
		}
	
		if c.Count > 0 {
			count = int64(c.Count)
		}

		for i := range dataInfo.TicketOptions {
			dataInfo.TicketOptions[i].RegisterNumber = fmt.Sprintf("%05d", int64(count+1))
			dataInfo.TicketOptions[i].UserOption.CreatedAt = time.Now()
			dataInfo.TicketOptions[i].UserOption.UpdatedAt = time.Now()
			if dataInfo.TicketOptions[i].Tickets.Category == "team" {
				dataInfo.IsTeamLead = true
				dataInfo.ParentRegID = dataInfo.ID
			}
		}
		update := bson.M{"$push": bson.M{"regs": dataInfo}}
		filter = bson.D{primitive.E{Key: "event_code", Value: register.EventCode}}
		_, err = registerMongo.ConnectionDB.Collection(registerCollection).UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Println(err.Error())
			return regModel, err
		}
		SaveRegister(dataInfo, register.EventCode)

	} else {
		mod := mongo.IndexModel{
			Keys: bson.M{
				"event_code": 1, // index in ascending order
				"ref2":       1,
			}, Options: nil,
		}

		// Create an Index using the CreateOne() method
		_, err := db.DB.Collection(registerCollection).Indexes().CreateOne(context.Background(), mod)
		if err != nil {
			log.Println(err)
		}
		event, err := DetailEventOwnerByCode(register.EventCode)
		var arrRegs []model.Regs
		dataInfo.ID = primitive.NewObjectID()
		dataInfo.CreatedAt = time.Now()
		dataInfo.UpdatedAt = time.Now()
		dataInfo.OrderID = utils.OrderIDGenerate()
		dataInfo.EventCode = register.EventCode
		for i := range dataInfo.TicketOptions {
			dataInfo.TicketOptions[i].RegisterNumber = fmt.Sprintf("%05d", int64(count+1))
			dataInfo.TicketOptions[i].UserOption.CreatedAt = time.Now()
			dataInfo.TicketOptions[i].UserOption.UpdatedAt = time.Now()

			if dataInfo.TicketOptions[i].Tickets.Category == "team" {
				dataInfo.IsTeamLead = true
				dataInfo.ParentRegID = dataInfo.ID
			}
		}
		arrRegs = append(arrRegs, dataInfo)

		registerModel := model.RegisterV2{
			EventCode: register.EventCode,
			Regs:      arrRegs,
			OwnerID:   event.UserID,
			UserCode:  event.UserID,
			Ref2:      utils.Ref2Generate(),
			Event:     event,
		}

		_, err = registerMongo.ConnectionDB.Collection(registerCollection).InsertOne(context.TODO(), registerModel)
		if err != nil {
			//log.Fatal(res)
			log.Println(err)
			return registerModel, err
		}
		SaveRegister(dataInfo, register.EventCode)
	}
	go registerMongo.SendMailRegisterNew(dataInfo.ID.Hex(), register.EventCode)
	type refVal struct {
		Ref2        string             `json:"ref2" bson:"ref2"`
	}
	var ref2 refVal
	err = registerMongo.ConnectionDB.Collection(registerCollection).FindOne(context.TODO(),bson.M{"event_code": register.EventCode}).Decode(&ref2)
	if err != nil {
		ref2.Ref2 = strings.ToUpper(register.EventCode)
	}

	var regModel = model.RegisterV2{
		Regs:      []model.Regs{dataInfo},
		EventCode: register.EventCode,
		Ref2: ref2.Ref2,
	}

	return regModel, err
}

func CheckTeamLead(register model.AddTeamRequest, userID primitive.ObjectID) (bool, error) {
	matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"event_code": register.EventCode}}}
	projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs._id", register.ParentRegID}}}}}}}
	cur, err := db.DB.Collection(registerCollection).Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage})
	if err != nil {
		return false, err
	}
	var reg model.RegisterV2

	for cur.Next(context.TODO()) {
		var u model.RegisterV2
		// decode the document
		if err := cur.Decode(&u); err != nil {
			return false, err
		}
		//fmt.Printf("post: %+v\n", p)
		reg = u
	}
	if reg.Regs[0].UserID == userID {
		return true, err
	}
	return false, err
}

func CheckSameUser(register model.AddTeamRequest) (bool, error) {
	matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"event_code": register.EventCode}}}
	projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$and": []interface{}{bson.M{"$eq": bson.A{"$$regs.user_id", register.TeamUserID}}, bson.M{"$eq": bson.A{"$$regs.parent_reg_id", register.ParentRegID}}}}}}}}}
	projectStage2 := bson.D{bson.E{Key: "$project", Value: bson.M{"count": bson.M{"$size": "$regs"}}}}
	cur, err := db.DB.Collection(registerCollection).Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage, projectStage2})
	if err != nil {
		log.Println(err)
		return false, err
	}
	var c CountBiB
	for cur.Next(context.TODO()) {
		// decode the document
		if err := cur.Decode(&c); err != nil {
			log.Println(err)
			return false, err
		}
		//fmt.Printf("post: %+v\n", p)
	}

	if c.Count > 0 {
		return true, err
	}
	return false, err
}

func CheckRunerInTeam(register model.AddTeamRequest) (int, bool, error) {
	log.Println("check runer")
	matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"event_code": register.EventCode}}}
	projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs.parent_reg_id", register.ParentRegID}}}}}}}
	//projectStage := bson.D{bson.E{Key: "$project", Value:bson.M{"regs": bson.M{"count": bson.M{"$size": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs.parent_reg_id", register.ParentRegID}}}}}}}}
	cur, err := db.DB.Collection(registerCollection).Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage})
	if err != nil {
		log.Println(err)
		return 0, false, err
	}
	// var cBib = CountBiB{
	// 	Count: 0,
	// }
	var regs = []model.RegisterV2{}
	for cur.Next(context.TODO()) {
		var u model.RegisterV2
		// decode the document
		if err := cur.Decode(&u); err != nil {
			return 0, false, err
		}
		//fmt.Printf("post: %+v\n", p)
		regs = append(regs, u)
	}
	if len(register.Regs.TicketOptions) > 0 {
		log.Println(register.Regs.TicketOptions[0].Tickets)
		i, err := strconv.Atoi(register.Regs.TicketOptions[0].Tickets.RunnerInTeam)
		if err != nil {
			return len(regs), false, err
		}
		if len(regs) < i {
			return len(regs), true, err
		}
		return len(regs), false, err
	}

	return len(regs), false, err
}

//GetRegEventDashboard
func GetRegEventDashboard(register model.RegEventDashboardRequest) (model.RegisterV2, error) {
	matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"event_code": register.EventCode}}}
	projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs.parent_reg_id", register.ParentRegID}}}}, "event_code": 1, "ref2": 1, "event": 1}}}
	if register.ParentRegID.IsZero() {
		log.Println("zero")
		projectStage = bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs._id", register.RegID}}}}, "event_code": 1, "ref2": 1, "event": 1}}}
	}
	cur, err := db.DB.Collection(registerCollection).Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage})
	if err != nil {
		return model.RegisterV2{}, err
	}
	var r model.RegisterV2
	for cur.Next(context.TODO()) {
		// decode the document
		if err := cur.Decode(&r); err != nil {
			return r, err
		}
		//fmt.Printf("post: %+v\n", p)
	}
	return r, err
}

//AddRegister repo
func AddTeamRegister(register model.AddTeamRequest) (model.RegisterV2, error) {
	dataInfo := register.Regs

	matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"event_code": register.EventCode}}}
	projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs.event_code", register.EventCode}}}}}}}
	projectStage2 := bson.D{bson.E{Key: "$project", Value: bson.M{"count": bson.M{"$size": "$regs"}}}}
	cur, err := db.DB.Collection(registerCollection).Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage, projectStage2})
	if err != nil {
		log.Println(err)
		return model.RegisterV2{}, err
	}
	var c CountBiB
	for cur.Next(context.TODO()) {
		// decode the document
		if err := cur.Decode(&c); err != nil {
			log.Println(err)
			return model.RegisterV2{}, err
		}
		//fmt.Printf("post: %+v\n", p)
	}

	//projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"count": bson.M{"$size": bson.M{"input": "$regs", "as": "regs"}}}}}

	// cur, err := db.DB.Collection(registerCollection).Aggregate(context.TODO(), mongo.Pipeline{matchStage})
	// if err != nil {
	// 	return regs, err
	// }
	// var cBib = CountBiB{
	// 	Count: 0,
	// }

	// for cur.Next(context.TODO()) {
	// 	var u model.RegisterV2
	// 	// decode the document
	// 	if err := cur.Decode(&u); err != nil {
	// 		return regs, err
	// 	}
	// }

	dataInfo.ID = primitive.NewObjectID()
	dataInfo.CreatedAt = time.Now()
	dataInfo.UpdatedAt = time.Now()
	dataInfo.OrderID = utils.OrderIDGenerate()
	dataInfo.EventCode = register.EventCode
	dataInfo.Status = config.REGISTER_WAIT_CONFIRM
	dataInfo.ParentRegID = register.ParentRegID

	log.Println(c.Count)

	for i := range dataInfo.TicketOptions {
		dataInfo.TicketOptions[i].RegisterNumber = fmt.Sprintf("%05d", c.Count+1)
		dataInfo.TicketOptions[i].UserOption.CreatedAt = time.Now()
		dataInfo.TicketOptions[i].UserOption.UpdatedAt = time.Now()
	}
	update := bson.M{"$push": bson.M{"regs": dataInfo}}
	filter := bson.D{primitive.E{Key: "event_code", Value: register.EventCode}}
	_, err = db.DB.Collection(registerCollection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println(err.Error())
		return model.RegisterV2{}, err
	}
	SaveRegister(dataInfo, register.EventCode)

	var regModel = model.RegisterV2{
		Regs:      []model.Regs{dataInfo},
		EventCode: register.EventCode,
	}

	return regModel, err
}

// AddRaceRegister repo for race running
func (registerMongo RegisterRepositoryMongo) AddRaceRegister(register model.RegisterRequest) (model.RegisterV2, error) {
	// register.CreatedAt = time.Now()
	// register.UpdatedAt = time.Now()
	filter := bson.M{"event_code": register.EventCode}
	dataInfo := register.Regs
	count, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filter)
	//log.Printf("[info] count %s", count)
	if err != nil {
		log.Println(err)
	}

	filterBib := bson.D{bson.E{Key: "event_code", Value: register.EventCode}}
	var ebibEvent model.EbibEvent
	err = registerMongo.ConnectionDB.Collection(ebibCollection).FindOne(context.TODO(), filterBib).Decode(&ebibEvent)
	if err != nil {
		log.Println("find ebib in AddRaceRegister", err)
	}
	//dataInfo.TicketOptions.RegisterNumber = fmt.Sprintf("%05d", (ebibEvent.LastNo + 1))

	for index := range dataInfo.TicketOptions {
		log.Println(fmt.Sprintf("%05d", (ebibEvent.LastNo + int64(index+1))))
		// element.RegisterNumber = fmt.Sprintf("%05d", (ebibEvent.LastNo + int64(index + 1)))
		dataInfo.TicketOptions[index].RegisterNumber = fmt.Sprintf("%05d", (ebibEvent.LastNo + int64(index+1)))
	}

	if count > 0 {
		dataInfo.ID = primitive.NewObjectID()
		dataInfo.CreatedAt = time.Now()
		dataInfo.UpdatedAt = time.Now()
		update := bson.M{"$push": bson.M{"regs": dataInfo}}
		_, err = registerMongo.ConnectionDB.Collection(registerCollection).UpdateOne(context.TODO(), filter, update)
		if err != nil {
			//log.Fatal(res)
			log.Println(err)
		}
	} else {
		var arrRegs []model.Regs
		dataInfo.ID = primitive.NewObjectID()
		dataInfo.CreatedAt = time.Now()
		dataInfo.UpdatedAt = time.Now()
		arrRegs = append(arrRegs, dataInfo)
		registerModel := model.RegisterV2{
			EventCode: register.EventCode,
			Regs:      arrRegs,
		}

		_, err := registerMongo.ConnectionDB.Collection(registerCollection).InsertOne(context.TODO(), registerModel)
		if err != nil {
			//log.Fatal(res)
			log.Println(err)
		}
	}

	// res, err := registerMongo.ConnectionDB.Collection(registerCollection).InsertOne(context.TODO(), register)
	// if err != nil {
	// 	log.Fatal(res)
	// }
	// fmt.Println("Inserted a single document: ", res.InsertedID)
	regModel, err2 := registerMongo.GetRegEventByIDNew(dataInfo.ID.Hex(), register.EventCode)
	if err2 != nil {
		return regModel, err
	}
	//ebibEvent.LastNo += int64(len(regModel.TicketOptions))
	ebibEvent.LastNo++
	ebibEvent.UpdatedAt = time.Now()
	updated := bson.M{"$set": ebibEvent}
	_, err = registerMongo.ConnectionDB.Collection(ebibCollection).UpdateOne(context.TODO(), filterBib, updated)
	if err != nil {
		log.Print(err)
	}

	err3 := registerMongo.SendMailRegisterNew(dataInfo.ID.Hex(), register.EventCode)
	if err3 != nil {
		log.Println(err3)
	}
	// if err3 != nil {
	// 	log.Fatal(err3)
	// }
	return regModel, err
}

// AddMerChant repo for payment success
func (registerMongo RegisterRepositoryMongo) AddMerChant(userID string, eventCode string, regID string, charge omise.Charge, orderID string) error {
	uid, err := primitive.ObjectIDFromHex(userID)
	rid, err := primitive.ObjectIDFromHex(regID)
	typePay := ""
	if charge.Source != nil {
		if len(charge.Source.Type) > 0 {
			typePay = charge.Source.Type
		}
	}

	var merchant model.Merchant
	merchant.UserID = uid
	merchant.EventCode = eventCode
	merchant.RegID = rid
	merchant.PaymentType = typePay
	merchant.Status = string(charge.Status)
	merchant.OrderID = orderID
	merchant.OmiseID = charge.Transaction
	merchant.TotalPrice = charge.Amount / 100
	merchant.CreatedAt = time.Now()
	merchant.UpdatedAt = time.Now()
	//fmt.Println("Inserted a single document: ", merchant)
	res, err := registerMongo.ConnectionDB.Collection(merchantCollection).InsertOne(context.TODO(), merchant)
	if err != nil {
		log.Println(res)
	}
	if charge.Status == "failed" {
		err = fmt.Errorf("Payment failed")
		return err
	} else {
		if registerMongo.UpdatePaymentStatus(rid, eventCode, uid, config.PAYMENT_SUCCESS, config.PAYMENT_CREDIT_CARD) {
			
		}
	// var mailTemplate model.EmailTemplateData2
	// mailTemplate.RefID = orderID
	// mailTemplate.IdentificationNumber = user.CitycenID
	// mailTemplate.ContactPhone = user.Phone
	// mailTemplate.CompetitionType = register.TicketOptions[0].Tickets.Title
	// mailTemplate.RegisterNumber = register.TicketOptions[0].RegisterNumber
	// mailTemplate.TicketName = register.TicketOptions[0].Tickets.Title
	// mailTemplate.Status = register.Status
	// mailTemplate.PaymentType = register.PaymentType
	// mailTemplate.ShipingAddress = register.TicketOptions[0].UserOption.Address

	// fmt.Println("Mail Object :  ", mailTemplate)

	// if register.PaymentType == "PAYMENT_FREE" {
	// 	mail.SendRegFreeEventMail2(mailTemplate)
	// } else {
	// 	mail.SendRegEventMail2(mailTemplate)
	// }
		return nil
	}
}

// UpdatePaymentStatus repo
func (registerMongo RegisterRepositoryMongo) UpdatePaymentStatus(registerID primitive.ObjectID, eventCode string, userID primitive.ObjectID, status string, payType string) bool {
	filter := bson.M{"$and": []interface{}{bson.M{"event_code": eventCode}, bson.M{"regs.user_id": userID}, bson.M{"regs._id": registerID}}}
	update := bson.M{"$set": bson.M{"regs.$.status": status, "regs.$.payment_type": payType, "regs.$.payment_date": time.Now(), "regs.$.updated_at": time.Now()}}
	_, err := registerMongo.ConnectionDB.Collection(registerCollection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println("updating the Data", err)
		return false
	}
	return true
}

// SlipUpdatePaymentStatus repo
func SlipUpdatePaymentStatus(registerID primitive.ObjectID, eventCode string, userID primitive.ObjectID, status string, payType string, imagePath string) bool {
	filter := bson.M{"$and": []interface{}{bson.M{"event_code": eventCode}, bson.M{"regs.user_id": userID}, bson.M{"regs._id": registerID}}}
	update := bson.M{"$set": bson.M{"regs.$.status": status, "regs.$.payment_type": payType, "regs.$.payment_date": time.Now(), "regs.$.slip": imagePath, "regs.$.updated_at": time.Now()}}
	_, err := db.DB.Collection(registerCollection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println("updating the Data", err)
		return false
	}
	return true
}

// EditRegister repo
func (registerMongo RegisterRepositoryMongo) EditRegister(registerID string, register model.RegisterRequest) error {
	regObjectID, err := primitive.ObjectIDFromHex(registerID)
	eventObjectID := register.EventCode
	filterUpdate := bson.D{bson.E{Key: "event_code", Value: eventObjectID}, bson.E{Key: "regs._id", Value: regObjectID}}
	//register.UpdatedAt = time.Now()
	dataInfo := register.Regs
	dataInfo.UpdatedAt = time.Now()
	dataInfo.ID = regObjectID
	//updated := bson.M{"$set": register}
	updated := bson.M{"$set": bson.M{"regs.$": dataInfo}}
	_, err = registerMongo.ConnectionDB.Collection(registerCollection).UpdateOne(context.TODO(), filterUpdate, updated)
	if err != nil {
		//log.Fatal(res)
		//log.Printf("[info] err %s", res)
		return err

	}

	// if register.Event.Category.Name == "Run" {
	// 	go registerMongo.SendRaceMailRegister(registerID)
	// } else {
	// 	go registerMongo.SendMailRegister(registerID)
	// }

	return nil
}

// GetRegisterByEvent repo list data register event
func (registerMongo RegisterRepositoryMongo) GetRegisterByEvent(eventID string) ([]model.ReportRegisterV2, error) {

	var register = []model.ReportRegisterV2{}
	filter := bson.D{bson.E{Key: "event_code", Value: eventID}}
	cur, err := registerMongo.ConnectionDB.Collection(registerCollection).Find(context.TODO(), filter)
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {
		var u model.ReportRegisterV2
		// decode the document
		if err := cur.Decode(&u); err != nil {
			//log.Fatal(err)
		}
		//fmt.Printf("post: %+v\n", p)
		register = append(register, u)
	}

	return register, err
}

//GetRegisterByUserID repository get my event register
func (registerMongo RegisterRepositoryMongo) GetRegisterByUserID(userID string) ([]model.RegisterV2, error) {

	var register []model.RegisterV2
	objectID, err := primitive.ObjectIDFromHex(userID)
	//filter := bson.D{{"regs.user_id", objectID}}
	//unwindStage := bson.D{{"$unwind", "$regs"}}
	matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"regs.user_id": objectID}}}

	//matchSubStage := bson.D{{"$match", bson.M{"regs.user_id": bson.M{"$eq": objectID}}}}
	//groupStage := bson.D{{"_id", "$_id"}, {"event_id", "$event_id"}, {"regs", bson.M{"$push": "$regs"}}}
	//filterStage := bson.D{{"$project", bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs.user_id", objectID}}}}}}}
	projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs.user_id", objectID}}}}, "event_code": 1, "ref2": 1, "event": 1}}}
	unwindStage := bson.D{primitive.E{Key: "$unwind", Value: "$regs"}}
	cur, err := registerMongo.ConnectionDB.Collection(registerCollection).Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage, unwindStage})
	//cur, err := registerMongo.ConnectionDB.Collection(registerCollection).Find(context.TODO(), filter)
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}

	//u.Regs = []model.Regs{}
	for cur.Next(context.TODO()) {
		var u Reg
		var regs = []model.Regs{}
		// decode the document
		if err := cur.Decode(&u); err != nil {
			log.Print(err)
		}
		event, err := DetailEventOwnerByCode(u.EventCode)
		if err == nil {
			u.Event = event
		}
		r := model.RegisterV2{
			UserCode:  u.UserCode,
			OwnerID:   u.OwnerID,
			EventCode: u.EventCode,
			Event:     event,
			Ref2:      u.Ref2,
			Regs:      append(regs, u.Regs),
		}
		register = append(register, r)
	}

	return register, err
}

//GetRegisterByUserEventAndRegID repository get my event register with user id, event id, and reg id
func (registerMongo RegisterRepositoryMongo) GetRegisterByUserEventAndRegID(userID string, eventID string, regID string) (model.RegisterV2, error) {

	var register model.RegisterV2
	objectID, err := primitive.ObjectIDFromHex(userID)
	eID, err := primitive.ObjectIDFromHex(userID)
	rID, err := primitive.ObjectIDFromHex(userID)
	//filter := bson.D{{"regs.user_id", objectID}}
	//unwindStage := bson.D{{"$unwind", "$regs"}}
	matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"regs.user_id": objectID, "regs._id": rID, "event_id": eID}}}
	//unwindStage := bson.D{{"$unwind", "$regs"}}
	//matchSubStage := bson.D{{"$match", bson.M{"regs.user_id": bson.M{"$eq": objectID}}}}
	//groupStage := bson.D{{"_id", "$_id"}, {"event_id", "$event_id"}, {"regs", bson.M{"$push": "$regs"}}}
	//filterStage := bson.D{{"$project", bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs.user_id", objectID}}}}}}}
	projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs.user_id", objectID}}}}, "event_id": eventID}}}
	cur, err := registerMongo.ConnectionDB.Collection(registerCollection).Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage})
	//cur, err := registerMongo.ConnectionDB.Collection(registerCollection).Find(context.TODO(), filter)
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {
		var u model.RegisterV2

		// decode the document
		if err := cur.Decode(&u); err != nil {
			log.Print(err)
		}

		// var event model.EventRegV2
		// filter := bson.D{primitive.E{Key: "_id", Value: u.EventID}}
		// err := registerMongo.ConnectionDB.Collection(eventCollection).FindOne(context.TODO(), filter).Decode(&event)
		// event, err := DetailEventByCode(u.EventCode)
		// if err != nil {
		// 	log.Print(err)
		// }
		//u.Regs[0].Event = event
		// var event model.Event
		// registerMongo.ConnectionDB.Collection(eventCollection).FindOne(context.TODO(), bson.D{{"_id", u.EventID}}).Decode(&event)
		// //fmt.Printf("post: %+v\n", p)
		//u.Event = event
		register = u
	}

	return register, err
}

//GetRegisterByUserAndEvent detail reg
func (registerMongo RegisterRepositoryMongo) GetRegisterByUserAndEvent(userID string, id string) (model.RegisterV2, error) {

	var register model.RegisterV2
	objectID, err := primitive.ObjectIDFromHex(userID)
	regID, err := primitive.ObjectIDFromHex(id)
	filter := bson.D{primitive.E{Key: "user_id", Value: objectID}, primitive.E{Key: "_id", Value: regID}}
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

//NotifySlipRegister repo for send mail on user send slip
func (registerMongo RegisterRepositoryMongo) NotifySlipRegister(registerID string, slip model.SlipTransfer) error {
	objectID, err := primitive.ObjectIDFromHex(registerID)
	filter := bson.D{bson.E{Key: "_id", Value: objectID}}

	updated := bson.M{"$set": bson.M{"slip": slip, "status": config.PAYMENT_WAITING_APPROVE, "updated_at": time.Now()}}

	_, err = registerMongo.ConnectionDB.Collection(registerCollection).UpdateOne(context.TODO(), filter, updated)
	if err != nil {
		//log.Fatal(res)
		log.Printf("[info] err %s", err.Error())
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
	filter := bson.D{bson.E{Key: "_id", Value: objectID}}

	var register model.Register
	err = registerMongo.ConnectionDB.Collection(registerCollection).FindOne(context.TODO(), filter).Decode(&register)
	paymentType := config.PAYMENT_TRANSFER
	if err == nil {
		if register.PaymentType != "" {
			paymentType = register.PaymentType
		}
	}

	updated := bson.M{"$set": bson.M{"slip": slip, "status": config.PAYMENT_WAITING_APPROVE, "payment_type": paymentType, "updated_at": time.Now()}}

	_, err = registerMongo.ConnectionDB.Collection(registerCollection).UpdateOne(context.TODO(), filter, updated)
	if err != nil {
		//log.Printf("[info] err %s", res)
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

//CountByEvent count register by event
func (registerMongo RegisterRepositoryMongo) CountByEvent(eventID string) (int64, error) {

	objectID, err := primitive.ObjectIDFromHex(eventID)
	filter := bson.D{bson.E{Key: "event_id", Value: objectID}}

	count, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filter)
	//log.Printf("[info] count %s", count)
	if err != nil {
		log.Println(err)
	}

	return count, err
}

// GetRegEventByID repo
func (registerMongo RegisterRepositoryMongo) GetRegEventByID(regID string) (model.RegisterV2, error) {

	var regEvent model.RegisterV2
	id, err := primitive.ObjectIDFromHex(regID)
	if err != nil {
		log.Println(err)
	}
	filter := bson.M{"_id": id}
	err2 := registerMongo.ConnectionDB.Collection(registerCollection).FindOne(context.TODO(), filter).Decode(&regEvent)
	if err2 != nil {
		log.Println(err2)
	}

	return regEvent, err2
}

//GetRegEventByIDNew repo for register detail
func (registerMongo RegisterRepositoryMongo) GetRegEventByIDNew(regID string, eventCode string) (model.RegisterV2, error) {

	var regEvent model.RegisterV2
	regObjectID, err3 := primitive.ObjectIDFromHex(regID)
	if err3 != nil {
		log.Println(err3)
	}

	filter := bson.D{bson.E{Key: "event_code", Value: eventCode}, bson.E{Key: "regs._id", Value: regObjectID}}

	options := options.FindOne()
	options.SetProjection(bson.M{"regs.$": 1})
	err2 := registerMongo.ConnectionDB.Collection(registerCollection).FindOne(context.TODO(), filter, options).Decode(&regEvent)
	log.Printf("[info RegEvent] cur %s", err2)
	if err2 != nil {
		log.Fatal(err2)
	}
	regEvent.EventCode = eventCode
	return regEvent, err2
}

//CheckUserRegisterEvent check user register event
func (registerMongo RegisterRepositoryMongo) CheckUserRegisterEvent(eventCode string, userID string) (bool, error) {
	userObjectID, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.D{primitive.E{Key: "event_code", Value: eventCode}, primitive.E{Key: "regs.user_id", Value: userObjectID}}
	count, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filter)
	log.Printf("[info] count %d", count)
	if err != nil {
		log.Println(err)
	}
	if count > 0 {

		return true, nil
	}

	return false, nil
}

//CheckUserRegisterEventCode check user register event
func (registerMongo RegisterRepositoryMongo) CheckUserRegisterEventCode(code string, userID string) (bool, error) {
	userObjectID, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.D{primitive.E{Key: "event_code", Value: code}, primitive.E{Key: "regs.user_id", Value: userObjectID}}
	count, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Println(err)
	}
	if count > 0 {

		return true, nil
	}

	return false, nil
}

// SendMailRegister send email for register
func (registerMongo RegisterRepositoryMongo) SendMailRegister(registerID string) error {
	var register model.Register
	var user model.UserProvider
	var mailTemplate model.EmailTemplateData2
	objectID, _ := primitive.ObjectIDFromHex(registerID)
	filter := bson.D{bson.E{Key: "_id", Value: objectID}}
	err := registerMongo.ConnectionDB.Collection("register").FindOne(context.TODO(), filter).Decode(&register)
	if err != nil {
		log.Fatal(err)
	}

	filterUser := bson.D{bson.E{Key: "_id", Value: register.UserID}}
	err = registerMongo.ConnectionDB.Collection("user").FindOne(context.TODO(), filterUser).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	filterEvent := bson.D{bson.E{Key: "event_id", Value: register.EventID}}
	count, err := registerMongo.ConnectionDB.Collection("register").CountDocuments(context.TODO(), filterEvent)
	registerNumber := fmt.Sprintf("%05d", count)
	address := register.ShipingAddress.Address + " City " + register.ShipingAddress.City + " District " + register.ShipingAddress.District + " Province " + register.ShipingAddress.Province + " " + register.ShipingAddress.Zipcode

	updated := bson.M{"$set": bson.M{"register_number": registerNumber}}
	_, err = registerMongo.ConnectionDB.Collection("register").UpdateOne(context.TODO(), filter, updated)
	if err != nil {
		log.Printf("[info] err %s", err.Error())
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

// SendMailRegisterNew new version
func (registerMongo RegisterRepositoryMongo) SendMailRegisterNew(registerID string, eventCode string) error {
	var registerEvent model.RegisterV2
	var register model.Regs
	var user model.UserProvider
	var mailTemplate model.EmailTemplateData2
	regObjectID, _ := primitive.ObjectIDFromHex(registerID)
	//filter := bson.D{{"event_id", eventObjectID}}
	filter := bson.D{bson.E{Key: "event_code", Value: eventCode}, bson.E{Key: "regs._id", Value: regObjectID}}

	options := options.FindOne()
	options.SetProjection(bson.M{"regs.$": 1})
	err := registerMongo.ConnectionDB.Collection(registerCollection).FindOne(context.TODO(), filter, options).Decode(&registerEvent)
	if err != nil {
		log.Println(err)
	}
	if len(registerEvent.Regs) < 1 {
		return err
	}
	register = registerEvent.Regs[0]
	filterUser := bson.D{bson.E{Key: "_id", Value: register.UserID}}
	err = registerMongo.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filterUser).Decode(&user)
	if err != nil {
		log.Println(err)
	}
	// var registerEventAll model.RegisterV2
	// filterEvent := bson.D{bson.E{Key: "event_code", Value: eventCode}}
	// err = registerMongo.ConnectionDB.Collection(registerCollection).FindOne(context.TODO(), filterEvent).Decode(&registerEventAll)
	// count := len(registerEventAll.Regs)
	// fmt.Println("First Length:", count)
	// registerNumber := fmt.Sprintf("%05d", count)
	// address := ""
	// if register.TicketOptions[0].UserOption.Address != "" {
	// 	address = register.TicketOptions[0].UserOption.Address

	// }

	// filterUpdate := bson.D{bson.E{Key: "event_id", Value: eventObjectID}, bson.E{Key: "regs._id", Value: regObjectID}}
	// updated := bson.M{"$set": bson.M{"regs.$.register_number": registerNumber}}
	// _, err = registerMongo.ConnectionDB.Collection(registerCollection).UpdateOne(context.TODO(), filterUpdate, updated)
	// if err != nil {
	// 	log.Printf("[info] err %s", err.Error())
	// 	return err
	// }

	mailTemplate.Name = user.FirstName + " " + user.LastName
	mailTemplate.RefID = register.OrderID
	mailTemplate.Email = user.Email
	mailTemplate.Phone = user.Phone
	mailTemplate.IdentificationNumber = user.CitycenID
	mailTemplate.ContactPhone = user.Phone
	mailTemplate.CompetitionType = register.TicketOptions[0].Tickets.Title
	mailTemplate.RegisterNumber = register.TicketOptions[0].RegisterNumber
	mailTemplate.TicketName = register.TicketOptions[0].Tickets.Title
	mailTemplate.Status = register.Status
	mailTemplate.PaymentType = register.PaymentType
	mailTemplate.ShipingAddress = register.TicketOptions[0].UserOption.Address

	fmt.Println("Mail Object :  ", mailTemplate)

	if register.PaymentType == "PAYMENT_FREE" {
		mail.SendRegFreeEventMail2(mailTemplate)
	} else {
		mail.SendRegEventMail2(mailTemplate)
	}

	return err
}

// SendRaceMailRegister to race register
func (registerMongo RegisterRepositoryMongo) SendRaceMailRegister(registerID string) error {
	var register model.Register
	var user model.UserProvider
	var mailTemplate model.EmailTemplateData2
	objectID, _ := primitive.ObjectIDFromHex(registerID)
	filter := bson.D{bson.E{Key: "_id", Value: objectID}}
	err := registerMongo.ConnectionDB.Collection("register").FindOne(context.TODO(), filter).Decode(&register)
	if err != nil {
		log.Println(err)
	}

	filterUser := bson.D{bson.E{Key: "_id", Value: register.UserID}}
	err = registerMongo.ConnectionDB.Collection("user").FindOne(context.TODO(), filterUser).Decode(&user)
	if err != nil {
		log.Println(err)
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

// GetRegisterActivateEvent repository
func (registerMongo RegisterRepositoryMongo) GetRegisterActivateEvent(userID string) ([]model.RegisterV2, error) {
	var register = []model.RegisterV2{}
	objectID, err := primitive.ObjectIDFromHex(userID)
	matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"regs.user_id": objectID}}}

	//matchSubStage := bson.D{{"$match", bson.M{"regs.user_id": bson.M{"$eq": objectID}}}}
	//groupStage := bson.D{{"_id", "$_id"}, {"event_id", "$event_id"}, {"regs", bson.M{"$push": "$regs"}}}
	//filterStage := bson.D{{"$project", bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs.user_id", objectID}}}}}}}
	projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$and": []interface{}{bson.M{"$eq": bson.A{"$$regs.user_id", objectID}}, bson.M{"$eq": bson.A{"$$regs.status", config.PAYMENT_SUCCESS}}}}}}, "event_code": 1, "ref2": 1, "event": 1}}}
	unwindStage := bson.D{primitive.E{Key: "$unwind", Value: "$regs"}}
	//projectStage := bson.D{primitive.E{Key: "$project", Value: bson.M{"regs": bson.M{"$elemMatch": bson.M{"$$regs.user_id": objectID, "$$regs.status": config.PAYMENT_SUCCESS}}, "event_id": 1}}}
	cur, err := registerMongo.ConnectionDB.Collection(registerCollection).Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage, unwindStage})
	//cur, err := registerMongo.ConnectionDB.Collection(registerCollection).Find(context.TODO(), filter)
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}

	//u.Regs = []model.Regs{}
	for cur.Next(context.TODO()) {
		var u Reg
		var regs = []model.Regs{}
		// decode the document
		if err := cur.Decode(&u); err != nil {
			log.Print(err)
		}
		event, err := DetailEventOwnerByCode(u.EventCode)
		if err == nil {
			u.Event = event
		}
		r := model.RegisterV2{
			UserCode:  u.UserCode,
			OwnerID:   u.OwnerID,
			EventCode: u.EventCode,
			Event:     event,
			Ref2:      u.Ref2,
			Regs:      append(regs, u.Regs),
		}
		register = append(register, r)
	}

	return register, err
}

//GetRegisterReport all
func (registerMongo RegisterRepositoryMongo) GetRegisterReport(formRequest model.DataRegisterRequest) (model.ReportRegister, error) {
	var register []model.Register
	var report model.ReportRegister

	objectEventID, err := primitive.ObjectIDFromHex(formRequest.EventID)
	if err != nil {
		log.Println(err)
	}

	filterSuccess := bson.D{bson.E{Key: "status", Value: config.PAYMENT_SUCCESS}, bson.E{Key: "event_id", Value: objectEventID}}
	filterWaiting := bson.D{bson.E{Key: "status", Value: config.PAYMENT_WAITING}, bson.E{Key: "event_id", Value: objectEventID}}
	filterWaitingApprove := bson.D{bson.E{Key: "status", Value: config.PAYMENT_WAITING_APPROVE}, bson.E{Key: "event_id", Value: objectEventID}}

	countAll, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), bson.D{bson.E{Key: "event_id", Value: objectEventID}})

	countSuccess, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filterSuccess)

	countWaiting, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filterWaiting)

	countWaitingApprove, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filterWaitingApprove)

	if err != nil {
		log.Println(err)

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
	options.SetSort(bson.D{bson.E{Key: "created_at", Value: -1}})

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
		filterUser := bson.D{bson.E{Key: "_id", Value: userID}}
		err := registerMongo.ConnectionDB.Collection("user").FindOne(context.TODO(), filterUser).Decode(&user)

		if err != nil {
			log.Println(err)
		}
		eventUser := item.EventID.Hex() + "." + userID.Hex()
		var activity model.Activity
		filter := bson.D{bson.E{Key: "event_user", Value: eventUser}}
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

// GetRegisterReportAll to list all register
func (registerMongo RegisterRepositoryMongo) GetRegisterReportAll(formRequest model.DataRegisterRequest) (model.ReportRegister, error) {
	var register []model.Register
	var report model.ReportRegister

	objectEventID, err := primitive.ObjectIDFromHex(formRequest.EventID)
	if err != nil {
		log.Fatal(err)

	}

	filterSuccess := bson.D{bson.E{Key: "status", Value: config.PAYMENT_SUCCESS}, bson.E{Key: "event_id", Value: objectEventID}}
	filterWaiting := bson.D{bson.E{Key: "status", Value: config.PAYMENT_WAITING}, bson.E{Key: "event_id", Value: objectEventID}}
	filterWaitingApprove := bson.D{bson.E{Key: "status", Value: config.PAYMENT_WAITING_APPROVE}, bson.E{Key: "event_id", Value: objectEventID}}

	countAll, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), bson.D{primitive.E{Key: "event_id", Value: objectEventID}})

	countSuccess, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filterSuccess)

	countWaiting, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filterWaiting)

	countWaitingApprove, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filterWaitingApprove)

	if err != nil {
		log.Println(err)
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
	options.SetSort(bson.D{bson.E{Key: "created_at", Value: -1}})

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
		filterUser := bson.D{bson.E{Key: "_id", Value: userID}}
		err := registerMongo.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filterUser).Decode(&user)

		if err != nil {
			log.Println(err)
		}
		eventUser := item.EventID.Hex() + "." + userID.Hex()
		var activity model.Activity
		filter := bson.D{bson.E{Key: "event_user", Value: eventUser}}
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
		log.Println(err)

	}

	filterSuccess := bson.D{bson.E{Key: "status", Value: config.PAYMENT_SUCCESS}, bson.E{Key: "event_id", Value: objectEventID}}
	filterWaiting := bson.D{bson.E{Key: "status", Value: config.PAYMENT_WAITING}, bson.E{Key: "event_id", Value: objectEventID}}
	filterWaitingApprove := bson.D{bson.E{Key: "status", Value: config.PAYMENT_WAITING_APPROVE}, bson.E{Key: "event_id", Value: objectEventID}}

	countAll, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), bson.D{bson.E{Key: "event_id", Value: objectEventID}})

	countSuccess, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filterSuccess)

	countWaiting, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filterWaiting)

	countWaitingApprove, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filterWaitingApprove)

	if err != nil {
		log.Println(err)

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
		filterUser := bson.D{bson.E{Key: "_id", Value: userID}}
		err := registerMongo.ConnectionDB.Collection("user").FindOne(context.TODO(), filterUser).Decode(&user)

		if err != nil {
			log.Println(err)
		}
		eventUser := item.EventID.Hex() + "." + userID.Hex()
		var activity model.Activity
		filter := bson.D{bson.E{Key: "event_user", Value: eventUser}}
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

// FindPersonRegEvent to search person register
func (registerMongo RegisterRepositoryMongo) FindPersonRegEvent(formRequest model.DataRegisterRequest) (model.ReportRegister, error) {
	var register []model.Register
	var report model.ReportRegister

	objectEventID, err := primitive.ObjectIDFromHex(formRequest.EventID)
	if err != nil {
		log.Println(err)

	}

	filterSuccess := bson.D{bson.E{Key: "status", Value: config.PAYMENT_SUCCESS}, bson.E{Key: "event_id", Value: objectEventID}}
	filterWaiting := bson.D{bson.E{Key: "status", Value: config.PAYMENT_WAITING}, bson.E{Key: "event_id", Value: objectEventID}}
	filterWaitingApprove := bson.D{bson.E{Key: "status", Value: config.PAYMENT_WAITING_APPROVE}, bson.E{Key: "event_id", Value: objectEventID}}

	countAll, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), bson.D{bson.E{Key: "event_id", Value: objectEventID}})

	countSuccess, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filterSuccess)

	countWaiting, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filterWaiting)

	countWaitingApprove, err := registerMongo.ConnectionDB.Collection(registerCollection).CountDocuments(context.TODO(), filterWaitingApprove)

	if err != nil {
		log.Println(err)

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
			filter := bson.D{bson.E{Key: "event_user", Value: eventUser}}
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

//UpdateStatusRegister repo with register_id, status, user_id
func (registerMongo RegisterRepositoryMongo) UpdateStatusRegister(registerID string, status string, userID string) error {

	objectID, err := primitive.ObjectIDFromHex(registerID)

	filter := bson.D{bson.E{Key: "_id", Value: objectID}}

	updated := bson.M{"$set": bson.M{"status": status}}

	_, err = registerMongo.ConnectionDB.Collection(registerCollection).UpdateOne(context.TODO(), filter, updated)
	if err != nil {
		//log.Fatal(res)
		log.Printf("[info] err %s", err.Error())
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

//SaveRegister repo
func SaveRegister(data model.Regs, eventID string) error {
	filter := bson.D{bson.E{Key: "reg_id", Value: data.ID}, bson.E{Key: "event_code", Value: eventID}}
	count, err := db.DB.Collection(registerAllCollection).CountDocuments(context.TODO(), filter)
	if count == 0 {
		_, err := db.DB.Collection(registerAllCollection).InsertOne(context.TODO(), data)
		return err
	}
	return err
}

// WorkoutHook repository for insert workouts
func PaymentWithSCB(p model.SCBPayment) error {
	var reg model.RegisterV2

	matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"ref2": p.Billpaymentref2, "regs.order_id": p.Billpaymentref1}}}
	//unwindStage := bson.D{{"$unwind", "$regs"}}
	//matchSubStage := bson.D{{"$match", bson.M{"regs.user_id": bson.M{"$eq": objectID}}}}
	//groupStage := bson.D{{"_id", "$_id"}, {"event_id", "$event_id"}, {"regs", bson.M{"$push": "$regs"}}}
	//filterStage := bson.D{{"$project", bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs.user_id", objectID}}}}}}}
	projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs.order_id", p.Billpaymentref1}}}}, "ref2": 1, "event_code": 1}}}
	cur, err := db.DB.Collection(registerCollection).Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage})
	if err != nil {
		return err
	}
	for cur.Next(context.TODO()) {
		// decode the document
		if err := cur.Decode(&reg); err != nil {
			log.Print(err)
		}
	}

	price := 0
	price, err = strconv.Atoi(p.Amount)
	var merchant model.Merchant
	merchant.UserID = reg.Regs[0].UserID
	merchant.EventCode = reg.EventCode
	merchant.RegID = reg.Regs[0].ID
	merchant.PaymentType = config.PAYMENT_QRCODE
	merchant.Status = config.PAYMENT_SUCCESS
	merchant.OrderID = p.Billpaymentref1
	merchant.TransactionID = p.Transactionid
	merchant.TotalPrice = int64(price)
	merchant.CreatedAt = time.Now()
	merchant.UpdatedAt = time.Now()
	//fmt.Println("Inserted a single document: ", merchant)
	res, err := db.DB.Collection(merchantCollection).InsertOne(context.TODO(), merchant)
	if err != nil {
		log.Println(res)
	}

	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	res, err = db.DB.Collection(scbPaymentCollection).InsertOne(context.TODO(), p)

	if UpdatePaymentRegStatus(reg.Regs[0].ID, reg.EventCode, reg.Regs[0].UserID, config.PAYMENT_SUCCESS, config.PAYMENT_QRCODE) {
		token, err := GetFirebaseToken(reg.Regs[0].UserID)
		if err == nil {
			if len(token.FirebaseTokens) > 0 {
				fcm := firebase.InitializeServiceAccountID()
				client, _ := fcm.Messaging(context.Background())
				body := fmt.Sprintf("Thank you for payment %s", p.Amount)
				go firebase.SendMulticastAndHandleErrors(context.Background(), client, token.FirebaseTokens, "RUNEX Register Payment", body)
			}
		}
		return nil
	}
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// AdminUpdateSlipPaymentStatus repo
func AdminUpdateSlipPaymentStatus(registerID primitive.ObjectID, eventCode string, userID primitive.ObjectID, status string, payType string, imagePath string) bool {
	filter := bson.M{"$and": []interface{}{bson.M{"event_code": eventCode}, bson.M{"regs.user_id": userID}, bson.M{"regs._id": registerID}}}
	update := bson.M{"$set": bson.M{"regs.$.status": status, "regs.$.payment_type": payType, "regs.$.payment_date": time.Now(), "regs.$.slip": imagePath, "regs.$.updated_at": time.Now()}}
	result := db.DB.Collection(registerCollection).FindOneAndUpdate(context.TODO(), filter, update)
	if result.Err() != nil {
		fmt.Println("updating the Data", result.Err())
		return false
	}
	return true
}

//UpdatePaymentRegStatus repo with register_id, status, user_id
func UpdatePaymentRegStatus(registerID primitive.ObjectID, eventCode string, userID primitive.ObjectID, status string, payType string) bool {
	filter := bson.M{"$and": []interface{}{bson.M{"event_code": eventCode}, bson.M{"regs.user_id": userID}, bson.M{"regs._id": registerID}}}
	update := bson.M{"$set": bson.M{"regs.$.status": status, "regs.$.payment_type": payType, "regs.$.payment_date": time.Now(), "regs.$.updated_at": time.Now()}}
	_, err := db.DB.Collection(registerCollection).UpdateOne(context.TODO(), filter, update)
	
	if err != nil {
		fmt.Println("updating the Data", err)
		return false
	}
	return true
}

//IsOwner repo
func IsOwner(eventCode string, userID string) bool {

	filter := bson.D{primitive.E{Key: "owner_id", Value: userID}, primitive.E{Key: "event_code", Value: eventCode}}
	count, err := db.DB.Collection(registerCollection).CountDocuments(context.TODO(), filter)
	if err != nil {
		return false
	}

	if count > 0 {
		return true
	}

	return false
}

func UpdateRegisterUserInfo(userID primitive.ObjectID, req model.RegUpdateUserInfoRequest) error {
	filter := bson.M{"$and": []interface{}{bson.M{"event_code": req.EventCode}, bson.M{"regs._id": req.RegID}}}
	t := []model.TicketOptionV2{}
	t = append(t, req.TicketOption)
	log.Println(filter)
	update := bson.M{"$set": bson.M{"regs.$.status": config.PAYMENT_SUCCESS, "regs.$.payment_date": time.Now(), "regs.$.updated_at": time.Now(), "regs.$.ticket_options": t}}
	_, err := db.DB.Collection(registerCollection).UpdateOne(context.TODO(), filter, update)
	return err
}

func EditUserOption(userID primitive.ObjectID, req model.RegUpdateUserInfoRequest) error {
	filter := bson.M{"$and": []interface{}{bson.M{"event_code": req.EventCode}, bson.M{"regs._id": req.RegID}}}
	t := []model.TicketOptionV2{}
	t = append(t, req.TicketOption)
	log.Println(req.TicketOption.UserOption)
	update := bson.M{"$set": bson.M{"regs.$.updated_at": time.Now(), "regs.$.ticket_options": t}}
	result := db.DB.Collection(registerCollection).FindOneAndUpdate(context.TODO(), filter, update)
	return result.Err()
}
