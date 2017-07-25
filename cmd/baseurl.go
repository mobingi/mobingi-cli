package cmd

import (
	"github.com/mobingilabs/mocli/pkg/check"
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/constants"
	"github.com/spf13/cobra"
)

func BaseApiUrl(cmd *cobra.Command) string {
	base := cli.GetCliStringFlag(cmd, "url")
	if base == "" {
		if check.IsDevMode() {
			base = constants.DEV_API_BASE
		} else {
			base = constants.PROD_API_BASE
		}
	}

	return base
}
