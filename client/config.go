package client

import (
	"github.com/mobingilabs/mocli/pkg/check"
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/constants"
	"github.com/mobingilabs/mocli/pkg/credentials"
	"github.com/spf13/cobra"
)

type Config struct {
	RootUrl     string
	ApiVersion  string
	AccessToken string
}

func NewApiConfig(cmd *cobra.Command) *Config {
	token := cli.GetCliStringFlag(cmd, "token")
	baseurl := cli.GetCliStringFlag(cmd, "url")
	apiver := cli.GetCliStringFlag(cmd, "apiver")

	// if already logged in
	if token == "" {
		t, err := credentials.GetToken()
		if err == nil {
			token = string(t)
		}
	}

	if baseurl == "" {
		baseurl = constants.PROD_API_BASE
		if check.IsDevMode() {
			baseurl = constants.DEV_API_BASE
		}
	}

	return &Config{
		RootUrl:     baseurl,
		ApiVersion:  apiver,
		AccessToken: token,
	}
}
