package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"reflect"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
	qrcode "github.com/skip2/go-qrcode"
)

func getHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	//Check that messages is GET
	if request.Method != http.MethodGet {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var tmpStringClaim tmpStringClaim
	// Decode JSON from the request body into the Message struct
	err := json.NewDecoder(request.Body).Decode(&tmpStringClaim)
	if err != nil {
		fmt.Println("Get request failed", err)
		http.Error(writer, "Error reading body", http.StatusNotAcceptable)
		return
	}

	if tmpStringClaim.CID[0] == 81 { // checks if the first char is Q
		CID := "/ipfs/" + tmpStringClaim.CID
		Dpp := getPassport(CID, "")
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(Dpp))

	}

}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func test(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var passport interface{}
	body, err := io.ReadAll(request.Body)
	if err != nil {
		fmt.Println("Error reading body", err)
		http.Error(writer, "Error reading body", http.StatusNotAcceptable)
		return
	}
	err = json.Unmarshal(body, &passport)
	passportData, _ := passport.(map[string]interface{})
	madeby := passportData["Made_by"]
	delete(passportData, "Made_by")
	fmt.Println("made_bay", madeby)

	if madeby != "" {
		str := fmt.Sprintf("%v", madeby)
		fmt.Println("EFTER IF", str)
	}
}

func createPassportHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	//Check that messages is Put
	if request.Method != http.MethodPost {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		fmt.Println("Error reading body", err)
		http.Error(writer, "Error reading body", http.StatusNotAcceptable)
		return
	}

	var passport interface{}
	var dataToCA dataToCA
	err = json.Unmarshal(body, &passport)
	if err != nil {
		fmt.Println("Error Unmarshal body", err)
		http.Error(writer, "Error Unmarshal body", http.StatusNotAcceptable)
		return
	}
	passportData, _ := passport.(map[string]interface{})
	passType := passportData["Type"]
	delete(passportData, "Type")
	if passType == "complete" {

		pubKey, privKey := generateKey()
		passportData["remanufacturing_events"] = pubKey
		dataToCA.Remanufacturing_events.Privatekey = privKey
		dataToCA.Remanufacturing_events.Publickey = pubKey

		pubKey, privKey = generateKey()
		passportData["shipping"] = pubKey
		dataToCA.Shipping.Privatekey = privKey
		dataToCA.Shipping.Publickey = pubKey

		pubKey, privKey = generateKey()
		passportData["makes"] = pubKey
		dataToCA.Makes.Privatekey = privKey
		dataToCA.Makes.Publickey = pubKey

		pubKey, privKey = generateKey()
		passportData["made_from"] = pubKey
		dataToCA.Made_from.Privatekey = privKey
		dataToCA.Made_from.Publickey = pubKey

	} else if passType == "simple" {

	} else {
		fmt.Println("No valid type sent")
		http.Error(writer, "No valid type sent", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &passport)
	output, err := json.Marshal(passportData)

	madeby := passportData["Made_by"]
	delete(passportData, "Made_by")

	sh := shell.NewShell("localhost:5001")
	cid, err := addFile(sh, string(output))
	if err != nil {
		fmt.Println("Error uploading to IPFS", err)
		http.Error(writer, "Error uploading to IPFS", http.StatusInternalServerError)
		return
	}

	if passType == "complete" {
		dataToCA.Cid = cid
		postData, _ := json.Marshal(dataToCA)
		fmt.Println("dataToCA", string(postData))
		sendToCa(postData, "POST")
	}
	fmt.Println("madeby string grej")
	fmt.Println("madeby type:", reflect.TypeOf(madeby))
	if reflect.TypeOf(madeby) != nil {
		fmt.Println("borde inte vara här..")
		madeByString := fmt.Sprintf("%v", madeby)
		cid, err := addFile(sh, madeByString)
		if err != nil {
			fmt.Println("Error uploading to IPFS", err)
			http.Error(writer, "Error uploading made_by to IPFS", http.StatusInternalServerError)
			return
		}
		addDataToIPNS(sh, dataToCA.Remanufacturing_events.Privatekey, cid)
	}
	response, err := json.Marshal(cid)
	if err != nil {
		fmt.Println("Error uploading to IPFS", err)
		http.Error(writer, "Error uploading to IPFS", http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(response)
}

func sendToCa(body []byte, method string) string {
	fmt.Println("JSON SENT TO CA \n", string(body))
	var address string
	if method == "GET" {
		var addressToCA addressToCA
		json.Unmarshal(body, &addressToCA)
		fmt.Println(addressToCA.PublicKey)
		address = "https://d0020e-project-dpp.vercel.app/api/v1/CA/" + addressToCA.PublicKey
	} else {
		address = "https://d0020e-project-dpp.vercel.app/api/v1/CA"
	}
	req, err := http.NewRequest(method, address, bytes.NewBuffer(body))

	// req, err := http.NewRequest(method, "http://localhost:80/test", bytes.NewBuffer(body))

	if err != nil {
		fmt.Println("Error sending keys to CA", err)
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	fmt.Println("requesten", req)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error Do", err)
		return ""
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error readall", err)
		return ""
	}
	fmt.Println("Response from CA: ", string(responseBody))
	return string(responseBody)
}

func addMutableData(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	fmt.Println("REQUEST METHODE", request.Method)
	//Check that messages is Put
	if request.Method != http.MethodPost {
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
	fmt.Println("After first unmarshal", MutableData.Key)
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
	fmt.Println("CID before key check", cid)
	if err != nil {
		fmt.Println("Error adding file to IPNS: ", err)
	}
	fmt.Println("Before checkKey", MutableData.Key)
	if !checkKey(MutableData.Key) {
		addDataToIPNS(sh, MutableData.Key, cid)
		writer.WriteHeader(http.StatusOK)
		statusText := "Data added"
		_, _ = writer.Write([]byte(statusText))
	} else {
		success, message := retrievePrivateKey(MutableData.Key)
		fmt.Println("Success", success)
		fmt.Println("Successmessage", message)

		if success == "true" {
			addDataToIPNS(sh, MutableData.Key, cid)
			writer.WriteHeader(http.StatusOK)
			statusText := "Data added"
			_, _ = writer.Write([]byte(statusText))
		} else {
			if err != nil {
				fmt.Println("Put request failed", err)
				http.Error(writer, message, http.StatusNotAcceptable)
				return
			}
		}

	}
}

func retriveEvent(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	fmt.Println("REQUEST METHODE", request.Method)
	//Check that messages is Get
	if request.Method != http.MethodGet {
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
		// if needed it can be implemented

	} else if chooseEvent.Type == "LastEvent" {
		remanData := catRemanContent(chooseEvent.Key)
		fmt.Println("remanData", remanData)
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
	if request.Method != http.MethodPost {
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

	json.Unmarshal(tmpByte, &record)
	found := false
	if MutableData.CIDToReplace != "" {
		for i := 0; i < len(record2); i++ {
			if record2[i].CID == MutableData.CIDToReplace {
				record2[i].CID = MutableData.CID
				record2[i].Datetime = MutableData.Datetime
				record2[i].ProductType = MutableData.ProductType
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
		sh := shell.NewShell("localhost:5001")
		cid, _ := addFile(sh, string(result))
		if checkKey(MutableData.Key) {
			addDataToIPNS(sh, MutableData.Key, cid)
			writer.WriteHeader(http.StatusOK)
			statusText := "Data added"
			_, _ = writer.Write([]byte(statusText))
		} else {
			success, message := retrievePrivateKey(MutableData.Key)
			if success == "true" {
				addDataToIPNS(sh, MutableData.Key, cid)
				writer.WriteHeader(http.StatusOK)
				statusText := "Data added"
				_, _ = writer.Write([]byte(statusText))
			} else {
				if err != nil {
					fmt.Println("Put request failed", err)
					http.Error(writer, message, http.StatusNotAcceptable)
					return
				}
			}

		}
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
		if checkKey(MutableData.Key) {
			addDataToIPNS(sh, MutableData.Key, cid)
			writer.WriteHeader(http.StatusOK)
			statusText := "Data added"
			_, _ = writer.Write([]byte(statusText))
		} else {
			success, message := retrievePrivateKey(MutableData.Key)
			if success == "true" {
				addDataToIPNS(sh, MutableData.Key, cid)
				writer.WriteHeader(http.StatusOK)
				statusText := "Data added"
				_, _ = writer.Write([]byte(statusText))
			} else {
				if err != nil {
					fmt.Println("Put request failed", err)
					http.Error(writer, message, http.StatusNotAcceptable)
					return
				}
			}

		}
	}
}

func retrieveMutableLog(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	fmt.Println("REQUEST METHODE", request.Method)
	//Check that messages is Get
	if request.Method != http.MethodGet {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		fmt.Println("Error reading body", err)
		http.Error(writer, "Error reading body", http.StatusNotAcceptable)
		return
	}

	var MutableLog MutableLog

	err = json.Unmarshal(body, &MutableLog)
	if err != nil {
		fmt.Println("Put request failed", err)
	}

	log := catRemanContent(MutableLog.Key)

	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte(log))

}

func generateQrCode(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	//Check that messages is GET
	if request.Method != http.MethodGet {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		fmt.Println("Error reading body", err)
		http.Error(writer, "Error reading body", http.StatusNotAcceptable)
		return
	}

	var tmpStringClaim tmpStringClaim

	// Decode JSON from the request body into the Message struct
	err = json.Unmarshal(body, &tmpStringClaim)
	if err != nil {
		fmt.Println("Error unmarshal body", err)
		http.Error(writer, "Error unmarshal body", http.StatusNotAcceptable)
		return
	}
	CID := "/ipfs/" + tmpStringClaim.CID
	Dpp := getPassport(CID, "")

	var QrCode QrCode
	var QrCodeImage QrCodeImage
	QrCode.CID = tmpStringClaim.CID
	json.Unmarshal([]byte(Dpp), &QrCode)
	qrString, err := json.Marshal(QrCode)
	if err != nil {
		fmt.Println("Error marshaling QrCode", err)
		http.Error(writer, "Error marshaling QrCode", http.StatusNotAcceptable)
		return
	}
	png, err := qrcode.Encode(string(qrString), qrcode.Medium, 256)
	if err != nil {
		fmt.Println("Error encoding QrCode", err)
		http.Error(writer, "Error encoding QrCode", http.StatusNotAcceptable)
		return
	}
	base64Encoded := base64.StdEncoding.EncodeToString(png)
	QrCodeImage.Filename = tmpStringClaim.CID
	QrCodeImage.Content = base64Encoded
	marshalQRImg, err := json.Marshal(QrCodeImage)
	if err != nil {
		fmt.Println("Error marshaling QrCodeImage", err)
		http.Error(writer, "Error marshaling QrCodeImage", http.StatusNotAcceptable)
		return
	}
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(marshalQRImg)
}
