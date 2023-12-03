package terminal

import (
	"bufio"
	"fmt"
	"syscall"

	"golang.org/x/term"
)

func GetPassword(r *bufio.Reader) string {
	fmt.Print("-> ")
	bytePassword, err := term.ReadPassword(syscall.Stdin)

	if err != nil {
		fmt.Println("Error reading file:", err)
	}

	return string(bytePassword)
}
