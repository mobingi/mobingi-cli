package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient(nil)
	if c == nil {
		t.Errorf("Expected client even with nil config")
	}
}

func TestGetTagDigest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected method %s; got %q", http.MethodGet, r.Method)
		}

		accpt := r.Header.Get("Accept")
		if accpt != "application/vnd.docker.distribution.manifest.v2+json" {
			t.Errorf("Expected 'application/vnd.docker.distribution.manifest.v2+json'; got %q", accpt)
		}

		w.Header().Add("Etag", "\"sha256:testdigest\"")
	}))

	defer ts.Close()
	c := NewClient(&Config{RootUrl: ts.URL})
	digest, err := c.GetTagDigest("/test")
	if err != nil {
		t.Errorf("Expected success; got %v", err)
	}

	if digest == "" {
		t.Errorf("Expected a digest value; got none")
	}
}
