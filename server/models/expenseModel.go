package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Expense struct {
	ID             primitive.ObjectID `bson:"_id"`
	Expense_id     string             `json:"expense_id" bson:"expense_id"`
	User_id        string             `json:"user_id" bson:"user_id"`
	Expense_label  string             `json:"expense_label" bson:"expense_label"`
	Expense_type   string             `json:"expense_type" bson:"expense_type"`
	Expense_amount int64              `json:"expense_amount" bson:"expense_amount"`
	Expense_time   time.Time          `json:"expense_time" bson:"expense_time"`
	Created_at     time.Time          `json:"created_at" bson:"created_at"`
	Updated_at     time.Time          `json:"updated_at" bson:"updated_at"`
}
