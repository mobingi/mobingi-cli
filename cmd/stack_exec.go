package cmd

import (
	"encoding/json"
	"io/ioutil"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/sesha3"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/spf13/cobra"
)

var targets []string

func StackExecCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec",
		Short: "execute a script in stack instance(s)",
		Long: `Execute a script in stack instance(s). You can use the --target flag more than once.

Format for --target flag:

  stack-id|ip1,ip2,...:flag

The format above will work if multiple ip's have the same flag. If not, you can use the format below:

  --target stack-id|ip1:flag1 --target stack-id|ip2:flag2

Examples:

  $ ` + cmdline.Args0() + ` stack exec --target "stackid|ip1,ip2:flag" --script /path/to/script`,
		Run: stackExec,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().StringSliceVar(&targets, "target", targets, "`fmt`: stackid|ip1,ip2,...:flag")
	cmd.Flags().String("script", "", "your script path")
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
		Targets: targets,
		Script:  scriptdata,
	}

	_, body, err := svc.ExecScript(in)
	cli.ErrorExit(err, 1)

	var res []sesha3.ExecScriptStackResponse
	err = json.Unmarshal(body, &res)
	cli.ErrorExit(err, 1)

	d.Info(string(body))
	d.Info(string(res[0].Outputs[0].CmdOut))
}
