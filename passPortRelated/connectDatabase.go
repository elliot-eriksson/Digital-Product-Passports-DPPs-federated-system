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
	var collection string = "LKAB_MainTable"

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://" + username + ":" + password + "@digital-product-passpor.mjd4fio.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}

	coll := client.Database(database).Collection(collection)
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

		insertResult, err := coll.InsertOne(ctx, doc)
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
		//remanufacture

	default:
		fmt.Println("xdd")
	}

}

/*package main

	// Convert the string representation of ObjectID to primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex("your_object_id_as_string")
	if err != nil {
		log.Fatal(err)
	}

	// Define a filter to find a document by its ObjectID
	filter := bson.D{{"_id", objectID}}

	// Perform the find operation to get a single document
	var result Passport
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Println("No matching document found.")
		return
	} else if err != nil {
		log.Fatal(err)
	}

	// Print the retrieved document
	fmt.Printf("ID: %s, ItemID: %d, ItemName: %s, Origin: %s\n", result.ID.Hex(), result.ItemID, result.ItemName, result.Origin)

	// Close the MongoDB connection
	err = client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Disconnected from MongoDB!")
}

}*/
