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
	UpdateUser(u model.User, userID string) (model.User, error)
	UpdateUserStrava(u model.UserStravaSyncRequest, userID string) (model.User, error)
	GetUserWithProvider(provider string, providerID string) (model.User, error)

	//fcm
	FirebaseRegister(token string, userID string)
	FirebaseRemove(u model.RegisterTokenRequest, userID string)
	GetFirebaseTokenAll() ([]model.FirebaseUser, error)
	GetFirebaseToken() (model.FirebaseUser, error)
}

//RepoUserDB db connection struct
type RepoUserDB struct {
	ConnectionDB *mongo.Database
}

const (
	userConlection = "user"
	fcmConlection  = "fcm"
)

// Signin repository for login api call
func (db RepoUserDB) Signin(u model.UserProviderRequest) (model.User, error) {
	var user model.User
	filter := bson.D{
		primitive.E{Key: "email", Value: u.Email},
		primitive.E{Key: "provider_id", Value: u.ProviderID},
		primitive.E{Key: "provider", Value: u.Provider},
	}
	err := db.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		count, err := db.ConnectionDB.Collection(userConlection).CountDocuments(context.TODO(), filter)
		if count < 1 {
			user, err = db.AddUser(u)
			return user, err
		}
		
	}
	return user, err
}

//CheckEmail Email in db
func (db RepoUserDB) CheckEmail(u model.UserProviderRequest) (model.User, bool, bool) {
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
func (db RepoUserDB) CheckProvider(provider model.UserProviderRequest) bool {
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
func (db RepoUserDB) GetUserByProvider(u model.UserProviderRequest) (model.User, error) {
	var user model.User
	filter := bson.D{primitive.E{Key: "provider", Value: u.Provider}, primitive.E{Key: "provider_id", Value: u.ProviderID}}
	err := db.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Println(err)
	}
	return user, err
}

// GetUser with user id and return user info
func (db RepoUserDB) GetUser(id primitive.ObjectID) (model.User, error) {
	var user model.User
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	err := db.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Println(err)
	}
	return user, err
}

// GetUserWithProvider with user id and return user info
func (db RepoUserDB) GetUserWithProvider(provider string, providerID string) (model.User, error) {
	var user model.User
	filter := bson.D{primitive.E{Key: "provider", Value: provider}, primitive.E{Key: "provider_id", Value: providerID}}
	err := db.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Println(err)
	}
	return user, err
}

// AddUser with provider api add user from social provider
func (db RepoUserDB) AddUser(user model.UserProviderRequest) (model.User, error) {
	var u model.User
	u = model.User{
		Email:      user.Email,
		FullName:   user.FullName,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Provider:   user.Provider,
		ProviderID: user.ProviderID,
		PF:         user.PF,
		Avatar:     user.Avatar,
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
func (db RepoUserDB) UpdateUser(u model.User, userID string) (model.User, error) {
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
		return u, result.Err()
	}
	return u, result.Err()
}

// FirebaseCreate repo initial
func (db RepoUserDB) FirebaseCreate(userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	filter := bson.D{primitive.E{Key: "user_id", Value: id}}
	fcm := model.FirebaseUser{
		UserID:         id,
		FirebaseTokens: []string{},
	}
	count, err := db.ConnectionDB.Collection(fcmConlection).CountDocuments(context.TODO(), filter)
	if err != nil {
		_, err = db.ConnectionDB.Collection(fcmConlection).InsertOne(context.TODO(), fcm)
		return err
	}
	if count == 0 {
		_, err = db.ConnectionDB.Collection(fcmConlection).InsertOne(context.TODO(), fcm)
		return err
	}
	return err
}

// FirebaseRegister repo register firebase token
func (db RepoUserDB) FirebaseRegister(token string, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	filter := bson.D{primitive.E{Key: "user_id", Value: id}}
	var fcm model.FirebaseUser
	count, err := db.ConnectionDB.Collection(fcmConlection).CountDocuments(context.TODO(), filter)
	if count == 0 {
		fcm := model.FirebaseUser{
			UserID:         id,
			FirebaseTokens: []string{ token },
		}
		_, err = db.ConnectionDB.Collection(fcmConlection).InsertOne(context.TODO(), fcm)
		return err
		
	}
	err = db.ConnectionDB.Collection(fcmConlection).FindOne(context.TODO(), filter).Decode(&fcm)
	if err != nil {
		return err
	}
	check := false
	for _, s := range fcm.FirebaseTokens {
		if token == s {
			check = true
		}
	}
	if check {
		return nil
	}
	fcm.FirebaseTokens = append(fcm.FirebaseTokens, token)
	isUpsert := true
	clientOptions := options.FindOneAndUpdateOptions{Upsert: &isUpsert}
	update := bson.M{"$set": fcm}
	result := db.ConnectionDB.Collection(fcmConlection).FindOneAndUpdate(context.TODO(), filter, update, &clientOptions)
	return result.Err()
}

// GetFirebaseToken repo register firebase token
func (db RepoUserDB) GetFirebaseToken(userID primitive.ObjectID) (model.FirebaseUser, error) {
	filter := bson.D{primitive.E{Key: "user_id", Value: userID}}
	var fcm model.FirebaseUser
	err := db.ConnectionDB.Collection(fcmConlection).FindOne(context.TODO(), filter).Decode(&fcm)
	return fcm, err
}

// GetFirebaseTokenAll repo register firebase token
func (db RepoUserDB) GetFirebaseTokenAll() ([]model.FirebaseUser, error) {
	filter := bson.D{}
	var fcms []model.FirebaseUser
	cur, err := db.ConnectionDB.Collection(fcmConlection).Find(context.TODO(), filter)
	if err != nil {
		return []model.FirebaseUser{}, err
	}
	for cur.Next(context.TODO()) {
		var a model.FirebaseUser
		// decode the document
		if err := cur.Decode(&a); err != nil {
			log.Println(err)
		}

		fcms = append(fcms, a)
	}
	return fcms, err
}

// FirebaseRemove api remove firebase token
func (db RepoUserDB) FirebaseRemove(u model.RegisterTokenRequest, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	filter := bson.M{"user_id":id}
	update := bson.M{"$pull": bson.M{"firebase_tokens": u.FirebaseToken }}
	_, err = db.ConnectionDB.Collection(fcmConlection).UpdateOne(context.TODO(), filter, update)
	return err
}

// UpdateUserStrava api update account profile sync strava
func (db RepoUserDB) UpdateUserStrava(u model.UserStravaSyncRequest, userID string) (model.User, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: userID},
		primitive.E{Key: "provider_id", Value: u.ProviderID},
		primitive.E{Key: "provider", Value: u.Provider}}
	//isUpsert := true
	//clientOptions := options.FindOneAndUpdateOptions{Upsert: &isUpsert}
	update := bson.M{"$set": bson.M{
		"strava_id":        u.StravaID,
		"strava_avatar":    u.StravaAvatar,
		"strava_firstname": u.StravaFirstname,
		"strava_lastname":  u.StravaLastname,
		"updated_at":       time.Now(),
	}}
	_, err := db.ConnectionDB.Collection(userConlection).UpdateOne(context.TODO(), filter, update)
	var user model.User
	if err == nil {
		db.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(&user)
		return user, err
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
