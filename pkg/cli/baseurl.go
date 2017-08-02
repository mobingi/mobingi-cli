package cli

import (
	"github.com/mobingilabs/mocli/pkg/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func BaseApiUrl(cmd *cobra.Command) string {
	base := GetCliStringFlag(cmd, "url")
	if base == "" {
		switch viper.GetString("runenv") {
		case "dev":
			base = constants.DEV_API_BASE
		case "qa":
			base = constants.QA_API_BASE
		default:
			base = constants.PROD_API_BASE
		}
	}

	return base
}

func BaseRegUrl(cmd *cobra.Command) string {
	base := GetCliStringFlag(cmd, "rurl")
	if base == "" {
		switch viper.GetString("runenv") {
		case "dev":
			base = constants.DEV_REG_BASE
		case "qa":
			base = constants.QA_REG_BASE
		default:
			base = constants.PROD_REG_BASE
		}
	}

	return base
}
