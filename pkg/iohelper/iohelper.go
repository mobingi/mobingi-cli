package iohelper

import (
	"io/ioutil"

	"github.com/pkg/errors"
)

func WriteToFile(f string, contents []byte) error {
	err := ioutil.WriteFile(f, contents, 0644)
	if err != nil {
		return errors.Wrap(err, "write file failed")
	}

	return nil
}
