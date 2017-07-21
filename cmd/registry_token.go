package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/mobingilabs/mocli/pkg/util"
	"github.com/spf13/cobra"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "get registry token",
	Long:  `Get registry token.`,
	Run:   token,
}

func init() {
	registryCmd.AddCommand(tokenCmd)
	tokenCmd.Flags().String("account", "", "subuser name")
	tokenCmd.Flags().String("username", "", "username (account subuser)")
	tokenCmd.Flags().String("password", "", "password (account subuser)")
	tokenCmd.Flags().String("service", "Mobingi Docker Registry", "service for authentication")
	tokenCmd.Flags().String("scope", "", "scope for authentication")
}

func token(cmd *cobra.Command, args []string) {
	passin := false
	user := util.GetCliStringFlag(cmd, "username")
	pass := util.GetCliStringFlag(cmd, "password")
	if user == "" {
		user = util.Username()
	}

	if user == "" {
		util.CheckErrorExit("username cannot be empty", 1)
	}

	if pass == "" {
		pass = util.Password()
		passin = true
	}

	if pass == "" {
		util.CheckErrorExit("password cannot be empty", 1)
	}

	if passin {
		fmt.Println("\n") // new line after the password input
	}

	var Url *url.URL
	Url, err := url.Parse(util.GetCliStringFlag(cmd, "url"))
	if err != nil {
		util.CheckErrorExit(err, 1)
	}

	Url.Path += "/" + util.GetCliStringFlag(cmd, "apiver") + "/docker/token"
	parameters := url.Values{}
	acct := util.GetCliStringFlag(cmd, "account")
	if acct == "" {
		acct = user
	}

	svc := util.GetCliStringFlag(cmd, "service")
	scope := util.GetCliStringFlag(cmd, "scope")
	parameters.Add("account", acct)
	parameters.Add("service", svc)
	parameters.Add("scope", scope)
	Url.RawQuery = parameters.Encode()

	client := &http.Client{}
	req, err := http.NewRequest("GET", Url.String(), nil)
	if err != nil {
		util.CheckErrorExit(err, 1)
	}

	req.SetBasicAuth(user, pass)
	d.Info(fmt.Sprintf("Get token for subuser '%s' with service '%s' and scope '%s'.", user, svc, scope))
	resp, err := client.Do(req)
	if err != nil {
		util.CheckErrorExit(err, 1)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.CheckErrorExit(err, 1)
	}

	pfmt := util.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	case "raw":
		// output raw for now
		fmt.Println(string(body))
	default:
		var m map[string]interface{}
		err = json.Unmarshal(body, &m)
		if err != nil {
			util.CheckErrorExit(err, 1)
		}

		t, found := m["token"]
		if !found {
			// should not happen :)
			d.Error("cannot find token")
		}

		d.Info("token:", fmt.Sprintf("%s", t))
	}
}
