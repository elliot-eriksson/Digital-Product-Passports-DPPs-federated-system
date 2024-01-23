package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var username, password string = "passAdmin", "passAdmin"

	var database string = "LKAB_DB"
	var Collection string = "LKAB_MainTable"

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://" + username + ":" + password + "@digital-product-passpor.mjd4fio.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database(database)
	fmt.Println(db)
	Coll := db.Collection(Collection)
	fmt.Println(Coll)
	var i int
	fmt.Println("What do you want to do? 1: Createpassport, 2: Remanufacture events for passports")
	fmt.Scan(&i)
	switch i {
	case 1:

		fmt.Println("Enter item name : ")
		var ItemN string
		fmt.Scan(&ItemN)
		fmt.Println("Enter item origin : ")
		var OriginN string
		fmt.Scan(&OriginN)

		highestItemID, err := GetHighestItemID(client, database, Collection)
		if err != nil {
			log.Fatal("Error getting highest itemid:", err)
		}

		log.Println("Highest ItemID:", highestItemID)

		doc := Createpassport(ItemN, OriginN, highestItemID)
		fmt.Println(doc)
		var ctx = context.TODO()

		insertResult, err := Coll.InsertOne(ctx, doc)
		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			if err = client.Disconnect(context.TODO()); err != nil {
				panic(err)
			}
		}()
		fmt.Println(insertResult)

		// skapa string för att lägga kommentar för vad som hänt
		// passera objectid för specifik product
	/*case 2:
	fmt.Println("Enter what has been updated on this certain product :")
	var RemanEvent string
	fmt.Scan(&RemanEvent)
	RemanufactureEvent(Coll, "65ae752266c6e05505af226c", RemanEvent)
	*/
	default:
		fmt.Println("xdd")

	}

}
