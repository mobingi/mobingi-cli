package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/mobingilabs/mobingi-sdk-go/client"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
	"github.com/pkg/errors"
)

type registry struct {
	session *session.Session
	client  client.HttpClient
}

type GetUserCatalogInput struct {
	Service string
	Scope   string
}

func (r *registry) GetUserCatalog(in *GetUserCatalogInput) (*client.Response, []byte, []string, error) {
	if in == nil {
		return nil, nil, nil, errors.New("input cannot be nil")
	}

	if in.Service == "" {
		in.Service = "Mobingi Docker Registry"
	}

	tokenIn := &GetRegistryTokenInput{
		Service: in.Service,
		Scope:   "registry:catalog:*",
	}

	resp, body, token, err := r.GetRegistryToken(tokenIn)
	if err != nil {
		return resp, body, nil, errors.Wrap(err, "get token failed")
	}

	r.session.AccessToken = token
	ep := r.session.RegistryEndpoint() + "/_catalog"
	req, err := http.NewRequest(http.MethodGet, ep, nil)
	req.Header.Add("Authorization", "Bearer "+r.session.AccessToken)
	resp, body, err = r.client.Do(req)
	if err != nil {
		return resp, body, nil, errors.Wrap(err, "client do failed")
	}

	type catalog struct {
		Repositories []string `json:"repositories"`
	}

	var ct catalog
	err = json.Unmarshal(body, &ct)
	if err != nil {
		return resp, nil, nil, errors.Wrap(err, "unmarshal failed")
	}

	ret := make([]string, 0)
	for _, v := range ct.Repositories {
		pair := strings.Split(v, "/")
		if len(pair) == 2 {
			if pair[0] == r.session.Config.Username {
				ret = append(ret, v)
			}
		}
	}

	return resp, nil, ret, nil
}

type GetRegistryTokenInput struct {
	Service string
	Scope   string
}

func (r *registry) GetRegistryToken(in *GetRegistryTokenInput) (*client.Response, []byte, string, error) {
	var token string

	if in == nil {
		return nil, nil, token, errors.New("input cannot be nil")
	}

	if in.Service == "" {
		in.Service = "Mobingi Docker Registry"
	}

	values := url.Values{}
	values.Add("service", in.Service)
	values.Add("scope", in.Scope)
	ep := r.session.ApiEndpoint() + "/docker/token"
	req, err := http.NewRequest(http.MethodGet, ep, nil)
	req.SetBasicAuth(r.session.Config.Username, r.session.Config.Password)
	req.URL.RawQuery = values.Encode()
	resp, body, err := r.client.Do(req)
	if err != nil {
		return resp, body, token, errors.Wrap(err, "client do failed")
	}

	var m map[string]interface{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		return resp, body, token, errors.Wrap(err, "unmarshal failed")
	}

	t, found := m["token"]
	if !found {
		return resp, body, token, errors.New("cannot find token")
	}

	token = fmt.Sprintf("%s", t)
	return resp, body, token, nil
}

type GetTagDigestInput struct {
	ImageName string
	TagName   string
}

func (r *registry) GetTagDigest(in *GetTagDigestInput) (*client.Response, []byte, string, error) {
	var digest string

	/*
		if in == nil {
			return nil, nil, digest, errors.New("input cannot be nil")
		}

		ep := r.session.RegistryEndpoint() + "/alm/pem?stack_id=" + in.StackId
		path := fmt.Sprintf("/%s/%s/manifests/%s", userpass.Username, pair[0], pair[1])
		req, err := http.NewRequest(http.MethodGet, ep, nil)
	*/

	return nil, nil, digest, nil
}

func New(s *session.Session) *registry {
	if s == nil {
		return nil
	}

	var c client.HttpClient
	if s.Config.HttpClientConfig != nil {
		c = client.NewSimpleHttpClient(s.Config.HttpClientConfig)
	} else {
		c = client.NewSimpleHttpClient()
	}

	return &registry{session: s, client: c}
}
