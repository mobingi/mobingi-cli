package cmd

import (
	"fmt"
	"time"

	"github.com/mobingi/mobingi-cli/client/timeout"
	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingi/mobingi-cli/pkg/cli/confmap"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/pretty"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// main parent (root) command
	rootCmd = &cobra.Command{
		Use:   "mobingi-cli",
		Short: "Mobingi API command line interface.",
		Long:  `Command line interface for Mobingi API and services.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			startTime = time.Now()
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if viper.GetBool(confmap.ConfigKey("verbose")) {
				delta := int64(time.Now().Sub(startTime) / time.Millisecond)
				d.Info(fmt.Sprintf("Elapsed time: %vms", delta))
			}
		},
	}

	startTime time.Time
)

func Execute() {
	err := rootCmd.Execute()
	cli.ErrorExit(err, 1)
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.PersistentFlags().String("token", "", "access token")
	rootCmd.PersistentFlags().String("url", "", "base url for API")
	rootCmd.PersistentFlags().String("rurl", "", "base url for Docker Registry")
	rootCmd.PersistentFlags().String("apiver", "v3", "API version")
	rootCmd.PersistentFlags().StringP("fmt", "f", "", "output format (values depends on command)")
	rootCmd.PersistentFlags().StringP("out", "o", "", "full file path to write the output")
	rootCmd.PersistentFlags().IntVar(&pretty.Pad, "indent", 2, "indent padding when fmt is 'json'")
	rootCmd.PersistentFlags().Int64Var(&timeout.Timeout, "timeout", 120, "timeout in seconds")
	rootCmd.PersistentFlags().BoolVar(&cli.Verbose, "verbose", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&cli.Debug, "debug", false, "debug mode when error occurs")
	rootCmd.SetHelpCommand(HelpCmd())

	viper.BindPFlag(confmap.ConfigKey("token"), rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag(confmap.ConfigKey("url"), rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag(confmap.ConfigKey("rurl"), rootCmd.PersistentFlags().Lookup("rurl"))
	viper.BindPFlag(confmap.ConfigKey("apiver"), rootCmd.PersistentFlags().Lookup("apiver"))
	viper.BindPFlag(confmap.ConfigKey("indent"), rootCmd.PersistentFlags().Lookup("indent"))
	viper.BindPFlag(confmap.ConfigKey("timeout"), rootCmd.PersistentFlags().Lookup("timeout"))
	viper.BindPFlag(confmap.ConfigKey("verbose"), rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag(confmap.ConfigKey("debug"), rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.AddCommand(
		LoginCmd(),
		StackCmd(),
		TemplateCmd(),
		RbacCmd(),
		ServerConfigCmd(),
		CredentialsCmd(),
		RegistryCmd(),
		VersionCmd(),
		ResetCmd(),
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
		cli.ErrorExit(err, 1)

		viper.SetConfigFile(f)
	}

	err = viper.ReadInConfig()
	cli.ErrorExit(err, 1)
}
