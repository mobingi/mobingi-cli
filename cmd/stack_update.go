package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/alm"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/filetype"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/pretty"
	"github.com/spf13/cobra"
)

func StackUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update a stack",
		Long: `Update a stack. You can get stack id from the command:

  $ ` + cmdline.Args0() + ` stack list

As an example for --spot-range, if you have a total of 20 instances running
in the autoscaling group and your spot range is set to 50 (50%), then there
will be a fleet of 10 spot instances and 10 on-demand instances.

Example(s):

  $ ` + cmdline.Args0() + ` stack update --id=mo-58c2297d25645-TEXlvYRBQ-tk --min=5 --max=20
  $ ` + cmdline.Args0() + ` stack update --id=mo-58c2297d25645-TEXlvYRBQ-tk --spot-range=25`,
		Run: updateStack,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().String("alm-template", "", "`path` to alm template file")
	cmd.Flags().String("id", "", "stack id to update")
	cmd.Flags().String("type", "m3.medium", "server type")
	cmd.Flags().Int("min", 2, "min auto scale group instance when arch is art_elb")
	cmd.Flags().Int("max", 10, "max auto scale group instance when arch is art_elb")
	cmd.Flags().Int("spot-range", 50, "spot instance percentage to deploy")
	return cmd
}

func updateStack(cmd *cobra.Command, args []string) {
	almt := cli.GetCliStringFlag(cmd, "alm-template")
	if almt != "" {
		updateAlmStack(cmd)
		return
	}

	var modified bool
	type updatet struct {
		Configurations string `json:"configurations,omitempty"`
	}

	id := cli.GetCliStringFlag(cmd, "id")
	if id == "" {
		cli.ErrorExit("stack id required", 1)
	}

	cnf := alm.StackCreateConfig{}
	if cmd.Flag("type").Changed {
		cnf.Type = cli.GetCliStringFlag(cmd, "type")
		modified = true
	}

	if cmd.Flag("min").Changed {
		cnf.Min = cli.GetCliIntFlag(cmd, "min")
		modified = true
	}

	if cmd.Flag("max").Changed {
		cnf.Max = cli.GetCliIntFlag(cmd, "max")
		modified = true
	}

	if cmd.Flag("spot-range").Changed {
		cnf.SpotRange = cli.GetCliIntFlag(cmd, "spot-range")
		modified = true
	}

	if !modified {
		d.Info("nothing to update")
		os.Exit(0)
	}

	mi, err := json.Marshal(&cnf)
	cli.ErrorExit(err, 1)

	p := updatet{}
	p.Configurations = string(mi)

	// for pretty print
	mi, err = json.MarshalIndent(&p, "", pretty.Indent(2))
	cli.ErrorExit(err, 1)

	d.Info("[update stack payload]")
	fmt.Println(string(mi))

	// for actual payload (smaller)
	mi, err = json.Marshal(&p)
	cli.ErrorExit(err, 1)

	sess, err := clisession()
	cli.ErrorExit(err, 1)

	svc := alm.New(sess)
	in := &alm.StackUpdateInput{
		StackId:        id,
		Configurations: cnf,
	}

	resp, body, err := svc.Update(in)
	cli.ErrorExit(err, 1)
	exitOn401(resp)

	var success bool
	var m map[string]interface{}
	err = json.Unmarshal(body, &m)
	cli.ErrorExit(err, 1)

	_, ok := m["status"]
	if ok {
		d.Info(fmt.Sprintf("[%s] %s", resp.Status, m["status"]))
		success = true
	}

	if !success {
		d.Info(string(body))
		return
	}
}

func updateAlmStack(cmd *cobra.Command) {
	id := cli.GetCliStringFlag(cmd, "id")
	if id == "" {
		cli.ErrorExit("stack id required", 1)
	}

	tf := cli.GetCliStringFlag(cmd, "alm-template")
	b, err := ioutil.ReadFile(tf)
	cli.ErrorExit(err, 1)

	if !filetype.IsJSON(string(b)) {
		cli.ErrorExit("invalid json", 1)
	}

	sess, err := clisession()
	cli.ErrorExit(err, 1)

	svc := alm.New(sess)
	in := &alm.StackUpdateInput{
		AlmTemplate: &alm.AlmTemplate{
			ContentType: "json",
			Contents:    string(b),
		},
		StackId: id,
	}

	resp, body, err := svc.Update(in)
	cli.ErrorExit(err, 1)
	exitOn401(resp)

	if strings.Contains(string(body), "success") {
		res := pretty.JSON(string(body), 2)
		d.Info(fmt.Sprintf("[%s] return payload:", resp.Status))
		fmt.Println(res)
		return
	}

	if (resp.StatusCode / 100) == 2 {
		d.Info(fmt.Sprintf("[%s] %s", resp.Status, string(body)))
	} else {
		d.Error(fmt.Sprintf("[%s] %s", resp.Status, string(body)))
	}
}
