package repository

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"thinkdev.app/think/runex/runexapi/model"
)

type EventRepository interface {
	// AddEvent(event model.Event) (string, error)
	// GetEventByStatus(status string) ([]model.Event, error)
	// GetEventAll() ([]model.Event, error)
	// GetEventActive() ([]model.Event, error)
	// ExistByName(name string) (bool, error)
	// ExistByNameForEdit(name string, eventID string) (bool, error)
	// GetEventByID(eventID string) (model.Event, error)
	// EditEvent(eventID string, event model.Event) error
	// DeleteEventByID(eventID string) error
	// UploadCoverEvent(eventID string, path string) error
	// AddProductEvent(eventID string, product model.ProduceEvent) (string, error)
	// EditProductEvent(eventID string, product model.ProduceEvent) error
	// GetProductByEventID(eventID string) ([]model.ProduceEvent, error)
	// DeleteProductEvent(eventID string, productID string) error
	// AddTicketEvent(eventID string, ticket model.TicketEvent) (string, error)
	// EditTicketEvent(eventID string, ticket model.TicketEvent) error
	// GetTicketByEventID(eventID string) ([]model.TicketEvent, error)
	// DeleteTicketEvent(eventID string, ticketID string) error
	// GetUserEvent(userID string) (model.UserEvent, error)
	// GetEventByUser(userID string) ([]model.Event, error)
	// SearchEvent(term string) ([]model.Event, error)
	// ValidateBySlug(slug string) (bool, error)
	// GetEventBySlug(slug string) (model.Event, error)
}
type EventRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

const (
	eventCollection = "event"
	ebibCollection  = "ebib"
)

// func (eventMongo EventRepositoryMongo) AddEvent(event model.Event) (string, error) {

// 	event.Product = []model.ProduceEvent{}
// 	event.Ticket = []model.TicketEvent{}
// 	event.CreatedTime = time.Now()
// 	event.UpdatedTime = time.Now()
// 	res, err := eventMongo.ConnectionDB.Collection(eventCollection).InsertOne(context.TODO(), event)
// 	if err != nil {
// 		log.Fatal(res)
// 	}
// 	var ebib model.EbibEvent
// 	ebib.EventID = event.ID
// 	ebib.LastNo = 0
// 	ebib.CreatedAt = time.Now()
// 	ebib.UpdatedAt = time.Now()
// 	eventMongo.ConnectionDB.Collection(ebibCollection).InsertOne(context.TODO(), ebib)

// 	fmt.Println("Inserted a single document: ", res.InsertedID)
// 	return res.InsertedID.(primitive.ObjectID).Hex(), err
// }

// func (eventMongo EventRepositoryMongo) EditEvent(eventID string, event model.Event) error {

// 	objectID, err := primitive.ObjectIDFromHex(eventID)
// 	filter := bson.D{{"_id", objectID}}
// 	event.UpdatedTime = time.Now()
// 	updated := bson.M{"$set": event}
// 	res, err := eventMongo.ConnectionDB.Collection(eventCollection).UpdateOne(context.TODO(), filter, updated)
// 	if err != nil {
// 		//log.Fatal(res)
// 		log.Printf("[info] err %s", res)
// 		return err
// 	}

// 	return nil
// }

// func (eventMongo EventRepositoryMongo) GetEventByStatus(status string) ([]model.Event, error) {
// 	var events []model.Event
// 	filter := bson.D{{"status", status}}
// 	cur, err := eventMongo.ConnectionDB.Collection(eventCollection).Find(context.TODO(), filter)
// 	//log.Printf("[info] cur %s", cur)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	for cur.Next(context.TODO()) {
// 		var u model.Event
// 		// decode the document
// 		if err := cur.Decode(&u); err != nil {
// 			log.Fatal(err)
// 		}
// 		//fmt.Printf("post: %+v\n", p)
// 		events = append(events, u)
// 	}

// 	return events, err
// }

// func (eventMongo EventRepositoryMongo) GetEventAll() ([]model.Event, error) {
// 	var events []model.Event
// 	options := options.Find()
// 	options.SetSort(bson.D{{"created_time", -1}})
// 	cur, err := eventMongo.ConnectionDB.Collection(eventCollection).Find(context.TODO(), bson.D{{}}, options)
// 	//log.Printf("[info] cur %s", cur)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	for cur.Next(context.TODO()) {
// 		var u model.Event
// 		// decode the document
// 		if err := cur.Decode(&u); err != nil {
// 			log.Fatal(err)
// 		}
// 		//fmt.Printf("post: %+v\n", p)
// 		events = append(events, u)
// 	}

// 	return events, err
// }

func (eventMongo EventRepositoryMongo) GetEventActive() ([]model.Event, error) {
	var events []model.Event
	// Sort by _id field descending
	options := options.Find()
	options.SetSort(bson.D{{"created_time", -1}})
	filter := bson.D{{"is_active", true}}
	cur, err := eventMongo.ConnectionDB.Collection(eventCollection).Find(context.TODO(), filter, options)
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {
		var u model.Event
		// decode the document
		if err := cur.Decode(&u); err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("post: %+v\n", p)
		events = append(events, u)
	}

	return events, err
}

