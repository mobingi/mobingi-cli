package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/alm"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/pretty"
	"github.com/spf13/cobra"
)

func TemplateDescribeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "describe a specific template version",
		Long: `Describe a specific template version.

To get template version ids, you can use the command:

  $ ` + cmdline.Args0() + ` template versions --id=foo

Example:

  $ ` + cmdline.Args0() + ` template describe --id=foo --ver=bar`,
		Run: tmplDescribe,
	}

	cmd.Flags().String("id", "", "stack id")
	cmd.Flags().String("ver", "", "template version id, can be empty or 'latest'")
	return cmd
}

func tmplDescribe(cmd *cobra.Command, args []string) {
	sess, err := clisession()
	d.ErrorExit(err, 1)

	svc := alm.New(sess)
	in := &alm.DescribeTemplateInput{
		StackId:   cli.GetCliStringFlag(cmd, "id"),
		VersionId: cli.GetCliStringFlag(cmd, "ver"),
	}

	resp, body, err := svc.DescribeTemplate(in)
	d.ErrorExit(err, 1)
	exitOn401(resp)

	out := cli.GetCliStringFlag(cmd, "out")
	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	case "raw":
		fmt.Println(string(body))
		if out != "" {
			err = ioutil.WriteFile(out, body, 0644)
			d.ErrorExit(err, 1)
		}
	default:
		if pfmt == "json" || pfmt == "" {
			indent := cli.GetCliIntFlag(cmd, "indent")
			js := pretty.JSON(string(body), indent)
			fmt.Println(js)

			// write to file option
			if out != "" {
				err = ioutil.WriteFile(out, []byte(js), 0644)
				d.ErrorExit(err, 1)
			}
		}
	}
}
