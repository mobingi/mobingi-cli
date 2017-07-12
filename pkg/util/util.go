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
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	credFolder = ".mocli"      // folder name for config file(s), created in home folder
	credFile   = "credentials" // we store access token here
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
	folder := filepath.Join(hd, credFolder)
	_ = os.Mkdir(folder, os.ModePerm)
	cred := filepath.Join(folder, credFile)
	return ioutil.WriteFile(cred, []byte(token), 0644)
}

func GetToken() ([]byte, error) {
	hd, _ := homedir.Dir()
	folder := filepath.Join(hd, credFolder)
	cred := filepath.Join(folder, credFile)
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

func ErrorExit(err string, code int) {
	log.Println("Error:", err)
	os.Exit(code)
}

func ResponseError(r gorequest.Response, m map[string]interface{}) string {
	var err string
	if r.StatusCode != 200 {
		err = err + "[" + r.Status + "]"
	}

	if c, found := m["code"]; found {
		err = err + "[" + fmt.Sprintf("%s", c) + "]"
	}

	if e, found := m["error"]; found {
		err = err + fmt.Sprintf(" %s:", e)
	}

	if s, found := m["message"]; found {
		err = err + fmt.Sprintf(" %s", s)
	}

	return err
}

func Indent(count int) string {
	pad := ""
	for i := 0; i < count; i++ {
		pad += " "
	}

	return pad
}
