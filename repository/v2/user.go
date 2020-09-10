package repository

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"thinkdev.app/think/runex/runexapi/config"
	model2 "thinkdev.app/think/runex/runexapi/model/v2"
	"thinkdev.app/think/runex/runexapi/model"
)

//UserRepository interface user repository
type UserRepository interface {
	CheckEmail(u model2.UserProviderRequest) (model2.User, bool, bool)
	CheckProvider(provider model2.UserProviderRequest) bool
	UpdateProvider(u model2.User, p model2.UserProviderRequest) error
	GetUserByProvider(u model2.UserProviderRequest) (model2.User, error)
	GetUser(id primitive.ObjectID) (model2.User, error)
	AddUser(user model2.UserProviderRequest) (model2.User, error)
}

//RepoDB db connection struct
type RepoDB struct {
	ConnectionDB *mongo.Database
}

const (
	userConlection = "user"
)

//CheckEmail Email in db
func (db RepoDB) CheckEmail(u model2.UserProviderRequest) (model2.User, bool, bool) {
	var user model2.User
	isMail := false
	isProvider := false
	filter := bson.D{
		primitive.E{Key: "email", Value: u.Email},
		primitive.E{Key: "provider.provider_id", Value: u.ProviderID},
		primitive.E{Key: "provider.provider_name", Value: u.Provider},
	}
	err := db.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		isProvider = false
		filter = bson.D{
			primitive.E{Key: "email", Value: u.Email},
		}
		err = db.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(&user)
		if err != nil {
			isMail = false
			return user, isMail, isProvider
		}
		isMail = true
		return user, isMail, isProvider
	}
	isProvider = true
	isMail = true

	return user, isMail, isProvider
}

//CheckProvider check provider with email user
func (db RepoDB) CheckProvider(provider model2.UserProviderRequest) bool {
	var user model2.User
	filter := bson.D{primitive.E{Key: "provider", Value: provider.Provider}, primitive.E{Key: "provider_id", Value: provider.ProviderID}}
	err := db.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// GetUserByProvider api login by provider
func (db RepoDB) GetUserByProvider(u model2.UserProviderRequest) (model2.User, error) {
	var user model2.User
	filter := bson.D{primitive.E{Key: "provider", Value: u.Provider}, primitive.E{Key: "provider_id", Value: u.ProviderID}}
	err := db.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Println(err)
	}
	return user, err
}

// GetUser with user id and return user info
func (db RepoDB) GetUser(id primitive.ObjectID) (model2.User, error) {
	var user model2.User
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	err := db.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Println(err)
	}
	return user, err
}

// AddUser with provider api add user from social provider
func (db RepoDB) AddUser(user model2.UserProviderRequest) (model2.User, error) {
	var u model2.User
	u = model2.User{
		Email:     user.Email,
		FullName:  user.FullName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		PF:        user.PF,
		Avatar:    user.Avatar,
	}
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	if u.Role == "" {
		u.Role = config.MEMBER
	}
	var p = model.Provider{
		ProviderID:   user.ProviderID,
		ProviderName: user.Provider,
	}

	u.Provider = append(u.Provider, p)
	u.Address = []model.Address{}
	result, err := db.ConnectionDB.Collection(userConlection).InsertOne(context.TODO(), u)
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		u.UserID = oid
		return u, err
	}
	return u, err

}

// UpdateProvider api update account login without provider
func (db RepoDB) UpdateProvider(u model2.User, p model2.UserProviderRequest) error {
	filter := bson.D{primitive.E{Key: "_id", Value: u.UserID}}
	var up = model.Provider{
		ProviderID:   p.ProviderID,
		ProviderName: p.Provider,
	}
	u.UpdatedAt = time.Now()
	u.Provider = append(u.Provider, up)
	isUpsert := true
	clientOptions := options.FindOneAndUpdateOptions{Upsert: &isUpsert}
	update := bson.M{"$set": u}
	result := db.ConnectionDB.Collection(userConlection).FindOneAndUpdate(context.TODO(), filter, update, &clientOptions)
	return result.Err()
}
