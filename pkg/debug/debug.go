package debug

import (
	"fmt"
	"log"

	"github.com/fatih/color"
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
}
