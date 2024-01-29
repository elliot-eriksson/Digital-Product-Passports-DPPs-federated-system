package main

import (
	"fmt"
	"os/exec"

	shell "github.com/ipfs/go-ipfs-api"
)

func main() {

	//SELF := "k51qzi5uqu5dgpie7j0flapmw67becwedlv5vjsvrsp634va9pl4pl3oe0yvyn"
	//LKAB := "k51qzi5uqu5dk0lknwezqu0hrcbgpbbrpynp3r5nh9typbj861k79bu8bud64t"
	//SSAB := "k51qzi5uqu5dlhsqq2mlmroidrca8vuautxhmbcmb5bvmb4g1lvljpj4fanf3x"
	//VOLVO := "k51qzi5uqu5dhqmsy1voi1wegln7cvehdqt7o2n485j451j5mqxpm1rccpzyga"

	//cid := "QmUbd3ZArm3fkLYK37oh17yAML218j4XuVnK4rGbG1b8Sz"

	// Initialize IPFS shell
	sh := shell.NewShell("127.0.0.1:5001")
	//addDataToIPNS(sh, VOLVO, cid)
	fmt.Println(keyGenerator(sh, "samuelsnyckel"))

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