func (eventMongo EventRepositoryMongo) ExistByName(name string) (bool, error) {

	filter := bson.D{{"name", name}}
	count, err := eventMongo.ConnectionDB.Collection(eventCollection).CountDocuments(context.TODO(), filter)
	log.Printf("[info] count %s", count)
	if err != nil {
		log.Println(err)
	}
	if count > 0 {
		return true, nil
	}

	return false, nil
}

// func (eventMongo EventRepositoryMongo) ExistByNameForEdit(name string, eventID string) (bool, error) {

// 	var event model.Event
// 	id, err := primitive.ObjectIDFromHex(eventID)
// 	filter := bson.M{"_id": id}
// 	err2 := eventMongo.ConnectionDB.Collection(eventCollection).FindOne(context.TODO(), filter).Decode(&event)
// 	log.Printf("[info] event %s", err2)
// 	if err2 != nil {
// 		log.Fatal(err2)
// 		//return true, err2
// 	}
// 	if event.Name == name {
// 		return false, nil
// 	}

// 	filter2 := bson.D{{"name", name}}
// 	count, err := eventMongo.ConnectionDB.Collection(eventCollection).CountDocuments(context.TODO(), filter2)
// 	log.Printf("[info] count %s", count)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	if count > 0 {
// 		return true, nil
// 	}

// 	return false, nil
// }

func (eventMongo EventRepositoryMongo) GetEventByID(eventID string) (model.Event, error) {

	var event model.Event
	id, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		log.Println(err)
	}
	filter := bson.M{"_id": id}
	err2 := eventMongo.ConnectionDB.Collection(eventCollection).FindOne(context.TODO(), filter).Decode(&event)
	log.Printf("[info Event] cur %s", err2)
	if err2 != nil {
		log.Println(err2)
	}

	return event, err2
}

func (eventMongo EventRepositoryMongo) DeleteEventByID(eventID string) error {

	id, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": id}
	deleteResult, err2 := eventMongo.ConnectionDB.Collection(eventCollection).DeleteOne(context.TODO(), filter)
	log.Printf("[info] cur %s", err2)
	if err2 != nil {
		return err2
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
	return nil
}

func (eventMongo EventRepositoryMongo) UploadCoverEvent(eventID string, path string) error {
	objectID, err := primitive.ObjectIDFromHex(eventID)
	filter := bson.D{{"_id", objectID}}

	updated := bson.M{"$set": bson.M{"cover": path, "updated_time": time.Now()}}
	res, err := eventMongo.ConnectionDB.Collection(eventCollection).UpdateOne(context.TODO(), filter, updated)
	if err != nil {
		//log.Fatal(res)
		log.Printf("[info] err %s", res)
		return err
	}

	return nil
}

func (eventMongo EventRepositoryMongo) AddProductEvent(eventID string, product model.ProduceEvent) (string, error) {
	objectID, err := primitive.ObjectIDFromHex(eventID)
	filter := bson.D{{"_id", objectID}}

	product.ProductID = primitive.NewObjectID()
	product.UpdatedAt = time.Now()
	product.CreatedAt = time.Now()

	log.Println(product)

	update := bson.M{"$push": bson.M{"product": product}}

	res, err := eventMongo.ConnectionDB.Collection(eventCollection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println(res)
		log.Printf("[info] err %s", err)
		// return err
	}

	return product.ProductID.Hex(), err
}

func (eventMongo EventRepositoryMongo) EditProductEvent(eventID string, product model.ProduceEvent) error {
	objectID, err := primitive.ObjectIDFromHex(eventID)
	//productObjectID, err := primitive.ObjectIDFromHex(product.ProductID)
	filter := bson.D{{"_id", objectID}, {"product._id", product.ProductID}}

	product.UpdatedAt = time.Now()

	update := bson.M{"$set": bson.M{"product.$": product}}

	res, err := eventMongo.ConnectionDB.Collection(eventCollection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		//log.Fatal(res)
		log.Printf("[info] err %s", res)
		return err
	}

	return nil
}

// func (eventMongo EventRepositoryMongo) GetProductByEventID(eventID string) ([]model.ProduceEvent, error) {

// 	objectID, err := primitive.ObjectIDFromHex(eventID)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	var event model.Event
// 	var productInfo []model.ProduceEvent
// 	filter := bson.D{{"_id", objectID}}

// 	err2 := eventMongo.ConnectionDB.Collection(eventCollection).FindOne(context.TODO(), filter).Decode(&event)

// 	if err2 != nil {
// 		log.Println(err2)
// 		return nil, err2
// 	}
// 	productInfo = event.Product
// 	return productInfo, err2
// }

func (eventMongo EventRepositoryMongo) DeleteProductEvent(eventID string, productID string) error {
	objectID, err := primitive.ObjectIDFromHex(eventID)
	productObjectID, err := primitive.ObjectIDFromHex(productID)

	//filter := bson.D{{"_id", objectID}}
	filter := bson.D{{"_id", objectID}}

	update := bson.M{"$pull": bson.M{"product": bson.M{"_id": productObjectID}}}

	res, err := eventMongo.ConnectionDB.Collection(eventCollection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		//log.Fatal(res)
		log.Printf("[info] err %s", res)
		return err
	}
	return nil
}

func (eventMongo EventRepositoryMongo) AddTicketEvent(eventID string, ticket model.TicketEvent) (string, error) {
	objectID, err := primitive.ObjectIDFromHex(eventID)
	filter := bson.D{{"_id", objectID}}

	ticket.TicketID = primitive.NewObjectID()
	ticket.UpdatedAt = time.Now()
	ticket.CreatedAt = time.Now()

	update := bson.M{"$push": bson.M{"ticket": ticket}}

	res, err := eventMongo.ConnectionDB.Collection(eventCollection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(res)
		// log.Printf("[info] err %s", res)
		// return err
	}

	return ticket.TicketID.Hex(), err
}

func (eventMongo EventRepositoryMongo) EditTicketEvent(eventID string, ticket model.TicketEvent) error {
	objectID, err := primitive.ObjectIDFromHex(eventID)
	//productObjectID, err := primitive.ObjectIDFromHex(product.ProductID)
	filter := bson.D{{"_id", objectID}, {"ticket._id", ticket.TicketID}}

	ticket.UpdatedAt = time.Now()

	update := bson.M{"$set": bson.M{"ticket.$": ticket}}

	res, err := eventMongo.ConnectionDB.Collection(eventCollection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		//log.Fatal(res)
		log.Printf("[info] err %s", res)
		return err
	}

	return nil
}

// func (eventMongo EventRepositoryMongo) GetTicketByEventID(eventID string) ([]model.TicketEvent, error) {

// 	objectID, err := primitive.ObjectIDFromHex(eventID)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	var event model.Event
// 	var ticketInfo []model.TicketEvent
// 	filter := bson.D{{"_id", objectID}}

// 	err2 := eventMongo.ConnectionDB.Collection(eventCollection).FindOne(context.TODO(), filter).Decode(&event)

// 	if err2 != nil {
// 		log.Println(err2)
// 		return nil, err2
// 	}
// 	ticketInfo = event.Ticket
// 	return ticketInfo, err2
// }

func (eventMongo EventRepositoryMongo) DeleteTicketEvent(eventID string, ticketID string) error {
	objectID, err := primitive.ObjectIDFromHex(eventID)
	ticketObjectID, err := primitive.ObjectIDFromHex(ticketID)

	//filter := bson.D{{"_id", objectID}}
	filter := bson.D{{"_id", objectID}}

	update := bson.M{"$pull": bson.M{"ticket": bson.M{"_id": ticketObjectID}}}

	res, err := eventMongo.ConnectionDB.Collection(eventCollection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		//log.Fatal(res)
		log.Printf("[info] err %s", res)
		return err
	}
	return nil
}

func (eventMongo EventRepositoryMongo) GetUserEvent(userID string) (model.UserEvent, error) {
	var user model.UserEvent
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
	}
	filter := bson.M{"_id": id}
	err2 := eventMongo.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(&user)
	if err2 != nil {
		log.Println(err2)
	}
	return user, err
}

