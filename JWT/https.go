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
		return
	}
	fmt.Println("--->CID ", testClaim.CID)

	if testClaim.CID[0] == 107 { // checks if the first char is k
		//fmt.Println("This is a public key ", key)
		key := "/ipns/" + testClaim.CID
		output := lsIPNS(key)
		content, contentLenght := splitListContent(output)
		stringindex := catContent(content, contentLenght)
		for output := range stringindex {
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(stringindex[output]))
		}
	}
	if testClaim.CID[0] == 81 { // checks if the first char is Q
		//fmt.Println("This is a CID ", key)
		Dpp := getPassport(testClaim.CID, keyD)
		if err != nil {
			fmt.Println("Wrong CID", err)
			return
		}
		fmt.Println("GETHANDLER INNAN MARSHAL ", Dpp)
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
		http.Error(writer, "Error reading body", http.StatusMethodNotAllowed)
		return
	}

	//var apiCheck APICheck

	fmt.Println("Body", string(body))
	sh := shell.NewShell("localhost:5001")

	cid, err := addFile(sh, string(body))
	fmt.Println("cid----", cid)

	response, err := json.Marshal(cid)

	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(response)
	_, _ = writer.Write([]byte(cid))

}

type RemanEvent struct {
	Key  string `json: Key`
	Data string `json: Data`
}

func addRemanafactureEventHandler(writer http.ResponseWriter, request *http.Request) {
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
		http.Error(writer, "Error reading body", http.StatusMethodNotAllowed)
		return
	}

	var remanEvent RemanEvent
	fmt.Println(request.Body)
	fmt.Println("Bodddy", string(body))

	// Decode JSON from the request body into the remanEvent struct
	//err = json.NewDecoder(request.Body).Decode(&remanEvent)
	err = json.Unmarshal(body, &remanEvent)
	if err != nil {
		fmt.Println("Put request failed", err)
	}
	// fmt.Println("remanevent", remanEvent)
	// fmt.Println("data :", remanEvent.Data)
	// fmt.Println("wowowow", remanEvent.Key)

	sh := shell.NewShell("localhost:5001")
	cid, err := addFile(sh, remanEvent.Data)
	fmt.Println("samuels print cid: ", cid)
	output, _ := addDataToIPNS(sh, remanEvent.Key, cid)

	fmt.Println("publised", output)

	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte(cid))
}
