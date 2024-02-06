package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

func readFile(sh *shell.Shell, cid string) (*string, error) {
	reader, err := sh.Cat(fmt.Sprintf("/ipfs/%s", cid))
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
	data := decryptIt([]byte(*text), key)
	err = json.Unmarshal(data, &target)
	target["cid"] = cid
	return target
}

func getPassport(cid, key string) string {
	result := passportFromCID(cid, key)
	jsonStr, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	fmt.Println("JSON STRING \n", string(jsonStr))
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
