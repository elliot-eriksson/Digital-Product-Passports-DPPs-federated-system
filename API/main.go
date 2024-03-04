package main

import (
	"log"
	"net/http"
)

// Message is a struct representing the response format

func main() {

	getURL := "/retrieveData"
	createPassportURL := "/publishPassport"
	addRemanafactureEventURL := "/addEvent"
	retrieveEventURL := "/retrieveEvent"
	addMutableProductURL := "/addMutableProduct"
	retrieveMutableLogURL := "/retrieveMutableLog"
	generateQcodeURL := "/getQrCode"
	testUTL := "/test"
	http.HandleFunc(testUTL, test)

	http.HandleFunc(getURL, getHandler)
	http.HandleFunc(createPassportURL, createPassportHandler)
	http.HandleFunc(addRemanafactureEventURL, addMutableData)
	http.HandleFunc(retrieveEventURL, retriveEvent)
	http.HandleFunc(addMutableProductURL, addMutableProduct)
	http.HandleFunc(retrieveMutableLogURL, retrieveMutableLog)
	http.HandleFunc(generateQcodeURL, generateQrCode)

	// Start the server on port 80
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Println("There was an error listening on port :80", err)
	}
}
