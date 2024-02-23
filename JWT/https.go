package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	shell "github.com/ipfs/go-ipfs-api"
)

type TestClaim struct {
	CID string `json: CID`
}

func getHandler(writer http.ResponseWriter, request *http.Request) {
	keyD := "hej"
	writer.Header().Set("Content-Type", "application/json")
	// var response []byte
	//Check that messages is GET
	if request.Method != http.MethodGet {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var testClaim TestClaim
	fmt.Println(request.Body)

	// Decode JSON from the request body into the Message struct
	err := json.NewDecoder(request.Body).Decode(&testClaim)
	if err != nil {
		fmt.Println("Get request failed", err)
		http.Error(writer, "Error reading body", http.StatusNotAcceptable)
		return
	}
	// fmt.Println("--->CID ", testClaim.CID)

	if testClaim.CID[0] == 107 { // checks if the first char is k
		//fmt.Println("This is a public key ", key)
		key := "/ipns/" + testClaim.CID
		output := getPassport(key, keyD)
		content, contentLenght := splitListContent(output)
		stringindex := catContent(content, contentLenght)
		for output := range stringindex {
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(stringindex[output]))
		}
	}
	if testClaim.CID[0] == 81 { // checks if the first char is Q
		//fmt.Println("This is a CID ", key)
		CID := "/ipfs/" + testClaim.CID
		Dpp := getPassport(CID, keyD)
		if err != nil {
			fmt.Println("Wrong CID", err)
			return
		}
		// fmt.Println("GETHANDLER INNAN MARSHAL ", Dpp)
		// response, err = json.Marshal(Dpp)
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(Dpp))

	}

}

func generateKey(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodGet {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	randomName := randSeq(10)
	keyname := keyGenerator(randomName)
	keyRename(randomName, keyname)
	// response, err := json.Marshal([]byte(keyname))
	// if err != nil {
	// 	fmt.Println("", err)
	// 	return
	// }
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte(keyname))
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type APICheck struct {
	APIKey string `json: api_key`
}

func createPassportHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	//Check that messages is Put
	if request.Method != http.MethodPut {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		fmt.Println("Error reading body", err)
		http.Error(writer, "Error reading body", http.StatusNotAcceptable)
		return
	}

	//var apiCheck APICheck

	fmt.Println("Body", string(body))
	sh := shell.NewShell("localhost:5001")

	cid, err := addFile(sh, string(body))
	if err != nil {
		fmt.Println("Error uploading to IPFS", err)
		http.Error(writer, "Error uploading to IPFS", http.StatusInternalServerError)
		return
	}
	fmt.Println("cid----", cid)

	response, err := json.Marshal(cid)

	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(response)
	// _, _ = writer.Write([]byte(cid))

}

// func errorHandling(field1 string){
// 		//Korrekt nyckel

// 	if field1  != nil{
// 		http.Error(writer, "Error uploading to IPFS", http.StatusInternalServerError)

// 	}

// }

type MutableData struct {
	Key  string `json: Key`
	Data string `json: Data`
}

type MutableData2 struct {
	Key  string `json: Key`
	Data string `json: Data`
	Name string `json: Name`
	Date string `json: Date`
}

type newData struct {
	CID  string `json: CID`
	Name string `json: Name`
	Date string `json: Date`
}

// type newData struct {
// 	CID  string
// 	Name string
// 	Date string
// }

func addMutableData(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	fmt.Println("REQUEST METHODE", request.Method)
	//Check that messages is Put
	if request.Method != http.MethodPut {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		fmt.Println("Error reading body", err)
		http.Error(writer, "Error reading body", http.StatusNotAcceptable)
		return
	}

	var mutableData MutableData
	fmt.Println(request.Body)
	fmt.Println("Bodddy", string(body))

	// Decode JSON from the request body into the remanEvent struct
	//err = json.NewDecoder(request.Body).Decode(&remanEvent)
	err = json.Unmarshal(body, &mutableData)
	if err != nil {
		fmt.Println("Put request failed", err)
	}

	sh := shell.NewShell("localhost:5001")
	cid, err := addFile(sh, mutableData.Data)

	var mutableData2 MutableData2
	sunesfulagrej := make(map[string]interface{})
	var oldData []newData
	var newData []newData

	key := "/ipns/" + mutableData.Key
	outputtest := getPassport(key, "")
	test, _ := json.Marshal([]byte(outputtest))
	json.Unmarshal(test, &oldData)

	err = json.Unmarshal(body, &mutableData2)
	fmt.Println("MUTABLE DATA 2", mutableData2)

	sunesfulagrej["CID"] = cid
	sunesfulagrej["Name"] = mutableData2.Name
	//sunesfulagrej["Data"] = mutableData2.Data
	sunesfulagrej["Date"] = mutableData2.Date

	fmt.Println("sunesfulagrej['CID']: ", sunesfulagrej["CID"])
	fmt.Println("sunesfulagrej['Name']: ", sunesfulagrej["Name"])
	// fmt.Println("sunesfulagrej['Data']: ", sunesfulagrej["Data"])
	fmt.Println("sunesfulagrej['Date']: ", sunesfulagrej["Date"])

	jsonAdd, err := json.Marshal(sunesfulagrej)
	fmt.Println("NEW DATA Marshal", string(jsonAdd))
	newString := "[" + string(jsonAdd) + "]"
	fmt.Println("NEW STRING", newString)

	// fmt.Println("type of jsonAdd %T\n: ", reflect.TypeOf(jsonAdd))

	err = json.Unmarshal([]byte(newString), &newData)
	if err != nil {
		fmt.Println("ERROR UNMARSHAL", err)
	}

	fmt.Println("NEW DATA Unmarshal", newData)

	oldData = append(oldData, newData...)

	// for _, v := range oldData {
	// 	fmt.Println(v)
	// }
	result, err := json.Marshal(oldData)
	fmt.Println("result: ", string(result))

	cid, err = addFile(sh, string(result))

	fmt.Println("samuels print cid: ", cid)
	output, _ := addDataToIPNS(sh, mutableData.Key, cid)
	fmt.Println("published", output)

	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte(cid))
}
