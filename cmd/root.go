package cmd

import (
	"fmt"
	"time"

	"github.com/mobingilabs/mocli/client/timeout"
	"github.com/mobingilabs/mocli/pkg/check"
	"github.com/mobingilabs/mocli/pkg/cli"
	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/mobingilabs/mocli/pkg/pretty"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// main parent (root) command
	rootCmd = &cobra.Command{
		Use:   "mocli",
		Short: "Mobingi API command line interface.",
		Long:  `Command line interface for Mobingi API and services.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			startTime = time.Now()
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if viper.GetBool(cli.ConfigKey("verbose")) {
				delta := int64(time.Now().Sub(startTime) / time.Millisecond)
				d.Info(fmt.Sprintf("Elapsed time: %vms", delta))
			}
		},
	}

	startTime time.Time
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		check.ErrorExit(err, 1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().String("token", "", "access token")
	rootCmd.PersistentFlags().String("url", "", "base url for API")
	rootCmd.PersistentFlags().String("rurl", "", "base url for Docker Registry")
	rootCmd.PersistentFlags().String("apiver", "v2", "API version")
	rootCmd.PersistentFlags().StringP("fmt", "f", "", "output format (values depends on command)")
	rootCmd.PersistentFlags().StringP("out", "o", "", "full file path to write the output")
	rootCmd.PersistentFlags().IntVar(&pretty.Pad, "indent", 2, "indent padding when fmt is 'text' or 'json'")
	rootCmd.PersistentFlags().Int64Var(&timeout.Timeout, "timeout", 120, "timeout in seconds")
	rootCmd.PersistentFlags().BoolVar(&d.Verbose, "verbose", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(cli.DbgMode(), "debug", false, "debug mode when error")
	rootCmd.SetHelpCommand(HelpCmd())

	viper.BindPFlag(cli.ConfigKey("token"), rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag(cli.ConfigKey("url"), rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag(cli.ConfigKey("rurl"), rootCmd.PersistentFlags().Lookup("rurl"))
	viper.BindPFlag(cli.ConfigKey("apiver"), rootCmd.PersistentFlags().Lookup("apiver"))
	viper.BindPFlag(cli.ConfigKey("indent"), rootCmd.PersistentFlags().Lookup("indent"))
	viper.BindPFlag(cli.ConfigKey("timeout"), rootCmd.PersistentFlags().Lookup("timeout"))
	viper.BindPFlag(cli.ConfigKey("verbose"), rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag(cli.ConfigKey("debug"), rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.AddCommand(
		LoginCmd(),
		StackCmd(),
		ServerConfigCmd(),
		CredentialsCmd(),
		RegistryCmd(),
		VersionCmd(),
	)
}

func initConfig() {
	cnf := cli.CliConfig{}
	f := cnf.ConfigFile()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(f)

	err := viper.ReadInConfig()
	if err != nil {
		d.Info("No config file found. Creating default.")
		err = cli.SetDefaultCliConfig()
		check.ErrorExit(err, 1)

		viper.SetConfigFile(f)
	}

	err = viper.ReadInConfig()
	check.ErrorExit(err, 1)
}
