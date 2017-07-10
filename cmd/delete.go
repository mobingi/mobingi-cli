package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mobingilabs/mocli/pkg/cli"
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
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().String("id", "", "stack id to delete")
}

func delete(cmd *cobra.Command, args []string) {
	token, err := util.GetToken()
	if err != nil {
		util.PrintErrorAndExit("Cannot read token. See `login` for information on how to login.", 1)
	}

	id := util.GetCliStringFlag(cmd, "id")
	if id == "" {
		util.PrintErrorAndExit("Stack id cannot be empty.", 1)
	}

	c := cli.New(util.GetCliStringFlag(cmd, "api-version"))
	ep := c.RootUrl + "/alm/stack/" + fmt.Sprintf("%s", id)
	resp, body, errs := c.DeleteSafe(ep, fmt.Sprintf("%s", token))
	if errs != nil {
		log.Println("Error(s):", errs)
		os.Exit(1)
	}

	var m map[string]interface{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		util.PrintErrorAndExit("Internal error.", 1)
	}

	status, found := m["status"]
	if !found {
		util.PrintErrorAndExit("Cannot read status.", 1)
	}

	log.Println(fmt.Sprintf("[%s] %s", resp.Status, status))
}
