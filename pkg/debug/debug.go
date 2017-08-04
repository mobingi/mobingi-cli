package debug

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/mobingilabs/mocli/pkg/cli/confmap"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var Verbose bool

// Info prints `v` into standard output (via log) with a green prefix "info:".
func Info(v ...interface{}) {
	green := color.New(color.FgGreen).SprintFunc()
	m := fmt.Sprintln(v...)
	log.Print(fmt.Sprintf("%s %s", green("info:"), m))
}

// Error prints `v` into standard output (via log) with a red prefix "error:".
func Error(v ...interface{}) {
	red := color.New(color.FgRed).SprintFunc()
	m := fmt.Sprintln(v...)
	log.Print(fmt.Sprintf("%s %s", red("error:"), m))

	// stack trace from 'errors'
	if viper.GetBool(confmap.ConfigKey("debug")) {
		err := fmt.Errorf(m)
		err = errors.WithStack(err)
		fmt.Printf("%+v\n", err)
	}
}
