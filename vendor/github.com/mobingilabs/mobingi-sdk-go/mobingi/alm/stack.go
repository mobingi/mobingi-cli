package alm

import (
	"net/http"

	"github.com/mobingilabs/mobingi-sdk-go/client"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
	"github.com/pkg/errors"
)

type stack struct {
	session *session.Session
	client  client.HttpClient
}

func (s *stack) List() (*http.Response, []byte, error) {
	req, err := http.NewRequest(http.MethodGet, s.session.ApiEndpoint()+"/alm/stack", nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "new request failed")
	}

	req.Header.Add("Authorization", "Bearer "+s.session.AccessToken)
	return s.client.Do(req)
}

func New(s *session.Session) *stack {
	if s == nil {
		return nil
	}

	var c client.HttpClient
	if s.Config.HttpClientConfig != nil {
		c = client.NewSimpleHttpClient(s.Config.HttpClientConfig)
	} else {
		c = client.NewSimpleHttpClient()
	}

	return &stack{session: s, client: c}
}
