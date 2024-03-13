package main

import (
	"fmt"
	"os/exec"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

// Imports a ipns key from a pem file, to the local Kubo node.
func importPEM(publicKey, filePath string) (string, error) {
	cmd := exec.Command("ipfs", "key", "import", publicKey, "-f", "pem-pkcs8-cleartext", filePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}
	return string(output), err
}

// Splitting the different files into their own string.
func splitListContent(Content string) ([]string, int) {
	temp := strings.Split(Content, "\n")
	lenvar := len(temp) - 1
	return temp, lenvar
}

// Retrieves the given IPNS public keys content.
func catRemanContent(key string) string {
	cmd := exec.Command("ipfs", "cat", ipnsKeyToCMD(key))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	return string(output)
}

// Printing and retriving the information from an IPNS pointer
func catContent(CID []string, length int) (contentIndex []string) {
	var splitIndex []string
	var appendvalue string
	// Trims unnecessary spaces and content from the CID-array
	for i := 0; i < length; i++ {
		splitIndex = append(splitIndex, strings.Split(string(CID[i]), " ")...)
	}
	// Splits the array to be able to print out the CID content
	if len(splitIndex) >= 2 {
		for i := 0; i < len(splitIndex); i += 3 {
			cmd := exec.Command("ipfs", "cat", splitIndex[i])
			output, err := cmd.CombinedOutput()
			appendvalue = string(output)
			contentIndex = append(contentIndex, appendvalue)

			if err != nil {
				fmt.Println(string(output))
				return
			}
		}
	}
	return contentIndex
}

// Helper function to find out if its and directory or just an file.
// Also retrieves the pointer data
func lsIPNS(key string) string {
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

// Renames the given IPNS key
func keyRename(oldAlias, newAlias string) {
	cmd := exec.Command("ipfs", "key", "rename", oldAlias, newAlias)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
	}
}

// Find out to what public key the CID is pointing to.
func resolveKeyPointer(sh *shell.Shell, key string) (string, error) {
	cmd := exec.Command("ipfs", "resolve", ipnsKeyToCMD(key))
	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("failed to resolve public key %v, output: %s", err, output)
	}
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
	return key
}

// Converts the IPFS CID to a string usable in the terminal
func cidToCMD(cid string) string {
	cid = "/ipfs/" + cid
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
