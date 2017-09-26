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

	if in.Scope == "" {
		in.Scope = "registry:catalog:*"
	}

	tokenIn := &GetRegistryTokenInput{
		Service: in.Service,
		Scope:   in.Scope,
	}

	resp, body, token, err := r.GetRegistryToken(tokenIn)
	if err != nil {
		return resp, body, nil, errors.Wrap(err, "get token failed")
	}

	r.session.AccessToken = token
	ep := r.session.RegistryEndpoint() + "/_catalog"
	req := r.session.SimpleAuthRequest(http.MethodGet, ep, nil)
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

type DescribeImageInput struct {
	Image string
}

func (r *registry) DescribeImage(in *DescribeImageInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

	if in.Image == "" {
		return nil, nil, errors.New("image name cannot be empty")
	}

	values := url.Values{}
	values.Add("targetKey", "repository")
	values.Add("targetValue", r.session.Config.Username+"/"+in.Image)
	ep := r.session.ApiEndpoint() + "/alm/registry"
	req := r.session.SimpleAuthRequest(http.MethodGet, ep, nil)
	req.URL.RawQuery = values.Encode()
	resp, body, err := r.client.Do(req)
	if err != nil {
		return resp, body, errors.Wrap(err, "client do failed")
	}

	return resp, body, nil
}

type GetTagsListInput struct {
	ManualOp bool // do not use api endpoint, other than login for token
	Service  string
	Scope    string
	Image    string
}

func (r *registry) GetTagsList(in *GetTagsListInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

	if in.Image == "" {
		return nil, nil, errors.New("image cannot be empty")
	}

	if !in.ManualOp {
		values := url.Values{}
		values.Add("account_id", r.session.Config.Username)
		values.Add("image_id", r.session.Config.Username+"/"+in.Image)
		ep := r.session.ApiEndpoint() + `/alm/registry/imagetags`
		req := r.session.SimpleAuthRequest(http.MethodGet, ep, nil)
		req.URL.RawQuery = values.Encode()
		resp, body, err := r.client.Do(req)
		if err != nil {
			return resp, body, errors.Wrap(err, "client do failed")
		}

		return resp, body, nil
	}

	if in.Service == "" {
		in.Service = "Mobingi Docker Registry"
	}

	if in.Scope == "" {
		in.Scope = fmt.Sprintf("repository:%s/%s:pull", r.session.Config.Username, in.Image)
	}

	tokenIn := &GetRegistryTokenInput{
		Service: in.Service,
		Scope:   in.Scope,
	}

	resp, body, token, err := r.GetRegistryToken(tokenIn)
	if err != nil {
		return resp, body, errors.Wrap(err, "get token failed")
	}

	r.session.AccessToken = token
	ep := r.session.RegistryEndpoint() + "/" + r.session.Config.Username + "/" + in.Image + "/tags/list"
	req := r.session.SimpleAuthRequest(http.MethodGet, ep, nil)
	resp, body, err = r.client.Do(req)
	if err != nil {
		return resp, body, errors.Wrap(err, "client do failed")
	}

	return resp, body, nil
}

type UpdateDescriptionInput struct {
	Image       string
	Description string
}

func (r *registry) UpdateDescription(in *UpdateDescriptionInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

	if in.Image == "" {
		return nil, nil, errors.New("image cannot be empty")
	}

	values := url.Values{}
	values.Add("image_id", r.session.Config.Username+"/"+in.Image)
	values.Add("description", in.Description)
	ep := r.session.ApiEndpoint() + `/alm/registry/description`
	req := r.session.SimpleAuthRequest(http.MethodPut, ep, nil)
	req.URL.RawQuery = values.Encode()
	resp, body, err := r.client.Do(req)
	if err != nil {
		return resp, body, errors.Wrap(err, "client do failed")
	}

	return resp, body, nil
}

type GetTagManifestInput struct {
	Service string
	Scope   string
	Image   string
	Tag     string
}

func (r *registry) GetTagManifest(in *GetTagManifestInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

	if in.Service == "" {
		in.Service = "Mobingi Docker Registry"
	}

	if in.Scope == "" {
		in.Scope = fmt.Sprintf("repository:%s/%s:pull", r.session.Config.Username, in.Image)
	}

	tokenIn := &GetRegistryTokenInput{
		Service: in.Service,
		Scope:   in.Scope,
	}

	resp, body, token, err := r.GetRegistryToken(tokenIn)
	if err != nil {
		return resp, body, errors.Wrap(err, "get token failed")
	}

	r.session.AccessToken = token
	ep := r.session.RegistryEndpoint() + "/" + r.session.Config.Username + "/" + in.Image + "/manifests/" + in.Tag
	req := r.session.SimpleAuthRequest(http.MethodGet, ep, nil)
	resp, body, err = r.client.Do(req)
	if err != nil {
		return resp, body, errors.Wrap(err, "client do failed")
	}

	return resp, body, nil
}

type UpdateVisibilityInput struct {
	Image      string
	Visibility string // public or private
}

func (r *registry) UpdateVisibility(in *UpdateVisibilityInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

	if in.Image == "" {
		return nil, nil, errors.New("image cannot be empty")
	}

	if in.Visibility == "" {
		return nil, nil, errors.New("visibility should be private or public")
	}

	values := url.Values{}
	values.Add("image_id", r.session.Config.Username+"/"+in.Image)
	values.Add("visibility", in.Visibility)
	ep := r.session.ApiEndpoint() + `/alm/registry/visibility`
	req := r.session.SimpleAuthRequest(http.MethodPut, ep, nil)
	req.URL.RawQuery = values.Encode()
	resp, body, err := r.client.Do(req)
	if err != nil {
		return resp, body, errors.Wrap(err, "client do failed")
	}

	return resp, body, nil
}

type DeleteImageInput struct {
	Image string
}

func (r *registry) DeleteImage(in *DeleteImageInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

	if in.Image == "" {
		return nil, nil, errors.New("image cannot be empty")
	}

	values := url.Values{}
	values.Add("image_id", r.session.Config.Username+"/"+in.Image)
	ep := r.session.ApiEndpoint() + `/alm/registry/image`
	req := r.session.SimpleAuthRequest(http.MethodDelete, ep, nil)
	req.URL.RawQuery = values.Encode()
	resp, body, err := r.client.Do(req)
	if err != nil {
		return resp, body, errors.Wrap(err, "client do failed")
	}

	return resp, body, nil
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
