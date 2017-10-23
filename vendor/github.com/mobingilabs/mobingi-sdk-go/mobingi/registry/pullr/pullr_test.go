package pullr

import (
	"os"
	"testing"
	"time"

	"github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
)

func TestPublishAndReadDevAcct(t *testing.T) {
	return
	if os.Getenv("AWS_ACCESS_KEY_ID") != "" && os.Getenv("AWS_SECRET_ACCESS_KEY") != "" {
		qc := QueueClient{}

		// publish first
		err := qc.Publish("message")
		if err != nil {
			t.Error(err)
		}

		// pause
		time.Sleep(time.Millisecond * 500)

		// read
		_, m, h, err := qc.Read(nil)
		if err != nil {
			t.Error(err)
		}

		debug.Info(m)
		debug.Info(h)
	}
}
