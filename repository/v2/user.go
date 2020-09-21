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
	"thinkdev.app/think/runex/runexapi/model"
)

//UserRepository interface user repository
type UserRepository interface {
	Signin(u model.UserProviderRequest) (model.User, error)
	CheckEmail(u model.UserProviderRequest) (model.User, bool, bool)
	CheckProvider(provider model.UserProviderRequest) bool
	//UpdateProvider(u model.User, p model.UserProviderRequest) error
	GetUserByProvider(u model.UserProviderRequest) (model.User, error)
	GetUser(id primitive.ObjectID) (model.User, error)
	AddUser(user model.UserProviderRequest) (model.User, error)
	UpdateUser(u model.User, userID string) ( model.User, error)
	UpdateUserStrava(u model.UserStravaSyncRequest, userID string) ( model.User, error)
}

//RepoDB db connection struct
type RepoDB struct {
	ConnectionDB *mongo.Database
}

const (
	userConlection = "user"
)

// Signin repository for login api call
func (db RepoDB) Signin(u model.UserProviderRequest) (model.User, error) {
	var user model.User
	filter := bson.D{
		primitive.E{Key: "email", Value: u.Email},
		primitive.E{Key: "provider_id", Value: u.ProviderID},
		primitive.E{Key: "provider", Value: u.Provider},
	}
	err := db.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		user, err = db.AddUser(u)
		return user, err
	}
	return user, err
}

//CheckEmail Email in db
func (db RepoDB) CheckEmail(u model.UserProviderRequest) (model.User, bool, bool) {
	var user model.User
	isMail := false
	isProvider := false
	filter := bson.D{
		primitive.E{Key: "email", Value: u.Email},
		primitive.E{Key: "provider_id", Value: u.ProviderID},
		primitive.E{Key: "provider", Value: u.Provider},
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
func (db RepoDB) CheckProvider(provider model.UserProviderRequest) bool {
	var user model.User
	filter := bson.D{primitive.E{Key: "provider", Value: provider.Provider}, primitive.E{Key: "provider_id", Value: provider.ProviderID}}
	err := db.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// GetUserByProvider api login by provider
func (db RepoDB) GetUserByProvider(u model.UserProviderRequest) (model.User, error) {
	var user model.User
	filter := bson.D{primitive.E{Key: "provider", Value: u.Provider}, primitive.E{Key: "provider_id", Value: u.ProviderID}}
	err := db.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Println(err)
	}
	return user, err
}

// GetUser with user id and return user info
func (db RepoDB) GetUser(id primitive.ObjectID) (model.User, error) {
	var user model.User
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	err := db.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Println(err)
	}
	return user, err
}

// AddUser with provider api add user from social provider
func (db RepoDB) AddUser(user model.UserProviderRequest) (model.User, error) {
	var u model.User
	u = model.User{
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
	u.Address = []model.Address{}
	u.Events = []primitive.ObjectID{}
	result, err := db.ConnectionDB.Collection(userConlection).InsertOne(context.TODO(), u)
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		u.UserID = oid
		return u, err
	}
	return u, err

}

// UpdateUser api update account profile
func (db RepoDB) UpdateUser(u model.User, userID string) ( model.User, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return u, err
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	u.UpdatedAt = time.Now()
	isUpsert := true
	clientOptions := options.FindOneAndUpdateOptions{Upsert: &isUpsert}
	update := bson.M{"$set": u}
	result := db.ConnectionDB.Collection(userConlection).FindOneAndUpdate(context.TODO(), filter, update, &clientOptions)
	if result.Err() == nil {
		u.UserID = id
		return u,result.Err()
	}
	return u,result.Err()
}

// UpdateUserStrava api update account profile sync strava
func (db RepoDB) UpdateUserStrava(u model.UserStravaSyncRequest, userID string) ( model.User, error){
	filter := bson.D{primitive.E{Key: "_id", Value: userID}, 
	primitive.E{Key: "provider_id", Value: u.ProviderID}, 
	primitive.E{Key: "provider", Value: u.Provider}}
	//isUpsert := true
	//clientOptions := options.FindOneAndUpdateOptions{Upsert: &isUpsert}
	update := bson.M{"$set": bson.M{
		"strava_id": u.StravaID,
		"strava_avatar": u.StravaAvatar,
		"strava_firstname": u.StravaFirstname,
		"strava_lastname": u.StravaLastname,
		"updated_at": time.Now(),
	  },}
	_, err := db.ConnectionDB.Collection(userConlection).UpdateOne(context.TODO(), filter, update)
	var user model.User
	if err == nil {
		db.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(&user)
		return user,err
	}
	return user, err
}

// UpdateProvider api update account login without provider
/*func (db RepoDB) UpdateProvider(u model.User, p model.UserProviderRequest) error {
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
}*/
