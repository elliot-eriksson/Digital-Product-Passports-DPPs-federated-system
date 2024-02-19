package main

import (
	"fmt"
	"os/exec"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

// Splitting the different files into their own string.
func splitListContent(Content string) ([]string, int) {
	temp := strings.Split(Content, "\n")
	lenvar := len(temp) - 1
	//fmt.Println("your splitted CIDs are: ", temp)
	//fmt.Println("you have ", lenvar, " different CIDs in this directory")
	return temp, lenvar
}

// Printing and retriving the information from an IPNS pointer
func catContent(CID []string, length int) (contentIndex []string) {
	var splitIndex []string
	var appendvalue string
	// Trims unnecessary spaces and content from the CID-array
	//print(CID)
	for i := 0; i < length; i++ {
		splitIndex = append(splitIndex, strings.Split(string(CID[i]), " ")...)
	}
	fmt.Println("--------->", splitIndex, " och längden av splitIndex är: ", len(splitIndex))
	// Splits the array to be able to print out the CID content

	if len(splitIndex) >= 2 {
		for i := 0; i < len(splitIndex); i += 3 {
			fmt.Println("File", splitIndex[i+2], " has CID :", splitIndex[i])
			cmd := exec.Command("ipfs", "cat", splitIndex[i])
			output, err := cmd.CombinedOutput()
			appendvalue = string(output)
			contentIndex = append(contentIndex, appendvalue)

			if err != nil {
				fmt.Println(string(output))
				return
			}
			fmt.Println("The content of the file", splitIndex[i+2], "is:", string(output))
		}
	}
	return contentIndex
}

// Helper function to find out if its and directory or just an file.
// Also retrieves the pointer data
func lsIPNS(key string) string {
	//simple check for if the sent link is an directory or a CID.
	//tempkey := key

	cmd := exec.Command("ipfs", "ls", key)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return ""
	}
	return string(output)

}

// Generates public key
func keyGenerator(keyAlias string) string {
	cmd := exec.Command("ipfs", "key", "gen", keyAlias)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return string(output)
	}
	output = []byte(strings.TrimSuffix(string(output), "\n"))
	return string(output)
}

func keyRename(oldAlias, newAlias string) {
	cmd := exec.Command("ipfs", "key", "rename", oldAlias, newAlias)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return
	}
	fmt.Println(string(output))
	return
}

// Find out to what public key the CID is pointing to.
func resolveKeyPointer(sh *shell.Shell, key string) (string, error) {
	cmd := exec.Command("ipfs", "resolve", ipnsKeyToCMD(key))
	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("failed to resolve public key %v, output: %s", err, output)
	}
	fmt.Println(string(output))
	return string(output), nil
}

// Helper function to get the IPNS key to a format for the terminal
func ipnsKeyToCMD(key string) string {
	key = "/ipns/" + key
	return key
}

// Converts a public key to a string usable in the terminal
func hashToCMD(key string) string {
	key = "--key=" + key
	fmt.Println(key)
	return key
}

// Converts the IPFS CID to a string usable in the terminal
func cidToCMD(cid string) string {
	cid = "/ipfs/" + cid
	fmt.Println(cid)
	return cid
}

// Uploads data to IPNS and return that adress, also does the same when you want to update information.
func addDataToIPNS(sh *shell.Shell, key string, cid string) (string, error) {

	fmt.Println(key)
	if key == "" {
		fmt.Println("ERROR: key is empty")
		return "", nil
	}
	cmd := exec.Command("ipfs", "name", "publish", hashToCMD(key), cidToCMD(cid))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to publish IPNS record: %v, output: %s", err, output)
	}
	return string(output), nil
}
