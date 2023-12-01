package main

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
)

// todo
// create data structure for pwd storing - DONE
// create ui for adding and removing data
// use storage for adding and removing data
// implement crypt and decrypt
// implement pwd
// datastructure
// password and id

var filePath = "storage"
var newValue = ""
var newKey = ""
var state = ""

const storageKey = "storageKey"

func Generate(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-=_+"

	randomBytes := make([]byte, length)
	for i := range randomBytes {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		randomBytes[i] = charset[n.Int64()]
	}

	return string(randomBytes), nil
}

func main() {
	passwordPairs := map[string]string{}

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create the file
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer file.Close()

		fmt.Println("Storage not found, creating new", filePath)
	} else {
		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		err = json.Unmarshal(fileContent, &passwordPairs)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return
		}

	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Simple Shell, press (L) for list and (A) for Add")
	fmt.Println("---------------------")

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if strings.Compare("q", text) == 0 {
			fmt.Println("Bye, have a nice day")
			os.Exit(0)
		}

		if strings.Compare("A", text) == 0 {
			println(state)
			fmt.Println("Enter new key")
			state = "A"
		}

		if strings.Compare("A", state) == 0 && strings.Compare("A", text) != 0 {
			println(state)
			fmt.Println("Enter password")
			newKey = text
			state = "B"
		}

		if strings.Compare("B", state) == 0 && strings.Compare("B", text) != 0 {
			println(state)
			newValue, err := Generate(8)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			_, ok := passwordPairs[newKey]

			// Check the result
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

				err = ioutil.WriteFile(filePath, jsonData, 0644)
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
