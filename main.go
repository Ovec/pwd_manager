package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Ovec/pwd_manager/crypt"
	"github.com/Ovec/pwd_manager/password"
	"github.com/Ovec/pwd_manager/random"
	"github.com/Ovec/pwd_manager/terminal"
)

var filePath = "storage"
var newKey = ""
var state = ""
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

	fmt.Println("Press (L) for list, (A) for Add, (Q) for Quit")
	fmt.Println("---------------------")

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if strings.Compare("Q", text) == 0 {
			fmt.Println("Bye, have a nice day")
			os.Exit(0)
		}

		if strings.Compare("A", text) == 0 {
			fmt.Println("Enter new key")
			state = "A"
		}

		if strings.Compare("A", state) == 0 && strings.Compare("A", text) != 0 {
			newKey = text
			state = "B"
		}

		if strings.Compare("B", state) == 0 && strings.Compare("B", text) != 0 {
			println(state)
			newValue, err := password.Generate(random.RandomNumber(8, 16))
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			_, ok := passwordPairs[newKey]

			if ok {
				fmt.Println("Already in, sorry bro")
				fmt.Printf("Id: %s - %s\n", newKey, passwordPairs[newKey])
			} else {
				fmt.Println(newKey)
				fmt.Println(newValue)

				passwordPairs[newKey] = newValue

				jsonData, err := json.Marshal(passwordPairs)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}

				cipherText, err := crypt.EncryptAES([]byte(jsonData), []byte(key))
				if err != nil {
					fmt.Println("Error encrypting plaintext:", err)
					return
				}

				fmt.Println(string(cipherText))

				err = os.WriteFile(filePath, cipherText, 0644)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}

			}

			state = ""
		}

		if strings.Compare("L", text) == 0 {
			println(state)
			for key, value := range passwordPairs {
				fmt.Printf("Id: %s - %s\n", key, value)
			}
		}
	}
}
