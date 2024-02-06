package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"io"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

// Paste here the local path of your computer where the file will be downloaded
const YourLocalPath = "C:/Users/boink/Desktop/test"

// // Paste here the local path of your computer where the file will be downloaded
// const YourLocalPath = "C:/Users/Ellio/Desktop/test"

func addFile(sh *shell.Shell, text string) (string, error) {
	return sh.Add(strings.NewReader(text))
}

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
	pinToIPFS(cid)

	return &text, nil
}

func downloadFile(sh *shell.Shell, cid string) error {
	return sh.Get(cid, YourLocalPath)
}

func pinToIPFS(cid string) {
	// URL of your IPFS node's API
	ipfsAPIURL := "http://localhost:5001/api/v0/pin/add?arg=" + cid

	// Make a POST request to pin the CID
	resp, err := http.Post(ipfsAPIURL, "application/json", bytes.NewBuffer([]byte{}))
	if err != nil {
		fmt.Println("Error pinning CID:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("CID: " + cid + " pinned successfully")
	} else {
		fmt.Println("Failed to pin CID:", resp.Status)
	}

	//this might work need testing
	// sh := shell.NewShell("localhost:5001")
	// sh.Pin(cid)

}

func passportFromCID(cid string) (target map[string]interface{}) {
	sh := shell.NewShell("localhost:5001")
	text, err := readFile(sh, cid)

	if err != nil {
		fmt.Println("Error reading the file:", err.Error())
		// return ""
	}
	data := decryptIt([]byte(*text), "hej")
	err = json.Unmarshal(data, &target)
	target["cid"] = cid
	return target
}

func ipfs(upploadString string) (string, error) {
	sh := shell.NewShell("localhost:5001")

	err := performChecks(sh)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	// 1. Add the "Hello from Launchpad!" text to IPFS
	// fmt.Println("Adding file to IPFS")
	cid, err := addFile(sh, upploadString)
	if err != nil {
		fmt.Println("Error adding file to IPFS:", err.Error())
		return "", err
	}
	fmt.Println("File added with CID:", cid)

	// cid = "QmUNGLqawa7dgDNBSt1yzR9sWphSCPppUYUFAkyKEDzyaH"
	separator()

	// 2. Read the file by using the generated CID
	// fmt.Println("Reading file")
	text, err := readFile(sh, cid)
	if err != nil {
		fmt.Println("Error reading the file:", err.Error())
		return "", err
	}
	fmt.Println("Content of the file:", *text)
	fmt.Println("Content of the file decrypt:", string(decryptIt([]byte(*text), "hej")))

	// separator()
	// // cid = "QmUNGLqawa7dgDNBSt1yzR9sWphSCPppUYUFAkyKEDzyaH"
	// // 3. Download the file to your computer
	// // fmt.Println("Downloading file")
	// err = downloadFile(sh, cid)
	// if err != nil {
	// 	fmt.Println("Error downloading file:", err.Error())
	// 	return "", err
	// }
	// fmt.Println("File downloaded")

	return cid, err
	// separator()
}
