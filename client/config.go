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

func NewConfig(cmd *cobra.Command) *Config {
	token := cli.GetCliStringFlag(cmd, "token")
	endpoint := cli.GetCliStringFlag(cmd, "url")
	apiver := cli.GetCliStringFlag(cmd, "apiver")

	// if already logged in
	if token == "" {
		t, err := credentials.GetToken()
		if err == nil {
			token = string(t)
		}
	}

	return &Config{
		RootUrl:     endpoint,
		ApiVersion:  apiver,
		AccessToken: token,
	}
}
