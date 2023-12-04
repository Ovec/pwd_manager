package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Ovec/pwd_manager/crypt"
	"github.com/Ovec/pwd_manager/terminal"
)

var filePath = "storage"
var key = ""
var salt = "some nice salt should be here"

func main() {
	reader := bufio.NewReader(os.Stdin)
	passwordPairs := map[string]string{}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("Storage not found, creating new")
		fmt.Println("Enter your master password")

		password := terminal.GetPassword(reader)
		key = string(crypt.GenerateAESKeyFromPassword([]byte(password), []byte(salt), 10000))

		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer file.Close()

	} else {
		fmt.Println("Enter your master password")

		password := terminal.GetPassword(reader)
		key = string(crypt.GenerateAESKeyFromPassword([]byte(password), []byte(salt), 10000))

		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		if len(fileContent) > 0 {
			plaintext, err := crypt.DecryptAES(fileContent, []byte(key))
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

	}

	terminal.Handle(reader, passwordPairs, filePath, key)
}
