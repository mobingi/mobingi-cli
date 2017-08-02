package iohelper

import (
	"fmt"
	"io/ioutil"

	d "github.com/mobingilabs/mocli/pkg/debug"
)

func WriteToFile(f string, contents []byte) error {
	err := ioutil.WriteFile(f, contents, 0644)
	if err != nil {
		return err
	}

	d.Info(fmt.Sprintf("write to file: %s", f))
	return nil
}
