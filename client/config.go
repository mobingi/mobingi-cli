package client

import (
	"github.com/mobingilabs/mocli/pkg/cli/confmap"
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

	apiver := viper.GetString(confmap.ConfigKey("apiver"))
	baseurl := viper.GetString(confmap.ConfigKey("url"))
	token := viper.GetString(confmap.ConfigKey("token"))

	return &Config{
		RootUrl:     baseurl,
		ApiVersion:  apiver,
		AccessToken: token,
	}
}
