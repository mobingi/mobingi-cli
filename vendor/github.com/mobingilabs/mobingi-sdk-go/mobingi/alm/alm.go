package alm

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/mobingilabs/mobingi-sdk-go/client"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/credentials"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
	"github.com/pkg/errors"
)

type StackCreateDb struct {
	Engine       string `json:"Engine,omitempty"`
	Type         string `json:"DBType,omitempty"`
	Storage      string `json:"DBStorage,omitempty"`
	ReadReplica1 bool   `json:"ReadReplica1,omitempty"`
	ReadReplica2 bool   `json:"ReadReplica2,omitempty"`
	ReadReplica3 bool   `json:"ReadReplica3,omitempty"`
	ReadReplica4 bool   `json:"ReadReplica4,omitempty"`
	ReadReplica5 bool   `json:"ReadReplica5,omitempty"`
}

type StackCreateElasticache struct {
	Engine    string `json:"ElastiCacheEngine,omitempty"`
	NodeType  string `json:"ElastiCacheNodeType,omitempty"`
	NodeCount string `json:"ElastiCacheNodes,omitempty"`
}

type StackCreateConfig struct {
	Region            interface{} `json:"region,omitempty"`
	Architecture      interface{} `json:"architecture,omitempty"`
	Type              interface{} `json:"type,omitempty"`
	Image             interface{} `json:"image,omitempty"`
	DockerHubUsername interface{} `json:"dockerHubUsername,omitempty"`
	DockerHubPassword interface{} `json:"dockerHubPassword,omitempty"`
	Min               interface{} `json:"min,omitempty"`
	Max               interface{} `json:"max,omitempty"`
	SpotRange         interface{} `json:"spotRange,omitempty"`
	Nickname          interface{} `json:"nickname,omitempty"`
	Code              interface{} `json:"code,omitempty"`
	GitReference      interface{} `json:"gitReference,omitempty"`
	GitPrivateKey     interface{} `json:"gitPrivateKey,omitempty"`
	Database          interface{} `json:"database,omitempty"`
	ElastiCache       interface{} `json:"elasticache,omitempty"`
}

type AlmTemplate struct {
	ContentType string // json, yaml
	Contents    string
}

type stack struct {
	session *session.Session
	client  client.HttpClient
}

func (s *stack) List() (*client.Response, []byte, error) {
	req, err := http.NewRequest(http.MethodGet, s.session.ApiEndpoint()+"/alm/stack", nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "new request failed")
	}

	req.Header.Add("Authorization", "Bearer "+s.session.AccessToken)
	return s.client.Do(req)
}

type StackDescribeInput struct {
	StackId string
}

func (s *stack) Describe(in *StackDescribeInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

	if in.StackId == "" {
		return nil, nil, errors.New("stack id cannot be empty")
	}

	ep := s.session.ApiEndpoint() + "/alm/stack/" + in.StackId
	req, err := http.NewRequest(http.MethodGet, ep, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "new request failed")
	}

	req.Header.Add("Authorization", "Bearer "+s.session.AccessToken)
	return s.client.Do(req)
}

type StackCreateInput struct {
	AlmTemplate    *AlmTemplate // if not nil, we use this for creation, discard others
	Vendor         string
	Region         string
	CredId         string
	Configurations interface{} // of type StackCreateConfig
}

func (s *stack) Create(in *StackCreateInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

	if in.AlmTemplate != nil {
		return s.createAlmStack(in)
	}

	return s.createStackV2(in)
}

type StackUpdateInput struct {
	AlmTemplate    *AlmTemplate // if not nil, we use this for update instead of Configurations
	StackId        string
	Configurations interface{} // of type StackCreateConfig
}

func (s *stack) Update(in *StackUpdateInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

	if in.AlmTemplate != nil {
		return s.updateAlmStack(in)
	}

	return s.updateStackV2(in)
}

type StackDeleteInput struct {
	StackId string
}

func (s *stack) Delete(in *StackDeleteInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

	if in.StackId == "" {
		return nil, nil, errors.New("stack id cannot be empty")
	}

	ep := s.session.ApiEndpoint() + "/alm/stack/" + in.StackId
	req, err := http.NewRequest(http.MethodDelete, ep, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "new request failed")
	}

	req.Header.Add("Authorization", "Bearer "+s.session.AccessToken)
	return s.client.Do(req)
}

