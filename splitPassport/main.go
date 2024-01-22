package main

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// query is user defined method used to query MongoDB,
// that accepts mongo.client,context, database name,
// collection name, a query and field.

//  database name and collection name is of type
// string. query is of type interface.
// field is of type interface, which limits
// the field being returned.

// query method returns a cursor and error.
func query(client *mongo.Client, ctx context.Context, dataBase, col string, query interface{}) (resultM bson.M, resultD bson.D, err error) {

	// select database and collection.
	collection := client.Database(dataBase).Collection(col)

	// collection has an method Find,
	// that returns a mongo.cursor
	// based on query and field.
	// result, err = collection.Find(ctx, query,
	// 	options.Find().SetProjection(field))
	// var result bson.M
	err = collection.FindOne(ctx, query).Decode(&resultM)
	err = collection.FindOne(ctx, query).Decode(&resultD)
	return
}

func uploadCID(client *mongo.Client, ctx context.Context, dataBase, col string, cid string, _id interface{}) {
	// upload CID to database by using filter on ObjectID
	coll := client.Database(dataBase).Collection(col)
	filter := bson.D{{"_id", _id}}
	update := bson.D{{"$set", bson.D{{"CID", cid}}}}

	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	fmt.Println("CID Upload successful", result)

}

func getSensitiveData(inputString string, resultD bson.D, sensetive int) (sensitiveArray2 bson.A) {
	re := regexp.MustCompile("\\d+")
	matches := re.FindAllString(inputString, -1)

	var intArray []int
	for _, match := range matches {
		num, err := strconv.Atoi(match)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		intArray = append(intArray, num)
	}

	var isSensitivePositionArray []int
	for i, value := range intArray {
		if value == sensetive {
			isSensitivePositionArray = append(isSensitivePositionArray, i)
		}
	}
	var sensitiveArray bson.A
	for _, value := range isSensitivePositionArray {
		fmt.Println(resultD[value])
		sensitiveArray = append(sensitiveArray, resultD[value])
	}
	return sensitiveArray
}

func main() {
	var username, password string = "TestComp1", ""
	// var database string = "Test"
	// var collection string = username
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://" + username + ":" + password + "@cluster0.qk8pnen.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	ctx := context.TODO()
	if err != nil {
		panic(err)
	}

	// Release resource when the main
	// function is returned.
	defer close(client, ctx)

	// create a filter  of type interface,
	// that stores bjson objects.
	var filter interface{}
	filter = bson.D{{"ItemID", 3}}

	// call the query method with client, context,
	// database name, collection  name, filter and option
	// This method returns two bson values and error if any.
	resultM, resultD, err := query(client, ctx, "Test", "TestComp1", filter)
	// handle the errors.
	if err != nil {
		panic(err)
	}

	// getSensitiveData(fmt.Sprintf("%v", resultM["isSensitive"]), resultD)
	sensitiveArray := getSensitiveData(fmt.Sprintf("%v", resultM["isSensitive"]), resultD, 1)
	// fmt.Println(sensitiveArray)
	jsonData, err := json.MarshalIndent(sensitiveArray, "", "   ")
	// fmt.Printf("json %s\n", jsonData)

	var upploadString string = string(encryptIt([]byte(jsonData), "hej"))
	cid, err := ipfs(upploadString)

	var ObjectID interface{} = resultM["_id"]

	uploadCID(client, ctx, "Test", "TestComp1", cid, ObjectID)
	ping(client, ctx)

}
