package cmd

import (
	"fmt"

	"github.com/mobingilabs/mobingi-sdk-go/client"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
)

func exitOn401(resp *client.Response) {
	if resp != nil {
		if resp.StatusCode == 401 {
			d.ErrorExit(fmt.Errorf(resp.Status), 1)
		}
	}
}
