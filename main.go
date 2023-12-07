package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	"github.com/Ovec/pwd_manager/crypt"
	"github.com/Ovec/pwd_manager/filesystem"
	"github.com/Ovec/pwd_manager/str"
	"github.com/Ovec/pwd_manager/terminal"
)

var validOss = []string{"linux", "windows", "darwin"}

const storeageDir = "/applications/pwd_manager"
const storageFile = "storage"
const saltFile = "salt"

func main() {
	if !str.ContainsString(validOss, runtime.GOOS) {
		fmt.Println("Operationg system", runtime.GOOS, "is not supported")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	passwordPairs := map[string]string{}

	dataDir, err := filesystem.DataDir(storeageDir)

	if err != nil {
		fmt.Println("Error finding home:", err)
		return
	}

	err = filesystem.CreateFiles(dataDir, []string{storageFile, saltFile})

	if err != nil {
		fmt.Println("Error creating storage:", err)
		return
	}

	salt, err := os.ReadFile(dataDir + "/" + saltFile)

	if err != nil {
		fmt.Println("Error reading salt file:", err)
		return
	}

	if len(salt) == 0 {
		newSalt, err := crypt.GenerateSalt(64)
		salt = []byte(newSalt)

		if err != nil {
			fmt.Println("Error creating salt:", err)
			return
		}

		err = os.WriteFile(dataDir+"/"+saltFile, []byte(salt), 0644)

		if err != nil {
			fmt.Println("Error string salt:", err)
			return
		}
	}

	fmt.Println("Enter your master password")

	password := terminal.GetPassword(reader)
	key := string(crypt.GenerateAESKeyFromPassword([]byte(password), salt, 10000))

	storage, err := os.ReadFile(dataDir + "/" + storageFile)
	if err != nil {
		fmt.Println("Error reading storage:", err)
		return
	}

	if len(storage) > 0 {
		plaintext, err := crypt.DecryptAES(storage, []byte(key))
		if err != nil {
			fmt.Println("Error decrypting plaintext:", err)
			return
		}

		err = json.Unmarshal(plaintext, &passwordPairs)
		if err != nil {
			fmt.Println("Wrong password")
			os.Exit(0)
		}
	}

	terminal.Handle(reader, passwordPairs, dataDir+"/"+storageFile, key)
}
