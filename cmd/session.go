package cmd

import (
	"strconv"
	"strings"
	"time"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingi/mobingi-cli/pkg/cli/confmap"
	"github.com/mobingilabs/mobingi-sdk-go/client"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/nativestore"
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
	verbose := viper.GetBool(confmap.ConfigKey("verbose"))
	dbg := viper.GetBool(confmap.ConfigKey("debug"))
	v := getApiVersionInt()
	if v < 0 {
		return nil, errors.New("cannot get api version")
	}

	// check if we have credentials in nativestore
	user, secret, err := nativestore.Get(cli.CliLabel, cli.CliUrl)
	if err == nil {
		cred1 := strings.Split(user, "|")
		cred2 := strings.Split(secret, "|")

		var id, sec, u, p string

		id = cred1[0]
		sec = cred2[0]
		if len(cred1) == 2 {
			u = cred1[1]
		}

		if len(cred2) == 2 {
			p = cred2[1]
		}

		if cred1[0] != "" && cred2[0] != "" {
			if verbose {
				str := "use credentials from native store: " + id
				if u != "" {
					str += "|" + u
				}
				d.Info(str)
			}

			return session.New(&session.Config{
				ClientId:        id,
				ClientSecret:    sec,
				Username:        u,
				Password:        p,
				ApiVersion:      v,
				BaseApiUrl:      viper.GetString(confmap.ConfigKey("url")),
				BaseRegistryUrl: viper.GetString(confmap.ConfigKey("rurl")),
				HttpClientConfig: &client.Config{
					Timeout: time.Second * time.Duration(viper.GetInt64(confmap.ConfigKey("timeout"))),
					Verbose: verbose,
				},
			})
		}
	}

	if verbose {
		d.Info("cannot access native store, use config file token")
		if dbg {
			d.ErrorD(err)
		}
	}

	return session.New(&session.Config{
		ApiVersion:      v,
		AccessToken:     viper.GetString(confmap.ConfigKey("token")),
		BaseApiUrl:      viper.GetString(confmap.ConfigKey("url")),
		BaseRegistryUrl: viper.GetString(confmap.ConfigKey("rurl")),
		HttpClientConfig: &client.Config{
			Timeout: time.Second * time.Duration(viper.GetInt64(confmap.ConfigKey("timeout"))),
			Verbose: verbose,
		},
	})
}
