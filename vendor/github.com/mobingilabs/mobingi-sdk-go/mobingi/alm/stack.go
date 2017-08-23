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

type StackDescribeInput struct {
	StackId string
}

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

type StackCreateInput struct {
	Vendor         string
	Region         string
	CredId         string
	Configurations interface{}
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

func (s *stack) Create(in *StackCreateInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

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

func (s *stack) Update(in *StackCreateInput) (*client.Response, []byte, error) {
	return nil, nil, nil
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
