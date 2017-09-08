package cmd

import (
	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/alm"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/spf13/cobra"
)

func StackGetPemCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pem",
		Short: "print stack pem file",
		Long: `Print your stack's pem file, if available.

Examples:

  $ ` + cmdline.Args0() + ` stack list
  $ ` + cmdline.Args0() + ` stack list --fmt=json --verbose
  $ ` + cmdline.Args0() + ` stack list --fmt=raw --out=/home/foo/tmp.txt`,
		Run: getpem,
	}

	cmd.Flags().StringP("id", "i", "", "stack id")
	return cmd
}

func getpem(cmd *cobra.Command, args []string) {
	sess, err := clisession()
	d.ErrorExit(err, 1)

	svc := alm.New(sess)
	in := &alm.GetPemInput{
		StackId: cli.GetCliStringFlag(cmd, "id"),
	}

	resp, body, err := svc.GetPem(in)
	d.ErrorExit(err, 1)
	exitOn401(resp)

	d.Info(resp, string(body))
}
