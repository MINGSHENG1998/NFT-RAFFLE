package database

import (
	"fmt"
	"nft-raffle/logger"
	"os"
	"regexp"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var RedisClient *redis.Client = RedisClientInstance()

func RedisClientInstance() *redis.Client {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))
	err := godotenv.Load(string(rootPath) + `/.env`)

	if err != nil {
		logger.Logger.Fatal("Error loading .env file in databaseConnection.go " + err.Error())
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
