package services

import (
	"context"
	"log"
	"nft-raffle/database"
	"nft-raffle/enums"
	"nft-raffle/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type VerificationMailService interface {
	CreateNewVerifcationMail(email string, randomSixDigits string, expires_at time.Time) error
	UpdateVerificationEmail(email string, randomSixDigits string, expires_at time.Time) error
}

type verificationMailServiceStruct struct{}

var mailCollection *mongo.Collection = database.OpenCollection(database.Client, "mail")

func NewVerificationMailService() VerificationMailService {
	return &verificationMailServiceStruct{}
}

func (v *verificationMailServiceStruct) CreateNewVerifcationMail(email string, randomSixDigits string, expires_at time.Time) error {
	var mail models.Mail
	var err error
	mail.ID = primitive.NewObjectID()
	mail.Mail_id = mail.ID.Hex()
	mail.Email = email

	mail.Code = randomSixDigits
	mail.Type = enums.MailVerification.String()

	mail.Created_at, err = time.Parse(time.RFC3339, time.Now().Local().Format(time.RFC3339))

	if err != nil {
		log.Println(err)
		return err
	}

	mail.Updated_at, err = time.Parse(time.RFC3339, time.Now().Local().Format(time.RFC3339))

	if err != nil {
		log.Println(err)
		return err
	}

	mail.Expires_at = expires_at

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	_, insertError := mailCollection.InsertOne(ctx, mail)
	defer cancel()

	if insertError != nil {
		log.Println(insertError)
		return insertError
	}

	return nil
}

func (v *verificationMailServiceStruct) UpdateVerificationEmail(email string, randomSixDigits string, expires_at time.Time) error {
	var updateObj bson.D

	updateObj = append(updateObj, bson.E{Key: "code", Value: randomSixDigits})
	updateObj = append(updateObj, bson.E{Key: "expires_at", Value: expires_at})

	Updated_at, err := time.Parse(time.RFC3339, time.Now().Local().Format(time.RFC3339))

	if err != nil {
		return err
	}

	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at})

	upsert := true
	filter := bson.D{
		{Key: "email", Value: email},
		{Key: "type", Value: enums.MailVerification.String()},
	}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	_, updateError := mailCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)
	defer cancel()

	if updateError != nil {
		return updateError
	}

	return nil
}
