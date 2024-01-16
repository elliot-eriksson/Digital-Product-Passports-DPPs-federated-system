package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Define a struct to represent the class

type PassPort struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ItemID       int                `bson:"ItemID"`
	ItemName     string             `bson:"ItemName"`
	Origin       string             `bson:"Origin"`
	IsNew        bool               `bson:"IsNew"`
	LinkMadeFrom []string           `bson:"LinkMadeFrom"`
	LinkMakes    []string           `bson:"LinkMakes"`
	CreationDate string             `bson:"CreationDate"`
}
