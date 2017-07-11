package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"

	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/stack"
	"github.com/mobingilabs/mocli/pkg/util"
	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "display stack details",
	Long: `Display stack details. If you specify the '--out=[filename]' option,
make sure you provide the full path of the file. If the path has
space(s) in it, make sure to surround it with double quotes.`,
	Run: describe,
}

func init() {
	rootCmd.AddCommand(describeCmd)
	describeCmd.Flags().StringP("id", "i", "", "stack id")
	describeCmd.Flags().StringP("fmt", "f", "text", "output format (valid values: text, json)")
	describeCmd.Flags().StringP("out", "o", "", "full file path to write the output")
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

	switch util.GetCliStringFlag(cmd, "fmt") {
	case "text":
		/*
			for _, s := range stacks {
				fmt.Fprintf(os.Stdout, "Stack ID       : %s\n", s.StackId)
				fmt.Fprintf(os.Stdout, "Stack name     : %s\n", s.Nickname)
				fmt.Fprintf(os.Stdout, "Stack type     : %s\n", s.Configuration.Type)
				fmt.Fprintf(os.Stdout, "Region         : %s\n", s.Configuration.Region)
				fmt.Fprintf(os.Stdout, "Architecture   : %s\n", s.Configuration.Architecture)
				fmt.Fprintf(os.Stdout, "Code           : %s\n", s.Configuration.Code)
				fmt.Fprintf(os.Stdout, "Image          : %s\n", s.Configuration.Image)
				fmt.Fprintf(os.Stdout, "Instances count: %d\n", len(s.Instances))
				for i, v := range s.Instances {
					fmt.Fprintf(os.Stdout, "  Index              : [%d]\n", i)
					fmt.Fprintf(os.Stdout, "  Instance ID        : %s\n", v.InstanceId)
					fmt.Fprintf(os.Stdout, "  Instance type      : %s\n", v.InstanceType)
					fmt.Fprintf(os.Stdout, "  Virtualization type: %s\n", v.VirtualizationType)
					fmt.Fprintf(os.Stdout, "  Public IP          : %s\n", v.PublicIpAddress)
					fmt.Fprintf(os.Stdout, "  Public DNS name    : %s\n", v.PublicDnsName)
					fmt.Fprintf(os.Stdout, "  Private IP         : %s\n", v.PrivateIpAddress)
					fmt.Fprintf(os.Stdout, "  Private DNS name   : %s\n", v.PrivateDnsName)
					fmt.Fprintf(os.Stdout, "  Architecture       : %s\n", v.Architecture)
					fmt.Fprintf(os.Stdout, "  Hypervisor         : %s\n", v.Hypervisor)
					fmt.Fprintf(os.Stdout, "  Image ID           : %s\n", v.ImageId)
					fmt.Fprintf(os.Stdout, "  Monitoring state   : %s\n", v.Monitoring.State)
					fmt.Fprintf(os.Stdout, "  State              : [%s], %s\n", v.State.Code, v.State.Name)
					fmt.Fprintf(os.Stdout, "  Availability zone  : %s\n", v.Placement.AvailabilityZone)
					fmt.Fprintf(os.Stdout, "  Root device name   : %s\n", v.RootDeviceName)
					fmt.Fprintf(os.Stdout, "  Root device type   : %s\n", v.RootDeviceType)
					fmt.Fprintf(os.Stdout, "  VPC ID             : %s\n", v.VpcId)
					fmt.Fprintf(os.Stdout, "\n")
				}

				fmt.Fprintf(os.Stdout, "Status         : %s\n", s.StackStatus)
				timestr := s.CreateTime
				t, err := time.Parse(time.RFC3339, s.CreateTime)
				if err == nil {
					timestr = t.Format(time.RFC1123)
				}

				fmt.Fprintf(os.Stdout, "Time created   : %s\n", timestr)
			}
		*/

		display(&stacks[0], 0)
	case "json":
		mi, err := json.MarshalIndent(stacks, "", "  ")
		if err != nil {
			util.PrintErrorAndExit(err.Error(), 1)
		}

		// this should be a prettified JSON output
		fmt.Println(string(mi))

		f := util.GetCliStringFlag(cmd, "out")
		if f != "" {
			err = ioutil.WriteFile(f, mi, 0644)
			if err != nil {
				util.PrintErrorAndExit(err.Error(), 1)
			}

			log.Println(fmt.Sprintf("Output written to %s.", f))
		}
	}
}

// display prints the field: value of the input struct recursively. Recursion level
// is provided for indention in printing.
func display(s interface{}, lvl int) {
	cnt := lvl * 2
	pad := ""
	for x := 0; x < cnt; x++ {
		pad += " "
	}

	rt := reflect.TypeOf(s).Elem()
	rv := reflect.ValueOf(s).Elem()

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i).Name
		value := rv.Field(i).Interface()

		switch rv.Field(i).Kind() {
		case reflect.String:
			fmt.Printf("%s%s: %s\n", pad, field, value)
		case reflect.Int32:
			fmt.Printf("%s%s: %i\n", pad, field, value)
		case reflect.Struct:
			fmt.Printf("%s[%s]\n", pad, field)
			v := rv.Field(i).Addr()
			display(v.Interface(), lvl+1)
		case reflect.Slice:
			fmt.Printf("%s[%s]\n", pad, field)
			slices, ok := value.([]stack.Instance)
			if ok {
				for _, slice := range slices {
					display(&slice, lvl+1)
				}
			} else {
				fmt.Printf("%s*** Not yet supported ***\n", pad)
			}
		}
	}
}
