package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/alm"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/pretty"
	"github.com/spf13/cobra"
)

var details bool

func StackListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all stacks",
		Long: `List all stacks. If you specify the '--out=[filename]' option,
make sure you provide the full path of the file. If the path has
space(s) in it, make sure to surround it with double quotes.

Valid format values: min (default), json, raw

For now, the 'min' format option cannot yet write to a file using the '--out=[filename]'
option. You need to specify either 'json', or 'raw'.

Examples:

  $ ` + cmdline.Args0() + ` stack list
  $ ` + cmdline.Args0() + ` stack list --fmt=json --verbose
  $ ` + cmdline.Args0() + ` stack list --fmt=raw --out=/home/foo/tmp.txt`,
		Run: stackList,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().BoolVar(&details, "details", false, "describe all stacks")
	return cmd
}

func stackList(cmd *cobra.Command, args []string) {
	sess, err := clisession()
	cli.ErrorExit(err, 1)

	svc := alm.New(sess)

	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	case "raw":
		resp, body, err := svc.List()
		cli.ErrorExit(err, 1)
		exitOn401(resp)

		fmt.Println(string(body))

		// write to file option
		f := cli.GetCliStringFlag(cmd, "out")
		if f != "" {
			err = ioutil.WriteFile(f, body, 0644)
			cli.ErrorExit(err, 1)
			d.Info(fmt.Sprintf("Output written to %s.", f))
		}
	case "json":
		resp, body, err := svc.List()
		cli.ErrorExit(err, 1)
		exitOn401(resp)

		indent := cli.GetCliIntFlag(cmd, "indent")
		js := pretty.JSON(string(body), indent)
		fmt.Println(js)

		// write to file option
		f := cli.GetCliStringFlag(cmd, "out")
		if f != "" {
			err = ioutil.WriteFile(f, []byte(js), 0644)
			cli.ErrorExit(err, 1)
			d.Info(fmt.Sprintf("Output written to %s.", f))
		}
	default:
		w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
		fmt.Fprintf(w, "STACK ID\tSTACK NAME\tPLATFORM\tSTATUS\tREGION\tLAUNCHED\n")

		// we don't need instance walking here, only stack
		in := alm.WalkerCtx{
			Data: w,
			StackCallback: func(i int, data interface{}, body []byte, ls *alm.ListStack) error {
				pw := data.(*tabwriter.Writer)
				timestr := ls.CreateTime
				t, err := time.Parse(time.RFC3339, ls.CreateTime)
				if err == nil {
					timestr = t.Format(time.RFC1123)
				}

				platform := "?"
				if ls.Configuration.Vendor != nil {
					var vm map[string]interface{}
					err = json.Unmarshal(ls.Configuration.Vendor, &vm)
					cli.ErrorExit(err, 1)

					for k, _ := range vm {
						switch k {
						case "aws", "alicloud":
							platform = k
						}
					}
				}

				if platform == "?" {
					if ls.Configuration.AWS != "" {
						platform = "aws"
					}
				}

				type cnf_t struct {
					Configuration json.RawMessage `json:"configuration"`
				}

				// if still invalid, find via regexp
				if platform == "?" {
					var cnfs []cnf_t
					err = json.Unmarshal(body, &cnfs)
					if err == nil {
						re := regexp.MustCompile(`"vendor":\{"aws":`)
						pltfm := re.FindString(string(cnfs[i].Configuration))
						if pltfm != "" {
							platform = "aws"
						}
					}
				}

				region := ls.Configuration.Region

				// if empty, extract the `"region:"xxxxxx"` part via regexp
				if region == "" {
					var cnfs []cnf_t
					err = json.Unmarshal(body, &cnfs)
					if err == nil {
						re := regexp.MustCompile(`"region":\s*".+"`)
						mi := pretty.JSON(cnfs[i].Configuration, 2)
						if mi != "" {
							rgn := re.FindString(mi)
							rgnkv := strings.Split(rgn, ":")
							if len(rgnkv) == 2 {
								r1 := strings.TrimSpace(rgnkv[1])
								region = strings.TrimRight(strings.TrimPrefix(r1, "\""), "\"")
							}
						}
					}
				}

				fmt.Fprintf(pw, "%s\t%s\t%s\t%s\t%s\t%s\n",
					ls.StackId,
					ls.Nickname,
					platform,
					ls.StackStatus,
					region,
					timestr)

				return nil
			},
		}

		err = svc.Walker(&in)
		cli.ErrorExit(err, 1)

		w.Flush()
	}
}
