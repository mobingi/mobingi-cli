package alm

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
)

func TestNew(t *testing.T) {
	sess, _ := session.NewSession(&session.Config{})
	alm := New(sess)
	if alm == nil {
		t.Errorf("Expecting non-nil")
	}
}

func TestList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}))

	defer ts.Close()
	sess, _ := session.NewSession(&session.Config{
		BaseApiUrl: ts.URL,
		ApiVersion: 2,
	})

	alm := New(sess)
	_, body, _ := alm.List()
	if string(body) != "hello" {
		t.Errorf("Expecting 'hello', received %s", string(body))
	}
}
