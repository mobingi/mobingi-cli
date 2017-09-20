package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/registry"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/pretty"
	"github.com/spf13/cobra"
)

type tags struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func RegistryTagsList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "list image tags",
		Long: `List image tags. At the very least, you only have to provide 'username', 'password',
and image name. Other values will be built based on inputs and command type.

Example:

  $ ` + cmdline.Args0() + ` registry tags --image=hello`,
		Run: tagsList,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().String("username", "", "username (account subuser)")
	cmd.Flags().String("password", "", "password (account subuser)")
	cmd.Flags().String("service", "Mobingi Docker Registry", "service for authentication")
	cmd.Flags().String("scope", "", "scope for authentication")
	cmd.Flags().String("image", "", "image name to query")
	return cmd
}

func tagsList(cmd *cobra.Command, args []string) {
	service := cli.GetCliStringFlag(cmd, "service")
	scope := cli.GetCliStringFlag(cmd, "scope")
	image := cli.GetCliStringFlag(cmd, "image")
	if image == "" {
		d.ErrorExit("image name cannot be empty", 1)
	}

	sess, err := clisession()
	d.ErrorExit(err, 1)

	ensureUserPass(cmd, sess)
	svc := registry.New(sess)
	in := &registry.GetTagsListInput{
		Service: service,
		Scope:   scope,
		Image:   image,
	}

	resp, body, err := svc.GetTagsList(in)
	d.ErrorExit(err, 1)
	exitOn401(resp)

	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	case "raw":
		fmt.Println(string(body))
	case "json":
		indent := cli.GetCliIntFlag(cmd, "indent")
		js := pretty.JSON(string(body), indent)
		fmt.Println(js)
	default:
		var t tags
		err = json.Unmarshal(body, &t)
		d.ErrorExit(err, 1)

		// write table
		w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
		fmt.Fprintf(w, "IMAGE\tTAG\n")
		for _, v := range t.Tags {
			fmt.Fprintf(w, "%s\t%s\n", t.Name, v)
		}

		w.Flush()
	}
}
