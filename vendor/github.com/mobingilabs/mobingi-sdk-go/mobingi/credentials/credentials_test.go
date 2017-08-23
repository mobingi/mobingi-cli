package credentials

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
)

func TestNew(t *testing.T) {
	sess, _ := session.New(&session.Config{})
	creds := New(sess)
	if creds == nil {
		t.Errorf("Expecting non-nil")
	}
}

func TestList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.URL.String()))
	}))

	defer ts.Close()
	sess, _ := session.New(&session.Config{
		BaseApiUrl: ts.URL,
		ApiVersion: 2,
	})

	creds := New(sess)
	_, body, _ := creds.List(&CredentialsListInput{})
	if string(body) != "/v2/credentials/aws" {
		t.Errorf("Expecting '/v2/credentials/aws', got %s", string(body))
	}
}
