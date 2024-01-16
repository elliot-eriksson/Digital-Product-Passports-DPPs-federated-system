package main

import (
	"gopkg.in/mgo.v2/bson"
)

// Define a struct to represent the class

type PassPort struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	ItemID       int           `bson:"ItemID"`
	ItemName     string        `bson:"ItemName"`
	Origin       string        `bson:"Origin"`
	IsNew        bool          `bson:"IsNew"`
	LinkMadeFrom []string      `bson:"LinkMadeFrom"`
	LinkMakes    []string      `bson:"LinkMakes"`
	CreationDate string        `bson:"CreationDate"`
}
