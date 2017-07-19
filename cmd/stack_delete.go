package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mobingilabs/mocli/api"
	"github.com/mobingilabs/mocli/pkg/util"
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
	id := util.GetCliStringFlag(cmd, "id")
	if id == "" {
		util.CheckErrorExit("stack id cannot be empty", 1)
	}

	c := api.NewClient(api.NewConfig(cmd))
	resp, body, errs := c.Del("/alm/stack/" + fmt.Sprintf("%s", id))
	util.CheckErrorExit(errs, 1)

	var m map[string]interface{}
	err := json.Unmarshal(body, &m)
	util.CheckErrorExit(err, 1)
	serr := util.ResponseError(resp, body)
	util.CheckErrorExit(serr, 1)
	status, found := m["status"]
	if !found {
		util.CheckErrorExit("cannot read status", 1)
	}

	log.Println(fmt.Sprintf("[%s] %s", resp.Status, status))
}
