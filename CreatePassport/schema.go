package main

import (
	"context"
	"fmt"
	"log"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
	"go.mongodb.org/mongo-driver/bson"
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
func Createpassport(ItemN string, OriginN string, client *mongo.Client, database, collection string, SensitiveArray []string, LinkMadeFromN []map[string]interface{}, LinkMakesN []string, ipnskey string) (itemID int) {
	//funktionsanrop för att hämta det nuvarande högsta mongodb passport _id i databasen
	highestItemID, err := GetHighestItemID(client, database, collection)
	if err != nil {
		log.Fatal("Error getting highest itemid:", err)
	}
	log.Println("Highest ItemID:", highestItemID)
	now := time.Now()
	newItemID := highestItemID + 1

	//Hämtar PassPort struct i models och ger den värden
	Passport := PassPort{
		ItemID:       newItemID,
		ItemName:     ItemN,
		Origin:       OriginN,
		LinkMadeFrom: LinkMadeFromN, //Ska matas in länk från IPFS som ska stores
		LinkMakes:    LinkMakesN,    //Samma här gäller det.
		Sensitive:    SensitiveArray,
		CreationDate: now.Format("01-02-2006"),
		Reman:        ipnskey,
	}

	//skickar det nyskapade passport till databas
	Coll := client.Database(database).Collection(collection)
	var ctx = context.TODO()
	insertResult, err := Coll.InsertOne(ctx, Passport)
	if err != nil {
		log.Fatal(err)
	}

	//ser till att vi disconnectar från databasen även om ett error skulle förekomma vid insert till databas
	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()
	fmt.Println(insertResult)
	return newItemID

}

// Creates a way to split data from database to implement diffrent access levels with different encryption
// TODO Need to change this to be fully dynamic with new rows of information
func sensetiveArray() (sensitiveArray []string) {
	var input string
	sensitiveArray = []string{"0", "0", "0", "0"}
	fmt.Println("Enter sensetivite value 0 not sensetive : 1 sensetive")
	fmt.Print("LinkMadeFrom: ")
	fmt.Scan(&input)
	for input != "0" && input != "1" {
		fmt.Println("Input must be 0 or 1")
		fmt.Print("LinkMadeFrom: ")
		fmt.Scan(&input)
	}
	sensitiveArray = append(sensitiveArray, input)
	fmt.Print("LinkMakes: ")
	fmt.Scan(&input)
	for input != "0" && input != "1" {
		fmt.Println("Input must be 0 or 1")
		fmt.Print("LinkMakes: ")
		fmt.Scan(&input)
	}
	sensitiveArray = append(sensitiveArray, input)
	// Sensetive
	sensitiveArray = append(sensitiveArray, "1")
	// CreationDate
	sensitiveArray = append(sensitiveArray, "0")
	// CID_sen
	sensitiveArray = append(sensitiveArray, "0")
	// Reman events special
	sensitiveArray = append(sensitiveArray, "0")
	return sensitiveArray
}

// Retrieves the passport infromation from the CID/products the passport is created from
// TODO Needs to retrieve the key the given CID is encrypted with
func LinkMadeFrom() (LinkMadeFrom []map[string]interface{}) {
	var CID, inputMore string
	var linkPassport map[string]interface{}
	// var data []map[string]interface{}
	fmt.Println("Press 1 to start entering CIDs for LinkMadeFrom: ")
	fmt.Scan(&inputMore)
	for inputMore == "1" {
		fmt.Println("Enter CID (Enter 0 if no more): ")
		fmt.Scan(&CID)
		if CID != "0" {
			linkPassport = passportFromCID(CID)
			sensetiveCID := fmt.Sprintf("%v", linkPassport["CID_sen"])
			fmt.Println("CID_SEN                        ", sensetiveCID)
			if sensetiveCID != "" {
				pinToIPFS(sensetiveCID)
			}
			LinkMadeFrom = append(LinkMadeFrom, linkPassport)
		} else {
			inputMore = "0"
		}
	}
	return LinkMadeFrom
}

func passportMenu(client *mongo.Client, database, collection string) (itemID int) {
	//temporär input för test ändamål, ska ändras framöver för att kunna göras via hemsida/program etc
	var i int
	fmt.Println("What do you want to do? 1: Createpassport, 2: Regenerate QR code")
	fmt.Scan(&i)
	switch i {
	case 1:

		//testinput av item name samt item origin
		var ItemN, OriginN string
		fmt.Println("Enter item name : ")
		fmt.Scan(&ItemN)
		fmt.Println("Enter item origin : ")
		fmt.Scan(&OriginN)
		// var LinkMadeFrom []string
		LinkMadeFrom := LinkMadeFrom()
		sensitiveArray := sensetiveArray()

		LinkMakes := []string{}
		sh := shell.NewShell("localhost:5001")
		ipnsKey := keyGenerator(sh, "tempAlias")

		//funktionsanrop för att skapa passport.
		//TODO: ska kunna hantera querys senare
		return Createpassport(ItemN, OriginN, client, database, collection, sensitiveArray, LinkMadeFrom, LinkMakes, ipnsKey)
	case 2:
		var cid string
		fmt.Println("Enter CID to regenerate QR-Code:")
		fmt.Scan(&cid)
		generateQRCode(cid)
		return 0

	default:
		fmt.Println("Error")

	}
	return 0
}
