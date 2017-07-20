package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

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
	tokenCmd.Flags().String("username", "", "username (account subuser)")
	tokenCmd.Flags().String("password", "", "password (account subuser)")
}

func token(cmd *cobra.Command, args []string) {
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
	}

	if pass == "" {
		util.CheckErrorExit("password cannot be empty", 1)
	}

	fmt.Println("\n") // new line after the password input
	log.Println(user, pass)

	/*
		data := url.Values{}
		data.Add("account", "chewsubuser1")
		// data.Add("username", "chewsubuser1")
		// data.Add("password", "mobingi")
		data.Add("service", "Mobingi Docker Registry")
		data.Add("scope", "registry:chewsubuser1:catalog:*")
	*/

	var Url *url.URL
	Url, err := url.Parse("https://apidev.mobingi.com")
	if err != nil {
		util.CheckErrorExit(err, 1)
	}

	Url.Path += "/v2/docker/token"
	parameters := url.Values{}
	parameters.Add("account", "chewsubuser1")
	parameters.Add("service", "Mobingi Docker Registry")
	parameters.Add("scope", "repository:chewsubuser1/hello:*")
	// parameters.Add("scope", "registry:catalog:*")
	Url.RawQuery = parameters.Encode()

	fmt.Printf("Encoded URL is %q\n", Url.String())

	client := &http.Client{}
	// req, err := http.NewRequest("GET", "https://apidev.mobingi.com/v2/docker/token", strings.NewReader(data.Encode()))
	req, err := http.NewRequest("GET", Url.String(), nil)
	req.SetBasicAuth(user, pass)
	log.Println(req.URL.String())

	resp, err := client.Do(req)
	if err != nil {
		util.CheckErrorExit(err, 1)
	}

	defer resp.Body.Close()
	log.Println(resp)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.CheckErrorExit(err, 1)
	}

	fmt.Println("body:", string(body))
}
