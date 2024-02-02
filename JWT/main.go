package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// Message is a struct representing the response format

func main() {

	getURL := "/get-endpoint"
	createPassportURL := "/put-createPassport"
	addRemanafactureEventURL := "/put-addRemanafactureEvent"

	http.HandleFunc(getURL, getHandler)
	http.HandleFunc(createPassportURL, createPassportHandler)
	http.HandleFunc(addRemanafactureEventURL, addRemanafactureEventHandler)

	// Define a route handler for the "/home" endpoint
	// http.HandleFunc("/home", handlePage)

	// Start the server on port 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("There was an error listening on port :8080", err)
	}

	// ------------ SUNES GREJER --------------
	// cid := "QmbjnEtna7T7hLN3CmaVYPNwwkQGxUoEZsGZJfNVVusmJB"
	// key := "hej"
	// getPassport(cid, key)
	// getSensetive(cid, key, key)
}

// handlePage is the handler function for the "/home" endpoint
func handlePage(writer http.ResponseWriter, request *http.Request) {
	// Set Content-Type header to indicate JSON response
	writer.Header().Set("Content-Type", "application/json")

	// Create an instance of the Message struct
	var userClaim userClaim

	// Decode JSON from the request body into the Message struct
	err := json.NewDecoder(request.Body).Decode(&userClaim)
	if err != nil {
		// If there is an error decoding, you can handle it as needed
		log.Println("Error decoding JSON:", err)
		return
	}

	// Encode the Message struct back to JSON and write it to the response
	err = json.NewEncoder(writer).Encode(userClaim)
	if err != nil {
		// If there is an error encoding, you can handle it as needed
		log.Println("Error encoding JSON:", err)
		return
	}

	fmt.Println(userClaim.Username)
	fmt.Println(userClaim.Hash)
	//Finns user ?
	// if authenticateRole(userClaim.hasedVlaue){
	rights := checkAccessRights(userClaim.Hash)
	fmt.Println("RIGHTS----------------", rights)
	jwtToken, err := createJWT(userClaim.Username, rights, userClaim.Hash)
	if err != nil {
		// If there is an error encoding, you can handle it as needed
		log.Println("Error creating JWT token:", err)
		return
	}
	// 	fmt.Printf("JWT Token: %s\n", jwtToken)
	// } else {
	// 	return
	// }
	if isValidJWT(jwtToken, key) {
		fmt.Printf("JWT Token: %s\n", jwtToken)
		fmt.Println("JWT is valid!")
	} else {
		fmt.Println("JWT is not valid.")
	}
	//------------Extracting payload without verifying signature
	payload, err := extractPayload(jwtToken)
	payloadJSON, err := extractPayloadJSON(jwtToken)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("-----Payload", payload)
	fmt.Println("-----PayloadJSON", payloadJSON)
}

// TODO: Generera, auth, validate

const key = "your-256-bit-secret"

type userClaim struct {
	jwt.RegisteredClaims
	Username         string
	Hash             string
	isAdmin          bool
	isRemanufactorer bool
	isUser           bool
}

type AccessRights struct {
	isAdmin          bool
	isRemanufactorer bool
	isUser           bool
}

func createJWT(username string, rights AccessRights, hash string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// "RegisteredClaims": jwt.RegisteredClaims{},
		"Username":         username,
		"Hash":             hash,
		"isAdmin":          rights.isAdmin,
		"isRemanufactorer": rights.isRemanufactorer,
		"isUser":           rights.isUser,
	})

	fmt.Println("TOKEN.................", token)
	//creation of JWT
	signedString, err := token.SignedString([]byte(key))
	fmt.Println("signedString", signedString)

	if err != nil {
		return "", fmt.Errorf("error creating signedString: %v", err)
	}
	return signedString, nil
}

func checkAccessRights(userHash string) AccessRights {
	//Connection till mongoDb för att kolla om user existerar
	userHashAdmin := "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	userHashUser := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	userHashGuest := "7890abcdef1234567890abcdef1234567890abcdef1234567890abcdef123456"

	//header + userHashAdmin
	//payload
	//secret
	// adminRights := userClaim{isAdmin: true, isRemanufactorer: false, isUser: false}
	// manufacturerRights := userClaim{isAdmin: false, isRemanufactorer: true, isUser: false}
	// userRights := userClaim{isAdmin: false, isRemanufactorer: false, isUser: true}
	// noRights := userClaim{isAdmin: false, isRemanufactorer: false, isUser: false}

	adminRights := AccessRights{isAdmin: true, isRemanufactorer: false, isUser: false}
	manufacturerRights := AccessRights{isAdmin: false, isRemanufactorer: true, isUser: false}
	userRights := AccessRights{isAdmin: false, isRemanufactorer: false, isUser: true}
	noRights := AccessRights{isAdmin: false, isRemanufactorer: false, isUser: false}

	if userHash == userHashAdmin {
		//Non sensitive & sensitive
		fmt.Println("Admin ", adminRights.isAdmin, "\n", "manu ", adminRights.isRemanufactorer, "\n", "user ", adminRights.isUser)
		return adminRights
	} else if userHash == userHashUser {
		//Möjlighet att inserta till databas
		//Skapa addresser till databas som uppladdas till IPFS
		//Non sensitive + lägga till remanufactor event
		fmt.Println(manufacturerRights.isRemanufactorer)
		return manufacturerRights
	} else if userHash == userHashGuest {
		//Non sensitive
		fmt.Println(userRights.isUser)
		return userRights
	} else {
		fmt.Println("norights")
		return noRights
	}

	//Clinet
	//ENcrypt hasahd verision av key
	//--> decrypta
	//key
}
func extractPayload(jwtToken string) (map[string]interface{}, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(jwtToken, jwt.MapClaims{})

	if err != nil {
		return nil, err
	}

	// Extract the payload
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("failed to extract payload")
}

func extractPayloadJSON(jwtToken string) (string, error) {
	//Without checking the
	token, _, err := new(jwt.Parser).ParseUnverified(jwtToken, jwt.MapClaims{})
	if err != nil {
		return "", err
	}

	// Extract the payload
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Convert the payload to JSON
		payloadJSON, err := json.Marshal(claims)
		if err != nil {
			return "", err
		}

		return string(payloadJSON), nil
	}

	return "", fmt.Errorf("failed to extract payload")
}

func isValidJWT(signedString string, key string) bool {
	token, err := jwt.Parse(signedString,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})

	if err != nil {
		fmt.Println("Error parsing token:", err)
		return false
	}

	//Validation
	if !token.Valid {
		fmt.Println("Invalid token")
		return false
	}
	return true

}
