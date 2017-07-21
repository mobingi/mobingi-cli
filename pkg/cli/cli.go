package cli

import (
	"strconv"

	"github.com/spf13/cobra"
)

func GetCliStringFlag(cmd *cobra.Command, f string) string {
	s := cmd.Flag(f).DefValue
	if cmd.Flag(f).Changed {
		s = cmd.Flag(f).Value.String()
	}

	return s
}

func GetCliIntFlag(cmd *cobra.Command, f string) int {
	s := cmd.Flag(f).DefValue
	if cmd.Flag(f).Changed {
		s = cmd.Flag(f).Value.String()
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return v
}

/*
func WriteToFile(f string, contents []byte) error {
	err := ioutil.WriteFile(f, contents, 0644)
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("output written to %s", f))
	return nil
}
*/
