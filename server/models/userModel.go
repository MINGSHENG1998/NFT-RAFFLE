package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	First_name    string             `json:"first_name" bson:"first_name" validate:"required,min=2,max=30"`
	Last_name     string             `json:"last_name" bson:"last_name" validate:"required,min=2,max=30"`
	Password      string             `json:"password" bson:"password" validate:"required"`
	Email         string             `json:"email" bson:"email" validate:"required"`
	Phone         string             `json:"phone" bson:"phone" validate:"required"`
	Access_token  string             `json:"access_token" bson:"access_token"`
	Refresh_token string             `json:"refresh_token" bson:"refresh_token"`
	User_role     string             `json:"user_role" bson:"user_role" validate:"required,eq=ADMIN|eq=USER"`
	Created_at    time.Time          `json:"created_at" bson:"created_at"`
	Updated_at    time.Time          `json:"updated_at" bson:"updated_at"`
	User_id       string             `json:"user_id" bson:"user_id"`
}
