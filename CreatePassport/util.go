package main

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func separator() {
	fmt.Println("--------------------------")
}

func performChecks(sh *shell.Shell) error {
	if YourLocalPath == "" {
		return fmt.Errorf("'YourLocalPath' constant is NOT defined. Please, provide a local path in your computer where the file will be downloaded")
	}

	if YourPublicKey == "" {
		return fmt.Errorf("'YourPublicKey' constant is NOT defined. Please, provide the public key of your IPFS node")
	}

	if !sh.IsUp() {
		return fmt.Errorf("You do not have an IPFS node running at port 5001")
	}

	return nil
}

// query function that retrieves a passport and returns the result as two differtent datatypes.
func queryPassport(client *mongo.Client, ctx context.Context, dataBase, col string, query interface{}) (resultM bson.M, resultD bson.D, err error) {

	collection := client.Database(dataBase).Collection(col)

	err = collection.FindOne(ctx, query).Decode(&resultM)
	err = collection.FindOne(ctx, query).Decode(&resultD)
	return
}

// Takes from ipfs and reformats it to jsonFormating
func jsonFormat(data primitive.A) (newJsonData string) {
	jsonData, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		panic(err)
	}
	newJsonData = strings.Replace(string(jsonData), "\"Key\": ", "", -1)
	newJsonData = strings.Replace(newJsonData, ",\n      \"Value\"", "", -1)
	newJsonData = strings.Replace(newJsonData, "[", "", 1)
	newJsonData = strings.Replace(newJsonData, "\n   },\n   {", ",", -1)
	sz := len(newJsonData)
	newJsonData = newJsonData[:sz-1]
	return newJsonData
}

func updateDatabase(client *mongo.Client, ctx context.Context, dataBase, col string, filter interface{}, update interface{}) {
	// update database with the data from the update by using filter to choose what to update
	coll := client.Database(dataBase).Collection(col)
	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	fmt.Println("CID Upload successful", result)

}

// Split passport into two seperate parts of either sensitive or non sensitive
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

	var splitSensitiveAndNonSensitive []int
	for i, value := range intArray {
		//Takes out either 0 or 1 depending if we want the sensitive or non-sensitive fields
		if value == sensetive {
			splitSensitiveAndNonSensitive = append(splitSensitiveAndNonSensitive, i)
		}
	}

	// Array with the values of either sensitive or nonsensitive
	var sensitiveArray bson.A
	fmt.Println("sensitiveArray 1 : ", sensitiveArray)
	for _, value := range splitSensitiveAndNonSensitive {
		sensitiveArray = append(sensitiveArray, resultD[value])
	}
	return sensitiveArray
}

// Uploads the passport to IPFS and calls updateDatabase(), where the filter is the CID
func uploadAndUpdateCID(intSensitive int, resultM primitive.M, resultD primitive.D, client *mongo.Client, database, collection string) {
	sensitiveArray := getSensitiveData(fmt.Sprintf("%v", resultM["Sensitive"]), resultD, intSensitive)
	jsonData := jsonFormat(sensitiveArray)
	fmt.Printf("json %s\n", jsonData)

	var filter primitive.D
	var update interface{}
	if json.Valid([]byte(jsonData)) {

		var upploadString string = string(encryptIt([]byte(jsonData), "hej"))
		cid, err := ipfs(upploadString)
		if err != nil {
			panic(err)
		}

		if intSensitive == 0 {
			filter = bson.D{{"_id", resultM["_id"]}}
			update = bson.D{{"$set", bson.D{{"CID", cid}}}}
		} else {
			filter = bson.D{{"_id", resultM["_id"]}}
			update = bson.D{{"$set", bson.D{{"CID_sen", cid}}}}
		}

		updateDatabase(client, context.TODO(), database, collection, filter, update)
	} else {
		fmt.Println("Not json valid")
	}
}
