package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"thinkdev.app/think/runex/runexapi/model"
)

//EventRepositoryMongo mongo ref
type EventRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

const (
	eventCollection = "event_v2"
	ebibCollection  = "ebib"
)

//DetailEventByCode go doc
//Description get Get event detail API calls to event runex
func DetailEventByCode(code string) (model.Event, error) {
	urlS := fmt.Sprintf("%s/event/%s",viper.GetString("events.api"), code)
	var bearer = "Bearer olcgZVpqDXQikRDG"
	//reqURL, _ := url.Parse(urlS)
	req, err := http.NewRequest("GET", urlS, nil)
	req.Header.Add("Authorization", bearer)
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded, charset=UTF-8")

	timeout := time.Duration(6 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	client.CheckRedirect = checkRedirectFunc

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	var eventRes model.EventResponse
	var event model.Event

	if resp.StatusCode >= 200 || resp.StatusCode < 300 {

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return event, err
		}

		err = json.Unmarshal(body, &eventRes)
		if err != nil {
			log.Println(err)
			return event, err
		}
		return eventRes.Event, err
	}

	return event, err

}

//DetailEventOwnerByCode go doc
//Description get Get event detail API calls to event runex
func DetailEventOwnerByCode(code string) (model.Event, error) {
	urlS := fmt.Sprintf("%s/event-detail/%s", viper.GetString("events.api"), code)
	var bearer = "5Dk2o03a4hPglRSUuNeEFXEah57CGQfMVjQ7f"
	//reqURL, _ := url.Parse(urlS)
	req, err := http.NewRequest("GET", urlS, nil)
	req.Header.Add("token", bearer)
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded, charset=UTF-8")

	timeout := time.Duration(6 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	client.CheckRedirect = checkRedirectFunc

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	var eventRes model.EventResponse
	var event model.Event

	if resp.StatusCode >= 200 || resp.StatusCode < 300 {

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return event, err
		}

		err = json.Unmarshal(body, &eventRes)
		if err != nil {
			log.Println(err)
			return event, err
		}
		return eventRes.Event, err
	}

	return event, err

}

// ExistByName repo
func (eventMongo EventRepositoryMongo) ExistByName(name string) (bool, error) {

	filter := bson.D{primitive.E{Key: "name", Value: name}}
	count, err := eventMongo.ConnectionDB.Collection(eventCollection).CountDocuments(context.TODO(), filter)
	log.Printf("[info] count %d", count)
	if err != nil {
		log.Println(err)
	}
	if count > 0 {
		return true, nil
	}

	return false, nil
}

//IsOwner repo
func (eventMongo EventRepositoryMongo) IsOwner(eventID string, userID string) bool {
	eID, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		return false
	}
	uID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false
	}
	filter := bson.D{primitive.E{Key: "_id", Value: eID}, primitive.E{Key: "user_id", Value: uID}}
	count, err := eventMongo.ConnectionDB.Collection(eventCollection).CountDocuments(context.TODO(), filter)
	if err != nil {
		return false
	}

	if count > 0 {
		return true
	}

	return false
}

func checkRedirectFunc(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", via[0].Header.Get("Authorization"))
	return nil
}
