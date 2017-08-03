package iohelper

import "io/ioutil"

func WriteToFile(f string, contents []byte) error {
	err := ioutil.WriteFile(f, contents, 0644)
	if err != nil {
		return err
	}

	return nil
}
