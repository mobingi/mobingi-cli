package sesha3

import (
	"os"
	"testing"

	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
)

func TestNewClient(t *testing.T) {
	url := &SeshaClientInput{URL: "https://sesha3.labs.mobingi.com:8568/d3aiwuxow4mxnsgc4j7usvcpw0bjh27kg94c/"}
	cli, err := NewClient(url)
	if err != nil {
		t.Fatal("expected nil error")
	}

	if cli == nil {
		t.Fatal("expected a valid client object")
	}
}

func TestNew(t *testing.T) {
	svc := New(nil)
	if svc != nil {
		t.Fatal("expected nil")
	}
}

func TestGetSessionUrlDevAcct(t *testing.T) {
	return
	if os.Getenv("MOBINGI_CLIENT_ID") != "" && os.Getenv("MOBINGI_CLIENT_SECRET") != "" {
		sess, _ := session.New(&session.Config{
			BaseApiUrl: "https://apidev.mobingi.com",
			ApiVersion: 2,
		})

		svc := New(sess)
		in := &GetSessionUrlInput{
			StackId:  "mo-58c2297d25645-Sd2aHRDq0-tk",
			IpAddr:   "54.238.234.202",
			InstUser: "ec2-user",
		}

		resp, body, u, err := svc.GetSessionUrl(in)
		if err != nil {
			t.Errorf("expecting nil error, received %v", err)
		}

		_, _, _ = resp, body, u
	}
}
