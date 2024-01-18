package main

import (
	"context"
	"fmt"
	"log"

	//"log"
	"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: Ändra så att funktionen tar query parametrar istället för hårdkodad data
func Createpassport() PassPort {
	Passport := PassPort{
		ItemName:     "Your Item Name",
		Origin:       "Your Origin",
		IsNew:        true,
		LinkMadeFrom: []string{"Link1", "Link2"},
		LinkMakes:    []string{"Link3", "Link4"},
		LinkEvents:   []string{"Link5", "Link6"},
		CreationDate: "2024-01-16",
	}
	return Passport
}

func RemanufactureEvent(Coll *mongo.Collection, mongoid string, RemanEvent string) {
	id, _ := primitive.ObjectIDFromHex(mongoid)
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "LinkEvents", Value: RemanEvent}}}}

	result, err := Coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Documents matched: %v\n", result.MatchedCount)
	fmt.Printf("Documents updated: %v\n", result.ModifiedCount)
}
