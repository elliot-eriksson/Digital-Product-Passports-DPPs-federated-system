package main

import (
	"context"
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// connectToDatabase2()
	var database string
	var collection string
	var itemID int
	// fmt.Println("Database: ")
	// fmt.Scanln(&database)
	// fmt.Println("Collection: ")
	// fmt.Scanln(&collection)
	fmt.Println("ItemID: ")
	fmt.Scanln(&itemID)
	database = "Test"
	collection = "TestComp1"

	client, err := connectDatabase(database, collection, itemID)
	fmt.Println("client efter", client)

	coll := client.Database(database).Collection(collection)
	var result bson.M
	err = coll.FindOne(context.TODO(), bson.D{{"ItemID", itemID}}).Decode(&result)
	//Error handling
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No item with item %s\n", itemID)
	}
	if err != nil {
		fmt.Printf("hej2")
		panic(err)

	}
	jsonData, err := json.MarshalIndent(result, "", "   ")
	if err != nil {
		fmt.Printf("hej3")
		panic(err)
	}
	fmt.Printf("json %s\n", jsonData)
	fmt.Printf("bson %s\n", result["Origin"])

}
