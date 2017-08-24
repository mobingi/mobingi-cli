package svrconf

import (
	"net/http"

	"github.com/mobingilabs/mobingi-sdk-go/client"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
	"github.com/pkg/errors"
)

type ServerConfigGetInput struct {
	StackId string
}

type cnf struct {
	session *session.Session
	client  client.HttpClient
}

func (s *cnf) Get(in *ServerConfigGetInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

	if in.StackId == "" {
		return nil, nil, errors.New("stack id cannot be empty")
	}

	ep := s.session.ApiEndpoint() + `/alm/serverconfig?stack_id=` + in.StackId
	req, err := http.NewRequest(http.MethodGet, ep, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "new request failed")
	}

	req.Header.Add("Authorization", "Bearer "+s.session.AccessToken)
	return s.client.Do(req)
}

func New(s *session.Session) *cnf {
	if s == nil {
		return nil
	}

	var c client.HttpClient
	if s.Config.HttpClientConfig != nil {
		c = client.NewSimpleHttpClient(s.Config.HttpClientConfig)
	} else {
		c = client.NewSimpleHttpClient()
	}

	return &cnf{session: s, client: c}
}
