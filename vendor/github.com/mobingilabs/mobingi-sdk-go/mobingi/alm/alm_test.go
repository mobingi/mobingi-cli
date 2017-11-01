package alm

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
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

	// v2 api
	sess, _ := session.New(&session.Config{
		BaseApiUrl: ts.URL,
		ApiVersion: 2,
	})

	alm := New(sess)
	_, body, _ := alm.List()
	if string(body) != "hello" {
		t.Errorf("Expecting 'hello', received %s", string(body))
	}

	// latest api
	sess, _ = session.New(&session.Config{BaseApiUrl: ts.URL})
	alm = New(sess)
	_, body, _ = alm.List()
	if string(body) != "hello" {
		t.Errorf("Expecting 'hello', received %s", string(body))
	}
}

// local test for dev; requires the following environment variables:
// MOBINGI_CLIENT_ID, MOBINGI_CLIENT_SECRET (dev accounts only)
func TestListDevAcct(t *testing.T) {
	return
	if os.Getenv("MOBINGI_CLIENT_ID") != "" && os.Getenv("MOBINGI_CLIENT_SECRET") != "" {
		// v2 api
		sess, _ := session.New(&session.Config{
			BaseApiUrl: "https://apidev.mobingi.com",
			ApiVersion: 2,
		})

		alm := New(sess)
		resp, body, err := alm.List()
		if err != nil {
			t.Errorf("Expecting nil error, received %v", err)
		}

		log.Println(resp)
		log.Println(string(body))

		// latest api
		sess, _ = session.New(&session.Config{BaseApiUrl: "https://apidev.mobingi.com"})
		alm = New(sess)
		resp, body, err = alm.List()
		if err != nil {
			t.Errorf("Expecting nil error, received %v", err)
		}

		log.Println(resp)
		log.Println(string(body))
		// _, _ = resp, body
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

func TestCreateAlmStackDevAcct(t *testing.T) {
	return
	if os.Getenv("MOBINGI_CLIENT_ID") != "" && os.Getenv("MOBINGI_CLIENT_SECRET") != "" && os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		var AWSSingleEC2JSON = `{
  "version": "2017-03-03",
  "label": "template version label #1",
  "description": "This template creates a sample stack with EC2 instance on AWS",
  "vendor": {
    "aws": {
      "cred": "` + os.Getenv("AWS_ACCESS_KEY_ID") + `",
      "region": "ap-northeast-1"
    }
  },
  "configurations": [
    {
      "role": "web",
      "flag": "Single1",
      "provision": {
        "instance_type": "t2.micro",
        "instance_count": 1,
        "keypair": false,
        "subnet": {
          "cidr": "10.0.1.0/24",
          "public": true,
          "auto_assign_public_ip": true
        },
        "availability_zone": "ap-northeast-1c"
      }
    }
  ]
}`

		sess, _ := session.New(&session.Config{BaseApiUrl: "https://apidev.mobingi.com"})
		alm := New(sess)

		in := &StackCreateInput{
			AlmTemplate: &AlmTemplate{
				Contents: AWSSingleEC2JSON,
			},
		}

		resp, body, err := alm.Create(in)
		if err != nil {
			t.Errorf("Expecting nil error, received %v", err)
		}

		_, _ = resp, body
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

// Local test for dev; requires the following environment variables:
// MOBINGI_CLIENT_ID, MOBINGI_CLIENT_SECRET (dev accounts only)
func TestGetTemplateVersionsDevAcct(t *testing.T) {
	return
	if os.Getenv("MOBINGI_CLIENT_ID") != "" && os.Getenv("MOBINGI_CLIENT_SECRET") != "" {
		sess, _ := session.New(&session.Config{
			BaseApiUrl: "https://apidev.mobingi.com",
		})

		alm := New(sess)
		in := &GetTemplateVersionsInput{StackId: "mo-58c2297d25645-ASERav0N1-tk"}
		resp, body, err := alm.GetTemplateVersions(in)
		if err != nil {
			t.Errorf("Expecting nil error, received %v", err)
		}

		_, _ = resp, body
	}
}

func TestCompareTemplateDevAcct(t *testing.T) {
	return
	if os.Getenv("MOBINGI_CLIENT_ID") != "" && os.Getenv("MOBINGI_CLIENT_SECRET") != "" && os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		sess, _ := session.New(&session.Config{
			BaseApiUrl: "https://apidev.mobingi.com",
		})

		alm := New(sess)
		in := &CompareTemplateInput{
			SourceStackId:   "mo-58c2297d25645-PxviFSJQV-tk",
			SourceVersionId: "jbyW_PxMAauQmOS31dUhij4KIqHAtqW2",
			TargetVersionId: "1xoPd.cg3juHK94vC8IdUh1bexx7sQ1T",
		}

		resp, body, err := alm.CompareTemplate(in)
		if err != nil {
			t.Errorf("Expecting nil error, received %v", err)
		}

		var testbody string = `{
  "version": "2017-03-03",
  "label": "template layer #1",
  "description": "Wayland Test. This template creates a sample stack with EC2 instance on AWS",
  "vendor": {
	"aws": {
	  "cred": "` + os.Getenv("AWS_ACCESS_KEY_ID") + `",
	  "region": "ap-northeast-1"
    }
  },
  "configurations": [
    {
	  "role": "web",
	  "flag": "Layer1",
	  "provision": {
	    "instance_type": "t2.micro",
		"instance_count": 1,
		"keypair": false,
		"subnet": {
		  "cidr": "10.0.1.0/24",
		  "public": true,
		  "auto_assign_public_ip": true
	    },
		"availability_zone": "ap-northeast-1c"
	  }
    }
  ]
}`

		in = &CompareTemplateInput{
			SourceStackId:   "mo-58c2297d25645-PxviFSJQV-tk",
			SourceVersionId: "jbyW_PxMAauQmOS31dUhij4KIqHAtqW2",
			TargetBody:      testbody,
		}

		resp, body, err = alm.CompareTemplate(in)
		if err != nil {
			t.Errorf("Expecting nil error, received %v", err)
		}

		_, _ = resp, body
	}
}

func TestGetPemDevAcct(t *testing.T) {
	return
	if os.Getenv("MOBINGI_CLIENT_ID") != "" && os.Getenv("MOBINGI_CLIENT_SECRET") != "" {
		/*
			sess, _ := session.New(&session.Config{
				BaseApiUrl: "https://apidev.mobingi.com",
				ApiVersion: 2,
			})

			alm := New(sess)
			in := &GetPemInput{StackId: "mo-58c2297d25645-Sd2aHRDq0-tk"}
			resp, body, pem, err := alm.GetPem(in)
			if err != nil {
				t.Errorf("Expecting nil error, received %v", err)
			}
		*/

		sess, _ := session.New(&session.Config{
			BaseApiUrl: "https://apidev.mobingi.com",
		})

		alm := New(sess)
		in := &GetPemInput{
			StackId: "mo-58c2297d25645-HVhUlcmM-tk",
			Flag:    "fweb",
		}

		resp, body, pem, err := alm.GetPem(in)
		if err != nil {
			t.Errorf("Expecting nil error, received %v", err)
		}

		// log.Println(resp, string(body), string(pem))
		_, _, _ = resp, body, pem
	}
}

func TestWalkerDevAcct(t *testing.T) {
	// return
	if os.Getenv("MOBINGI_CLIENT_ID") != "" && os.Getenv("MOBINGI_CLIENT_SECRET") != "" {
		sess, _ := session.New(&session.Config{
			BaseApiUrl: "https://apidev.mobingi.com",
		})

		type data_t struct {
			Data string
		}

		data := &data_t{
			Data: "hello",
		}

		alm := New(sess)
		in := WalkerCtx{
			Data: data,
			StackCallback: func(data interface{}, ls *ListStack) error {
				_data := data.(*data_t)
				if _data.Data != "hello" {
					t.Error("should be hello")
				}

				debug.Info("stack-callback:", ls.StackId)
				return nil
			},
			InstanceCallback: func(data interface{}, ls *ListStack, flag string, inst *Instance, err error) error {
				debug.Info("instance-callback:", ls.StackId, inst.PublicDnsName)
				return nil
			},
		}

		err := alm.Walker(&in)
		if err != nil {
			t.Error(err)
		}
	}
}
