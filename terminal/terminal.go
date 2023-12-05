package terminal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"

	"github.com/Ovec/pwd_manager/crypt"
	"github.com/Ovec/pwd_manager/password"
	"github.com/Ovec/pwd_manager/random"
)

func GetPassword(r *bufio.Reader) string {
	fmt.Print("-> ")
	bytePassword, err := term.ReadPassword(syscall.Stdin)

	if err != nil {
		fmt.Println("Error reading file:", err)
	}

	return string(bytePassword)
}

func Handle(r *bufio.Reader, passwordPairs map[string]string, filePath string, key string) {
	var newKey = ""
	var state = ""
	fmt.Println("Press (L) for list, (A) for Add, (Q) for Quit")
	fmt.Println("---------------------")

	for {
		fmt.Print("-> ")
		text, _ := r.ReadString('\n')
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
			newValue, err := password.Generate(random.RandomNumber(8, 16))
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			_, ok := passwordPairs[newKey]

			if ok {
				fmt.Println("Key exists already")
				fmt.Printf("Id: %s - %s\n", newKey, passwordPairs[newKey])
			} else {
				fmt.Printf("Id: %s - %s\n", newKey, newValue)

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

				err = os.WriteFile(filePath, cipherText, 0644)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}

			}

			state = ""
		}

		if strings.Compare("L", text) == 0 {
			for key, value := range passwordPairs {
				fmt.Printf("Id: %s - %s\n", key, value)
			}
		}
	}
}
