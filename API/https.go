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
	"time"

	shell "github.com/ipfs/go-ipfs-api"
	qrcode "github.com/skip2/go-qrcode"
)

// Function that handles the retrieval of passports from IPFS requires a CID
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
		Dpp := getPassport(CID)
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(Dpp))

	}

}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Used to give generated private keys a temporary random name before being renamed to the responding public key
func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Creates a DPP with the information contained in the call to the endpoint. It can create either simple/barebone passports without keys or a complete with keys.
// If the passport is created with keys it uploads the public and private keys to the CA.
// If the call contains the field Made_by it uploads that to the responding key and updates the makes log in the responding passports.
func createPassportHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	//Check that messages is Post
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
	// Generates the 4 keypairs needed for a complete passport
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
		// If the passport is a simple/barebone passport there is no need for key generation
	} else {
		fmt.Println("No valid type sent")
		http.Error(writer, "No valid type sent", http.StatusInternalServerError)
		return
	}

	madeby := make(map[string]interface{})
	madeby["Made_by"] = passportData["Made_by"]
	delete(passportData, "Made_by")
	output, err := json.Marshal(passportData)
	if err != nil {
		fmt.Println("Error Marshal passportData createPassportHandler:124", err)
		http.Error(writer, "Error in passportData", http.StatusInternalServerError)
		return
	}

	sh := shell.NewShell("localhost:5001")
	// Uploads the passport to IPFS, with all information contained in the call except the Made_by and Type fields.
	// This gives the program possibility to be handle different types of passport containing different data fields.
	cid, err := addFile(sh, string(output))
	if err != nil {
		fmt.Println("Error uploading to IPFS", err)
		http.Error(writer, "Error uploading to IPFS", http.StatusInternalServerError)
		return
	}

	// Uploads the generated keys to the CA created by https://github.com/LTU-D0020E/d0020e-project-dpp
	// If the project is run without the CA from https://github.com/LTU-D0020E/d0020e-project-dpp this code can be removed
	// Could be modified to use a different CA by changin the address here and in retrievePrivateKey
	if passType == "complete" {
		dataToCA.Cid = cid
		postData, _ := json.Marshal(dataToCA)
		outboundCalls(postData, "POST", "https://d0020e-project-dpp.vercel.app/api/v1/CA/")
	}

	// Populates the Made_by key with the data contained in the Made_by field
	if reflect.TypeOf(madeby) != nil {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(body), &data); err != nil {
			fmt.Println("Error unmarshal made_by Main 157")
			http.Error(writer, "Error unmarshal Made_by", http.StatusInternalServerError)
		}
		madeByJSON, err := json.Marshal(data["Made_by"])
		if err != nil {
			fmt.Println("Error marshal made_by Main 163")
			http.Error(writer, "Error in Made_by", http.StatusInternalServerError)
		}
		madeByString := string(madeByJSON)
		// Uploads Made_by to IPFS
		cidMade, err := addFile(sh, madeByString)
		if err != nil {
			fmt.Println("Error uploading to IPFS", err)
			http.Error(writer, "Error uploading made_by to IPFS", http.StatusInternalServerError)
			return
		}
		// Points the Made_by private key to the newly created log containing the Made_by data (Publish to IPNS)
		addDataToIPNS(sh, dataToCA.Made_from.Publickey, cidMade)

		// Goes through the incoming made_by and retrieves the key from each of the CID contained within.
		if madeBy, ok := madeby["Made_by"].([]interface{}); ok {
			for _, item := range madeBy {
				if m, ok := item.(map[string]interface{}); ok {
					makesData := make(map[string]interface{})
					makesData["CID"] = m["CID"]
					jsonData, _ := json.Marshal(makesData)
					dataFromCall := outboundCalls(jsonData, "GET", "http://localhost:80/retrieveData")
					var makesKey makesKey
					json.Unmarshal([]byte(dataFromCall), &makesKey)
					if reflect.TypeOf(makesKey) != nil {
						newString := fmt.Sprintf("%v", makesKey)
						newString = newString[1 : len(newString)-1]
						makesData["Key"] = newString
						makesData["CID"] = cid
						makesData["ProductType"] = passportData["ProductType"]
						makesData["Datetime"] = time.Now().Format("YYYY-MM-DD")
						jsonData, _ := json.Marshal(makesData)
						// Updates the Make data for each of the made_by passes
						if makesData["Key"] != "" {
							response := outboundCalls(jsonData, "POST", "http://localhost:80/addMutableProduct")
							fmt.Println("RESPONS 180", response)
						}
					}
				}
			}
		}
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

