package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var username, password string = "TestComp1", "TestComp1"
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://" + username + ":" + password + "@cluster0.qk8pnen.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	ctx := context.TODO()
	if err != nil {
		panic(err)
	}

	// Release resource when the main function is returned
	defer close(client, ctx)

	// create a filter of type interface,
	// that stores bson objects.
	var filter interface{}
	var update interface{}
	filter = bson.D{{"ItemID", 2}}

	// Queries database to retrieve passports
	// This method returns two bson values and error if any.
	resultM, resultD, err := queryPassport(client, ctx, "Test", "TestComp1", filter)
	if err != nil {
		panic(err)
	}

	// 1 for Sensitive Passport
	uploadAndUpdateCID(1, resultM, resultD, filter, update, client)
	// 0 for Non Sensitive Passport
	uploadAndUpdateCID(0, resultM, resultD, filter, update, client)

	ping(client, ctx)

}
