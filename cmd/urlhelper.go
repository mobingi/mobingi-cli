package cmd

import (
	"github.com/mobingi/mobingi-cli/pkg/cli/confmap"
	"github.com/spf13/viper"
)

func buildUrl(path string) string {
	apiver := viper.GetString(confmap.ConfigKey("apiver"))
	baseurl := viper.GetString(confmap.ConfigKey("url"))
	return baseurl + "/" + apiver + path
}
