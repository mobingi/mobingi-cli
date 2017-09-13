package sesha3

import (
	"os"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	url := &SeshaClientInput{URL: "https://sesha3.labs.mobingi.com:8568/d3aiwuxow4mxnsgc4j7usvcpw0bjh27kg94c/"}
	sesha3cli, _ := NewClient(url)
	go func() {
		time.Sleep(time.Second)
		os.Exit(0)
	}()
	err := sesha3cli.Run()
	_ = err
}
