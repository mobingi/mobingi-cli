package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mobingilabs/mocli/api"
	"github.com/mobingilabs/mocli/pkg/check"
	"github.com/mobingilabs/mocli/pkg/cli"
	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a stack",
	Long:  `Delete a stack.`,
	Run:   delete,
}

func init() {
	stackCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().String("id", "", "stack id to delete")
}

func delete(cmd *cobra.Command, args []string) {
	id := cli.GetCliStringFlag(cmd, "id")
	if id == "" {
		check.ErrorExit("stack id cannot be empty", 1)
	}

	c := api.NewClient(api.NewConfig(cmd))
	resp, body, errs := c.Del("/alm/stack/" + fmt.Sprintf("%s", id))
	check.ErrorExit(errs, 1)

	var m map[string]interface{}
	err := json.Unmarshal(body, &m)
	check.ErrorExit(err, 1)
	serr := check.ResponseError(resp, body)
	check.ErrorExit(serr, 1)
	status, found := m["status"]
	if !found {
		check.ErrorExit("cannot read status", 1)
	}

	d.Info(fmt.Sprintf("[%s] %s", resp.Status, status))
}
