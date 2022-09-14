package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Mail struct {
	ID         primitive.ObjectID `bson:"_id"`
	Mail_id    string             `json:"mail_id" bson:"mail_id"`
	Email      string             `json:"email" bson:"email"`
	Code       string             `json:"code" bson:"code"`
	Type       string             `json:"type" bson:"type"`
	Created_at time.Time          `json:"created_at" bson:"created_at"`
	Updated_at time.Time          `json:"updated_at" bson:"updated_at"`
	Expires_at time.Time          `json:"expires_at" bson:"expires_at"`
}
