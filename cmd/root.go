package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/mobingilabs/mocli/client"
	"github.com/mobingilabs/mocli/pkg/check"
	"github.com/mobingilabs/mocli/pkg/cli"
	d "github.com/mobingilabs/mocli/pkg/debug"
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
			if d.Verbose {
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
	rootCmd.PersistentFlags().IntP("indent", "n", 4, "indent padding when fmt is 'text' or 'json'")
	rootCmd.PersistentFlags().String("runenv", "", "run in environment (dev, qa, prod)")
	rootCmd.PersistentFlags().BoolVar(&d.Verbose, "verbose", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(cli.DbgMode(), "debug", false, "debug mode when error")
	rootCmd.PersistentFlags().Int64Var(&client.Timeout, "timeout", 120, "timeout in seconds")
	rootCmd.SetHelpCommand(HelpCmd())

	viper.BindPFlag("runenv", rootCmd.PersistentFlags().Lookup("runenv"))

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
	home, err := homedir.Dir()
	check.ErrorExit(err, 1)

	cfgpath := filepath.Join(home, "."+cli.BinName())
	f := filepath.Join(cfgpath, "config")
	viper.SetConfigType("yaml")
	viper.SetConfigFile(f)

	err = viper.ReadInConfig()
	if err != nil {
		d.Info("Creating default config file...")
		err = cli.SetDefaultCliConfig(f)
		check.ErrorExit(err, 1)

		viper.SetConfigFile(f)
	}

	err = viper.ReadInConfig()
	check.ErrorExit(err, 1)
}
