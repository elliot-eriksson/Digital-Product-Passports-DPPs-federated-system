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
	var i int
	var database string
	var collection string

	//temporär
	fmt.Println("What database and table do you want? 1: LKAB, 2: SSAB, 3: VOLVO")
	fmt.Scan(&i)
	switch i {
	case 1:

		database = "LKAB_DB"
		collection = "LKAB_MainTable"
	case 2:
		database = "SSAB_DB"
		collection = "SSAB_MainTable"
	case 3:
		database = "VOLVO_DB"
		collection = "VOLVO_MainTable"
	default:
		fmt.Println("incorrect input")
		return
	}

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://" + username + ":" + password + "@digital-product-passpor.mjd4fio.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)
	// Skapar en client och koppling till servern
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}

	//Coll := client.Database(database).collection(collection)

	//funktionsanrop för passport meny. Presenterar en med 2 stycken val just nu. Antingen skapa ett nytt passport eller uppdatera ett passport med remanafacture event
	passportMenu(client, database, collection)

}
