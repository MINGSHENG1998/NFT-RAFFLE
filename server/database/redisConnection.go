package database

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

type RedisConnection interface {
	RedisClient() *redis.Client
}

type redisConnectionStruct struct{}

func NewRedisConenction() RedisConnection {
	return &redisConnectionStruct{}
}

func (r *redisConnectionStruct) RedisClient() *redis.Client {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))
	err := godotenv.Load(string(rootPath) + `/.env`)

	if err != nil {
		log.Fatal("Error loading .env file in databaseConnection.go ", err.Error())
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisUsername := os.Getenv("REDIS_USERNAME")
	redisUserPassword := os.Getenv("REDIS_USER_PASSWORD")

	redisAddress := fmt.Sprintf("%s:%s", redisHost, redisPort)

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Username: redisUsername,
		Password: redisUserPassword,
		DB:       0,
	})

	return client
}
