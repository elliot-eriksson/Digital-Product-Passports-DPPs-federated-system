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

	var filter interface{}
	filter = bson.D{{"ItemID", 2}}
	resultM, resultD, err := queryPassport(client, ctx, "Test", "TestComp1", filter)
	if err != nil {
		panic(err)
	}
	// 1 for Sensitive Passport
	uploadAndUpdateCID(1, resultM, resultD, client)
	// 0 for Non Sensitive Passport
	uploadAndUpdateCID(0, resultM, resultD, client)

}
