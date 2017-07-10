package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/util"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login to Mobingi API",
	Long:  `Login to Mobingi API server. If 'grant_type' is set to 'password', you will be prompted to enter your username and password.`,
	Run:   login,
}

type authPayload struct {
	ClientId     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	GrantType    string `json:"grant_type,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringP("client-id", "i", "", "client id (required)")
	loginCmd.Flags().StringP("client-secret", "s", "", "client secret (required)")
	loginCmd.Flags().StringP("grant-type", "g", "client_credentials", "grant type (valid values: 'client_credentials', 'password')")
	loginCmd.Flags().StringP("username", "u", "", "user name")
	loginCmd.Flags().StringP("password", "p", "", "password")
}

func login(cmd *cobra.Command, args []string) {
	id := util.GetCliStringFlag(cmd, "client-id")
	secret := util.GetCliStringFlag(cmd, "client-secret")

	if id == "" {
		id = util.GetClientId()
	}

	if id == "" {
		log.Println("Client id cannot be empty.")
		os.Exit(1)
	}

	if secret == "" {
		secret = util.GetClientSecret()
	}

	if secret == "" {
		log.Println("Client secret cannot be empty.")
		secret = util.GetClientSecret()
	}

	var m map[string]interface{}
	var user, pass string
	var p *authPayload
	c := cli.New(util.GetCliStringFlag(cmd, "api-version"))

	if util.GetCliStringFlag(cmd, "grant-type") == "client_credentials" {
		p = &authPayload{
			ClientId:     id,
			ClientSecret: secret,
			GrantType:    util.GetCliStringFlag(cmd, "grant-type"),
		}
	}

	if util.GetCliStringFlag(cmd, "grant-type") == "password" {
		// user, pass = util.GetUserPassword()
		user = util.GetCliStringFlag(cmd, "username")
		if user == "" {
			user = util.GetUsername()
		}

		if user == "" {
			log.Println("Username cannot be empty.")
			os.Exit(1)
		}

		pass = util.GetCliStringFlag(cmd, "password")
		if pass == "" {
			pass = util.GetPassword()
		}

		if pass == "" {
			log.Println("Password cannot be empty.")
			os.Exit(1)
		}

		fmt.Println("\n") // new line after the password input
		p = &authPayload{
			ClientId:     id,
			ClientSecret: secret,
			GrantType:    util.GetCliStringFlag(cmd, "grant-type"),
			Username:     user,
			Password:     pass,
		}
	}

	// should not be nil when `grant_type` is valid
	if p == nil {
		util.PrintErrorAndExit("Invalid argument(s). See `help` for more information.", 1)
	}

	resp, body, errs := c.PostJSON(c.RootUrl+"/access_token", p)
	if errs != nil {
		log.Println("Error(s):", errs)
		os.Exit(1)
	}

	err := json.Unmarshal(body, &m)
	if err != nil {
		util.PrintErrorAndExit(err.Error(), 1)
	}

	serr := util.BuildRequestError(resp, m)
	if serr != "" {
		util.PrintErrorAndExit(serr, 1)
	}

	token, found := m["access_token"]
	if !found {
		util.PrintErrorAndExit("Internal error.", 1)
	}

	// always overwrite file
	util.SaveToken(fmt.Sprintf("%s", token))
	log.Println("Login successful.")
}
