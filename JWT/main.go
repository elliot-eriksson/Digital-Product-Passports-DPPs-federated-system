package main

import (
	"log"
	"net/http"
)

// Message is a struct representing the response format

func main() {

	getURL := "/retrieveData"
	// getAllIPNS := "/retrieveAllMutableData"
	getKey := "/get-publicKey"
	createPassportURL := "/publishPassport"
	addRemanafactureEventURL := "/addMutableData"
	retriveEventURL := "/retriveEvent"

	http.HandleFunc(getURL, getHandler)
	http.HandleFunc(getKey, generateKey)
	http.HandleFunc(createPassportURL, createPassportHandler)
	http.HandleFunc(addRemanafactureEventURL, addMutableData)
	http.HandleFunc(retriveEventURL, retriveEvent)

	// Define a route handler for the "/home" endpoint
	// http.HandleFunc("/home", handlePage)

	// Start the server on port 8080
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Println("There was an error listening on port :8081", err)
	}
}
