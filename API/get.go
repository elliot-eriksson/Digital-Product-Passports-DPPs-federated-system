package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

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

	return &text, nil
}

func passportFromCID(cid, key string) (target map[string]interface{}) {
	sh := shell.NewShell("localhost:5001")
	text, err := readFile(sh, cid)
	if err != nil {
		fmt.Println("Error reading the file:", err.Error())
		return
	}
	err = json.Unmarshal([]byte(*text), &target)
	return target
}

func getPassport(cid, key string) string {
	result := passportFromCID(cid, key)
	jsonStr, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return string(jsonStr)
}

func addFile(sh *shell.Shell, text string) (string, error) {
	return sh.Add(strings.NewReader(text))
}

func generateKey() (publicKey, privatekey string) {
	randomName := randSeq(10)
	publicKey = keyGenerator(randomName)
	keyRename(randomName, publicKey)

	keynamePriv := ".\\PrivateKeys\\" + publicKey + ".pem"

	// Create the directory if it doesn't exist
	err := os.MkdirAll(".\\PrivateKeys\\", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	cmd := exec.Command("ipfs", "key", "export", publicKey, "--format=pem-pkcs8-cleartext", "-o", keynamePriv)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
		fmt.Println("Command output:", string(output))
		return
	}

	data, err := os.ReadFile(keynamePriv)
	privatekey = string(data)
	return publicKey, privatekey
}

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

func retrievePrivateKey(publicKey string) (success, message string) {
	remanEventData := make(map[string]interface{})
	remanEventData["publicKey"] = publicKey
	jsonToCA, err := json.Marshal(remanEventData)
	if err != nil {
		fmt.Println("Error unmarshaling jsonAdd, error code: ", err)
	}
	// fmt.Println("Data till CA", string(test))
	response := sendToCa(jsonToCA, "GET")

	var dataFromCa dataFromCa
	// fmt.Println("Responsen", response)
	json.Unmarshal([]byte(response), &dataFromCa)
	if dataFromCa.Success == "true" {
		filePath := ".\\PrivateKeys\\" + publicKey + ".pem"
		// fmt.Println("data from CA NY:", dataFromCa.PrivateKey)
		err := os.WriteFile(filePath, []byte(dataFromCa.PrivateKey), 0644)
		if err != nil {
			fmt.Println("Error writing to file, error code: ", err)
		}
		// fmt.Println("file wirthe")
		message, err := importPEM(publicKey, filePath)
		// fmt.Println("message:", message)
		if err != nil {
			fmt.Println("Error:", err)
			return "false", message
		}
		return "true", message
	} else {
		return "false", dataFromCa.Message
	}
}
