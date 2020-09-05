package db

import (
	"context"
	"fmt"
	"log"
	"os"

	redis "github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"thinkdev.app/think/runex/runexapi/config"
	//"thinkdev.app/think/runex/runexapi/logger"
)

const db_host = "mongodb://localhost:27017"

// const db_host = "mongodb://farmme.in.th:27017"

// const db_host = "mongodb://178.128.85.151:27017"

// const db_host = "mongodb://mongodb:27017"
const db_user = "idever"
const db_pass = "idever@987"

var (
	RedisClient *redis.Client
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

func GetDBCollection() (*mongo.Database, error) {

	clientOptions := options.Client().SetAuth(options.Credential{
		AuthSource: "admin", Username: db_user,
		Password: db_pass, PasswordSet: true,
	}).ApplyURI(db_host)
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
	db := client.Database("runex_v2")
	return db, nil
}

func connectDB(ctx context.Context) (*mongo.Database, error) {
	uri := fmt.Sprintf(db_host, db_user, db_pass, db_host, config.DB_NAME)
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to mongo: %v", err)
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("mongo client couldn't connect with background context: %v", err)
	}
	db := client.Database(config.DB_NAME)
	return db, nil
}

func ConnectRedisDB() *redis.Client {

	redisHost := GetEnv("REDIS_HOST", "localhost")
	redisPort := GetEnv("REDIS_PORT", "6379")
	//redisPassword := GetEnv("REDIS_PASSWORD", "__@redis__P@ss")

	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "__@redis__P@ss",
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
