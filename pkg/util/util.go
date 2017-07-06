package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func GetUserPassword() (string, string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	user, _ := reader.ReadString('\n')
	fmt.Print("Password: ")
	pass, _ := terminal.ReadPassword(int(syscall.Stdin))
	return strings.TrimSpace(user), strings.TrimSpace(string(pass))
}
