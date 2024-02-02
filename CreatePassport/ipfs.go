package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"io"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

// Paste here the local path of your computer where the file will be downloaded
const YourLocalPath = "C:/Users/boink/Desktop/test"

// Paste here your public key
// Go to powershell run the following command: ipfs key list -l
const YourPublicKey = "k51qzi5uqu5djc71p3quno2nypbts7k7t14el81gwjpxsjksp25kbbl22n70rh"

// // Paste here the local path of your computer where the file will be downloaded
// const YourLocalPath = "C:/Users/Ellio/Desktop/test"

// // Paste here your public key
// // Go to powershell run the following command: ipfs key list -l
// const YourPublicKey = "k51qzi5uqu5dk93h4gqv2fml1x95vc92hbh6tr5atdo5fgc2623qvfak4o6qe3"

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

func addToIPNS(sh *shell.Shell, cid string) error {
	var lifetime time.Duration = 50 * time.Hour
	var ttl time.Duration = 1 * time.Microsecond

	_, err := sh.PublishWithDetails(cid, YourPublicKey, lifetime, ttl, true)
	return err
}

func resolveIPNS(sh *shell.Shell) (string, error) {
	return sh.Resolve(YourPublicKey)
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
		fmt.Println("CID pinned successfully")
	} else {
		fmt.Println("Failed to pin CID:", resp.Status)
	}
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

	separator()
	// cid = "QmUNGLqawa7dgDNBSt1yzR9sWphSCPppUYUFAkyKEDzyaH"
	// 3. Download the file to your computer
	// fmt.Println("Downloading file")
	err = downloadFile(sh, cid)
	if err != nil {
		fmt.Println("Error downloading file:", err.Error())
		return "", err
	}
	// fmt.Println("File downloaded")

	return cid, err
	// separator()

	// // 4. Publish the file to IPNS
	// fmt.Println("Adding file to IPNS")
	// err = addToIPNS(sh, cid)
	// if err != nil {
	// 	fmt.Println("Error publishing to IPNS:", err.Error())
	// 	return
	// }
	// fmt.Println("File added to IPNS")

	// separator()

	// // 5. Resolve IPNS based on your public key
	// fmt.Println("Resolving file in IPNS")
	// result, err := resolveIPNS(sh)
	// if err != nil {
	// 	fmt.Println("Error resolving IPNS:", err.Error())
	// 	return
	// }

	// fmt.Println("IPNS is pointing to:", result)
}
