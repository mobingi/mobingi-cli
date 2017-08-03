package client

import (
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	apiver := cli.GetCliStringFlag(cmd, "apiver")
	baseurl := cli.BaseApiUrl(cmd)
	token := viper.GetString("access_token")

	return &Config{
		RootUrl:     baseurl,
		ApiVersion:  apiver,
		AccessToken: token,
	}
}
