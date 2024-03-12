package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

// Reads the file with the given CID from IPFS
// Pin all retrieved files to the local IPFS Node
func readFile(sh *shell.Shell, cid string) (*string, error) {
	reader, err := sh.Cat(cid)
	if err != nil {
		return nil, fmt.Errorf("Error reading the file %s", err.Error())
	}

	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("Error reading the file %s", err.Error())
	}

	text := string(bytes)

	// Pin the given CID to the local IPFS Node
	pinToIPFS(cid)

	return &text, nil
}

// Creates the shell needed for accessing the local Kubo node
// Calls readFile and returns a map of the content
func passportFromCID(cid string) (target map[string]interface{}) {
	sh := shell.NewShell("localhost:5001")
	text, err := readFile(sh, cid)
	if err != nil {
		fmt.Println("Error reading the file:", err.Error())
		return
	}
	err = json.Unmarshal([]byte(*text), &target)
	return target
}

// Converts the map from passportFromCID to a string
func getPassport(cid string) string {
	result := passportFromCID(cid)
	jsonStr, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return string(jsonStr)
}

// Adds a file to IPFS
func addFile(sh *shell.Shell, text string) (string, error) {
	return sh.Add(strings.NewReader(text))
}

// Generates a IPNS key (ed25519) and stores it on the local node.
// Names the key after its public key, Returns both the public and private key as text where the private key is formated as a pem file.
func generateKey() (publicKey, privatekey string) {
	// Generates random temp name.
	randomName := randSeq(10)
	// Generates key with temp name.
	publicKey = keyGenerator(randomName)
	// Renames the key to the public key.
	keyRename(randomName, publicKey)

	// Generates the filepath to the PrivateKeys folder
	filePath2 := filepath.Join("PrivateKeys", publicKey+".pem")

	// Exports the private key to a pem file, stores it in the PrivateKeys folder.
	cmd := exec.Command("ipfs", "key", "export", publicKey, "--format=pem-pkcs8-cleartext", "-o", filePath2)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
		fmt.Println("Command output:", string(output))
		return
	}

	// Read the pem file in the PrivateKeys folder to get a string with its content.
	data, err := os.ReadFile(filePath2)
	privatekey = string(data)
	return publicKey, privatekey
}

// Checks if the local node has the given key.
func checkKey(key string) (hasKey bool) {
	cmd := exec.Command("ipfs", "key", "list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
		fmt.Println("Command output:", string(output))
		return
	}
	s := strings.Split(string(output), "\n")
	for _, cKey := range s {
		if key == cKey {
			return true
		}
	}
	return false
}

// Retrieves the given private key from the CA created by https://github.com/LTU-D0020E/d0020e-project-dpp
// If the project is run without a CA this code can be removed.
// If the project is run with a different CA this code needs to be modified.
func retrievePrivateKey(publicKey string) (success, message string) {
	remanEventData := make(map[string]interface{})
	remanEventData["publicKey"] = publicKey
	jsonToCA, err := json.Marshal(remanEventData)
	if err != nil {
		fmt.Println("Error unmarshaling jsonAdd, error code: ", err)
	}
	fmt.Println("Data till CA", string(jsonToCA))
	// Retrieves the private key from the CA
	address := "https://d0020e-project-dpp.vercel.app/api/v1/CA/" + publicKey
	response := outboundCalls(jsonToCA, "GET", address)

	var dataFromCa dataFromCa
	json.Unmarshal([]byte(response), &dataFromCa)
	// Saves the private key to the PrivateKeys folder as a pem file
	if dataFromCa.Success {
		filePath2 := filepath.Join("PrivateKeys", publicKey+".pem")
		err := os.WriteFile(filePath2, []byte(dataFromCa.PrivateKey), 0644)
		if err != nil {
			fmt.Println("Error writing to file, error code: ", err)
		}
		// fmt.Println("\nfilepath2 is :", filePath2, "\nand datafromca is: ", dataFromCa.PrivateKey, "\n")
		// Imports the private key to the local node
		message, err := importPEM(publicKey, filePath2)
		fmt.Println("retrievePrivateKey Message: ", message)
		if err != nil {
			fmt.Println("Error:", err)
			return "false", message
		}
		return "true", message
	} else {
		return "false", dataFromCa.Message
	}
}

// Pin the CID to the local IPFS node
func pinToIPFS(cid string) {
	sh := shell.NewShell("localhost:5001")
	err := sh.Pin(cid)
	if err != nil {
		fmt.Println("Error pinToIPFS : ", err)
	}
}
