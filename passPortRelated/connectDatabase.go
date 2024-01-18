package main

//passAdmin passAdmin
//
import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var username, password string = "passAdmin", "passAdmin"
	var itemID int
	// fmt.Println("Username: ")
	// fmt.Scanln(&username)
	// fmt.Println("Password: ")
	// fmt.Scanln(&password)
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

		fmt.Println("ItemID: ")
		fmt.Scanln(&itemID)
		doc := Createpassport()
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
	case 2:
		RemanufactureEvent(Coll, "65a8f9c45a1a8a3ddf32c503", "den e cool naijs")
	default:
		fmt.Println("xdd")
	}

}
