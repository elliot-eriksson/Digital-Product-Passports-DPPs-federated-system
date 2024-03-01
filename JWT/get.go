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
		// return ""
	}
	// fmt.Println("passportFromCID file read", text)
	// data := decryptIt([]byte(*text), key)
	err = json.Unmarshal([]byte(*text), &target)
	// fmt.Println("passportFromCID Unmarsal", target)
	// target["cid"] = cid
	return target
}

func getPassport(cid, key string) string {
	// fmt.Println("GetPASSPORT ENTERD")
	result := passportFromCID(cid, key)
	jsonStr, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	// fmt.Println("JSON STRING \n", string(jsonStr))
	return string(jsonStr)
}

func getSensetive(cid, key, keySen string) {
	result := passportFromCID(cid, key)
	sensetive := passportFromCID(fmt.Sprintf("%v", result["CID_sen"]), keySen)
	jsonStr, err := json.Marshal(sensetive)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("JSON STRING sensetive \n", string(jsonStr))
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

	if err != nil {
		fmt.Println(string(output))
		return
	}

	data, err := os.ReadFile(keynamePriv)
	fmt.Println("DATA GREJS", string(data))
	privatekey = string(data)

	return publicKey, privatekey
}
