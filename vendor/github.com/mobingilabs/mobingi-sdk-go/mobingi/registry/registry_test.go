package registry

import (
	"os"
	"testing"

	"github.com/mobingilabs/mobingi-sdk-go/client"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
)

func TestGetRegistryTokenDevAcct(t *testing.T) {
	return
	if os.Getenv("MOBINGI_CLIENT_ID") != "" && os.Getenv("MOBINGI_CLIENT_SECRET") != "" &&
		os.Getenv("MOBINGI_USERNAME") != "" && os.Getenv("MOBINGI_PASSWORD") != "" {
		sess, _ := session.New(&session.Config{
			BaseApiUrl: "https://apidev.mobingi.com",
		})

		reg := New(sess)
		in := &GetRegistryTokenInput{
			Scope: "repository:" + os.Getenv("MOBINGI_USERNAME") + "/hello:*",
		}

		resp, body, token, err := reg.GetRegistryToken(in)
		if err != nil {
			t.Errorf("expecting nil error, received %v", err)
		}

		// log.Println(resp.Status, string(body), token)
		_, _, _ = resp, body, token
	}
}

func TestGetUserCatalogDevAcct(t *testing.T) {
	return
	if os.Getenv("MOBINGI_CLIENT_ID") != "" && os.Getenv("MOBINGI_CLIENT_SECRET") != "" &&
		os.Getenv("MOBINGI_USERNAME") != "" && os.Getenv("MOBINGI_PASSWORD") != "" {
		sess, _ := session.New(&session.Config{
			BaseApiUrl:      "https://apidev.mobingi.com",
			BaseRegistryUrl: "https://dockereg2.labs.mobingi.com",
		})

		reg := New(sess)
		in := &GetUserCatalogInput{}
		resp, body, list, err := reg.GetUserCatalog(in)
		if err != nil {
			t.Errorf("expecting nil error, received %v", err)
		}

		// log.Println(resp, string(body), list)
		_, _, _ = resp, body, list
	}
}

func TestGetTagsListDevAcct(t *testing.T) {
	return
	if os.Getenv("MOBINGI_CLIENT_ID") != "" && os.Getenv("MOBINGI_CLIENT_SECRET") != "" &&
		os.Getenv("MOBINGI_USERNAME") != "" && os.Getenv("MOBINGI_PASSWORD") != "" {
		sess, _ := session.New(&session.Config{
			BaseApiUrl:      "https://apidev.mobingi.com",
			BaseRegistryUrl: "https://dockereg2.labs.mobingi.com",
		})

		reg := New(sess)
		in := &GetTagsListInput{
			Image: "hello",
		}

		resp, body, err := reg.GetTagsList(in)
		if err != nil {
			t.Errorf("expecting nil error, received %v", err)
		}

		// log.Println(resp, string(body))
		_, _ = resp, body
	}
}

func TestGetTagDigestDevAcct(t *testing.T) {
	return
	if os.Getenv("MOBINGI_CLIENT_ID") != "" && os.Getenv("MOBINGI_CLIENT_SECRET") != "" &&
		os.Getenv("MOBINGI_USERNAME") != "" && os.Getenv("MOBINGI_PASSWORD") != "" {
		sess, _ := session.New(&session.Config{
			BaseApiUrl:      "https://apidev.mobingi.com",
			BaseRegistryUrl: "https://dockereg2.labs.mobingi.com",
		})

		reg := New(sess)
		in := &GetTagManifestInput{
			Image: "hello",
			Tag:   "latest",
		}

		resp, body, err := reg.GetTagManifest(in)
		if err != nil {
			t.Errorf("expecting nil error, received %v", err)
		}

		// log.Println(resp, string(body))
		_, _ = resp, body
	}
}

func TestDeleteImageDevAcct(t *testing.T) {
	return
	if os.Getenv("MOBINGI_CLIENT_ID") != "" && os.Getenv("MOBINGI_CLIENT_SECRET") != "" &&
		os.Getenv("MOBINGI_USERNAME") != "" && os.Getenv("MOBINGI_PASSWORD") != "" {
		sess, _ := session.New(&session.Config{
			BaseApiUrl: "https://apidev.mobingi.com",
			HttpClientConfig: &client.Config{
				Verbose: true,
			},
		})

		reg := New(sess)
		in := &DeleteImageInput{
			Image: "hello",
		}

		resp, body, err := reg.DeleteImage(in)
		if err != nil {
			t.Errorf("expecting nil error, received %v", err)
		}

		// log.Println(resp, string(body))
		_, _ = resp, body
	}
}
