package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// Message is a struct representing the response format
type userClaim struct {
	jwt.RegisteredClaims
	Username string
	Password string
}

func main() {
	// Define a route handler for the "/home" endpoint
	http.HandleFunc("/home", handlePage)

	//Start the server on port 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("There was an error listening on port :8080", err)
	}

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
	fmt.Println(userClaim.Password)
	//Finns user ?
	// if authenticateRole(userClaim.hasedVlaue){
	jwtToken, err := createJWT(userClaim.Username, userClaim.Password)
	// 	if err != nil {
	// 		// If there is an error encoding, you can handle it as needed
	// 		log.Println("Error creating JWT token:", err)
	// 		return
	// 	}
	// 	fmt.Printf("JWT Token: %s\n", jwtToken)
	// } else {
	// 	log.Println("This user does not exist")
	// 	return
	// }
	if isValidJWT(jwtToken, key) {
		fmt.Printf("JWT Token: %s\n", jwtToken)
		fmt.Println("JWT is valid!")
	} else {
		fmt.Println("JWT is not valid.")
	}
}

// TODO: Generera, auth, validate

const key = "your-256-bit-secret"

func createJWT(username string, password string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim{
		RegisteredClaims: jwt.RegisteredClaims{},
		Username:         username,
		Password:         password,
	})

	//creation of JWT
	signedString, err := token.SignedString([]byte(key))
	//fmt.Println("signedString", signedString)

	if err != nil {
		return "", fmt.Errorf("error creating signedString: %v", err)
	}
	return signedString, nil
}

func checkAccessRights(userHash string) string {
	//Connection till mongoDb för att kolla om user existerar
	userHashAdmin := "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	userHashUser := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	userHashGuest := "7890abcdef1234567890abcdef1234567890abcdef1234567890abcdef123456"

	//header + userHashAdmin
	//payload
	//secret
	if userHash == userHashAdmin {
		return "admin"
	}
	if userHash == userHashUser {
		return "user"
	}
	if userHash == userHashGuest {
		return "guest"
	}
	return "no access"
	//Clinet
	//ENcrypt hasahd verision av key
	//--> decrypta
	//key
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