type GetTemplateVersionsInput struct {
	StackId string
}

func (s *stack) GetTemplateVersions(in *GetTemplateVersionsInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

	if in.StackId == "" {
		return nil, nil, errors.New("stack id cannot be empty")
	}

	ep := s.session.ApiEndpoint() + "/alm/template?stack_id=" + in.StackId
	req, err := http.NewRequest(http.MethodGet, ep, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "new request failed")
	}

	req.Header.Add("Authorization", "Bearer "+s.session.AccessToken)
	return s.client.Do(req)
}

type DescribeTemplateInput struct {
	StackId   string
	VersionId string // can be empty or 'latest'
}

func (s *stack) DescribeTemplate(in *DescribeTemplateInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

	if in.StackId == "" {
		return nil, nil, errors.New("stack id cannot be empty")
	}

	var param string
	if in.VersionId == "" || in.VersionId == "latest" {
		param = "?version_id=latest"
	}

	if in.VersionId != "" && in.VersionId != "latest" {
		param = "?version_id=" + in.VersionId
	}

	ep := s.session.ApiEndpoint() + "/alm/template/" + in.StackId + param
	req, err := http.NewRequest(http.MethodGet, ep, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "new request failed")
	}

	req.Header.Add("Authorization", "Bearer "+s.session.AccessToken)
	return s.client.Do(req)
}

type CompareTemplateInput struct {
	SourceStackId   string // required
	SourceVersionId string // required
	TargetStackId   string // optional, if empty, use SourceStackId
	TargetVersionId string // optional, can be this or TargetBody
	TargetBody      string // optional, can be this or TargetVersionId
}

func (s *stack) CompareTemplate(in *CompareTemplateInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

	if in.SourceStackId == "" {
		return nil, nil, errors.New("source stack id cannot be empty")
	}

	if in.SourceVersionId == "" {
		return nil, nil, errors.New("source version id cannot be empty")
	}

	if in.TargetStackId == "" {
		in.TargetStackId = in.SourceStackId
	}

	if in.TargetVersionId == "" && in.TargetBody == "" {
		return nil, nil, errors.New("should provide either version id or body as target")
	}

	type payload_t struct {
		Id   []json.RawMessage `json:"id,omitempty"`
		Body []string          `json:"body,omitempty"`
	}

	var set bool
	var payload payload_t

	payload.Id = make([]json.RawMessage, 0)
	srcid := json.RawMessage(`{"` + in.SourceStackId + `":{"version":"` + in.SourceVersionId + `"}}`)
	payload.Id = append(payload.Id, srcid)
	if in.TargetVersionId != "" {
		tgtid := json.RawMessage(`{"` + in.TargetStackId + `":{"version":"` + in.TargetVersionId + `"}}`)
		payload.Id = append(payload.Id, tgtid)
		set = true
	}

	if !set {
		if in.TargetBody != "" {
			payload.Body = make([]string, 0)
			payload.Body = append(payload.Body, in.TargetBody)
		}
	}

	p, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, errors.Wrap(err, "marshal failed")
	}

	ep := s.session.ApiEndpoint() + "/alm/template/compare"
	req, err := http.NewRequest(http.MethodPost, ep, bytes.NewBuffer(p))
	if err != nil {
		return nil, nil, errors.Wrap(err, "new request failed")
	}

	req.Header.Add("Authorization", "Bearer "+s.session.AccessToken)
	req.Header.Add("Content-Type", "application/json")
	return s.client.Do(req)
}

type GetPemInput struct {
	StackId string
}

func (s *stack) GetPem(in *GetPemInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

	if in.StackId == "" {
		return nil, nil, errors.New("stack id cannot be empty")
	}

	ep := s.session.ApiEndpoint() + "/alm/pem?stack_id=" + in.StackId
	req, err := http.NewRequest(http.MethodGet, ep, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "new request failed")
	}

	req.Header.Add("Authorization", "Bearer "+s.session.AccessToken)
	return s.client.Do(req)
}

func (s *stack) getCredsList(vendor string) ([]credentials.VendorCredentials, error) {
	creds := credentials.New(s.session)
	_, body, err := creds.List(&credentials.CredentialsListInput{
		Vendor: vendor,
	})

	var list []credentials.VendorCredentials
	err = json.Unmarshal(body, &list)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal failed")
	}

	_ = err
	return list, nil
}

