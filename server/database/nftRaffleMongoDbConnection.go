package database

import (
	"context"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const projectDirName = "server"

type NftRaffleMongoDbConnection interface {
	OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection
	DBClient() *mongo.Client
}

type nftRaffleMongoDbConnectionStruct struct{}

func NewNftRaffleMongoDbConnection() NftRaffleMongoDbConnection {
	return &nftRaffleMongoDbConnectionStruct{}
}

func (n *nftRaffleMongoDbConnectionStruct) DBClient() *mongo.Client {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))
	err := godotenv.Load(string(rootPath) + `/.env`)

	if err != nil {
		log.Fatal("Error loading .env file in databaseConnection.go ", err.Error())
	}

	MongoDb := os.Getenv("MONGODB_URL")

	// MongoDB connection
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	defer cancel()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!!!")

	return client
}

func (n *nftRaffleMongoDbConnectionStruct) OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("nft_raffle_db").Collection(collectionName)
	return collection
}
