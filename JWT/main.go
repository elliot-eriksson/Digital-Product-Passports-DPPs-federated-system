package main

import (
	"log"
	"net/http"
)

// Message is a struct representing the response format

func main() {

	getURL := "/retrieveData"
	// getAllIPNS := "/retrieveAllMutableData"
	// getKey := "/get-publicKey"
	createPassportURL := "/publishPassport"
	addRemanafactureEventURL := "/addEvent"
	retriveEventURL := "/retriveEvent"
	addMutableProductURL := "/addMutableProduct"
	retrieveMutableLogURL := "/retrieveMutableLog"
	generateQcodeURL := "/getQrCode"
	// retriveMutableProductURL := "/retriveMutableProduct"
	testURL := "/test"

	http.HandleFunc(getURL, getHandler)
	// http.HandleFunc(getKey, generateKey)
	http.HandleFunc(createPassportURL, createPassportHandler)
	http.HandleFunc(addRemanafactureEventURL, addMutableData)
	http.HandleFunc(retriveEventURL, retriveEvent)
	http.HandleFunc(addMutableProductURL, addMutableProduct)
	http.HandleFunc(retrieveMutableLogURL, retrieveMutableLog)
	http.HandleFunc(generateQcodeURL, generateQrCode)
	http.HandleFunc(testURL, test)
	// http.HandleFunc(retriveMutableProductURL, retriveMutableProduct)

	// Define a route handler for the "/home" endpoint
	// http.HandleFunc("/home", handlePage)

	// Start the server on port 8080
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Println("There was an error listening on port :8081", err)
	}
}
