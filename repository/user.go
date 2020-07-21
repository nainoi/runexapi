package repository

import (
	"context"
	"log"
	"time"

	"bitbucket.org/suthisakch/runex/config"
	"bitbucket.org/suthisakch/runex/model"
	"bitbucket.org/suthisakch/runex/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetAllUser() ([]model.User, error)
	GetMyUser(userID string) (model.User, error)
	GetUserByProvider(provider string, providerID string) (model.UserAuth, error)
	GetUserByEmail(email string, password string) (model.UserAuth, error)
	CheckEmail(email string) bool
	CheckProvider(provider model.UserProvider) bool
	UpdateProvider(user model.UserProvider) error
	//GetLastUser() (model.User, error)
	GetUser(userID string) (model.User, error)
	AddUserPD(user model.UserProvider) error
	AddUserEP(user model.UserMail) error
	AddUserAdmin(user model.UserMail) error
	EditUser(userID string, user model.User) (model.User, error)
	Confirm(userID string) error
	DeleteUserByID(userID string) error
	GetUserEvent(userID string) (model.UserEvent, error)
	UploadAvatar(eventID string, path string) error

	AddAddress(userID string, address model.Address) (model.User, error)
	UpdateAddress(userID string, address model.Address) (model.User, error)
	ChangePassword(userID string, oldPassword string, newPassword string) error
	ForgotPassword(email string,) (model.UserForgot ,error)
	UpdatePassword(userID string, newPassword string) error
}

type UserRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

const (
	userConlection = "user"
)

func (userMongo UserRepositoryMongo) GetAllUser() ([]model.User, error) {
	var users []model.User
	cursor, err := userMongo.ConnectionDB.Collection(userConlection).Find(context.TODO(), bson.D{})
	if err != nil {
		log.Println(err)
	}

	for cursor.Next(context.Background()) {
		var u model.User
		// decode the document
		if err := cursor.Decode(&u); err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("post: %+v\n", p)
		users = append(users, u)
	}

	return users, err
}

// GetUserByProvider api login by provider
func (userMongo UserRepositoryMongo) GetUserByProvider(provider string, providerID string) (model.UserAuth, error) {
	var user model.UserAuth
	filter := bson.D{{"provider", provider}, {"provider_id", providerID}}
	err := userMongo.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Println(err)
	}
	return user, err
}

//GetUserByEmail api login by email
func (userMongo UserRepositoryMongo) GetUserByEmail(email string, password string) (model.UserAuth, error) {
	var user model.UserAuth
	//pass := utils.HashAndSalt([]byte(password))
	filter := bson.D{{"email", email}}
	err := userMongo.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Println(err)
	}
	err2 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err2 != nil {
		return user, err2
	}
	return user, err
}

func (userMongo UserRepositoryMongo) CheckEmail(email string) bool {
	var user model.UserAuth
	filter := bson.D{{"email", email}}
	err := userMongo.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(&user)
	//log.Println("-----------------------")
	//log.Println(email)
	if err != nil {
		//log.Println(err)
		return false
	}
	//log.Println(user)
	return true
}

func (userMongo UserRepositoryMongo) CheckProvider(provider model.UserProvider) bool {
	var user model.UserAuth
	filter := bson.D{{"provider", provider.Provider}, {"provider_id", provider.ProviderID}}
	err := userMongo.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(&user)
	//log.Println("-----------------------")
	//log.Println(email)
	if err != nil {
		log.Println(err)
		return false
	}
	log.Println(user)
	return true
}

func (userMongo UserRepositoryMongo) GetMyUser(userID string) (model.User, error) {
	var user model.User
	filter := bson.D{{"_id", userID}}
	err := userMongo.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(user)
	if err != nil {
		log.Println(err)
	}
	return user, err
}

func (userMongo UserRepositoryMongo) GetUserEvent(userID string) (model.UserEvent, error) {
	var user model.UserEvent
	log.Println(userID)
	filter := bson.D{{"_id", userID}}
	err := userMongo.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(user)
	if err != nil {
		log.Println(err)
	}
	return user, err
}

func (userMongo UserRepositoryMongo) GetUser(userID string) (model.User, error) {
	var user model.User
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
	}
	filter := bson.M{"_id": id}
	err2 := userMongo.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(&user)
	if err2 != nil {
		log.Println(err2)
	}
	return user, err
}

// AddUserEP api add user from email and password
func (userMongo UserRepositoryMongo) AddUserEP(user model.UserMail) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	if user.Role == "" {
		user.Role = config.MEMBER
	}
	_, err := userMongo.ConnectionDB.Collection(userConlection).InsertOne(context.TODO(), user)
	return err
}

// AddUserPD api add user from social provider
func (userMongo UserRepositoryMongo) AddUserPD(user model.UserProvider) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	if user.Role == "" {
		user.Role = config.MEMBER
	}

	_, err := userMongo.ConnectionDB.Collection(userConlection).InsertOne(context.TODO(), user)
	return err
}

// AddUserAdmin api add user admin
func (userMongo UserRepositoryMongo) AddUserAdmin(user model.UserMail) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	if user.Role == "" {
		user.Role = config.ADMIN
	}
	_, err := userMongo.ConnectionDB.Collection(userConlection).InsertOne(context.TODO(), user)
	return err
}

// UpdateProvider api update account login without provider
func (userMongo UserRepositoryMongo) UpdateProvider(provider model.UserProvider) error {
	filter := bson.D{{"email", provider.Email}}
	update := bson.D{{"$set",
		bson.D{
			{"updated_time", time.Now()}, {"provider", provider.Provider}, {"provider_id", provider.ProviderID},
		},
	}}
	//newProvider := bson.M{"$set": bson.M{"provider": provider.Provider}}
	res, err := userMongo.ConnectionDB.Collection(userConlection).UpdateOne(context.TODO(), filter, update)
	log.Println(res)
	return err
}

