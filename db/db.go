package db

import (
	"FunNow/url-shortener/constants"
	"log"

	"github.com/go-redis/redis/v8"
)

type Database struct {
	RedisClient *redis.Client
}

var DatabaseClient = ConstructDatabase()

func ConstructDatabase() *Database {
	database := new(Database)
	err := database.BuildConnection()
	if err != nil {
		log.Fatalln(err)
	}
	return database
}

func (database *Database) BuildConnection() error {

	database.RedisClient = redis.NewClient(&redis.Options{
		Addr:     constants.RedisAddr,
		Password: "",
		DB:       0,
	})

	if err := database.RedisClient.Ping(constants.Ctx).Err(); err != nil {
		return err
	}

	return nil

}
