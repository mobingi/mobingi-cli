package rbac

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/mobingilabs/mobingi-sdk-go/client"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
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

	rb, err := json.Marshal(in.Scope)
	if err != nil {
		return nil, nil, errors.Wrap(err, "marshal role failed")
	}

	v := url.Values{}
	v.Set("name", in.Name)
	v.Set("scope", string(rb))
	payload := []byte(v.Encode())
	ep := r.session.ApiEndpoint() + "/role"
	req, err := http.NewRequest(http.MethodPost, ep, bytes.NewBuffer(payload))
	if err != nil {
		return nil, nil, errors.Wrap(err, "new request failed")
	}

	if r.session.Config.HttpClientConfig.Verbose {
		debug.Info("[BODY]", string(payload))
	}

	req.Header.Add("Authorization", "Bearer "+r.session.AccessToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	return r.client.Do(req)
}

type AttachRoleToUserInput struct {
	Username string
	RoleId   string
}

func (r *rbac) AttachRoleToUser(in *AttachRoleToUserInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

	if in.Username == "" {
		return nil, nil, errors.New("username cannot be empty")
	}

	if in.RoleId == "" {
		return nil, nil, errors.New("role id cannot be empty")
	}

	v := url.Values{}
	v.Set("username", in.Username)
	v.Set("role_id", in.RoleId)
	payload := []byte(v.Encode())
	ep := r.session.ApiEndpoint() + "/user/role"
	req, err := http.NewRequest(http.MethodPost, ep, bytes.NewBuffer(payload))
	if err != nil {
		return nil, nil, errors.Wrap(err, "new request failed")
	}

	req.Header.Add("Authorization", "Bearer "+r.session.AccessToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	return r.client.Do(req)
}

type DeleteRoleInput struct {
	RoleId string
}

func (r *rbac) DeleteRole(in *DeleteRoleInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

	if in.RoleId == "" {
		return nil, nil, errors.New("role id cannot be empty")
	}

	ep := r.session.ApiEndpoint() + "/role/" + in.RoleId
	req, err := http.NewRequest(http.MethodDelete, ep, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "new request failed")
	}

	req.Header.Add("Authorization", "Bearer "+r.session.AccessToken)
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
