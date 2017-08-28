package cmd

import (
	"strconv"
	"strings"
	"time"

	"github.com/mobingi/mobingi-cli/client/timeout"
	"github.com/mobingi/mobingi-cli/pkg/cli/confmap"
	"github.com/mobingi/mobingi-cli/pkg/dbg"
	"github.com/mobingilabs/mobingi-sdk-go/client"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func clisession() (*session.Session, error) {
	v := 2
	vparam := viper.GetString(confmap.ConfigKey("apiver"))
	in, err := strconv.Atoi(strings.TrimLeft(vparam, "v"))
	if err != nil {
		return nil, errors.Wrap(err, "cannot setup input api version")
	}

	v = in

	return session.New(&session.Config{
		ApiVersion:      v,
		AccessToken:     viper.GetString(confmap.ConfigKey("token")),
		BaseApiUrl:      viper.GetString(confmap.ConfigKey("url")),
		BaseRegistryUrl: viper.GetString(confmap.ConfigKey("rurl")),
		HttpClientConfig: &client.Config{
			Timeout: time.Second * time.Duration(timeout.Timeout),
			Verbose: dbg.Verbose,
		},
	})
}
