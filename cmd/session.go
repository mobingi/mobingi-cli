package cmd

import (
	"strconv"
	"strings"
	"time"

	"github.com/mobingi/mobingi-cli/client/timeout"
	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingi/mobingi-cli/pkg/cli/confmap"
	"github.com/mobingilabs/mobingi-sdk-go/client"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func getApiVersionInt() int {
	v := 3
	vparam := viper.GetString(confmap.ConfigKey("apiver"))
	in, err := strconv.Atoi(strings.TrimLeft(vparam, "v"))
	if err != nil {
		return -1
	}

	v = in
	return v
}

func clisession() (*session.Session, error) {
	v := getApiVersionInt()
	if v < 0 {
		return nil, errors.New("cannot get api version")
	}

	return session.New(&session.Config{
		ApiVersion:      v,
		AccessToken:     viper.GetString(confmap.ConfigKey("token")),
		BaseApiUrl:      viper.GetString(confmap.ConfigKey("url")),
		BaseRegistryUrl: viper.GetString(confmap.ConfigKey("rurl")),
		HttpClientConfig: &client.Config{
			Timeout: time.Second * time.Duration(timeout.Timeout),
			Verbose: cli.Verbose,
		},
	})
}
