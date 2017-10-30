package pullr

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
)

func TestPublishAndReadDevAcct(t *testing.T) {
	return
	if os.Getenv("AWS_ACCESS_KEY_ID") != "" && os.Getenv("AWS_SECRET_ACCESS_KEY") != "" {
		qc := QueueClient{}

		for i := 0; i < 100; i++ {
			// publish first
			err := qc.Publish(fmt.Sprintf("message%03d", i))
			if err != nil {
				t.Error(err)
			}

			// pause
			time.Sleep(time.Millisecond * 50)

			// read
			_, m, h, err := qc.Read(nil)
			if err != nil {
				t.Error(err)
			}

			debug.Info(m)
			debug.Info(h)
		}
	}
}
