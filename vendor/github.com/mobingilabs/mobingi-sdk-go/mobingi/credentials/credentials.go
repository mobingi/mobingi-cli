package credentials

import (
	"net/http"

	"github.com/mobingilabs/mobingi-sdk-go/client"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
	"github.com/pkg/errors"
)

type CredentialsListInput struct {
	Vendor string
}

type creds struct {
	session *session.Session
	client  client.HttpClient
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
	req, err := http.NewRequest(http.MethodGet, ep, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "new request failed")
	}

	req.Header.Add("Authorization", "Bearer "+c.session.AccessToken)
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
