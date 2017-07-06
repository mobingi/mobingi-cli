package cmd

import (
	"io/ioutil"
	"log"

	"github.com/mitchellh/go-homedir"
	"github.com/mobingilabs/mocli/pkg/util"
	"github.com/spf13/cobra"
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
	log.Println("login here")
	hd, _ := homedir.Dir()
	log.Println("home:", hd)
	cred := hd + `/.mocli/credentials`

	token, err := ioutil.ReadFile(cred)
	if err != nil {
		log.Println(err)
		ioutil.WriteFile(hd+`/.mocli/credentials`, []byte("hello"), 0644)
	}

	log.Println(string(token))

	user, pass := util.GetUserPassword()
	log.Println(user, pass)
}
