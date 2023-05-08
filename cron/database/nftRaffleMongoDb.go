package database

import (
	"context"
	"nft-raffle-cron/logger"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	nftRaffleMongoDb     *NftRaffleMongoDb
	nftRaffleMongoDbOnce sync.Once

	nftRaffleMongoDbClient     *mongo.Client
	nftRaffleMongoDbClientOnce sync.Once
)

type NftRaffleMongoDb struct{}

const projectDirName = "cron"

func GetNftRaffleMongoDb() *NftRaffleMongoDb {
	if nftRaffleMongoDb == nil {
		nftRaffleMongoDbOnce.Do(func() {
			nftRaffleMongoDb = &NftRaffleMongoDb{}
		})
	}
	return nftRaffleMongoDb
}

func (db *NftRaffleMongoDb) GetClient() *mongo.Client {
	if nftRaffleMongoDbClient == nil {
		nftRaffleMongoDbClientOnce.Do(func() {
			projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
			currentWorkDirectory, _ := os.Getwd()
			rootPath := projectName.Find([]byte(currentWorkDirectory))
			err := godotenv.Load(string(rootPath) + `/.env`)

			if err != nil {
				logger.Logger.Fatal("Error loading .env file in databaseConnection.go " + err.Error())
			}

			MongoDb := os.Getenv("MONGODB_URL")

			// MongoDB connection
			client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))

			if err != nil {
				logger.Logger.Fatal(err.Error())
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			err = client.Connect(ctx)

			if err != nil {
				logger.Logger.Fatal(err.Error())
			}

			logger.Logger.Info("Connected to MongoDB!!!")

			nftRaffleMongoDbClient = client
		})
	}
	return nftRaffleMongoDbClient
}

func (db *NftRaffleMongoDb) OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("nft_raffle_db").Collection(collectionName)
	return collection
}
