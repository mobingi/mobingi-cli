package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingi/mobingi-cli/pkg/cli/confmap"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/alm"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func StackGetPemCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pem",
		Short: "print stack pem file",
		Long: `Print your stack's pem file, if available.

Example:

  $ ` + cmdline.Args0() + ` stack pem --id foo`,
		Run: getpem,
	}

	cmd.Flags().StringP("id", "i", "", "stack id")
	return cmd
}

func getpem(cmd *cobra.Command, args []string) {
	verbose := viper.GetBool(confmap.ConfigKey("verbose"))
	sess, err := clisession()
	d.ErrorExit(err, 1)

	svc := alm.New(sess)
	in := &alm.GetPemInput{
		StackId: cli.GetCliStringFlag(cmd, "id"),
	}

	resp, body, pem, err := svc.GetPem(in)
	exitOn401(resp)
	if err != nil {
		d.Error(err)
		if verbose {
			d.Error("req payload:")
			fmt.Println(string(body))
			d.Error("pem payload:")
			fmt.Println(string(pem))
		}

		return
	}

	if verbose {
		d.Info("req payload:")
		fmt.Println(string(body))
	}

	d.Info("payload:")
	fmt.Println(string(pem))

	// write to file if requested
	out := cli.GetCliStringFlag(cmd, "out")
	if out != "" {
		err = ioutil.WriteFile(out, pem, 0644)
		d.ErrorExit(err, 1)
	}
}
