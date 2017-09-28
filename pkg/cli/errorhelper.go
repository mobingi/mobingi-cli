package cli

import (
	"github.com/mobingi/mobingi-cli/pkg/cli/confmap"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/spf13/viper"
)

func ErrorExit(err interface{}, code int) {
	dbg := viper.GetBool(confmap.ConfigKey("debug"))
	if dbg {
		debug.ErrorTraceExit(err, code)
	} else {
		debug.ErrorExit(err, code)
	}
}
