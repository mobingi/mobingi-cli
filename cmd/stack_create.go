package cmd

import (
	"github.com/spf13/cobra"
)

func StackCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a stack",
		Long:  `Create a stack.`,
		Run:   createStack,
	}

	return cmd
}

func createStack(cmd *cobra.Command, args []string) {
	/*
		c := client.NewClient(client.NewApiConfig(cmd))
		body, err := c.AuthGet("/alm/stack")
		check.ErrorExit(err, 1)

		var stacks []stack.ListStack
		err = json.Unmarshal(body, &stacks)
		check.ErrorExit(err, 1)

		pfmt := cli.GetCliStringFlag(cmd, "fmt")
		switch pfmt {
		case "raw":
			fmt.Println(string(body))

			// write to file option
			f := cli.GetCliStringFlag(cmd, "out")
			if f != "" {
				err = ioutil.WriteFile(f, body, 0644)
				check.ErrorExit(err, 1)
				d.Info(fmt.Sprintf("Output written to %s.", f))
			}
		case "text":
			indent := cli.GetCliIntFlag(cmd, "indent")
			stack.PrintR(os.Stdout, &stacks[0], 0, indent)

			// write to file option
			f := cli.GetCliStringFlag(cmd, "out")
			if f != "" {
				fp, err := os.Create(f)
				check.ErrorExit(err, 1)

				defer fp.Close()
				w := bufio.NewWriter(fp)
				defer w.Flush()
				stack.PrintR(w, &stacks[0], 0, indent)
				d.Info(fmt.Sprintf("Output written to %s.", f))
			}
		case "json":
			indent := cli.GetCliIntFlag(cmd, "indent")
			mi, err := json.MarshalIndent(stacks, "", pretty.Indent(indent))
			check.ErrorExit(err, 1)

			fmt.Println(string(mi))

			// write to file option
			f := cli.GetCliStringFlag(cmd, "out")
			if f != "" {
				err = ioutil.WriteFile(f, mi, 0644)
				check.ErrorExit(err, 1)
				d.Info(fmt.Sprintf("Output written to %s.", f))
			}
		default:
			if pfmt == "min" || pfmt == "" {
				stbl := term.NewTable(0, 10, 5, ' ', 0)
				fmt.Fprintf(stbl, "STACK ID\tSTACK NAME\tPLATFORM\tSTATUS\tREGION\tLAUNCHED\n")
				for _, s := range stacks {
					timestr := s.CreateTime
					t, err := time.Parse(time.RFC3339, s.CreateTime)
					if err == nil {
						timestr = t.Format(time.RFC1123)
					}

					platform := "?"
					if s.Configuration.AWS != "" {
						platform = "AWS"
					}

					fmt.Fprintf(stbl, "%s\t%s\t%s\t%s\t%s\t%s\n",
						s.StackId,
						s.Nickname,
						platform,
						s.StackStatus,
						s.Configuration.Region,
						timestr)
				}

				term.Print(stbl)
				term.Flush()
			}
		}
	*/
}
