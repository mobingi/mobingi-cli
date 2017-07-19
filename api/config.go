package api

import (
	"github.com/mobingilabs/mocli/pkg/util"
	"github.com/spf13/cobra"
)

type Config struct {
	RootUrl     string
	ApiVersion  string
	AccessToken string
}

func NewConfig(cmd *cobra.Command) *Config {
	token := util.GetCliStringFlag(cmd, "token")
	endpoint := util.GetCliStringFlag(cmd, "url")
	apiver := util.GetCliStringFlag(cmd, "apiver")

	// if already logged in
	if token == "" {
		t, err := util.GetToken()
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
