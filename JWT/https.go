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
	fmt.Println("------>NYTT TEST SAMUEL REMANCONTENTCAT:", dataOnIPNS)
	tmpByte, _ = json.Marshal([]byte(dataOnIPNS))
	json.Unmarshal(tmpByte, &record)

	err = json.Unmarshal(body, &MutableData)
	if err != nil {
		fmt.Println("Error unmarshaling body, error code: ", err)
	}
	remanEventData["CID"] = cid
	remanEventData["EventType"] = MutableData.Eventtype
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
	newRecord := append([]byte(dataOnIPNS), recordJson...)

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
		remanData := catRemanContent(chooseEvent.Key)
		err = json.Unmarshal([]byte(remanData), &getEvent)
		if err != nil {
			fmt.Println("Error unmarshaling remanData, error code: ", err)
		}

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

func addMutableProduct(writer http.ResponseWriter, request *http.Request) {
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

	var MutableData httpDataProduct

	err = json.Unmarshal(body, &MutableData)
	if err != nil {
		fmt.Println("Put request failed", err)
	}

	remanEventData := make(map[string]interface{})
	var record []appendEntryProduct
	var record2 []appendEntryProduct
	var appendEntry []appendEntryProduct
	dataOnIPNS := catRemanContent(MutableData.Key)

	tmpByte, _ := json.Marshal([]byte(dataOnIPNS))
	err = json.Unmarshal([]byte(dataOnIPNS), &record2)

	fmt.Println("EFTER UNMARSHAL ", record2)
	json.Unmarshal(tmpByte, &record)
	found := false
	fmt.Println("CID TO REPLACEDCED ", MutableData.CIDToReplace)
	if MutableData.CIDToReplace != "" {
		for i := 0; i < len(record2); i++ {
			if record2[i].CID == MutableData.CIDToReplace {
				fmt.Println("HIttade en matchande CID: ", record2[i].CID)
				record2[i].CID = MutableData.CID
				record2[i].Datetime = MutableData.Datetime
				record2[i].ProductType = MutableData.ProductType
				fmt.Println("RECORD 2", record2)
				found = true
				break
			}
		}
		if !found {
			fmt.Println("Error CID not found in event log: ", MutableData.CIDToReplace)
			http.Error(writer, "Error CID not found in event log: ", http.StatusNotFound)
			return
		}
		result, _ := json.Marshal(record2)
		fmt.Println("RECORD 2 MARSHAL", string(result))
		sh := shell.NewShell("localhost:5001")
		cid, _ := addFile(sh, string(result))
		addDataToIPNS(sh, MutableData.Key, cid)

		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(cid))
	} else {
		err = json.Unmarshal(body, &MutableData)
		if err != nil {
			fmt.Println("Error unmarshaling body, error code: ", err)
		}
		remanEventData["CID"] = MutableData.CID
		remanEventData["Datetime"] = MutableData.Datetime
		remanEventData["ProductType"] = MutableData.ProductType

		jsonAdd, err := json.Marshal(remanEventData)
		if err != nil {
			fmt.Println("Error unmarshaling jsonAdd, error code: ", err)
		}
		newString := "[" + string(jsonAdd) + "]"

		err = json.Unmarshal([]byte(newString), &appendEntry)
		if err != nil {
			fmt.Println("Error unmarshaling newString, error code: ", err)
		}

		if err != nil {
			fmt.Println("Error adding file to IPNS: ", err)
		}
		fmt.Println("reman event------------->", newString)
		record = append(record, appendEntry...)
		recordJson, err := json.Marshal(record)
		fmt.Println("recordJson event------------->", string(recordJson))
		newRecord := append([]byte(dataOnIPNS), recordJson...)
		if err != nil {
			fmt.Println("Error marshalling record: ", err)
		}

		record2 := strings.Replace(string(newRecord), "[", "", -1)
		record2 = strings.Replace(record2, "]", "", -1)
		record2 = strings.Replace(record2, "}{", "},{", -1)
		record2 = "[" + record2 + "]"
		fmt.Println("DATAN SOM BLIR UPPLADDADDD:) \n", record2)

		sh := shell.NewShell("localhost:5001")
		cid, err := addFile(sh, record2)
		addDataToIPNS(sh, MutableData.Key, cid)

		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(cid))
	}

	//fmt.Println("Result SOM BLIR UPPLADDADDD:) \n", string(result))

}
