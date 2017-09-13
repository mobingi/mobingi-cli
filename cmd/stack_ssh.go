package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/sesha3"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/pretty"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

var linkOnly bool

func StackSshCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssh",
		Short: "ssh to your instance",
		Long: `Try to establish an ssh connection to an instance.

Examples:

  $ ` + cmdline.Args0() + ` stack ssh --id=58c2297d25645-Y6NSE4VjP-tk --ip=1.1.1.1`,
		Run: stackSsh,
	}

	cmd.Flags().String("id", "", "stack id")
	cmd.Flags().String("ip", "", "instance ip address")
	cmd.Flags().String("user", "ec2-user", "ssh username")
	cmd.Flags().BoolVar(&linkOnly, "url-only", false, "true if you want the url only")
	return cmd
}

func stackSsh(cmd *cobra.Command, args []string) {
	sess, err := clisession()
	d.ErrorExit(err, 1)

	svc := sesha3.New(sess)
	in := &sesha3.GetSessionUrlInput{
		StackId:  cli.GetCliStringFlag(cmd, "id"),
		IpAddr:   cli.GetCliStringFlag(cmd, "ip"),
		InstUser: cli.GetCliStringFlag(cmd, "user"),
	}

	resp, body, u, err := svc.GetSessionUrl(in)
	d.ErrorExit(err, 1)
	exitOn401(resp)

	out := cli.GetCliStringFlag(cmd, "out")
	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	case "raw":
		fmt.Println(string(body))
		if out != "" {
			err = ioutil.WriteFile(out, body, 0644)
			d.ErrorExit(err, 1)
		}
	case "json":
		indent := cli.GetCliIntFlag(cmd, "indent")
		js := pretty.JSON(string(body), indent)
		fmt.Println(js)

		// write to file option
		if out != "" {
			err = ioutil.WriteFile(out, []byte(js), 0644)
			d.ErrorExit(err, 1)
		}
	default:
		if linkOnly {
			d.Info("open link with a browser:", u)
			_ = open.Run(u)
			return
		}

		sshcli, err := sesha3.NewClient(&sesha3.SeshaClientInput{URL: u})
		d.ErrorExit(err, 1)

		err = sshcli.Run()
		if err != nil {
			d.Error("session return:", err)
		}
	}
}
