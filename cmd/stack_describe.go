package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/alm"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/pretty"
	"github.com/spf13/cobra"
)

func StackDescribeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "display stack details",
		Long: `Display stack details. If you specify the '--out=[filename]' option,
make sure you provide the full path of the file. If the path has
space(s) in it, make sure to surround it with double quotes.

Valid format values: min (default), json, raw

Examples:

  $ ` + cmdline.Args0() + ` stack describe --id=58c2297d25645-Y6NSE4VjP-tk
  $ ` + cmdline.Args0() + ` stack describe --id=58c2297d25645-Y6NSE4VjP-tk --fmt=json`,
		Run: describe,
	}

	cmd.Flags().String("id", "", "stack id")
	return cmd
}

func describe(cmd *cobra.Command, args []string) {
	sess, err := clisession()
	cli.ErrorExit(err, 1)

	svc := alm.New(sess)
	in := &alm.StackDescribeInput{
		StackId: cli.GetCliStringFlag(cmd, "id"),
	}

	resp, body, err := svc.Describe(in)
	cli.ErrorExit(err, 1)
	exitOn401(resp)

	// we process `--fmt=raw` option first
	out := cli.GetCliStringFlag(cmd, "out")
	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	/*
		if sess.Config.ApiVersion == 3 {
			if pfmt == "min" || pfmt == "" {
				pfmt = "json"
			}
		}
	*/

	switch pfmt {
	case "raw":
		fmt.Println(string(body))
		if out != "" {
			err = ioutil.WriteFile(out, body, 0644)
			cli.ErrorExit(err, 1)
		}
	case "json":
		indent := cli.GetCliIntFlag(cmd, "indent")
		js := pretty.JSON(string(body), indent)
		fmt.Println(js)

		// write to file option
		if out != "" {
			err = ioutil.WriteFile(out, []byte(js), 0644)
			cli.ErrorExit(err, 1)
		}
	default:
		if pfmt == "min" || pfmt == "" {
			/*
				if sess.Config.ApiVersion >= 3 {
					err = v3DescribeStack(cmd, body)
					cli.ErrorExit(err, 1)
					return
				}
			*/

			type Configuration struct {
				// v3
				Description string          `json:"description,omitempty"`
				Label       string          `json:"label,omitempty"`
				Version     string          `json:"version,omitempty"`
				Vendor      json.RawMessage `json:"vendor,omitempty"`
				// v2
				AWS                 string      `json:"AWS,omitempty"`
				AWSAccountName      string      `json:"AWS_ACCOUNT_NAME,omitempty"`
				AssociatePublicIp   string      `json:"AssociatePublicIP,omitempty"`
				ELBOpen443Port      string      `json:"ELBOpen443Port,omitempty"`
				ELBOpen80Port       string      `json:"ELBOpen80Port,omitempty"`
				SpotInstanceMaxSize int         `json:"SpotInstanceMaxSize,omitempty"`
				SpotInstanceMinSize int         `json:"SpotInstanceMinSize,omitempty"`
				SpotPrice           string      `json:"SpotPrice,omitempty"`
				Architecture        string      `json:"architecture,omitempty"`
				Code                string      `json:"code,omitempty"`
				Image               string      `json:"image,omitempty"`
				Max                 interface{} `json:"max,omitempty"`
				MaxOrigin           interface{} `json:"maxOrigin,omitempty"`
				Min                 interface{} `json:"min,omitempty"`
				MinOrigin           interface{} `json:"minOrigin,omitempty"`
				Nickname            string      `json:"nickname,omitempty"`
				Region              string      `json:"region,omitempty"`
				Type                string      `json:"type,omitempty"`
			}

			type State struct {
				Code string `json:"Code,omitempty"`
				Name string `json:"Name,omitempty"`
			}

			type Instance struct {
				AmiLaunchIndex        string      `json:"AmiLaunchIndex,omitempty"`
				Architecture          string      `json:"Architecture,omitempty"`
				BlockDeviceMappings   interface{} `json:"BlockDeviceMappings,omitempty"`
				ClientToken           string      `json:"ClientToken,omitempty"`
				EbsOptimized          bool        `json:"EbsOptimized,omitempty"`
				Hypervisor            string      `json:"Hypervisor,omitempty"`
				ImageId               string      `json:"ImageId,omitempty"`
				InstanceId            string      `json:"InstanceId,omitempty"`
				InstanceType          string      `json:"InstanceType,omitempty"`
				InstanceLifecycle     string      `json:"InstanceLifecycle,omitempty"`
				SpotInstanceRequestId string      `json:"SpotInstanceRequestId,omitempty"`
				KeyName               string      `json:"KeyName,omitempty"`
				LaunchTime            string      `json:"LaunchTime,omitempty"`
				Monitoring            interface{} `json:"Monitoring,omitempty"`
				NetworkInterfaces     interface{} `json:"NetworkInterfaces,omitempty"`
				Placement             interface{} `json:"Placement,omitempty"`
				PrivateDnsName        string      `json:"PrivateDnsName,omitempty"`
				PrivateIpAddress      string      `json:"PrivateIpAddress,omitempty"`
				ProductCodes          []string    `json:"ProductCodes,omitempty"`
				PublicDnsName         string      `json:"PublicDnsName,omitempty"`
				PublicIpAddress       string      `json:"PublicIpAddress,omitempty"`
				Reservation           interface{} `json:"Reservation,omitempty"`
				RootDeviceName        string      `json:"RootDeviceName,omitempty"`
				RootDeviceType        string      `json:"RootDeviceType,omitempty"`
				SecurityGroups        interface{} `json:"SecurityGroups,omitempty"`
				SourceDestCheck       bool        `json:"SourceDestCheck,omitempty"`
				State                 State       `json:"State,omitempty"`
				StateTransitionReason string      `json:"StateTransitionReason,omitempty"`
				SubnetId              string      `json:"SubnetId,omitempty"`
				Tags                  interface{} `json:"Tags,omitempty"`
				VirtualizationType    string      `json:"VirtualizationType,omitempty"`
				VpcId                 string      `json:"VpcId,omitempty"`
				EnaSupport            string      `json:"enaSupport,omitempty"`
			}

			type DescribeStack struct {
				AuthToken     string        `json:"auth_token,omitempty"`
				Configuration Configuration `json:"configuration,omitempty"`
				CreateTime    string        `json:"create_time,omitempty"`
				Instances     []Instance    `json:"Instances,omitempty"`
				Nickname      string        `json:"nickname,omitempty"`
				StackId       string        `json:"stack_id,omitempty"`
				StackOutputs  interface{}   `json:"stack_outputs,omitempty"`
				StackStatus   string        `json:"stack_status,omitempty"`
				UserId        string        `json:"user_id,omitempty"`
			}

			var stacks []alm.DescribeStack
			var stack alm.DescribeStack

			switch sess.Config.ApiVersion {
			case 3:
				err = json.Unmarshal(body, &stack)
				cli.ErrorExit(err, 1)
			default:
				err = json.Unmarshal(body, &stacks)
				cli.ErrorExit(err, 1)
				stack = stacks[0]
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
			fmt.Fprintf(w, "INSTANCE ID\tINSTANCE TYPE\tINSTANCE MODEL\tPUBLIC IP\tPRIVATE IP\tSTATUS\n")
			for _, inst := range stack.Instances {
				instype := "on-demand"
				if inst.InstanceLifecycle == "spot" {
					instype = inst.InstanceLifecycle
				}

				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
					inst.InstanceId,
					instype,
					inst.InstanceType,
					inst.PublicIpAddress,
					inst.PrivateIpAddress,
					inst.State.Name)
			}

			w.Flush()
		}
	}
}

func v3DescribeStack(cmd *cobra.Command, body []byte) error {
	return nil
}
