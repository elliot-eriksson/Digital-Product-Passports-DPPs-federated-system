package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Definerar en struct för en produkt. ID är mongoDBs _id

type PassPort struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ItemID       int                `bson:"ItemID"`
	ItemName     string             `bson:"ItemName"`
	Origin       string             `bson:"Origin"`
	IsNew        bool               `bson:"IsNew"`
	LinkMadeFrom []string           `bson:"LinkMadeFrom"`
	LinkMakes    []string           `bson:"LinkMakes"`
	LinkEvents   []string           `bson:"LinkEvents"`
	CreationDate string             `bson:"CreationDate"`
}
