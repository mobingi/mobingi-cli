package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	term "github.com/buger/goterm"
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/stack"
	"github.com/mobingilabs/mocli/pkg/util"
	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "display stack details",
	Long:  `Display stack details.`,
	Run:   describe,
}

func init() {
	rootCmd.AddCommand(describeCmd)
	describeCmd.Flags().String("id", "", "stack id")
}

func describe(cmd *cobra.Command, args []string) {
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
	resp, body, errs := c.GetSafe(ep, fmt.Sprintf("%s", token))
	if errs != nil {
		log.Println("Error(s):", errs)
		os.Exit(1)
	}

	var stacks []stack.DescribeStack
	err = json.Unmarshal(body, &stacks)
	if err != nil {
		log.Println(err)
		var m map[string]interface{}
		err = json.Unmarshal(body, &m)
		if err != nil {
			util.PrintErrorAndExit("Internal error.", 1)
		}

		serr := util.BuildRequestError(resp, m)
		if serr != "" {
			util.PrintErrorAndExit(serr, 1)
		}
	}

	for _, s := range stacks {
		l1 := term.NewTable(0, 10, 1, ' ', 0)
		fmt.Fprintf(l1, "Stack ID:\t%s\n", s.StackId)
		fmt.Fprintf(l1, "Stack name:\t%s\n", s.Nickname)
		fmt.Fprintf(l1, "Stack type:\t%s\n", s.Configuration.Type)
		fmt.Fprintf(l1, "Region:\t%s\n", s.Configuration.Region)
		fmt.Fprintf(l1, "Architecture:\t%s\n", s.Configuration.Architecture)
		fmt.Fprintf(l1, "Code:\t%s\n", s.Configuration.Code)
		fmt.Fprintf(l1, "Image:\t%s\n", s.Configuration.Image)
		fmt.Fprintf(l1, "Instances count:\t%d\n", len(s.Instances))
		for i, v := range s.Instances {
			fmt.Fprintf(l1, "  Index:\t[%d]\n", i)
			fmt.Fprintf(l1, "  Instance ID:\t%s\n", v.InstanceId)
			fmt.Fprintf(l1, "  Instance type:\t%s\n", v.InstanceType)
			fmt.Fprintf(l1, "  Virtualization type:\t%s\n", v.VirtualizationType)
			fmt.Fprintf(l1, "  Public IP:\t%s\n", v.PublicIpAddress)
			fmt.Fprintf(l1, "  Public DNS name:\t%s\n", v.PublicDnsName)
			fmt.Fprintf(l1, "  Private IP:\t%s\n", v.PrivateIpAddress)
			fmt.Fprintf(l1, "  Private DNS name:\t%s\n", v.PrivateDnsName)
			fmt.Fprintf(l1, "  Architecture:\t%s\n", v.Architecture)
			fmt.Fprintf(l1, "  Hypervisor:\t%s\n", v.Hypervisor)
			fmt.Fprintf(l1, "  Image ID:\t%s\n", v.ImageId)
			fmt.Fprintf(l1, "  Monitoring state:\t%s\n", v.Monitoring.State)
			fmt.Fprintf(l1, "  State:\t[%s], %s\n", v.State.Code, v.State.Name)
			fmt.Fprintf(l1, "  Availability zone:\t%s\n", v.Placement.AvailabilityZone)
			fmt.Fprintf(l1, "  Root device name:\t%s\n", v.RootDeviceName)
			fmt.Fprintf(l1, "  Root device type:\t%s\n", v.RootDeviceType)
			fmt.Fprintf(l1, "  VPC ID:\t%s\n", v.VpcId)
			fmt.Fprintf(l1, "\t\n")
		}

		fmt.Fprintf(l1, "Status:\t%s\n", s.StackStatus)
		timestr := s.CreateTime
		t, err := time.Parse(time.RFC3339, s.CreateTime)
		if err == nil {
			timestr = t.Format(time.RFC1123)
		}

		fmt.Fprintf(l1, "Time created:\t%s\n", timestr)
		term.Print(l1)
		term.Flush()
	}
}