func (eventMongo EventRepositoryMongo) GetEventByUser(userID string) ([]model.Event, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	var events []model.Event
	filter := bson.D{{"owner_id", objectID}}
	cur, err := eventMongo.ConnectionDB.Collection(eventCollection).Find(context.TODO(), filter)
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {
		var u model.Event
		// decode the document
		if err := cur.Decode(&u); err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("post: %+v\n", p)
		events = append(events, u)
	}

	return events, err
}

func (eventMongo EventRepositoryMongo) SearchEvent(term string) ([]model.Event, error) {

	filter := bson.D{{"name", bson.D{{"$regex", strings.TrimSpace(term)}}}}
	//filter := bson.D{{"$text", bson.D{{"$search", strings.TrimSpace(term)}}}}
	//index := bson.D{{"name", "text"}, {"description", "text"}}
	var events []model.Event
	cur, err := eventMongo.ConnectionDB.Collection(eventCollection).Find(context.TODO(), filter)
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {
		var u model.Event
		// decode the document
		if err := cur.Decode(&u); err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("post: %+v\n", p)
		events = append(events, u)
	}

	return events, err
}

func (eventMongo EventRepositoryMongo) ValidateBySlug(slugText string) (bool, error) {

	filter := bson.D{{"slug", slugText}}
	count, err := eventMongo.ConnectionDB.Collection(eventCollection).CountDocuments(context.TODO(), filter)
	log.Printf("[info] count %s", count)
	if err != nil {
		log.Println(err)
	}
	if count > 0 {
		return false, nil
	}

	return true, nil
}

func (eventMongo EventRepositoryMongo) GetEventBySlug(slug string) (model.Event, error) {
	var event model.Event
	filter := bson.D{{"slug", slug}}
	err2 := eventMongo.ConnectionDB.Collection(eventCollection).FindOne(context.TODO(), filter).Decode(&event)
	log.Printf("[info Event] cur %s", err2)
	if err2 != nil {
		log.Fatal(err2)
	}

	return event, err2
}
