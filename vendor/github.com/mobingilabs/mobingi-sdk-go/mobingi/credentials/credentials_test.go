package credentials

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mobingilabs/mobingi-sdk-go/client"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
)

func TestNew(t *testing.T) {
	sess, _ := session.New(&session.Config{})
	creds := New(sess)
	if creds == nil {
		t.Errorf("Expecting non-nil")
	}
}

func TestUserDevAcct(t *testing.T) {
	return
	if os.Getenv("MOBINGI_ROOT_USERNAME") != "" && os.Getenv("MOBINGI_ROOT_PASSWORD") != "" {
		// if os.Getenv("MOBINGI_USERNAME") != "" && os.Getenv("MOBINGI_PASSWORD") != "" {
		sess, _ := session.New(&session.Config{
			Username: os.Getenv("MOBINGI_ROOT_USERNAME"),
			Password: os.Getenv("MOBINGI_ROOT_PASSWORD"),
			// Username:         os.Getenv("MOBINGI_USERNAME"),
			// Password:         os.Getenv("MOBINGI_PASSWORD"),
			BaseApiUrl:       "https://apidev.mobingi.com",
			BaseRegistryUrl:  "https://dockereg2.labs.mobingi.com",
			HttpClientConfig: &client.Config{Verbose: true},
		})

		svc := New(sess)
		resp, body, err := svc.User()
		if err != nil {
			t.Errorf("expecting nil error, received %v", err)
		}

		var u UserDetails
		err = json.Unmarshal(body, &u)
		if err != nil {
			t.Fatal("error:", err)
		}

		if u.UserId == "" {
			t.Fatal("user_id should at least be set")
		}

		_ = resp
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

func TestAddVendorCredentialsDevAcct(t *testing.T) {
	return
	if os.Getenv("MOBINGI_TESTADD_CLIENT_ID") != "" && os.Getenv("MOBINGI_TESTADD_CLIENT_SECRET") != "" {
		sess, _ := session.New(&session.Config{
			Username:         os.Getenv("MOBINGI_ROOT_USERNAME"),
			Password:         os.Getenv("MOBINGI_ROOT_PASSWORD"),
			BaseApiUrl:       "https://apidev.mobingi.com",
			BaseRegistryUrl:  "https://dockereg2.labs.mobingi.com",
			HttpClientConfig: &client.Config{Verbose: true},
		})

		svc := New(sess)
		in := &AddVendorCredentialsInput{
			Vendor:       "alicloud",
			ClientId:     os.Getenv("MOBINGI_TESTADD_CLIENT_ID"),
			ClientSecret: os.Getenv("MOBINGI_TESTADD_CLIENT_SECRET"),
			AcctName:     "sdktest",
		}

		resp, body, err := svc.AddVendorCredentials(in)
		if err != nil {
			t.Errorf("expecting nil error, received %v", err)
		}

		log.Println(resp, string(body))
		// _, _ = resp, body
	}
}
