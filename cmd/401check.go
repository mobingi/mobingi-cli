package cmd

import (
	"fmt"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/client"
)

func exitOn401(resp *client.Response) {
	if resp != nil {
		if resp.StatusCode == 401 {
			cli.ErrorExit(fmt.Errorf(resp.Status), 1)
		}
	}
}
