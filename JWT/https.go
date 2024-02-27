package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

type tmpStringClaim struct {
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

	var tmpStringClaim tmpStringClaim
	fmt.Println(request.Body)

	// Decode JSON from the request body into the Message struct
	err := json.NewDecoder(request.Body).Decode(&tmpStringClaim)
	if err != nil {
		fmt.Println("Get request failed", err)
		http.Error(writer, "Error reading body", http.StatusNotAcceptable)
		return
	}
	// fmt.Println("--->CID ", tmpStringClaim.CID)

	if tmpStringClaim.CID[0] == 107 { // checks if the first char is k
		//fmt.Println("This is a public key ", key)
		key := "/ipns/" + tmpStringClaim.CID
		output := getPassport(key, keyD)
		content, contentLenght := splitListContent(output)
		stringindex := catContent(content, contentLenght)
		for output := range stringindex {
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(stringindex[output]))
		}
	}
	if tmpStringClaim.CID[0] == 81 { // checks if the first char is Q
		//fmt.Println("This is a CID ", key)
		CID := "/ipfs/" + tmpStringClaim.CID
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

// type MutableData struct {
// 	Key  string `json: Key`
// 	Data string `json: Data`
// 	Name string `json: Name`
// 	Date string `json: Date`
// }

// type MutableDataToUpload struct {
// 	Data string `json: Data`
// 	Name string `json: Name`
// 	Date string `json: Date`
// }

type httpData struct {
	Key       string `json: Key`
	Eventtype string `json: Eventtype`
	Datetime  string `json: Datetime`
	Data      string `json: Data`
}

type ledgerData struct {
	Eventtype string `json: Eventtype`
	Data      string `json: Data`
	Datetime  string `json: Datetime`
}

type appendEntry struct {
	CID       string `json: CID`
	Eventtype string `json: Eventtype`
	// Name string `json: Name`
	Datetime string `json: Datetime`
}

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

	var MutableData httpData

	err = json.Unmarshal(body, &MutableData)
	if err != nil {
		fmt.Println("Put request failed", err)
	}

	var MutableDataToUpload ledgerData
	tmpByte, _ := json.Marshal(MutableData)
	_ = json.Unmarshal(tmpByte, &MutableDataToUpload)
	uploadData, _ := json.Marshal(MutableDataToUpload)

	sh := shell.NewShell("localhost:5001")
	cid, err := addFile(sh, string(uploadData))
	if err != nil {
		fmt.Println("Failed to add file to IPNS", err)
	}

	remanEventData := make(map[string]interface{})
	var record []appendEntry
	var appendEntry []appendEntry

	dataOnIPNS := catRemanContent(MutableData.Key)

	tmpByte, _ = json.Marshal([]byte(dataOnIPNS))
	json.Unmarshal(tmpByte, &record)

	err = json.Unmarshal(body, &MutableData)
	if err != nil {
		fmt.Println("Error unmarshaling body, error code: ", err)
	}
	remanEventData["CID"] = cid
	remanEventData["Eventtype"] = MutableData.Eventtype
	remanEventData["Datetime"] = MutableData.Datetime

	jsonAdd, err := json.Marshal(remanEventData)
	if err != nil {
		fmt.Println("Error unmarshaling jsonAdd, error code: ", err)
	}
	newString := "[" + string(jsonAdd) + "]"

	err = json.Unmarshal([]byte(newString), &appendEntry)
	if err != nil {
		fmt.Println("Error unmarshaling newString, error code: ", err)
	}

	record = append(record, appendEntry...)
	recordJson, err := json.Marshal(record)
	newRecord := append(recordJson, dataOnIPNS...)

	if err != nil {
		fmt.Println("Error marshalling record: ", err)
	}

	newJsonData := strings.Replace(string(newRecord), "[", "", -1)
	newJsonData = strings.Replace(newJsonData, "]", "", -1)
	newJsonData = strings.Replace(newJsonData, "}{", "},{", -1)
	newJsonData = "[" + newJsonData + "]"

	cid, err = addFile(sh, newJsonData)
	if err != nil {
		fmt.Println("Error adding file to IPNS: ", err)
	}

	addDataToIPNS(sh, MutableData.Key, cid)

	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte(cid))
}

type chooseEvent struct {
	Key  string `json: Key`
	Type string `json: Type`
	CID  string `json: CID`

	// "Key":"k51qzi5uqu5dl8vkvhdrynmw3blxw6r2rx43ui0nhybad5nbvikmq04nd7gzb0", "type" : "last" ,"CID" : "QmVcvZu5N7VRyuarfZ2bAz6KkwdnsaEuQNFD8wdX1xmgJG"

}

type getEvent struct {
	CID string `json:"CID"`
	// Name string `json:"Name"`
	// Date string `json:"Date"`
	//Data      string `json: Data`

	// "Key":"k51qzi5uqu5dl8vkvhdrynmw3blxw6r2rx43ui0nhybad5nbvikmq04nd7gzb0", "type" : "last" ,"CID" : "QmVcvZu5N7VRyuarfZ2bAz6KkwdnsaEuQNFD8wdX1xmgJG"

}

func retriveEvent(writer http.ResponseWriter, request *http.Request) {
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
	var chooseEvent chooseEvent
	var getEvent []getEvent

	err = json.Unmarshal(body, &chooseEvent)
	if err != nil {
		fmt.Println("Put request failed", err)
		http.Error(writer, "Error not Json valid", http.StatusNotAcceptable)
		return
	}

	if chooseEvent.Type == "AllEvent" {
		remanData := catRemanContent(chooseEvent.Key)
		fmt.Println("remanData", remanData)
		err = json.Unmarshal([]byte(remanData), &getEvent)
		// fmt.Println("getEvent", getEvent.CID)
		if err != nil {
			fmt.Println("Error unmarshaling remanData, error code: ", err)
		}
		fmt.Println("CID: ", getEvent)

		for i, _ := range getEvent {

			fmt.Sprintf("%v", getEvent[i])

			newJsonData := strings.Replace(fmt.Sprintf("%v", getEvent[i]), "{", "", -1)
			newJsonData = strings.Replace(newJsonData, "}", "", -1)

			data := getPassport(newJsonData, "")
			data = data + "\n"
			fmt.Println("data: ", data)
			fmt.Println("CID: ", getEvent[i])

			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(data))
		}

	} else if chooseEvent.Type == "SpecificEvent" {
		// if needed

	} else if chooseEvent.Type == "LastEvent" {
		// 1. hämta event logg
		// 2. Ta ur data från sista append i event loggbody
		// 3. return eventtype, date, data
		remanData := catRemanContent(chooseEvent.Key)
		//fmt.Println("remanData", remanData)
		err = json.Unmarshal([]byte(remanData), &getEvent)
		// fmt.Println("getEvent", getEvent.CID)
		if err != nil {
			fmt.Println("Error unmarshaling remanData, error code: ", err)
		}

		// getLast := getEvent[len(getEvent)-1]
		getLast := strings.Replace(fmt.Sprintf("%v", getEvent[len(getEvent)-1]), "{", "", -1)
		getLast = strings.Replace(getLast, "}", "", -1)

		data := getPassport(getLast, "")

		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(data))

	} else {
		http.Error(writer, "Error no type selected", http.StatusNotAcceptable)
		return
	}

}
