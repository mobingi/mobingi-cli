package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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
		printStackText(os.Stdout, &stacks[0], 0)
		f := util.GetCliStringFlag(cmd, "out")
		if f != "" {
			fp, err := os.Create(f)
			if err != nil {
				util.PrintErrorAndExit(err.Error(), 1)
			}

			defer fp.Close()
			w := bufio.NewWriter(fp)
			defer w.Flush()
			printStackText(w, &stacks[0], 0)
			log.Println(fmt.Sprintf("Output written to %s.", f))
		}
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

// printStackText prints the field: value of the input struct recursively. Recursion level
// is provided for indention in printing.
func printStackText(w io.Writer, s interface{}, lvl int) {
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
			fmt.Fprintf(w, "%s%s: %s\n", pad, field, value)
		case reflect.Int32:
			fmt.Fprintf(w, "%s%s: %i\n", pad, field, value)
		case reflect.Struct:
			fmt.Fprintf(w, "%s[%s]\n", pad, field)
			v := rv.Field(i).Addr()
			printStackText(w, v.Interface(), lvl+1)
		case reflect.Slice:
			fmt.Fprintf(w, "%s[%s]\n", pad, field)
			slices, ok := value.([]stack.Instance)
			if ok {
				for _, slice := range slices {
					printStackText(w, &slice, lvl+1)
					fmt.Fprintf(w, "\n")
				}
			} else {
				fmt.Fprintf(w, "%s*** Not yet supported ***\n", pad)
			}
		}
	}
}
