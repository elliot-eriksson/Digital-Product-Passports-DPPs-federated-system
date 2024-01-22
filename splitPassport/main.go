package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// This is a user defined method to close resources.
// This method closes mongoDB connection and cancel context.
func close(client *mongo.Client, ctx context.Context) {

	// CancelFunc to cancel to context

	// client provides a method to close
	// a mongoDB connection.
	defer func() {

		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

// This is a user defined method that returns mongo.Client,
// context.Context, context.CancelFunc and error.
// mongo.Client will be used for further database operation.
// context.Context will be used set deadlines for process.
// context.CancelFunc will be used to cancel context and
// resource associated with it.

func connect(uri string) (*mongo.Client, context.Context,
	context.CancelFunc, error) {

	// ctx will be used to set deadline for process, here
	// deadline will of 30 seconds.
	ctx, cancel := context.WithTimeout(context.Background(),
		30*time.Second)

	// mongo.Connect return mongo.Client method
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

// This is a user defined method that accepts
// mongo.Client and context.Context
// This method used to ping the mongoDB, return error if any.
func ping(client *mongo.Client, ctx context.Context) error {

	// mongo.Client has Ping to ping mongoDB, deadline of
	// the Ping method will be determined by cxt
	// Ping method return error if any occurred, then
	// the error can be handled.
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("connected successfully")
	return nil
}

// query is user defined method used to query MongoDB,
// that accepts mongo.client,context, database name,
// collection name, a query and field.

//  database name and collection name is of type
// string. query is of type interface.
// field is of type interface, which limits
// the field being returned.

// query method returns a cursor and error.
func query(client *mongo.Client, ctx context.Context,
	dataBase, col string, query, field interface{}) (resultM bson.M, resultD bson.D, err error) {

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

func main() {
	var username, password string = "TestComp1", "TestComp1"
	// var database string = "Test"
	// var collection string = username
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://" + username + ":" + password + "@cluster0.qk8pnen.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	ctx := context.TODO()
	// Get Client, Context, CancelFunc and
	// err from connect method.
	// client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	// Release resource when the main
	// function is returned.
	defer close(client, ctx)

	// create a filter an option of type interface,
	// that stores bjson objects.
	var filter, option interface{}

	// filter  gets all document,
	// with maths field greater that 70
	filter = bson.D{{"ItemID", 2}}
	// filter = bson.D{{"ItemID", bson.D{{"$gt", 1}}}}
	// filter = bson.D{{"itemID", bson.D{{"$gt", 70}}},
	// }

	//  option remove id field from all documents
	option = bson.D{{"_id", 0}}

	// call the query method with client, context,
	// database name, collection  name, filter and option
	// This method returns momngo.cursor and error if any.
	resultM, resultD, err := query(client, ctx, "Test", "TestComp1", filter, option)
	// handle the errors.
	if err != nil {
		panic(err)
	}

	fmt.Println(resultD)
	fmt.Println(resultM)

	test3 := fmt.Sprintf("%v", resultM["isSensitive"])

	fmt.Println(test3)

	// var results []bson.M

	// to get bson object  from cursor,
	// returns error if any.
	// if err := cursor.All(ctx, &results); err != nil {

	// 	// handle the error
	// 	panic(err)
	// }

	// jsonData, err := json.MarshalIndent(results, "", "   ")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("json %s\n", jsonData)
	// fmt.Printf("bson %s\n", results["Origin"])

	// printing the result of query.
	// fmt.Println("Query Result")
	// var test primitive.M
	// for _, doc := range results {
	// 	fmt.Println(doc)
	// 	test = doc
	// 	fmt.Println("----------->result: ", results)
	// }

	// fmt.Printf("bson %s\n", test["isSensitive"])
	// Ping mongoDB with Ping method

	inputString := test3

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

	fmt.Println(intArray)
	fmt.Println(reflect.TypeOf(intArray))
	var isSensitivePositionArray []int

	for i, value := range intArray {
		if value == 1 {
			isSensitivePositionArray = append(isSensitivePositionArray, i)
		}
	}
	fmt.Println(isSensitivePositionArray)

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
