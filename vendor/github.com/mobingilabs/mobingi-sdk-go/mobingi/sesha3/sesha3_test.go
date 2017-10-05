package sesha3

import (
	"os"
	"testing"

	"github.com/mobingilabs/mobingi-sdk-go/client"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
)

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
			// ApiVersion: 2,
			HttpClientConfig: &client.Config{
				Verbose: true,
			},
		})

		svc := New(sess)
		in := &GetSessionUrlInput{
			// StackId:  "mo-58c2297d25645-Sd2aHRDq0-tk",
			// IpAddr:   "54.238.234.202",
			StackId:  "mo-58c2297d25645-thHNtg0YS-tk",
			IpAddr:   "13.114.34.121",
			Flag:     "fweb",
			InstUser: "ec2-user",
		}

		resp, body, u, err := svc.GetSessionUrl(in)
		if err != nil {
			t.Errorf("expecting nil error, received %v", err)
		}

		// log.Println(resp, string(body), u)
		_, _, _ = resp, body, u
	}
}
