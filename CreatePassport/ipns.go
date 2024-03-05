package main

import (
	"fmt"
	"os/exec"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

func importPrivKey(CID string) {
	// download privkey from an IPFS node
	cmd := exec.Command("ipfs", "get", CID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return
	}
	fmt.Println("Private key successfully downloaded under local path: ./")
	// import private key locally, creating the option to use it yourself
	cmd = exec.Command("ipfs", "key", "import", "newtestkey", "-f", "pem-pkcs8-cleartext", "privkey.pem")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return
	}
	fmt.Println(string(output))

}

func uploadPrivKey() string {
	cmd := exec.Command("ipfs", "add", "privkey.pem")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return ""
	}
	fmt.Println(string(output))
	return string(output)
}

func exportPrivKey(keyAlias string) {
	// ipfs key export testkey --format=pem-pkcs8-cleartext -o privkey.pem
	cmd := exec.Command("ipfs", "key", "export", keyAlias, "--format=pem-pkcs8-cleartext", "-o", "privkey.pem")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return
	}
	fmt.Println("Private key successfully exported to path: ./privkey.pem")
}

func LinkMakesAdd(str string) string {
	// echo "new test string to upload" | ipfs add -w
	sh := shell.NewShell("localhost:5001")
	res, _ := addFile(sh, str)
	return res

}

// Splitting the different files into their own string.
func splitListContent(Content string) ([]string, int) {
	temp := strings.Split(Content, "\n")
	lenvar := len(temp) - 1
	//fmt.Println("your splitted CIDs are: ", temp)
	//fmt.Println("you have ", lenvar, " different CIDs in this directory")
	return temp, lenvar
}

// Printing and retriving the information from an IPNS pointer
func catContent(CID []string, length int) {
	var splitIndex []string
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
			if err != nil {
				fmt.Println(string(output))
				return
			}
			fmt.Println("The content of the file", splitIndex[i+2], "is:", string(output))
		}
	} else {
		tempString := splitIndex[0]
		tempString = tempString[6:]
		cmd := exec.Command("ipfs", "cat", tempString)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(output))
			return
		}
		fmt.Println("The content of the file", tempString, "is:", string(output))
	}

}

// Helper function to find out if its and directory or just an file.
// Also retrieves the pointer data
func lsIPNS(sh *shell.Shell, key string) (string, error) {
	//simple check for if the sent link is an directory or a CID.
	tempkey := key
	if key[0] == 107 { // checks if the first char is k
		//fmt.Println("This is a public key ", key)
		key = "/ipns/" + key
	}
	if key[0] == 81 { // checks if the first char is Q
		//fmt.Println("This is a CID ", key)
		key = "/ipfs/" + key
	}
	cmd := exec.Command("ipfs", "ls", key)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return "", err
	}
	if string(output) == "" {
		fmt.Println("this is a file")
		newOut, err := resolveKeyPointer(sh, tempkey)
		if err != nil {
			fmt.Println(string(newOut))
			return "", err
		}
		fmt.Println(newOut)
		return string(newOut), err
	} else {
		fmt.Println("this is a directory")
	}

	fmt.Println(string(output))
	return string(output), err

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
	fmt.Println("The public key value is: ", string(output))
	return string(output)
}

func keyRename(newAlias string) {
	cmd := exec.Command("ipfs", "key", "rename", "tempAlias", newAlias)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return
	}
	fmt.Println(string(output))
	return
}

func keyRenameLinkMakes(newAlias string, input string) {
	cmd := exec.Command("ipfs", "key", "rename", input, "LinkMakes_"+newAlias)
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
func addDataToIPNS(key string, cid string) string {
	cmd := exec.Command("ipfs", "name", "publish", hashToCMD(key), cidToCMD(cid))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output)
		//return fmt.Errorf("failed to publish IPNS record: %v, output: %s", err, output)
	}
	return string(output)
}
