package main

import (
	"encoding/json"
	"fmt"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func separator() {
	fmt.Println("--------------------------")
}

func performChecks(sh *shell.Shell) error {
	if YourLocalPath == "" {
		return fmt.Errorf("'YourLocalPath' constant is NOT defined. Please, provide a local path in your computer where the file will be downloaded")
	}

	if YourPublicKey == "" {
		return fmt.Errorf("'YourPublicKey' constant is NOT defined. Please, provide the public key of your IPFS node")
	}

	if !sh.IsUp() {
		return fmt.Errorf("You do not have an IPFS node running at port 5001")
	}

	return nil
}

func jsonFormat(data primitive.A) (newJsonData string) {
	jsonData, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		panic(err)
	}
	newJsonData = strings.Replace(string(jsonData), "\"Key\": ", "", -1)
	newJsonData = strings.Replace(newJsonData, ",\n      \"Value\"", "", -1)
	newJsonData = strings.Replace(newJsonData, "[", "", 1)
	newJsonData = strings.Replace(newJsonData, "\n   },\n   {", ",", -1)
	sz := len(newJsonData)
	newJsonData = newJsonData[:sz-1]
	return newJsonData
}
