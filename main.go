// Package mobingi-cli is a command line interface client for Mobingi API.
package main

import (
	"log"

	"github.com/mobingi/mobingi-cli/cmd"
)

func main() {
	log.SetPrefix("[mobingi-cli]: ")
	log.SetFlags(0)
	cmd.Execute()
}