// EditUser api edit user
func (userMongo UserRepositoryMongo) EditUser(userID string, user model.User) (model.User, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	filter := bson.D{{"_id", objectID}}
	if err != nil {
		log.Println(err)
	}
	user.UpdatedAt = time.Now()
	user.UserID = objectID
	for index := 0; index < len(user.Address); index++ {
		if user.Address[index].ID == primitive.NilObjectID {
			user.Address[index].ID = primitive.NewObjectID()
		}
		user.Address[index].UpdatedAt = time.Now()
		if user.Address[index].CreatedAt.IsZero() {
			user.Address[index].CreatedAt = time.Now()
		}
	}
	isUpsert := true
	clientOptions := options.FindOneAndUpdateOptions{Upsert: &isUpsert}
	update := bson.M{"$set": user}
	var userData model.User
	err = userMongo.ConnectionDB.Collection(userConlection).FindOneAndUpdate(context.TODO(), filter, update, &clientOptions).Decode(&userData)
	//res, err := userMongo.ConnectionDB.Collection(userConlection).UpdateOne(context.TODO(), filter, newName)
	return userData, err
}

// Confirm api for confirm Email
func (userMongo UserRepositoryMongo) Confirm(userID string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	filter := bson.D{{"user_id", objectID}}
	if err != nil {
		log.Println(err)
	}
	data := bson.M{"$set": bson.M{"confirm": true}}
	res, err := userMongo.ConnectionDB.Collection(userConlection).UpdateOne(context.TODO(), filter, data)
	log.Println(res)
	return err
}

func (userMongo UserRepositoryMongo) DeleteUserByID(userID string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	filter := bson.D{{"user_id", objectID}}
	if err != nil {
		log.Println(err)
	}
	res, err := userMongo.ConnectionDB.Collection(userConlection).DeleteOne(context.TODO(), filter)
	log.Println(res)
	return err
}

// Get user on auth
func (userMongo UserRepositoryMongo) Get(userID string) (model.UserAuth, error) {
	var user model.UserAuth
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
	}
	filter := bson.M{"_id": id}
	err2 := userMongo.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(&user)
	if err2 != nil {
		log.Println(err2)
	}
	return user, err
}

//UploadAvatar update image profile
func (userMongo UserRepositoryMongo) UploadAvatar(userID string, path string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	filter := bson.D{{"_id", objectID}}

	updated := bson.M{"$set": bson.M{"avatar": path, "updated_time": time.Now()}}
	res, err := userMongo.ConnectionDB.Collection(userConlection).UpdateOne(context.TODO(), filter, updated)
	if err != nil {
		//log.Fatal(res)
		log.Printf("[info] err %s", res)
		return err
	}

	return nil
}

// AddAddress api add user admin
func (userMongo UserRepositoryMongo) AddAddress(userID string, address model.Address) (model.User, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	filter := bson.D{{"_id", objectID}}
	address.ID = primitive.NewObjectID()
	address.UpdatedAt = time.Now()
	address.CreatedAt = time.Now()
	update := bson.M{"$push": bson.M{"address": address}}

	res, err := userMongo.ConnectionDB.Collection(userConlection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(res)
		// log.Printf("[info] err %s", res)
		// return err
	}
	user, err2 := userMongo.GetUser(userID)
	return user, err2
}

// UpdateAddress api add edit address
func (userMongo UserRepositoryMongo) UpdateAddress(userID string, address model.Address) (model.User, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	filter := bson.D{{"_id", objectID}, {"product._id", address.ID}}
	address.UpdatedAt = time.Now()
	update := bson.M{"$push": bson.M{"address": address}}

	res, err := userMongo.ConnectionDB.Collection(userConlection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(res)
		// log.Printf("[info] err %s", res)
		// return err
	}
	user, err2 := userMongo.GetUser(userID)
	return user, err2
}

//ChangePassword api change password
func (userMongo UserRepositoryMongo) ChangePassword(userID string, oldPassword string, newPassword string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return err
	}
	pass := utils.HashAndSalt([]byte(oldPassword))
	newpass := utils.HashAndSalt([]byte(newPassword))
	filter := bson.D{{"_id", objectID}, {"password", pass}}
	update := bson.M{"$set": bson.M{"password": newpass}}
	result := userMongo.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter)
	if result.Err() != nil {
		log.Println(result.Err())
		return result.Err()
	}
	res, err2 := userMongo.ConnectionDB.Collection(userConlection).UpdateOne(context.TODO(), filter, update)
	log.Printf("[info] err %s", res)
	return err2
}

//ChangePassword api change password
func (userMongo UserRepositoryMongo) UpdatePassword(userID string, newPassword string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return err
	}
	newpass := utils.HashAndSalt([]byte(newPassword))
	filter := bson.D{{"_id", objectID}}
	update := bson.M{"$set": bson.M{"password": newpass}}
	result := userMongo.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter)
	if result.Err() != nil {
		log.Println(result.Err())
		return result.Err()
	}
	res, err2 := userMongo.ConnectionDB.Collection(userConlection).UpdateOne(context.TODO(), filter, update)
	log.Printf("[info] err %s", res)
	return err2
}

//ForgotPassword api change password
func (userMongo UserRepositoryMongo) ForgotPassword(email string,) (model.UserForgot ,error) {
	var user model.UserForgot
	filter := bson.D{{"email", email}}
	err := userMongo.ConnectionDB.Collection(userConlection).FindOne(context.TODO(), filter).Decode(&user)
	return user,err
}