package db

import (
	"context"
	"fmt"
	"log"
	"os"

	redis "github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"thinkdev.app/think/runex/runexapi/config"
	//"thinkdev.app/think/runex/runexapi/logger"
)

// const db_host = "mongodb://localhost:27017"

// const db_host = "mongodb://farmme.in.th:27017"

// const db_host = "mongodb://178.128.85.151:27017"

// const db_host = "mongodb://mongodb:27017"
// const db_user = "idever"
// const db_pass = "idever@987"

var (
	// RedisClient redis variable
	RedisClient *redis.Client
	// DB variable
	DB *mongo.Database
)

// GetEnv accepts the ENV as key and a default string
// If the lookup returns false then it uses the default string else it leverages the value set in ENV variable
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	//logger.Logger.Info("Setting default values for ENV variable " + key)
	return fallback
}

// GetDBCollection mongo db
func GetDBCollection() (*mongo.Database, error) {
	host := viper.GetString("mongodb.connection")
	user := viper.GetString("mongodb.user")
	pass := viper.GetString("mongodb.pass")
	dbName := viper.GetString("mongodb.db")
	clientOptions := options.Client().SetAuth(options.Credential{
		AuthSource: "admin", Username: user,
		Password: pass, PasswordSet: true,
	}).ApplyURI(host)
	//clientOptions := options.Client().ApplyURI(db_host)
	//clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Print("can't connect database!!!")
		return nil, err
	}
	DB = client.Database(dbName)
	return DB, nil
}

func connectDB(ctx context.Context) (*mongo.Database, error) {
	host := viper.GetString("mongodb.connection")
	user := viper.GetString("mongodb.user")
	pass := viper.GetString("mongodb.pass")
	dbName := viper.GetString("mongodb.db")
	uri := fmt.Sprintf(host, user, pass, host, dbName)
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to mongo: %v", err)
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("mongo client couldn't connect with background context: %v", err)
	}
	DB = client.Database(dbName)
	return DB, nil
}

// ConnectRedisDB connect to redis
func ConnectRedisDB() *redis.Client {

	redisHost := viper.GetString("redis.host")
	redisPort := viper.GetString("redis.port")
	redisPassword := viper.GetString("redis.pass")

	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	_, err := RedisClient.Ping(context.TODO()).Result()
	//logger.Logger.Infof("Reply from Redis %s", pong)
	if err != nil {
		fmt.Println(err.Error())
		//logger.Logger.Fatalf("Failed connecting to redis db %s", err.Error())
		os.Exit(1)
	}
	//logger.Logger.Infof("Successfully connected to redis database")
	return RedisClient
}

// CloseDB accepst Session as input to close Connection to the database
func CloseDB(s *mongo.Database) {

	err := s.Client().Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}
