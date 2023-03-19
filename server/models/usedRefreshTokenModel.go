package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UsedRefreshToken struct {
	ID              primitive.ObjectID `bson:"_id"`
	Token_id        string             `json:"token_id" bson:"token_id"`
	Refresh_token   string             `json:"refresh_token" bson:"refresh_token"`
	Issued_at_unix  int64              `json:"issued_at_unix" bson:"issued_at_unix"`
	Expired_at_unix int64              `json:"expired_at_unix" bson:"expired_at_unix"`
}
