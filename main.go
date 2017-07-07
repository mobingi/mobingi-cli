// Package mocli is a command line interface client for Mobingi API.
package main

import (
	"log"

	"github.com/mobingilabs/mocli/cmd"
)

func main() {
	log.SetPrefix("[mocli]: ")
	log.SetFlags(0)
	cmd.Execute()
}
