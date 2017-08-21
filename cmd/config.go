package cmd

import (
	"time"

	"github.com/mobingi/mobingi-cli/client/timeout"
	"github.com/mobingi/mobingi-cli/pkg/cli/confmap"
	"github.com/mobingilabs/mobingi-sdk-go/client"
	"github.com/spf13/viper"
)

func httpClientConfig() *client.Config {
	return &client.Config{
		Timeout: time.Second * time.Duration(timeout.Timeout),
		Verbose: viper.GetBool(confmap.ConfigKey("verbose")),
	}
}
