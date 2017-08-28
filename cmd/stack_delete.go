package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/alm"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/spf13/cobra"
)

func StackDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete a stack",
		Long: `Delete a stack.
		
Example:

  $ ` + cmdline.Args0() + ` stack delete --id=58c2297d25645-Y6NSE4VjP-tk`,
		Run: delete,
	}

	cmd.Flags().String("id", "", "stack id to delete")
	return cmd
}

func delete(cmd *cobra.Command, args []string) {
	id := cli.GetCliStringFlag(cmd, "id")
	if id == "" {
		d.ErrorExit("stack id cannot be empty", 1)
	}

	sess, err := clisession()
	d.ErrorExit(err, 1)

	svc := alm.New(sess)
	in := &alm.StackDeleteInput{StackId: id}
	_, body, err := svc.Delete(in)
	d.ErrorExit(err, 1)

	var m map[string]interface{}
	err = json.Unmarshal(body, &m)
	d.ErrorExit(err, 1)

	status, found := m["status"]
	if !found {
		d.ErrorExit("cannot read status", 1)
	}

	d.Info(fmt.Sprintf("%s", status))
}
