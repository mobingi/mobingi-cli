package util

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/mobingilabs/mocli/pkg/constants"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

func ClientId() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Client ID: ")
	id, _ := reader.ReadString('\n')
	return strings.TrimSpace(id)
}

func ClientSecret() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Client secret: ")
	secret, _ := reader.ReadString('\n')
	return strings.TrimSpace(secret)
}

func Username() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	user, _ := reader.ReadString('\n')
	return strings.TrimSpace(user)
}

func Password() string {
	fmt.Print("Password: ")
	pass, _ := terminal.ReadPassword(int(syscall.Stdin))
	return string(pass)
}

func SaveToken(token string) error {
	hd, _ := homedir.Dir()
	folder := filepath.Join(hd, constants.CRED_FOLDER)
	_ = os.Mkdir(folder, os.ModePerm)
	cred := filepath.Join(folder, constants.CRED_FILE)
	return ioutil.WriteFile(cred, []byte(token), 0644)
}

func GetToken() ([]byte, error) {
	hd, _ := homedir.Dir()
	folder := filepath.Join(hd, constants.CRED_FOLDER)
	cred := filepath.Join(folder, constants.CRED_FILE)
	return ioutil.ReadFile(cred)
}

func GetCliStringFlag(cmd *cobra.Command, f string) string {
	s := cmd.Flag(f).DefValue
	if cmd.Flag(f).Changed {
		s = cmd.Flag(f).Value.String()
	}

	return s
}

func GetCliIntFlag(cmd *cobra.Command, f string) int {
	s := cmd.Flag(f).DefValue
	if cmd.Flag(f).Changed {
		s = cmd.Flag(f).Value.String()
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return v
}

func Indent(count int) string {
	pad := ""
	for i := 0; i < count; i++ {
		pad += " "
	}

	return pad
}

func WriteToFile(f string, contents []byte) error {
	err := ioutil.WriteFile(f, contents, 0644)
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("output written to %s", f))
	return nil
}
