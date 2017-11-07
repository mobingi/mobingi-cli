package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/client"
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

Examples:

  $ ` + cmdline.Args0() + ` stack exec --target "stackid|ip:flag" --script /path/to/script
  $ ` + cmdline.Args0() + ` stack exec --target "stackid|ip1:flag1" \
      --target "stackid|ip2:flag2" \
      --script /path/to/script`,
		Run: stackExec,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().StringSliceVar(&targets, "target", targets, "`fmt`: stackid|ip:flag")
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
		OutputCallback: func(idx int, cr *client.Response, body []byte, th *sesha3.TargetHeader, err error) {
			if err != nil {
				d.Error(th.StackId+"|"+th.VmUser+"@"+th.Ip+":"+th.Flag+" failed:", err)
			}

			d.Info("["+fmt.Sprintf("%d", idx)+"]output:", th.StackId+",", "instance:", th.VmUser+"@"+th.Ip+",", "flag:", th.Flag)
			var res sesha3.ExecScriptStackResponse
			err = json.Unmarshal(body, &res)
			if err != nil {
				d.Error(err)
			}

			fmt.Println(string(res.Outputs[0].CmdOut))
		},
		Script: scriptdata,
	}

	_, _, err = svc.ExecScript(in)
	cli.ErrorExit(err, 1)
}
