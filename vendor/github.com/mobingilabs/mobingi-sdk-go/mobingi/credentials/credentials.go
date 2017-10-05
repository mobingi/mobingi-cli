package credentials

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/mobingilabs/mobingi-sdk-go/client"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
	"github.com/pkg/errors"
)

type creds struct {
	session *session.Session
	client  client.HttpClient
}

func (c *creds) User() (*client.Response, []byte, error) {
	ep := c.session.ApiEndpoint() + `/user?username=` + c.session.Config.Username
	req := c.session.SimpleAuthRequest(http.MethodGet, ep, nil)
	return c.client.Do(req)
}

type CredentialsListInput struct {
	Vendor string
}

func (c *creds) List(in *CredentialsListInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

	// default to aws
	if in.Vendor == "" {
		in.Vendor = "aws"
	}

	ep := c.session.ApiEndpoint() + "/credentials/" + in.Vendor
	req := c.session.SimpleAuthRequest(http.MethodGet, ep, nil)
	return c.client.Do(req)
}

type AddVendorCredentialsInput struct {
	Vendor       string
	ClientId     string
	ClientSecret string
	AcctName     string
}

func (c *creds) AddVendorCredentials(in *AddVendorCredentialsInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

	if in.Vendor == "" {
		return nil, nil, errors.New("vendor cannot be empty")
	}

	if in.ClientId == "" {
		return nil, nil, errors.New("client id cannot be empty")
	}

	if in.ClientSecret == "" {
		return nil, nil, errors.New("client secret cannot be empty")
	}

	if in.AcctName == "" {
		return nil, nil, errors.New("acct name cannot be empty")
	}

	var payload []byte
	switch in.Vendor {
	case "aws":
		type add_t struct {
			Credentials AWSCredentials `json:"credentials,omitempty"`
		}

		creds := add_t{
			Credentials: AWSCredentials{
				Name:   in.AcctName,
				KeyId:  in.ClientId,
				Secret: in.ClientSecret,
			},
		}

		b, err := json.Marshal(creds)
		if err != nil {
			return nil, nil, errors.Wrap(err, "marshal aws payload failed")
		}

		payload = b
	case "alicloud":
		type add_t struct {
			Credentials AliCloudCredentials `json:"credentials,omitempty"`
		}

		creds := &add_t{
			Credentials: AliCloudCredentials{
				Name:   in.AcctName,
				KeyId:  in.ClientId,
				Secret: in.ClientSecret,
			},
		}

		log.Printf("%+v", creds)

		b, err := json.Marshal(creds)
		if err != nil {
			return nil, nil, errors.Wrap(err, "marshal aws payload failed")
		}

		payload = b
	default:
		return nil, nil, errors.New("vendor not supported")
	}

	ep := c.session.ApiEndpoint() + "/credentials/" + in.Vendor
	req := c.session.SimpleAuthRequest(http.MethodPost, ep, bytes.NewBuffer(payload))
	return c.client.Do(req)
}

func New(s *session.Session) *creds {
	if s == nil {
		return nil
	}

	var c client.HttpClient
	if s.Config.HttpClientConfig != nil {
		c = client.NewSimpleHttpClient(s.Config.HttpClientConfig)
	} else {
		c = client.NewSimpleHttpClient()
	}

	return &creds{session: s, client: c}
}
