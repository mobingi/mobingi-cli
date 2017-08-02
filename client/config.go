package client

import (
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/credentials"
	"github.com/spf13/cobra"
)

type Config struct {
	RootUrl     string
	ApiVersion  string
	AccessToken string
}

func NewApiConfig(cmd *cobra.Command) *Config {
	if cmd == nil {
		return nil
	}

	token := cli.GetCliStringFlag(cmd, "token")
	apiver := cli.GetCliStringFlag(cmd, "apiver")
	baseurl := cli.BaseApiUrl(cmd)

	// if already logged in
	if token == "" {
		t, err := credentials.GetToken()
		if err == nil {
			token = string(t)
		}
	}

	return &Config{
		RootUrl:     baseurl,
		ApiVersion:  apiver,
		AccessToken: token,
	}
}
