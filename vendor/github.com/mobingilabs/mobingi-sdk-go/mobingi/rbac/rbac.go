package rbac

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/mobingilabs/mobingi-sdk-go/client"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
	"github.com/pkg/errors"
)

type rbac struct {
	session *session.Session
	client  client.HttpClient
}

type DescribeRolesInput struct {
	User string
}

func (r *rbac) DescribeRoles(in *DescribeRolesInput) (*client.Response, []byte, error) {
	ep := r.session.ApiEndpoint() + "/role"
	if in != nil {
		if in.User != "" {
			ep = r.session.ApiEndpoint() + "/user/" + in.User + "/role"
		} else {
			return nil, nil, errors.New("user cannot be empty")
		}
	}

	req, err := http.NewRequest(http.MethodGet, ep, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "new request failed")
	}

	req.Header.Add("Authorization", "Bearer "+r.session.AccessToken)
	return r.client.Do(req)
}

type CreateRoleInput struct {
	Name  string `json:"name"`
	Scope Role   `json:"scope"`
}

func (r *rbac) CreateRole(in *CreateRoleInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

	if in.Name == "" {
		return nil, nil, errors.New("name cannot be empty")
	}

	p, err := json.Marshal(in)
	if err != nil {
		return nil, nil, errors.Wrap(err, "marshal failed")
	}

	ep := r.session.ApiEndpoint() + "/role"
	req, err := http.NewRequest(http.MethodPost, ep, bytes.NewBuffer(p))
	if err != nil {
		return nil, nil, errors.Wrap(err, "new request failed")
	}

	req.Header.Add("Authorization", "Bearer "+r.session.AccessToken)
	req.Header.Add("Content-Type", "application/json")
	return r.client.Do(req)
}

func New(s *session.Session) *rbac {
	if s == nil {
		return nil
	}

	var c client.HttpClient
	if s.Config.HttpClientConfig != nil {
		c = client.NewSimpleHttpClient(s.Config.HttpClientConfig)
	} else {
		c = client.NewSimpleHttpClient()
	}

	return &rbac{session: s, client: c}
}
