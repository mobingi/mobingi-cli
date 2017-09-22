package cmd

import (
	"fmt"
	"io/ioutil"

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
		/*
			if pfmt == "min" || pfmt == "" {
				w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
				fmt.Fprintf(w, "STACK ID\tSTACK NAME\tPLATFORM\tSTATUS\tREGION\tLAUNCHED\n")
				for i, s := range stacks {
					timestr := s.CreateTime
					t, err := time.Parse(time.RFC3339, s.CreateTime)
					if err == nil {
						timestr = t.Format(time.RFC1123)
					}

					platform := "?"
					if s.Configuration.AWS != "" {
						platform = "AWS"
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
								platform = "AWS"
							}
						}
					}

					region := s.Configuration.Region

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

					fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
						s.StackId,
						s.Nickname,
						platform,
						s.StackStatus,
						region,
						timestr)
				}

				w.Flush()
			}
		*/
	}
}
