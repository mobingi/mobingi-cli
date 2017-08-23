package client

import (
	"bufio"
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewSimpleHttpClient(t *testing.T) {
	c := NewSimpleHttpClient()
	if c == nil {
		t.Errorf("Expected an object, received nil")
	}
}

func TestDo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}))

	defer ts.Close()
	c := NewSimpleHttpClient()
	r, err := http.NewRequest(http.MethodGet, ts.URL+"/test", nil)
	if err != nil {
		t.Errorf("New request failed: %#v", err)
	}

	resp, body, err := c.Do(r)
	if err != nil {
		t.Errorf("Do failed: %#v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, received %v", resp.StatusCode)
	}

	if string(body) != "hello" {
		t.Errorf("Expected body 'hello', received %s", string(body))
	}
}

func TestDoWithTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 1)
		w.Write([]byte("hello"))
	}))

	defer ts.Close()
	c := NewSimpleHttpClient(&Config{Timeout: time.Second * 2})
	r, err := http.NewRequest(http.MethodGet, ts.URL+"/test", nil)
	if err != nil {
		t.Errorf("New request failed: %#v", err)
	}

	_, _, err = c.Do(r)
	if err != nil {
		t.Errorf("Do failed: %#v", err)
	}
}

func TestDoVerbose(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}))

	defer ts.Close()

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	logger := log.New(w, "TEST", 0)

	c := NewSimpleHttpClient(&Config{Verbose: true, Logger: logger})
	r, err := http.NewRequest(http.MethodGet, ts.URL+"/test", nil)
	if err != nil {
		t.Errorf("New request failed: %#v", err)
	}

	_, _, err = c.Do(r)
	if err != nil {
		t.Errorf("Do failed: %#v", err)
	}

	w.Flush()
	if b.Len() == 0 {
		t.Errorf("Buffer should not be empty")
	}
}
