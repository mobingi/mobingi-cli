package svrconf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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

func TestUpdateEnvVars(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Error reading body: %v", err)
		}

		var m map[string]interface{}
		_ = json.Unmarshal(body, &m)
		_, ok := m["envvars"]
		if ok {
			if !strings.HasPrefix(fmt.Sprintf("%v", m["envvars"]), "map") {
				t.Errorf("Should start with map, got %v", m["envvars"])
			}
		}

		w.Write([]byte(r.URL.String()))
	}))

	defer ts.Close()
	sess, _ := session.New(&session.Config{
		BaseApiUrl: ts.URL,
		ApiVersion: 2,
	})

	cnf := New(sess)
	_, body, err := cnf.UpdateEnvVars(&ServerConfigUpdateEnvVarsInput{
		StackId: "id",
		EnvVars: "KEY1:value1,KEY2:value2",
	})

	_, _ = body, err
}

// local test for dev; requires the following environment variables:
// MOBINGI_CLIENT_ID, MOBINGI_CLIENT_SECRET (dev accounts only)
func TestUpdateEnvVarsDevAcct(t *testing.T) {
	return
	if os.Getenv("MOBINGI_CLIENT_ID") != "" && os.Getenv("MOBINGI_CLIENT_SECRET") != "" {
		sess, _ := session.New(&session.Config{
			BaseApiUrl: "https://apidev.mobingi.com",
			ApiVersion: 2,
		})

		cnf := New(sess)
		_, body, err := cnf.UpdateEnvVars(&ServerConfigUpdateEnvVarsInput{
			StackId: "mo-58c2297d25645-NZvoZDVMg-tk",
			EnvVars: "KEY1:value1,KEY2:value2",
		})

		if err != nil {
			t.Errorf("Expecting nil error, received %v", err)
		}

		_ = body
	}
}

func TestUpdateFilePath(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Error reading body: %v", err)
		}

		var m map[string]interface{}
		_ = json.Unmarshal(body, &m)
		_, ok := m["filepath"]
		if ok {
			if fmt.Sprintf("%v", m["filepath"]) != "filepath" {
				t.Errorf("Expecting 'filepath', got %v", m["envvars"])
			}
		}

		w.Write([]byte(r.URL.String()))
	}))

	defer ts.Close()
	sess, _ := session.New(&session.Config{
		BaseApiUrl: ts.URL,
		ApiVersion: 2,
	})

	cnf := New(sess)
	_, body, err := cnf.UpdateFilePath(&ServerConfigUpdateFilePathInput{
		StackId:  "id",
		FilePath: "filepath",
	})

	_, _ = body, err
}

// local test for dev; requires the following environment variables:
// MOBINGI_CLIENT_ID, MOBINGI_CLIENT_SECRET (dev accounts only)
func TestUpdateFilePathDevAcct(t *testing.T) {
	return
	if os.Getenv("MOBINGI_CLIENT_ID") != "" && os.Getenv("MOBINGI_CLIENT_SECRET") != "" {
		sess, _ := session.New(&session.Config{
			BaseApiUrl: "https://apidev.mobingi.com",
			ApiVersion: 2,
		})

		cnf := New(sess)
		_, body, err := cnf.UpdateFilePath(&ServerConfigUpdateFilePathInput{
			StackId:  "mo-58c2297d25645-NZvoZDVMg-tk",
			FilePath: "git://github.com/mobingilabs/default1",
		})

		if err != nil {
			t.Errorf("Expecting nil error, received %v", err)
		}

		_ = body
	}
}
