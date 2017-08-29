package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingi/mobingi-cli/pkg/iohelper"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/alm"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/pretty"
	"github.com/spf13/cobra"
)

func TemplateVersionsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "versions",
		Short: "list all template versions for a stack",
		Long: `List all template versions for a stack.

Example:

  $ ` + cmdline.Args0() + ` template versions --id foo`,
		Run: tmplVersionsList,
	}

	cmd.Flags().StringP("id", "i", "", "stack id")
	return cmd
}

func tmplVersionsList(cmd *cobra.Command, args []string) {
	sess, err := clisession()
	d.ErrorExit(err, 1)

	svc := alm.New(sess)
	in := &alm.GetTemplateVersionsInput{
		StackId: cli.GetCliStringFlag(cmd, "id"),
	}

	resp, body, err := svc.GetTemplateVersions(in)
	d.ErrorExit(err, 1)
	exitOn401(resp)

	out := cli.GetCliStringFlag(cmd, "out")
	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	case "raw":
		fmt.Println(string(body))
		if out != "" {
			err = iohelper.WriteToFile(out, body)
			d.ErrorExit(err, 1)
		}
	case "json":
		indent := cli.GetCliIntFlag(cmd, "indent")
		js := pretty.JSON(string(body), indent)
		fmt.Println(js)

		// write to file option
		if out != "" {
			err = iohelper.WriteToFile(out, []byte(js))
			d.ErrorExit(err, 1)
		}
	default:
		var vers []alm.AlmTemplateVersion
		err = json.Unmarshal(body, &vers)
		d.ErrorExit(err, 1)

		w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
		fmt.Fprintf(w, "VERSION ID\tLATEST\tLAST MODIFIED\tSIZE\n")
		for _, ver := range vers {
			timestr := ver.LastModified
			t, err := time.Parse(time.RFC3339, ver.LastModified)
			if err == nil {
				timestr = t.Format(time.RFC1123)
			}

			fmt.Fprintf(w, "%s\t%v\t%s\t%s\n",
				ver.VersionId,
				ver.Latest,
				timestr,
				ver.Size)
		}

		w.Flush()
	}
}
