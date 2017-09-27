// Package mobingi-cli is a command line interface client for Mobingi API.
package main

import (
	"log"

	"github.com/mobingi/mobingi-cli/cmd"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
)

func main() {
	pfx := "[" + cmdline.Args0() + "]: "
	log.SetPrefix(pfx)
	log.SetFlags(0)
	cmd.Execute()
}
