package main

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
	qrcode "github.com/skip2/go-qrcode"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func separator() {
	fmt.Println("---------------------------------------------------------------------------------------------------------------")
}

func performChecks(sh *shell.Shell) error {
	if YourLocalPath == "" {
		return fmt.Errorf("'YourLocalPath' constant is NOT defined. Please, provide a local path in your computer where the file will be downloaded")
	}

	if !sh.IsUp() {
		return fmt.Errorf("You do not have an IPFS node running at port 5001")
	}

	return nil
}

// query function that retrieves a passport and returns the result as two differtent datatypes.
func queryPassport(client *mongo.Client, ctx context.Context, dataBase, col string, query interface{}) (passData map[string]interface{}, err error) {

	collection := client.Database(dataBase).Collection(col)
	err = collection.FindOne(ctx, query).Decode(&passData)
	fmt.Println("QUERYPASSPORT err", err)
	return
}

// Takes from ipfs and reformats it to jsonFormating
func jsonFormat(data primitive.A) (newJsonData string) {
	jsonData, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		fmt.Println("--------------ad")
		panic(err)
	}
	newJsonData = strings.Replace(string(jsonData), "\"Key\": ", "", -1)
	newJsonData = strings.Replace(newJsonData, ",\n      \"Value\"", "", -1)
	newJsonData = strings.Replace(newJsonData, ",\n         \"Value\"", "", -1)
	newJsonData = strings.Replace(newJsonData, ",\n            \"Value\"", "", -1)
	newJsonData = strings.Replace(newJsonData, ",\n               \"Value\"", "", -1)
	newJsonData = strings.Replace(newJsonData, "[", "", 1)
	newJsonData = strings.Replace(newJsonData, "\n   },\n   {", ",", -1)
	sz := len(newJsonData)
	newJsonData = newJsonData[:sz-1]
	return newJsonData
}

func updateDatabase(client *mongo.Client, ctx context.Context, dataBase, col string, filter interface{}, update interface{}) {
	// update database with the data from the update by using filter to choose what to update
	coll := client.Database(dataBase).Collection(col)
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println("updateDatabase")
		panic(err)
	}
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
	for _, value := range splitSensitiveAndNonSensitive {
		sensitiveArray = append(sensitiveArray, resultD[value])
	}
	// fmt.Println("sensitiveArray 1 : ", sensitiveArray)
	return sensitiveArray
}

func selectPassData(sensitiveLevel string, data map[string]interface{}) map[string]interface{} {
	dataArray := data[sensitiveLevel]
	test := fmt.Sprintf("%v", dataArray)
	selectData := make(map[string]interface{})

	fmt.Println("dataArray", test)

	for _, value := range test {
		key := string(value)
		fmt.Println("selectPassData key: ", key)
		selectData[key] = data[key]
	}
	fmt.Println("selectPassData", selectData)
	return selectData

}

// Uploads the passport to IPFS and calls updateDatabase(), where the filter is the CID
func uploadAndUpdateCID(sensitiveLevel string, passData map[string]interface{}, client *mongo.Client, database, collection string) (cid string) {
	fmt.Println("UPLOAD DATA FIRST LINE")
	// sensitiveArray := getSensitiveData(fmt.Sprintf("%v", resultM["Sensitive"]), resultD, intSensitive)
	uploadData := selectPassData(sensitiveLevel, passData)
	fmt.Println("UPLOAD DATA ", uploadData)

	jsonData, err := json.Marshal(uploadData)
	fmt.Println("JSON DATA", jsonData)
	// jsonData, err := json.format(uploadData)
	if err != nil {
		fmt.Println("Error marshalling object to JSON")
		// Handle the case where the sensetiveArray is not found or not of the expected type
	}

	fmt.Println("JSON DATA", jsonData)
	var ItemN string
	fmt.Scan(&ItemN)

	// fmt.Printf("json %s\n", jsonData)
	// var err error
	// var filter primitive.D
	// var update interface{}
	// if json.Valid([]byte(jsonData)) {

	// 	var upploadString string = string(encryptIt([]byte(jsonData), "hej"))
	// 	cid, err = ipfs(upploadString)

	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	if intSensitive == 0 {
	// 		filter = bson.D{{"_id", resultM["_id"]}}
	// 		update = bson.D{{"$set", bson.D{{"CID", cid}}}}
	// 	} else {
	// 		filter = bson.D{{"_id", resultM["_id"]}}
	// 		update = bson.D{{"$set", bson.D{{"CID_sen", cid}}}}
	// 	}

	// 	updateDatabase(client, context.TODO(), database, collection, filter, update)
	// } else {
	// 	fmt.Println("Not json valid")
	// }
	return cid
}

// Will be used for testing without the CA.
func testLinkMakes(cid string) string {
	target := passportFromCID(cid)
	// var newJson2 map[string]interface{}
	res := fmt.Sprintf("%v", target["LinkMakes"])
	return res[1 : len(res)-1]
}

func generateLinkMakesData(cid string) string {
	target := passportFromCID(cid)
	// var newJson2 map[string]interface{}
	newJson2 := make(map[string]interface{})
	newJson2["CID"] = cid
	newJson2["ItemName"] = target["ItemName"]
	newJson2["Origin"] = target["Origin"]
	newJson2["CreationDate"] = target["CreationDate"]
	test, _ := json.Marshal(newJson2)
	newjson := LinkMakesAdd(string(test))
	return newjson
}

func generateQRCode(cid string) {
	target := passportFromCID(cid)

	newjson := "{" + "\n      \"cid\": " + "\"" + cid + "\",\n" +
		"      \"ItemName\": " + "\"" + fmt.Sprintf("%v", target["ItemName"]) + "\",\n" +
		"      \"Origin\": " + "\"" + fmt.Sprintf("%v", target["Origin"]) + "\",\n" +
		"      \"CreationDate\": " + "\"" + fmt.Sprintf("%v", target["CreationDate"]) + "\"\n" +
		"}"
	qrcode.WriteFile(newjson, qrcode.Medium, 256, cid+".png")
}
