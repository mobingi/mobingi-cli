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
	Long:  `Login to Mobingi API server.`,
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
	loginCmd.Flags().StringP("client-id", "i", "", "client id")
	loginCmd.Flags().StringP("client-secret", "s", "", "client secret")
	loginCmd.Flags().StringP("grant-type", "g", "client_credentials", "grant type (valid values: 'client_credentials', 'password')")
}

func login(cmd *cobra.Command, args []string) {
	if util.GetCliStringFlag(cmd, "client-id") == "" {
		log.Println("Client id not provided. See `help login` for more information.")
		os.Exit(1)
	}

	if util.GetCliStringFlag(cmd, "client-secret") == "" {
		log.Println("Client secret not provided. See `help login` for more information.")
		os.Exit(1)
	}

	var m map[string]interface{}
	var user, pass string
	var p *authPayload
	c := cli.New()

	if util.GetCliStringFlag(cmd, "grant-type") == "client_credentials" {
		p = &authPayload{
			ClientId:     util.GetCliStringFlag(cmd, "client-id"),
			ClientSecret: util.GetCliStringFlag(cmd, "client-secret"),
			GrantType:    util.GetCliStringFlag(cmd, "grant-type"),
		}
	}

	if util.GetCliStringFlag(cmd, "grant-type") == "password" {
		user, pass = util.GetUserPassword()
		fmt.Println("\n")
		p = &authPayload{
			ClientId:     util.GetCliStringFlag(cmd, "client-id"),
			ClientSecret: util.GetCliStringFlag(cmd, "client-secret"),
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
