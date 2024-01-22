package ConnectDatabase

import (
	"context"
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectDatabase(database string, collection string, itemID int) (client *mongo.Client, err error) {

	// var username, password string = "TestComp1", "TestComp1"
	// fmt.Println("hej")

	// fmt.Println("Username: ")
	// fmt.Scanln(&username)
	// fmt.Println("Password: ")
	// fmt.Scanln(&password)
	// var database string = "Test"
	// var collection string = username
	// fmt.Println("ItemID: ")
	// fmt.Scanln(&itemID)
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	// serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	// // opts := options.Client().ApplyURI("mongodb+srv://" + username + ":" + password + "@digital-product-passpor.mjd4fio.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)
	// // gamla här under
	// opts := options.Client().ApplyURI("mongodb+srv://" + username + ":" + password + "@cluster0.qk8pnen.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)
	// // Create a new client and connect to the server
	// client2, err2 := mongo.Connect(context.TODO(), opts)
	// fmt.Println("client innan", client2)
	// if err2 != nil {
	// 	fmt.Println("hej")
	// 	panic(err2)
	// }
	// defer func() {
	// 	if err = client2.Disconnect(context.TODO()); err2 != nil {
	// 		panic(err2)
	// 	}
	// }()
	// return client2, err2

	// var keyPhrase string = "hej"
	// fmt.Println("-------------string format...............")
	// fmt.Println(string(encryptIt([]byte(jsonData), keyPhrase)))
	// fmt.Println("-------------Decrypt format...............")
	// fmt.Println(string(decryptIt(encryptIt([]byte(jsonData), keyPhrase), keyPhrase)))

	// // Send a ping to confirm a successful connection
	// if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
	//  panic(err)
	// }
	// fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}

func connectToDatabase2() {
	var username, password string = "TestComp1", "TestComp1"
	var itemID int
	// fmt.Println("Username: ")
	// fmt.Scanln(&username)
	// fmt.Println("Password: ")
	// fmt.Scanln(&password)
	var database string = "Test"
	var collection string = username
	fmt.Println("ItemID: ")
	fmt.Scanln(&itemID)
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	// opts := options.Client().ApplyURI("mongodb+srv://" + username + ":" + password + "@digital-product-passpor.mjd4fio.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)
	// gamla här under
	opts := options.Client().ApplyURI("mongodb+srv://" + username + ":" + password + "@cluster0.qk8pnen.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	print("client", client)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database(database).Collection(collection)
	var result bson.M
	err = coll.FindOne(context.TODO(), bson.D{{"ItemID", itemID}}).Decode(&result)
	//Error handling
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No item with item %s\n", itemID)
	}
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(result, "", "   ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("json %s\n", jsonData)
	fmt.Printf("bson %s\n", result["Origin"])
}

// func splitPassports() {

// }
