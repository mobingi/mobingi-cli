package cmd

import (
	"github.com/mobingi/mobingi-cli/pkg/cli/confmap"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
	"github.com/spf13/viper"
)

func sessionv2() (*session.Session, error) {
	return session.New(&session.Config{
		ApiVersion:      2,
		AccessToken:     viper.GetString(confmap.ConfigKey("token")),
		BaseApiUrl:      viper.GetString(confmap.ConfigKey("url")),
		BaseRegistryUrl: viper.GetString(confmap.ConfigKey("rurl")),
	})
}