// Handles outbound calls the API does to the CA as well as some calls to itself.
func outboundCalls(body []byte, method string, address string) string {

	fmt.Println("JSON outbound \n", string(body))
	fmt.Println("Outbound Address\n", address)

	req, err := http.NewRequest(method, address, bytes.NewBuffer(body))

	if err != nil {
		fmt.Println("Error sending keys to CA", err)
		return ""
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
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
	fmt.Println("Response from OutboundCall: ", string(responseBody))
	return string(responseBody)
}

// Adds event to Remanufacturing or Shipping
// Requires the responding key and the event data
func addMutableData(writer http.ResponseWriter, request *http.Request) {
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

	// Adds the new event data to IPFS
	cid, err := addFile(sh, string(uploadData))
	if err != nil {
		fmt.Println("Failed to add file to IPNS", err)
	}

	remanEventData := make(map[string]interface{})
	var record []appendEntry
	var appendEntry []appendEntry

	// Retrieves the current event log on the given key
	dataOnIPNS := catRemanContent(MutableData.Key)
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

	// Appends the new event to the event log
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

	// Upload the new event log to IPFS
	cid, err = addFile(sh, newJsonData)
	if err != nil {
		fmt.Println("Error adding file to IPNS: ", err)
	}

	// Check if local node already has the private key responding to the public key.
	// If it does updates the IPNS record to point to new event log.
	// Else retrieves the private key from the CA and then update IPNS
	// If the project is run without a CA all keys need to be local
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

// Retrieves either the latest or all events corresponding to a private key.
// Used for Remanufacturing or Shipping events
func retriveEvent(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

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
		// Retrieves all events corresponding to a private key.
		remanData := catRemanContent(chooseEvent.Key)
		err = json.Unmarshal([]byte(remanData), &getEvent)
		if err != nil {
			fmt.Println("Error unmarshaling remanData, error code: ", err)
		}
		// Loops though the event log and returns all the events one by one
		for i, _ := range getEvent {
			newJsonData := strings.Replace(fmt.Sprintf("%v", getEvent[i]), "{", "", -1)
			newJsonData = strings.Replace(newJsonData, "}", "", -1)
			data := getPassport(newJsonData)
			data = data + "\n"

			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(data))
		}

	} else if chooseEvent.Type == "SpecificEvent" {
		// if needed it can be implemented

	} else if chooseEvent.Type == "LastEvent" {
		// Retrieves last/latest event corresponding to a private key.
		remanData := catRemanContent(chooseEvent.Key)
		// fmt.Println("remanData", remanData)
		err = json.Unmarshal([]byte(remanData), &getEvent)
		if err != nil {
			fmt.Println("Error unmarshaling remanData, error code: ", err)
		}

		getLast := strings.Replace(fmt.Sprintf("%v", getEvent[len(getEvent)-1]), "{", "", -1)
		getLast = strings.Replace(getLast, "}", "", -1)

		data := getPassport(getLast)

		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(data))

	} else {
		http.Error(writer, "Error no type selected", http.StatusNotAcceptable)
		return
	}
}

// Adds a product to the made_by or makes logs for a passport
func addMutableProduct(writer http.ResponseWriter, request *http.Request) {
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

	var MutableData httpDataProduct

	err = json.Unmarshal(body, &MutableData)
	if err != nil {
		fmt.Println("Put request failed", err)
	}
	fmt.Println("AddMutableProduct Body 459 ", string(body))
	remanEventData := make(map[string]interface{})
	var record []appendEntryProduct
	var record2 []appendEntryProduct
	var appendEntry []appendEntryProduct

	// Retrieves current log
	dataOnIPNS := catRemanContent(MutableData.Key)

	tmpByte, _ := json.Marshal([]byte(dataOnIPNS))
	err = json.Unmarshal([]byte(dataOnIPNS), &record2)

	json.Unmarshal(tmpByte, &record)
	found := false

	// Checks if the new product already exists in current log.
	// If exists update current log by changing the product information to the new provided information.
	// Else adds the new information as a new product last in the log

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
		// Uploads new event log to IPFS
		cid, _ := addFile(sh, string(result))

		// Check if local node already has the private key responding to the public key.
		// If it does updates the IPNS record to point to new event log.
		// Else retrieves the private key from the CA and then update IPNS
		// If the project is run without a CA all keys need to be local

		if checkKey(MutableData.Key) {
			statusText, _ := addDataToIPNS(sh, MutableData.Key, cid)
			writer.WriteHeader(http.StatusOK)
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

		sh := shell.NewShell("localhost:5001")
		// Uploads new event log to IPFS
		cid, err := addFile(sh, record2)
		// Check if local node already has the private key responding to the public key.
		// If it does updates the IPNS record to point to new event log.
		// Else retrieves the private key from the CA and then update IPNS
		// If the project is run without a CA all keys need to be local
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

// Retrieves the whole event log from a public key
func retrieveMutableLog(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
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

// Generates a QR code from a CID by retrieving the passport.
// Returns the filename = CID and a base 64 encoding of the QR code
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
	// Retrieves passport from IPFS
	Dpp := getPassport(CID)

	var QrCode QrCode
	var QrCodeImage QrCodeImage
	QrCode.CID = tmpStringClaim.CID
	json.Unmarshal([]byte(Dpp), &QrCode)
	// Creates the QR code with styling
	qrString, err := json.MarshalIndent(QrCode, "", "    ")
	if err != nil {
		fmt.Println("Error marshal indenting http.go 636", err)
		http.Error(writer, "Error marshaling QrCode", http.StatusNotAcceptable)
		return
	}
	// Generate QR code with predetermined data fields such as Product name
	png, err := qrcode.Encode(string(qrString), qrcode.Medium, 256)
	if err != nil {
		fmt.Println("Error encoding QrCode", err)
		http.Error(writer, "Error encoding QrCode", http.StatusNotAcceptable)
		return
	}
	// Base64 encodes the QR code
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