func (s *stack) createStackV2(in *StackCreateInput) (*client.Response, []byte, error) {
	if in.Vendor == "" {
		in.Vendor = "aws"
	}

	if in.CredId == "" {
		list, err := s.getCredsList(in.Vendor)
		if err != nil {
			return nil, nil, errors.Wrap(err, "get creds list failed")
		}

		if len(list) > 0 {
			in.CredId = list[0].Id
		}
	}

	mi, err := json.Marshal(&in.Configurations)
	if err != nil {
		return nil, nil, errors.Wrap(err, "marshal failed")
	}

	v := url.Values{}
	v.Set("vendor", in.Vendor)
	v.Set("region", in.Region)
	v.Set("cred", in.CredId)
	v.Set("configurations", string(mi))
	payload := []byte(v.Encode())
	req, err := http.NewRequest(
		http.MethodPost,
		s.session.ApiEndpoint()+"/alm/stack",
		bytes.NewBuffer(payload))
	if err != nil {
		return nil, nil, errors.Wrap(err, "new request failed")
	}

	req.Header.Add("Authorization", "Bearer "+s.session.AccessToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	return s.client.Do(req)
}

func (s *stack) updateStackV2(in *StackUpdateInput) (*client.Response, []byte, error) {
	if in.StackId == "" {
		return nil, nil, errors.New("stack id cannot be empty")
	}

	type updatet struct {
		Configurations string `json:"configurations,omitempty"`
	}

	mi, err := json.Marshal(&in.Configurations)
	if err != nil {
		return nil, nil, errors.Wrap(err, "marshal config failed")
	}

	p := updatet{}
	p.Configurations = string(mi)

	mi, err = json.Marshal(&p)
	if err != nil {
		return nil, nil, errors.Wrap(err, "marshal payload failed")
	}

	ep := s.session.ApiEndpoint() + "/alm/stack/" + in.StackId
	req, err := http.NewRequest(http.MethodPut, ep, bytes.NewBuffer(mi))
	if err != nil {
		return nil, nil, errors.Wrap(err, "new request failed")
	}

	req.Header.Add("Authorization", "Bearer "+s.session.AccessToken)
	req.Header.Add("Content-Type", "application/json")
	return s.client.Do(req)
}

func (s *stack) createAlmStack(in *StackCreateInput) (*client.Response, []byte, error) {
	var ct string
	if in.AlmTemplate.Contents == "" {
		return nil, nil, errors.New("contents cannot be empty")
	}

	switch in.AlmTemplate.ContentType {
	case "json":
		ct = "application/json"
	case "yaml":
		ct = "application/x-yaml" // same with Ruby on Rails
	default:
		return nil, nil, errors.New("invalid content type; should be json or yaml")
	}

	ep := s.session.ApiEndpoint() + "/alm/template"
	req, err := http.NewRequest(http.MethodPost, ep, bytes.NewBuffer([]byte(in.AlmTemplate.Contents)))
	if err != nil {
		return nil, nil, errors.Wrap(err, "new request failed")
	}

	req.Header.Add("Authorization", "Bearer "+s.session.AccessToken)
	req.Header.Add("Content-Type", ct)
	return s.client.Do(req)
}

func (s *stack) updateAlmStack(in *StackUpdateInput) (*client.Response, []byte, error) {
	if in.StackId == "" {
		return nil, nil, errors.New("stack id cannot be empty")
	}

	var ct string
	if in.AlmTemplate.Contents == "" {
		return nil, nil, errors.New("contents cannot be empty")
	}

	switch in.AlmTemplate.ContentType {
	case "json":
		ct = "application/json"
	case "yaml":
		ct = "application/x-yaml" // same with Ruby on Rails
	default:
		return nil, nil, errors.New("invalid content type; should be json or yaml")
	}

	ep := s.session.ApiEndpoint() + "/alm/template/" + in.StackId
	req, err := http.NewRequest(http.MethodPut, ep, bytes.NewBuffer([]byte(in.AlmTemplate.Contents)))
	if err != nil {
		return nil, nil, errors.Wrap(err, "new request failed")
	}

	req.Header.Add("Authorization", "Bearer "+s.session.AccessToken)
	req.Header.Add("Content-Type", ct)
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
