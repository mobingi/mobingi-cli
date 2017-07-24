package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mobingilabs/mocli/client"
	"github.com/mobingilabs/mocli/pkg/check"
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/credentials"
	d "github.com/mobingilabs/mocli/pkg/debug"
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
	idsec := &credentials.ClientIdSecret{
		Id:     cli.GetCliStringFlag(cmd, "client-id"),
		Secret: cli.GetCliStringFlag(cmd, "client-secret"),
	}

	err := idsec.EnsureInput(false)
	if err != nil {
		check.ErrorExit(err, 1)
	}

	// id := cli.GetCliStringFlag(cmd, "client-id")
	// secret := cli.GetCliStringFlag(cmd, "client-secret")
	grant := cli.GetCliStringFlag(cmd, "grant-type")
	// user := cli.GetCliStringFlag(cmd, "username")
	// pass := cli.GetCliStringFlag(cmd, "password")

	var m map[string]interface{}
	var p *authPayload
	cnf := client.NewConfig(cmd)
	c := client.NewClient(cnf)

	if grant == "client_credentials" {
		p = &authPayload{
			ClientId:     idsec.Id,
			ClientSecret: idsec.Secret,
			GrantType:    grant,
		}
	}

	if grant == "password" {
		up := &credentials.UserPass{
			Username: cli.GetCliStringFlag(cmd, "username"),
			Password: cli.GetCliStringFlag(cmd, "password"),
		}

		in, err := up.EnsureInput(false)
		if err != nil {
			check.ErrorExit(err, 1)
		}

		if in[1] {
			fmt.Println("\n") // new line after the password input
		}

		p = &authPayload{
			ClientId:     idsec.Id,
			ClientSecret: idsec.Secret,
			Username:     up.Username,
			Password:     up.Password,
			GrantType:    grant,
		}
	}

	// should not be nil when `grant_type` is valid
	if p == nil {
		check.ErrorExit("Invalid argument(s). See `help` for more information.", 1)
	}

	payload, err := json.Marshal(p)
	check.ErrorExit(err, 1)

	resp, body, errs := c.PostU("/access_token", string(payload))
	check.ErrorExit(errs, 1)
	serr := check.ResponseError(resp, body)
	check.ErrorExit(serr, 1)

	err = json.Unmarshal(body, &m)
	check.ErrorExit(err, 1)
	token, found := m["access_token"]
	if !found {
		check.ErrorExit("cannot find access token", 1)
	}

	// always overwrite file
	err = credentials.SaveToken(fmt.Sprintf("%s", token))
	check.ErrorExit(err, 1)
	d.Info("Login successful.")
}
