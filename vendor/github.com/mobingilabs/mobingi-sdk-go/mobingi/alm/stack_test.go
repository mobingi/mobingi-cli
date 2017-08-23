package alm

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
)

func TestNew(t *testing.T) {
	sess, _ := session.New(&session.Config{})
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
	sess, _ := session.New(&session.Config{
		BaseApiUrl: ts.URL,
		ApiVersion: 2,
	})

	alm := New(sess)
	_, body, _ := alm.List()
	if string(body) != "hello" {
		t.Errorf("Expecting 'hello', received %s", string(body))
	}
}

// local test for dev; requires the following environment variables:
// MOBINGI_CLIENT_ID, MOBINGI_CLIENT_SECRET (dev accounts only)
func TestListDevAcct(t *testing.T) {
	if os.Getenv("MOBINGI_CLIENT_ID") != "" && os.Getenv("MOBINGI_CLIENT_SECRET") != "" {
		sess, _ := session.New(&session.Config{
			BaseApiUrl: "https://apidev.mobingi.com",
			ApiVersion: 2,
		})

		alm := New(sess)
		_, _, err := alm.List()
		if err != nil {
			t.Errorf("Expecting nil error, received %v", err)
		}
	}
}

func TestDescribe(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.URL.String()))
	}))

	defer ts.Close()
	sess, _ := session.New(&session.Config{
		BaseApiUrl: ts.URL,
		ApiVersion: 2,
	})

	alm := New(sess)
	_, body, _ := alm.Describe(&StackDescribeInput{StackId: "id"})
	if string(body) != "/v2/alm/stack/id" {
		t.Errorf("Expecting '/v2/alm/stack/id', got %s", string(body))
	}
}

// local test for dev; requires the following environment variables:
// MOBINGI_CLIENT_ID, MOBINGI_CLIENT_SECRET (dev accounts only)
func TestDescribeDevAcct(t *testing.T) {
	if os.Getenv("MOBINGI_CLIENT_ID") != "" && os.Getenv("MOBINGI_CLIENT_SECRET") != "" {
		sess, _ := session.New(&session.Config{
			BaseApiUrl: "https://apidev.mobingi.com",
			ApiVersion: 2,
		})

		alm := New(sess)
		_, _, err := alm.Describe(&StackDescribeInput{StackId: "id"})
		if err != nil {
			t.Errorf("Expecting nil error, received %v", err)
		}
	}
}

func TestCreate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Error reading body: %v", err)
		}

		sbody, _ := url.QueryUnescape(string(body))
		params := strings.Split(sbody, "&")
		for _, v := range params {
			kvs := strings.Split(v, "=")
			if len(kvs) == 2 {
				if kvs[0] == "configurations" {
					var in StackCreateConfig
					err = json.Unmarshal([]byte(kvs[1]), &in)
					if err != nil {
						t.Errorf("Error unmarshaling body: %v", err)
					}
				}
			}
		}

		w.Write([]byte(r.URL.String()))
	}))

	defer ts.Close()
	sess, _ := session.New(&session.Config{
		BaseApiUrl: ts.URL,
		ApiVersion: 2,
	})

	alm := New(sess)
	_, body, _ := alm.Create(&StackCreateInput{
		CredId: "id",
		Configurations: StackCreateConfig{
			Region:            "ap-northeast-1",
			Architecture:      "art_elb",
			Type:              "m3.medium",
			Image:             "mobingi/ubuntu-apache2-php7:7.1",
			DockerHubUsername: "",
			DockerHubPassword: "",
			Min:               2,
			Max:               10,
			SpotRange:         50,
			Nickname:          "stack_create_go_test",
			Code:              "github.com/mobingilabs/default-site-php",
			GitReference:      "master",
			GitPrivateKey:     "",
			// no database
			// no elasticache
		},
	})

	if string(body) != "/v2/alm/stack" {
		t.Errorf("Expecting '/v2/alm/stack', got %s", string(body))
	}
}

