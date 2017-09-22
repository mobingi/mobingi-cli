package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"
	"time"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/registry"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/pretty"
	"github.com/spf13/cobra"
)

func DescribeImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "describe an image",
		Long: `Describe an image.

Example:

  $ ` + cmdline.Args0() + ` registry descrube --image foo`,
		Run: describeImage,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().String("username", "", "username (account subuser)")
	cmd.Flags().String("password", "", "password (account subuser)")
	cmd.Flags().String("image", "", "image name")
	return cmd
}

func describeImage(cmd *cobra.Command, args []string) {
	sess, err := clisession()
	d.ErrorExit(err, 1)

	ensureUserPass(cmd, sess)
	svc := registry.New(sess)
	in := &registry.DescribeImageInput{
		Image: cli.GetCliStringFlag(cmd, "image"),
	}

	resp, body, err := svc.DescribeImage(in)
	d.ErrorExit(err, 1)
	exitOn401(resp)

	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	case "raw":
		fmt.Println(string(body))

		// write to file option
		f := cli.GetCliStringFlag(cmd, "out")
		if f != "" {
			err = ioutil.WriteFile(f, body, 0644)
			d.ErrorExit(err, 1)
			d.Info(fmt.Sprintf("Output written to %s.", f))
		}
	case "json":
		indent := cli.GetCliIntFlag(cmd, "indent")
		js := pretty.JSON(string(body), indent)
		fmt.Println(js)

		// write to file option
		f := cli.GetCliStringFlag(cmd, "out")
		if f != "" {
			err = ioutil.WriteFile(f, []byte(js), 0644)
			d.ErrorExit(err, 1)
			d.Info(fmt.Sprintf("Output written to %s.", f))
		}
	default:
		w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
		fmt.Fprintf(w, "REPO\tSIZE\tCREATED\tVISIBILITY\tTAGS\n")
		var repos []json.RawMessage
		err = json.Unmarshal(body, &repos)
		d.ErrorExit(err, 1)

		for _, item := range repos {
			var m map[string]interface{}
			err = json.Unmarshal(item, &m)
			d.ErrorExit(err, 1)

			var repo, size, created, vis, tags string
			if _, ok := m["repository"]; ok {
				repo = fmt.Sprintf("%s", m["repository"])
			}

			if _, ok := m["size"]; ok {
				size = fmt.Sprintf("%v", m["size"])
			}

			if _, ok := m["created_at"]; ok {
				created = fmt.Sprintf("%v", m["created_at"])
				t, err := time.Parse(time.RFC3339, created)
				if err == nil {
					created = t.Format(time.RFC1123)
				}
			}

			if _, ok := m["visibility"]; ok {
				vis = fmt.Sprintf("%v", m["visibility"])
			}

			// count tags
			if _, ok := m["tags"]; ok {
				var t1 map[string]json.RawMessage
				err = json.Unmarshal(item, &t1)
				if err == nil {
					if t2, ok := t1["tags"]; ok {
						var cnt map[string]interface{}
						err = json.Unmarshal(t2, &cnt)
						if err == nil {
							tags = fmt.Sprintf("%v", len(cnt))
						}
					}
				}
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				repo,
				size,
				created,
				vis,
				tags)
		}

		w.Flush()
	}
}
