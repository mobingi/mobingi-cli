package svrconf

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
)

func TestNew(t *testing.T) {
	sess, _ := session.New(&session.Config{})
	cnf := New(sess)
	if cnf == nil {
		t.Errorf("Expecting non-nil")
	}
}

func TestGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.URL.String()))
	}))

	defer ts.Close()
	sess, _ := session.New(&session.Config{
		BaseApiUrl: ts.URL,
		ApiVersion: 2,
	})

	cnf := New(sess)
	_, body, err := cnf.Get(&ServerConfigGetInput{StackId: "id"})

	if string(body) != "/v2/alm/serverconfig?stack_id=id" {
		t.Errorf("Expecting '/v2/alm/serverconfig?stack_id=id', received %v", string(body))
	}

	_ = err
}

// local test for dev; requires the following environment variables:
// MOBINGI_CLIENT_ID, MOBINGI_CLIENT_SECRET (dev accounts only)
func TestGetDevAcct(t *testing.T) {
	return
	if os.Getenv("MOBINGI_CLIENT_ID") != "" && os.Getenv("MOBINGI_CLIENT_SECRET") != "" {
		sess, _ := session.New(&session.Config{
			BaseApiUrl: "https://apidev.mobingi.com",
			ApiVersion: 2,
		})

		cnf := New(sess)
		_, body, err := cnf.Get(&ServerConfigGetInput{StackId: "mo-58c2297d25645-4t7SRL1P-tk"})
		if err != nil {
			t.Errorf("Expecting nil error, received %v", err)
		}

		_ = body
	}
}
