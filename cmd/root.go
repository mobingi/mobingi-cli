package cmd

import (
	"fmt"
	"time"

	"github.com/mobingilabs/mocli/client"
	"github.com/mobingilabs/mocli/pkg/check"
	"github.com/mobingilabs/mocli/pkg/cli"
	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/spf13/cobra"
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
	rootCmd.PersistentFlags().String("token", "", "access token")
	rootCmd.PersistentFlags().String("url", "", "base url for API")
	rootCmd.PersistentFlags().String("rurl", "", "base url for Docker Registry")
	rootCmd.PersistentFlags().String("apiver", "v2", "API version")
	rootCmd.PersistentFlags().StringP("fmt", "f", "", "output format (values depends on command)")
	rootCmd.PersistentFlags().StringP("out", "o", "", "full file path to write the output")
	rootCmd.PersistentFlags().IntP("indent", "n", 4, "indent padding when fmt is 'text' or 'json'")
	rootCmd.PersistentFlags().BoolVar(&d.Verbose, "verbose", false, "verbose output")
	rootCmd.PersistentFlags().StringVar(&cli.RunEnv, "runenv", "prod", "run in environment (dev, qa, prod)")
	rootCmd.PersistentFlags().BoolVar(cli.DbgMode(), "debug", false, "debug mode when error")
	rootCmd.PersistentFlags().Int64Var(&client.Timeout, "timeout", 120, "timeout in seconds")
	rootCmd.SetHelpCommand(HelpCmd())

	rootCmd.AddCommand(
		LoginCmd(),
		StackCmd(),
		ServerConfigCmd(),
		CredentialsCmd(),
		RegistryCmd(),
		VersionCmd(),
	)
}
