package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetHighestItemID(client *mongo.Client, dbName, collectionName string) (int, error) {
	collection := client.Database(dbName).Collection(collectionName)

	var result struct {
		ItemID int `bson:"itemid"`
	}

	options := options.FindOne().SetSort(bson.D{{"ItemID", -1}})

	err := collection.FindOne(context.TODO(), bson.D{}, options).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Handle case when no documents are found
			return 0, nil
		}
		log.Println("Error retrieving highest itemid:", err)
		return 0, err
	}

	return result.ItemID, nil
}

// TODO: Ändra så att funktionen tar query parametrar istället för hårdkodad data
func Createpassport(ItemN string, OriginN string, itemIDN int) PassPort {
	now := time.Now()
	Passport := PassPort{
		ItemID:       itemIDN + 1,
		ItemName:     ItemN,
		Origin:       OriginN,
		IsNew:        true,
		LinkMadeFrom: []string{"Link1", "Link2"},
		LinkMakes:    []string{"Link3", "Link4"},
		LinkEvents:   []string{},
		CreationDate: now.Format("01-02-2006"),
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
