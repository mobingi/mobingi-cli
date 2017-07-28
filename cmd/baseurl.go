package cmd

import (
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/constants"
	"github.com/spf13/cobra"
)

func BaseApiUrl(cmd *cobra.Command) string {
	base := cli.GetCliStringFlag(cmd, "url")
	if base == "" {
		base = constants.PROD_API_BASE
		if cli.IsDevMode() {
			base = constants.DEV_API_BASE
		}
	}

	return base
}

func BaseRegUrl(cmd *cobra.Command) string {
	base := cli.GetCliStringFlag(cmd, "rurl")
	if base == "" {
		base = constants.PROD_REG_BASE
		if cli.IsDevMode() {
			base = constants.DEV_REG_BASE
		}
	}

	return base
}
