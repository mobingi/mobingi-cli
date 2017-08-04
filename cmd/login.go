package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mobingilabs/mocli/client"
	"github.com/mobingilabs/mocli/client/timeout"
	"github.com/mobingilabs/mocli/pkg/check"
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/credentials"
	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/mobingilabs/mocli/pkg/pretty"
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
enter your username and password. Token will be saved in $HOME/.` + cli.BinName() + `/credentials.

Valid 'grant-type' values: client_credentials, password

Examples:

  $ ` + cli.BinName() + ` login --client-id=foo --client-secret=bar`,
		Run: login,
	}

	cmd.Flags().StringP("client-id", "i", "", "client id (required)")
	cmd.Flags().StringP("client-secret", "s", "", "client secret (required)")
	cmd.Flags().StringP("grant-type", "g", "client_credentials", "grant type")
	cmd.Flags().StringP("username", "u", "", "user name")
	cmd.Flags().StringP("password", "p", "", "password")
	return cmd
}

func login(cmd *cobra.Command, args []string) {
	idsec := &credentials.ClientIdSecret{
		Id:     cli.GetCliStringFlag(cmd, "client-id"),
		Secret: cli.GetCliStringFlag(cmd, "client-secret"),
	}

	err := idsec.EnsureInput(false)
	if err != nil {
		check.ErrorExit(err, 1)
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
		check.ErrorExit("Invalid argument(s). See `help` for more information.", 1)
	}

	cnf := cli.ReadCliConfig()
	if cnf == nil {
		check.ErrorExit("read config failed", 1)
	}

	apiurl := fmt.Sprint(fval(cmd, "url", cli.ProductionBaseApiUrl))
	cnf.BaseApiUrl = apiurl
	viper.Set(cli.ConfigKey("url"), apiurl)

	regurl := fmt.Sprint(fval(cmd, "rurl", cli.ProductionBaseRegistryUrl))
	cnf.BaseRegistryUrl = regurl
	viper.Set(cli.ConfigKey("rurl"), regurl)

	apiver := fmt.Sprint(fval(cmd, "apiver", pretty.Pad))
	cnf.ApiVersion = apiver
	viper.Set(cli.ConfigKey("apiver"), apiver)

	indent := fval(cmd, "indent", pretty.Pad)
	cnf.Indent = indent.(int)
	viper.Set(cli.ConfigKey("indent"), indent.(int))

	tm := fval(cmd, "timeout", timeout.Timeout)
	cnf.Timeout = tm.(int64)
	viper.Set(cli.ConfigKey("timeout"), tm.(int64))

	verbose := fval(cmd, "verbose", d.Verbose)
	cnf.Verbose = verbose.(bool)
	viper.Set(cli.ConfigKey("verbose"), verbose.(bool))

	dbg := fval(cmd, "debug", cli.IsDbgMode())
	cnf.Debug = dbg.(bool)
	viper.Set(cli.ConfigKey("debug"), dbg.(bool))

	payload, err := json.Marshal(p)
	check.ErrorExit(err, 1)

	c := client.NewClient(client.NewApiConfig(cmd))
	token, err := c.GetAccessToken(payload)
	check.ErrorExit(err, 1)

	cnf.AccessToken = token
	err = cnf.WriteToConfig()
	check.ErrorExit(err, 1)

	// reload updated config to viper
	err = viper.ReadInConfig()
	check.ErrorExit(err, 1)

	d.Info("Login successful.")
}

func fval(cmd *cobra.Command, flag string, defval interface{}) interface{} {
	var ret interface{}
	switch defval.(type) {
	case string:
		fvalue := cli.GetCliStringFlag(cmd, flag)
		if fvalue == "" {
			tmp := viper.Get(cli.ConfigKey(flag))
			if tmp == nil {
				return defval
			} else {
				ret = viper.GetString(cli.ConfigKey(flag))
			}
		} else {
			ret = fvalue
		}
	case int:
		if cmd.Flag(flag).Changed {
			return cli.GetCliIntFlag(cmd, flag)
		} else {
			tmp := viper.Get(cli.ConfigKey(flag))
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
			tmp := viper.Get(cli.ConfigKey(flag))
			if tmp == nil {
				return cli.GetCliInt64Flag(cmd, flag)
			} else {
				// viper's get returns int, not int64
				ret = viper.GetInt64(cli.ConfigKey(flag))
			}
		}
	case bool:
		if cmd.Flag(flag).Changed {
			return defval
		} else {
			tmp := viper.Get(cli.ConfigKey(flag))
			if tmp == nil {
				return defval
			} else {
				ret = tmp
			}
		}
	default:
		check.ErrorExit("defval type not supported", 1)
	}

	return ret
}
