package main

import (
	"fmt"
	"os/exec"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

func main() {

	// 107 = k
	// 81 = Q
	//tst := "string"
	//fmt.Println(reflect.TypeOf(tst))

	//SELF := "k51qzi5uqu5dgpie7j0flapmw67becwedlv5vjsvrsp634va9pl4pl3oe0yvyn"
	//LKAB := "k51qzi5uqu5dk0lknwezqu0hrcbgpbbrpynp3r5nh9typbj861k79bu8bud64t"
	//SSAB := "k51qzi5uqu5dlhsqq2mlmroidrca8vuautxhmbcmb5bvmb4g1lvljpj4fanf3x"
	//VOLVO := "k51qzi5uqu5dhqmsy1voi1wegln7cvehdqt7o2n485j451j5mqxpm1rccpzyga"

	//cid := "QmUbd3ZArm3fkLYK37oh17yAML218j4XuVnK4rGbG1b8Sz"

	// Initialize IPFS shell
	sh := shell.NewShell("127.0.0.1:5001")
	//addDataToIPNS(sh, VOLVO, cid)
	//fmt.Println(keyGenerator(sh, "samuelsnyckel"))

	thisvar, err := lsIPNS(sh, "k51qzi5uqu5dgpie7j0flapmw67becwedlv5vjsvrsp634va9pl4pl3oe0yvyn")
	if err != nil {
		fmt.Println("big error oh no")
	}
	content, contentlength := splitListContent(thisvar)
	fmt.Println(content, contentlength)
	// for x = 0, x < len(list)
	// 	fmt.Println("CID", x, content[x])

}

func separator() {
	fmt.Println("-------------------")
}

func splitListContent(Content string) ([]string, int) {
	temp := strings.Split(Content, "\n")
	lenvar := len(temp) - 1
	fmt.Println("your splitted CIDs are: ", temp)
	fmt.Println("you have ", lenvar, " different CIDs in this directory")
	return temp, lenvar
}

func catContent(CID []string, length int) (string, error) {

	cmd := exec.Command("ipfs", "cat")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return "", err
	}
	return string(output), err

}

// om den returnar "" vet vi att det inte är ett directory, utan bara en fil av något slag
func lsIPNS(sh *shell.Shell, key string) (string, error) {
	//simple check for if the sent link is an directory or a CID.
	if key[0] == 107 { // checks if the first char is k
		fmt.Println("This is a public key ", key)
		separator()
		key = "/ipns/" + key
	}
	if key[0] == 81 { // checks if the first char is Q
		fmt.Println("This is a CID ", key)
		separator()
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
		separator()
	} else {
		fmt.Println("this is a directory")
		separator()
	}

	fmt.Println(string(output))
	separator()
	return string(output), err

}

func keyGenerator(sh *shell.Shell, keyAlias string) (string, error) {
	cmd := exec.Command("ipfs", "key", "gen", keyAlias)
	output, err := cmd.CombinedOutput()
	if err != nil {
		//return "", fmt.Errorf("failed to generate a public key %v, output: %s", err, output)
		fmt.Println(string(output))
		return "", err
	}
	fmt.Println("The public key value is: ", string(output))
	return string(output), nil
}

func resolveKeyPointer(sh *shell.Shell, key string) (string, error) {
	cmd := exec.Command("ipfs", "resolve", ipnsKeyToCMD(key))
	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("failed to resolve public key %v, output: %s", err, output)
	}
	fmt.Println(string(output))
	return string(output), nil
}

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
	cmd := exec.Command("ipfs", "name", "publish", hashToCMD(key), cidToCMD(cid))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to publish IPNS record: %v, output: %s", err, output)
	}
	return string(output), nil
}
