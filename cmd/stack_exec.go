package cmd

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/sesha3"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/spf13/cobra"
)

func StackExecCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec",
		Short: "execute a script to your instance",
		Long: `execute a script to your instance.

Examples:

  $ ` + cmdline.Args0() + ` stack exec --target "stackid1|ip1,ip2,ip3,ipn:stackid2|ip1,ip2,ip3,ipn" --script /path/to/script`,
		Run: stackExec,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().String("target", "", "stackid1|ip1,ip2,ip3,ipn:stackid2|ip1,ip2,ip3,ipn")
	cmd.Flags().String("script", "", "your script path")
	cmd.Flags().String("flag", "", "configuration flag")
	cmd.Flags().String("user", "ec2-user", "ssh username")
	return cmd
}

func stackExec(cmd *cobra.Command, args []string) {
	sess, err := clisession()
	cli.ErrorExit(err, 1)

	svc := sesha3.New(sess)
	sfile := cli.GetCliStringFlag(cmd, "script")
	scriptdata, err := ioutil.ReadFile(sfile)
	cli.ErrorExit(err, 1)
	in := &sesha3.ExecScriptInput{
		Target:     cli.GetCliStringFlag(cmd, "target"),
		Script:     string(scriptdata),
		ScriptName: path.Base(sfile),
		InstUser:   cli.GetCliStringFlag(cmd, "user"),
		Flag:       cli.GetCliStringFlag(cmd, "flag"),
	}
	_, _, u, err := svc.ExecScript(in)
	if err != nil {
		d.Error("command execution failed", err)
	}
	if os.Stderr != nil {
		d.Error("\n" + u.Err)
	}
	d.Info("\n" + u.Out)
}
