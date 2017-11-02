package cmd

import (
	"encoding/json"
	"io/ioutil"
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
		Short: "execute a script in stack instance(s)",
		Long: `Execute a script in stack instance(s).

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

	_, body, err := svc.ExecScript(in)
	cli.ErrorExit(err, 1)

	var res []sesha3.ExecScriptStackResponse
	err = json.Unmarshal(body, &res)
	cli.ErrorExit(err, 1)

	d.Info(string(res[0].Outputs[0].CmdOut))
}
