package cli

import (
	"strconv"

	"github.com/spf13/cobra"
)

func flag(cmd *cobra.Command, f string) string {
	s := cmd.Flag(f).DefValue
	if cmd.Flag(f).Changed {
		s = cmd.Flag(f).Value.String()
	}

	return s
}

func GetCliStringFlag(cmd *cobra.Command, f string) string {
	return flag(cmd, f)
}

func GetCliIntFlag(cmd *cobra.Command, f string) int {
	v, err := strconv.Atoi(flag(cmd, f))
	if err != nil {
		return 0
	}

	return v
}

func GetCliInt64Flag(cmd *cobra.Command, f string) int64 {
	v, err := strconv.ParseInt(flag(cmd, f), 10, 64)
	if err != nil {
		return 0
	}

	return v
}
