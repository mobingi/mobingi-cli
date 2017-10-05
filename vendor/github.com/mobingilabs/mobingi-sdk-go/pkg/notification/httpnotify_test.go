package notification

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewSimpleHttpNotify(t *testing.T) {
	n := NewSimpleHttpNotify("test")
	if n == nil {
		t.Fatal("should not be nil")
	}
}

func TestNotify(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		if string(body) != "test" {
			t.Fatal("should be test")
		}
	}))

	defer ts.Close()
	n := NewSimpleHttpNotify(ts.URL)
	n.Notify([]byte("test"))

}
