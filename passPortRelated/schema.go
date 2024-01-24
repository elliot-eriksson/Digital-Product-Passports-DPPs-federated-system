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

// hämtar högsta mongodb passport _id
func GetHighestItemID(client *mongo.Client, dbName, collectionName string) (int, error) {
	collection := client.Database(dbName).Collection(collectionName)

	var result struct {
		ItemID int `bson:"itemid"`
	}

	options := options.FindOne().SetSort(bson.D{{"ItemID", -1}})

	err := collection.FindOne(context.TODO(), bson.D{}, options).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			//returnerar _id som 0 om det inte finns några existerande passports
			return 0, nil
		}
		//felhantering om queryn misslyckar att hämta _id
		log.Println("Error retrieving highest itemid:", err)
		return 0, err
	}

	//returnerar högsta _id samt nil om error
	return result.ItemID, nil
}

// TODO: Ändra så att funktionen tar query parametrar istället för hårdkodad data
func Createpassport(ItemN string, OriginN string, client *mongo.Client, database, Collection string, SensitiveArray []int) {
	//funktionsanrop för att hämta det nuvarande högsta mongodb passport _id i databasen
	highestItemID, err := GetHighestItemID(client, database, Collection)
	if err != nil {
		log.Fatal("Error getting highest itemid:", err)
	}
	log.Println("Highest ItemID:", highestItemID)
	now := time.Now()

	//Hämtar PassPort struct i models och ger den värden
	Passport := PassPort{
		ItemID:       highestItemID + 1,
		ItemName:     ItemN,
		Origin:       OriginN,
		IsNew:        true,
		LinkMadeFrom: []string{"Link1", "Link2"}, //Ska matas in länk från IPFS som ska stores
		LinkMakes:    []string{"Link3", "Link4"}, //Samma här gäller det.
		LinkEvents:   []string{},
		Sensitive:    SensitiveArray,
		CreationDate: now.Format("01-02-2006"),
	}

	//skickar det nyskapade passport till databas
	Coll := client.Database(database).Collection(Collection)
	var ctx = context.TODO()
	insertResult, err := Coll.InsertOne(ctx, Passport)
	if err != nil {
		log.Fatal(err)
	}

	//ser till att vi disconnectar från databasen även om ett error skulle förekomma vid insert till databas
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	fmt.Println(insertResult)

}

// Funtion som tar in hårdkodad objectid för tillfället och gör det möjligt att lägga till event som hänt med produkten.
// Behöver lägga till där man hämtar objectid för att välja vilken produkt som det ska uppdateras för
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
