package cmd

import (
	"time"

	"github.com/mobingi/mobingi-cli/client/timeout"
	"github.com/mobingi/mobingi-cli/pkg/cli/confmap"
	"github.com/mobingi/mobingi-cli/pkg/dbg"
	"github.com/mobingilabs/mobingi-sdk-go/client"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
	"github.com/spf13/viper"
)

func sessionv2() (*session.Session, error) {
	return session.New(&session.Config{
		ApiVersion:      2,
		AccessToken:     viper.GetString(confmap.ConfigKey("token")),
		BaseApiUrl:      viper.GetString(confmap.ConfigKey("url")),
		BaseRegistryUrl: viper.GetString(confmap.ConfigKey("rurl")),
		HttpClientConfig: &client.Config{
			Timeout: time.Second * time.Duration(timeout.Timeout),
			Verbose: dbg.Verbose,
		},
	})
}
