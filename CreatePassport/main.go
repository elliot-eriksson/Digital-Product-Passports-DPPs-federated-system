package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var username, password string = "passAdmin", "passAdmin"
	var i int
	var database string
	var collection string
	var cid string
	var lmArray []string

	// Lets you chose witch of the three test database create to uppdate.
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
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	//funktionsanrop för passport meny. Presenterar en med 2 stycken val just nu. Antingen skapa ett nytt passport eller återskapa en qr-kod
	itemID, randomName, lmArray := passportMenu(client, database, collection, lmArray)
	if itemID != 0 {
		var filter interface{}
		filter = bson.D{{"ItemID", itemID}}
		resultM, resultD, err := queryPassport(client, ctx, database, collection, filter)
		if err != nil {
			panic(err)
		}

		// 1 for Sensitive Passport
		cid = uploadAndUpdateCID(1, resultM, resultD, client, database, collection)
		// 0 for Non Sensitive Passport
		resultM2, resultD2, err := queryPassport(client, ctx, database, collection, filter)
		cid = uploadAndUpdateCID(0, resultM2, resultD2, client, database, collection)
		keyRename(cid)
		keyRenameLinkMakes(cid, randomName)
		//generateQRCode(cid)
		fmt.Println(lmArray)
		if len(lmArray) > 0 {
			LinkMakesPointerUpdate(lmArray)
		}
	}
}
