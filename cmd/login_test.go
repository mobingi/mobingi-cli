package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProcess(t *testing.T) {
	path := "/v2/access_token"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected method %q; got %q", "POST", r.Method)
		}

		if r.Header == nil {
			t.Error("Expected non-nil request Header")
		}

		if r.URL.Path != path {
			t.Errorf("Expected %v, got %v", path, r.URL.Path)
		}

		/*
			switch r.URL.Path {
			default:
				t.Errorf("No testing for this case yet : %q", r.URL.Path)
			case case1_empty:
				t.Logf("case %v ", case1_empty)
			case case2_set_header:
				t.Logf("case %v ", case2_set_header)
				if r.Header.Get("API-Key") != "fookey" {
					t.Errorf("Expected 'API-Key' == %q; got %q", "fookey", r.Header.Get("API-Key"))
				}
			}
		*/
	}))

	defer ts.Close()

	// c := cli.New("")
	// c.Requester.Post(ts.URL + path).End()
}
