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
	Long: `Get registry token. This command supports '--fmt=raw' option. By default,
it will only print the token value.`,
	Run: token,
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

	base := util.GetCliStringFlag(cmd, "url")
	apiver := util.GetCliStringFlag(cmd, "apiver")
	acct := util.GetCliStringFlag(cmd, "account")
	if acct == "" {
		acct = user
	}

	svc := util.GetCliStringFlag(cmd, "service")
	scope := util.GetCliStringFlag(cmd, "scope")
	body, token, err := getRegistryToken(base, apiver, user, pass, acct, svc, scope)
	if err != nil {
		util.CheckErrorExit(err, 1)
	}

	pfmt := util.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	case "raw":
		fmt.Println(string(body))
	default:
		d.Info("token:", token)

	}
}

func getRegistryToken(base, apiver, user, pass, subuser, svc, scope string) ([]byte, string, error) {
	var u *url.URL
	u, err := url.Parse(base)
	if err != nil {
		return nil, "", err
	}

	u.Path += "/" + apiver + "/docker/token"
	v := url.Values{}
	v.Add("account", subuser)
	v.Add("service", svc)
	v.Add("scope", scope)
	u.RawQuery = v.Encode()

	client := &http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, "", err
	}

	req.SetBasicAuth(user, pass)
	d.Info(fmt.Sprintf("Get token for subuser '%s' with service '%s' and scope '%s'.", subuser, svc, scope))
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	var m map[string]interface{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		return nil, "", err
	}

	t, found := m["token"]
	if !found {
		return nil, "", fmt.Errorf("cannot find token")
	}

	return body, fmt.Sprintf("%s", t), nil
}
