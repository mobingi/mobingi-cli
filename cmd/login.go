package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/mobingilabs/mocli/pkg/util"
	"github.com/spf13/cobra"
)

const (
	credFolder = ".mocli"
	credFile   = "credentials"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "",
	Long:  `Placeholder for the documentation.`,
	Run:   login,
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringP("client-id", "i", "", "client id")
	loginCmd.Flags().StringP("client-secret", "s", "", "client secret")
	loginCmd.Flags().StringP("grant-type", "g", "client_credentials", "grant type (valid values: 'client_credentials', 'password')")
}

func login(cmd *cobra.Command, args []string) {
	hd, _ := homedir.Dir()
	log.Println("home:", hd)
	folder := filepath.Join(hd, credFolder)
	cred := filepath.Join(folder, credFile)
	log.Println(folder, cred)

	token, err := ioutil.ReadFile(cred)
	if err != nil {
		os.MkdirAll(folder, os.ModePerm)
		ioutil.WriteFile(cred, []byte("hello"), 0644)
	}

	log.Println(string(token))

	user, pass := util.GetUserPassword()
	log.Println(user, pass)
}
