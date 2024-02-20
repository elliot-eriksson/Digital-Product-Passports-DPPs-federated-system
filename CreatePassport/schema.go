package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

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
	//log.Println("Highest ItemID:", highestItemID)
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
		CreationDate: now.Format("2006-01-02"),
		Reman:        ipnskey,
	}

	//skickar det nyskapade passport till databas
	Coll := client.Database(database).Collection(collection)
	var ctx = context.TODO()
	_, err = Coll.InsertOne(ctx, Passport)
	if err != nil {
		log.Fatal(err)
	}
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
	separator()
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
func LinkMadeFrom(lmArray []string) ([]map[string]interface{}, []string) {
	var CID, inputMore string
	var linkPassport map[string]interface{}
	var LinkMadeFrom []map[string]interface{}
	// var data []map[string]interface{}
	fmt.Println("Press 1 to start entering CIDs for LinkMadeFrom. Press 0 if your product is not made from something: ")
	fmt.Scan(&inputMore)
	for inputMore == "1" {
		fmt.Println("Enter CID (Enter 0 if no more): ")
		fmt.Scan(&CID)
		if CID != "0" {
			linkPassport = passportFromCID(CID)
			lmArray = LinkMakesAppend(lmArray, CID)
			//fmt.Println("\nHär börjar testprint\n", linkPassport, "\nHär slutar testprint\n")
			delete(linkPassport, "LinkMadeFrom")
			delete(linkPassport, "LinkMakes")
			sensetiveCID := fmt.Sprintf("%v", linkPassport["CID_sen"])
			if sensetiveCID != "" {
				pinToIPFS(sensetiveCID)
			}
			LinkMadeFrom = append(LinkMadeFrom, linkPassport)
		} else {
			inputMore = "0"
		}

	}
	// create linkMakes in referenced object
	return LinkMadeFrom, lmArray
}

// TODO: retrieve private key from CA
func LinkMakesPointerUpdate(lmArray []string) {
	for i := range lmArray {
		pkey := "k51qzi5uqu5dk2i4blnf7qwri0gf2he2cdyp10of13aqclrrdklhha1605lu0i" //temporary, key should be retrieved dynamically from CA
		out := addDataToIPNS(pkey, lmArray[i])
		fmt.Println(out)
	}
}

func LinkMakes(alias string) string {
	out := keyGenerator(alias)
	return out
}
func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func LinkMakesAppend(lmArray []string, CID string) []string {
	lmArray = append(lmArray, CID)
	fmt.Println("----------------------->", lmArray)
	return lmArray
}

func passportMenu(client *mongo.Client, database, collection string, lmArray []string) (int, string, []string) {
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
		LinkMadeFrom, lmArray := LinkMadeFrom(lmArray)
		randomName := randSeq(10) //Glöm ej!!!!
		LinkMakesTemp := LinkMakes(randomName)
		sensitiveArray := sensetiveArray()

		LinkMakes := []string{}
		LinkMakes = append(LinkMakes, LinkMakesTemp)
		//sh := shell.NewShell("localhost:5001")
		ipnsKey := keyGenerator("tempAlias")

		//funktionsanrop för att skapa passport.
		//TODO: ska kunna hantera querys senare
		return Createpassport(ItemN, OriginN, client, database, collection, sensitiveArray, LinkMadeFrom, LinkMakes, ipnsKey), randomName, lmArray
	case 2:
		var cid string
		fmt.Println("Enter CID to regenerate QR-Code:")
		fmt.Scan(&cid)
		generateQRCode(cid)
		return 0, "", lmArray

	default:
		fmt.Println("Error")

	}
	return 0, "", lmArray
}