// Local test for dev; requires the following environment variables:
// MOBINGI_CLIENT_ID, MOBINGI_CLIENT_SECRET (dev accounts only)
//
// NOTE: This is an example of how to actually create a stack in AWS. If you want to try it out,
// set the environment variables above with your dev account (local) then run the test.
// Don't forget to uncomment the first line `return` statement.
func TestCreateDevAcct(t *testing.T) {
	return
	if os.Getenv("MOBINGI_CLIENT_ID") != "" && os.Getenv("MOBINGI_CLIENT_SECRET") != "" {
		sess, _ := session.New(&session.Config{
			BaseApiUrl: "https://apidev.mobingi.com",
			ApiVersion: 2,
		})

		alm := New(sess)

		// these are the default values in `mobingi-cli`
		cnf := StackCreateConfig{
			Region:            "ap-northeast-1",
			Architecture:      "art_elb",
			Type:              "m3.medium",
			Image:             "mobingi/ubuntu-apache2-php7:7.1",
			DockerHubUsername: "",
			DockerHubPassword: "",
			Min:               2,
			Max:               10,
			SpotRange:         50,
			Nickname:          "stack_create_go_test",
			Code:              "github.com/mobingilabs/default-site-php",
			GitReference:      "master",
			GitPrivateKey:     "",
			// no database
			// no elasticache
		}

		in := &StackCreateInput{
			// Vendor not provided; will default to "aws"
			Region:         "ap-northeast-1",
			Configurations: cnf,
			// CredId not provided; sdk will attempt to retrieve thru API
		}

		_, body, err := alm.Create(in)
		if err != nil {
			t.Errorf("Expecting nil error, received %v", err)
		}

		_ = body
	}
}

func TestUpdate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Error reading body: %v", err)
		}

		var m map[string]interface{}
		var in StackCreateConfig
		err = json.Unmarshal(body, &m)
		if err != nil {
			t.Errorf("Expecting nil error, received %v", err)
		}

		_, ok := m["configurations"]
		if ok {
			err = json.Unmarshal([]byte(m["configurations"].(string)), &in)
			if err != nil {
				t.Errorf("Expecting nil error, received %v", err)
			}

			if in.SpotRange.(float64) != 40 {
				t.Errorf("Expecting a 40 spot range, got %v", in.SpotRange)
			}
		}

		w.Write([]byte(r.URL.String()))
	}))

	defer ts.Close()
	sess, _ := session.New(&session.Config{
		BaseApiUrl: ts.URL,
		ApiVersion: 2,
	})

	alm := New(sess)
	_, body, err := alm.Update(&StackUpdateInput{
		StackId: "id",
		Configurations: StackCreateConfig{
			SpotRange: 40,
		},
	})

	_, _ = body, err
}

// Local test for dev; requires the following environment variables:
// MOBINGI_CLIENT_ID, MOBINGI_CLIENT_SECRET (dev accounts only)
func TestUpdateDevAcct(t *testing.T) {
	return
	if os.Getenv("MOBINGI_CLIENT_ID") != "" && os.Getenv("MOBINGI_CLIENT_SECRET") != "" {
		sess, _ := session.New(&session.Config{
			BaseApiUrl: "https://apidev.mobingi.com",
			ApiVersion: 2,
		})

		alm := New(sess)

		cnf := StackCreateConfig{
			SpotRange: 40,
		}

		in := &StackUpdateInput{
			// this id is an actual stack id; if you want to test, use an actual id
			StackId:        "mo-58c2297d25645-4t7SRL1P-tk",
			Configurations: cnf,
		}

		_, body, err := alm.Update(in)
		if err != nil {
			t.Errorf("Expecting nil error, received %v", err)
		}

		_ = body
	}
}

func TestDelete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.URL.String()))
	}))

	defer ts.Close()
	sess, _ := session.New(&session.Config{
		BaseApiUrl: ts.URL,
		ApiVersion: 2,
	})

	alm := New(sess)
	_, body, err := alm.Delete(&StackDeleteInput{
		StackId: "id",
	})

	if string(body) != "/v2/alm/stack/id" {
		t.Errorf("Expecting '/v2/alm/stack/id', received %v", string(body))
	}

	_ = err
}

// Local test for dev; requires the following environment variables:
// MOBINGI_CLIENT_ID, MOBINGI_CLIENT_SECRET (dev accounts only)
func TestDeleteDevAcct(t *testing.T) {
	return
	if os.Getenv("MOBINGI_CLIENT_ID") != "" && os.Getenv("MOBINGI_CLIENT_SECRET") != "" {
		sess, _ := session.New(&session.Config{
			BaseApiUrl: "https://apidev.mobingi.com",
			ApiVersion: 2,
		})

		alm := New(sess)

		in := &StackDeleteInput{
			// this id is an actual stack id; if you want to test, use an actual id
			StackId: "mo-58c2297d25645-4t7SRL1P-tk",
		}

		_, body, err := alm.Delete(in)
		if err != nil {
			t.Errorf("Expecting nil error, received %v", err)
		}

		_ = body
	}
}
