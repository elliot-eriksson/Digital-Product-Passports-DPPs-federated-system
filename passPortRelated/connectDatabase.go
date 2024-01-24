package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

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
	// Skapar en client och koppling till servern
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}

	Coll := client.Database(database).Collection(Collection)

	//temporär input för test ändamål, ska ändras framöver för att kunna göras via hemsida/program etc
	var i int
	fmt.Println("What do you want to do? 1: Createpassport, 2: Remanufacture events for passports")
	fmt.Scan(&i)
	switch i {
	case 1:

		//testinput av item name samt item origin
		var ItemN, OriginN string
		fmt.Println("Enter item name : ")
		fmt.Scan(&ItemN)
		fmt.Println("Enter item origin : ")
		fmt.Scan(&OriginN)
		sensitiveArray := []int{1, 1, 1, 1, 1, 1, 1, 1, 1}

		//funktionsanrop för att skapa passport.
		//TODO: ska kunna hantera querys senare
		Createpassport(ItemN, OriginN, client, database, Collection, sensitiveArray)
	case 2:

		//testinput för att lägga till ett remanufacture event till en produkt
		fmt.Println("Enter what has been updated on this certain product:")
		fmt.Scan("")
		reader := bufio.NewReader(os.Stdin)
		RemanEvent, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		RemanEvent = RemanEvent[:len(RemanEvent)-1]

		//funktionsanrop för att uppdatera en produkt med ett remanufacture event
		//TODO: andra variabeln som skickas med i funktionen måste bytas ut med en dynamisk variabel "objectid" senare, är hårdkodad för nuvarandet med ett _id
		//TODO: ska kunna hantera querys
		RemanufactureEvent(Coll, "65b103b5c0ba3fc65303b998", RemanEvent)

	default:
		fmt.Println("xdd")

	}

}
