package cmd

import (
	"fmt"
	"time"

	"github.com/mobingi/mobingi-cli/client/timeout"
	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingi/mobingi-cli/pkg/cli/confmap"
	sdkclient "github.com/mobingilabs/mobingi-sdk-go/client"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/credentials"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/nativestore"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/pretty"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type authPayload struct {
	ClientId     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	GrantType    string `json:"grant_type,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
}

func LoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "login to Mobingi API",
		Long: `Login to Mobingi API server. If 'grant_type' is set to 'password', you will be prompted to
enter your username and password. Cli will attempt to store your credentials to the native
store, if supported. Otherwise, token will be saved in $HOME/.` + cmdline.Args0() + `/` + cli.ConfigFileName + `.

Valid 'grant-type' values: password, client_credentials

Examples:

  $ ` + cmdline.Args0() + ` login --client-id=foo --client-secret=bar`,
		Run: login,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().StringP("client-id", "i", "", "client id (required)")
	cmd.Flags().StringP("client-secret", "s", "", "client secret (required)")
	cmd.Flags().StringP("grant-type", "g", "password", "grant type")
	cmd.Flags().StringP("username", "u", "", "user name")
	cmd.Flags().StringP("password", "p", "", "password")
	cmd.Flags().String("endpoints", "prod", "set endpoints (dev, qa, prod)")
	return cmd
}

func login(cmd *cobra.Command, args []string) {
	idsec := &credentials.ClientIdSecret{
		Id:     cli.GetCliStringFlag(cmd, "client-id"),
		Secret: cli.GetCliStringFlag(cmd, "client-secret"),
	}

	err := idsec.EnsureInput(false)
	if err != nil {
		cli.ErrorExit(err, 1)
	}

	grant := cli.GetCliStringFlag(cmd, "grant-type")

	var p *authPayload
	if grant == "client_credentials" {
		p = &authPayload{
			ClientId:     idsec.Id,
			ClientSecret: idsec.Secret,
			GrantType:    grant,
		}
	}

	if grant == "password" {
		userpass := userPass(cmd)
		p = &authPayload{
			ClientId:     idsec.Id,
			ClientSecret: idsec.Secret,
			Username:     userpass.Username,
			Password:     userpass.Password,
			GrantType:    grant,
		}
	}

	// should not be nil when `grant_type` is valid
	if p == nil {
		cli.ErrorExit("Invalid argument(s). See `help` for more information.", 1)
	}

	cnf := cli.ReadCliConfig()
	if cnf == nil {
		cli.ErrorExit("read config failed", 1)
	}

	switch cli.GetCliStringFlag(cmd, "endpoints") {
	case "dev":
		cnf.BaseApiUrl = cli.DevelopmentBaseApiUrl
		cnf.BaseRegistryUrl = cli.DevelopmentBaseRegistryUrl
		viper.Set(confmap.ConfigKey("url"), cli.DevelopmentBaseApiUrl)
		viper.Set(confmap.ConfigKey("rurl"), cli.DevelopmentBaseRegistryUrl)
	case "qa":
		cnf.BaseApiUrl = cli.TestBaseApiUrl
		cnf.BaseRegistryUrl = cli.TestBaseRegistryUrl
		viper.Set(confmap.ConfigKey("url"), cli.TestBaseApiUrl)
		viper.Set(confmap.ConfigKey("rurl"), cli.TestBaseRegistryUrl)
	case "prod":
		cnf.BaseApiUrl = cli.ProductionBaseApiUrl
		cnf.BaseRegistryUrl = cli.ProductionBaseRegistryUrl
		viper.Set(confmap.ConfigKey("url"), cli.ProductionBaseApiUrl)
		viper.Set(confmap.ConfigKey("rurl"), cli.ProductionBaseRegistryUrl)
	default:
		err = fmt.Errorf("endpoint value not supported")
		err = errors.Wrap(err, "invalid flag")
		cli.ErrorExit(err, 1)
	}

	apiver := fmt.Sprint(fval(cmd, "apiver", cli.ApiVersion))
	cnf.ApiVersion = apiver
	viper.Set(confmap.ConfigKey("apiver"), apiver)

	indent := fval(cmd, "indent", pretty.Pad)
	cnf.Indent = indent.(int)
	viper.Set(confmap.ConfigKey("indent"), indent.(int))

	tm := fval(cmd, "timeout", timeout.Timeout)
	cnf.Timeout = tm.(int64)
	viper.Set(confmap.ConfigKey("timeout"), tm.(int64))

	verbose := fval(cmd, "verbose", cli.Verbose)
	cnf.Verbose = verbose.(bool)
	viper.Set(confmap.ConfigKey("verbose"), verbose.(bool))

	dbg := fval(cmd, "debug", cli.Debug)
	cnf.Debug = dbg.(bool)
	viper.Set(confmap.ConfigKey("debug"), dbg.(bool))

	// create our own config
	var sess *session.Session
	if grant == "password" {
		sess, err = session.New(&session.Config{
			ClientId:        p.ClientId,
			ClientSecret:    p.ClientSecret,
			Username:        p.Username,
			Password:        p.Password,
			ApiVersion:      getApiVersionInt(),
			BaseApiUrl:      viper.GetString(confmap.ConfigKey("url")),
			BaseRegistryUrl: viper.GetString(confmap.ConfigKey("rurl")),
			HttpClientConfig: &sdkclient.Config{
				Timeout: time.Second * time.Duration(viper.GetInt64(confmap.ConfigKey("timeout"))),
				Verbose: cnf.Verbose,
			},
		})

		cli.ErrorExit(err, 1)
	} else {
		sess, err = session.New(&session.Config{
			ClientId:        p.ClientId,
			ClientSecret:    p.ClientSecret,
			ApiVersion:      getApiVersionInt(),
			BaseApiUrl:      viper.GetString(confmap.ConfigKey("url")),
			BaseRegistryUrl: viper.GetString(confmap.ConfigKey("rurl")),
			HttpClientConfig: &sdkclient.Config{
				Timeout: time.Second * time.Duration(viper.GetInt64(confmap.ConfigKey("timeout"))),
				Verbose: cnf.Verbose,
			},
		})

		cli.ErrorExit(err, 1)
	}

	// prefer to store credentials to native store (keychain, wincred)
	nid := p.ClientId
	nsec := p.ClientSecret
	if grant == "password" {
		nid += "|" + p.Username
		nsec += "|" + p.Password
	}

	err = nativestore.Set(cli.CliLabel, cli.CliUrl, nid, nsec)
	if err != nil {
		if cnf.Verbose {
			d.Error("Error in accessing native store, will use config file.")
		}

		if cnf.Debug {
			d.ErrorD(err)
		}
	}

	if cnf.Verbose {
		d.Info("apiver:", "v"+fmt.Sprintf("%d", getApiVersionInt()))
		d.Info("token:", sess.AccessToken)
	}

	cnf.AccessToken = sess.AccessToken
	err = cnf.WriteToConfig()
	cli.ErrorExit(err, 1)

	// reload updated config to viper
	err = viper.ReadInConfig()
	cli.ErrorExit(err, 1)

	d.Info("Login successful.")
}

func fval(cmd *cobra.Command, flag string, defval interface{}) interface{} {
	var ret interface{}
	switch defval.(type) {
	case string:
		fvalue := cli.GetCliStringFlag(cmd, flag)
		if fvalue == "" {
			tmp := viper.Get(confmap.ConfigKey(flag))
			if tmp == nil {
				return defval
			} else {
				ret = viper.GetString(confmap.ConfigKey(flag))
			}
		} else {
			ret = fvalue
		}
	case int:
		if cmd.Flag(flag).Changed {
			return cli.GetCliIntFlag(cmd, flag)
		} else {
			tmp := viper.Get(confmap.ConfigKey(flag))
			if tmp == nil {
				return cli.GetCliIntFlag(cmd, flag)
			} else {
				ret = tmp
			}
		}
	case int64:
		if cmd.Flag(flag).Changed {
			return cli.GetCliInt64Flag(cmd, flag)
		} else {
			tmp := viper.Get(confmap.ConfigKey(flag))
			if tmp == nil {
				return cli.GetCliInt64Flag(cmd, flag)
			} else {
				// viper's get returns int, not int64
				ret = viper.GetInt64(confmap.ConfigKey(flag))
			}
		}
	case bool:
		if cmd.Flag(flag).Changed {
			return defval
		} else {
			tmp := viper.Get(confmap.ConfigKey(flag))
			if tmp == nil {
				return defval
			} else {
				ret = tmp
			}
		}
	default:
		cli.ErrorExit("defval type not supported", 1)
	}

	return ret
}
