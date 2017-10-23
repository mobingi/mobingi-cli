package sesha3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mobingilabs/mobingi-sdk-go/client"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/alm"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
	"github.com/moul/gotty-client"
	"github.com/pkg/errors"
)

type SeshaClientInput struct {
	URL string
}

type sesha3Client struct {
	client *gottyclient.Client
}

func NewClient(in *SeshaClientInput) (*sesha3Client, error) {
	var err error
	if len(in.URL) < 1 {
		err = errors.Wrap(err, "url should not be empty")
		return nil, err
	}

	client, err := gottyclient.NewClient(in.URL)
	if err != nil {
		err = errors.Wrap(err, "sesha3 client creation failed")
		return nil, err
	}

	return &sesha3Client{client: client}, err
}

func (c *sesha3Client) Run() error {
	err := c.client.Loop()
	if err != nil {
		err = errors.Wrap(err, "sesha3 run failed")
	}

	return err
}

type sesha3 struct {
	session *session.Session
	client  client.HttpClient
}

func (s *sesha3) GetToken() (*client.Response, []byte, string, error) {
	var token string

	tp := TokenPayload{
		Username: s.session.Config.Username,
		Passwd:   s.session.Config.Password,
	}

	b, err := json.Marshal(tp)
	if err != nil {
		return nil, nil, token, errors.New("payload token marshal failed")
	}

	ep := s.session.Sesha3Endpoint() + "/token"
	req, err := http.NewRequest(http.MethodGet, ep, bytes.NewBuffer(b))
	req.Header.Add("Content-Type", "application/json")
	resp, body, err := s.client.Do(req)
	if err != nil {
		return resp, body, token, errors.Wrap(err, "client do failed")
	}

	var m map[string]string
	err = json.Unmarshal(body, &m)
	if err != nil {
		return resp, body, token, errors.Wrap(err, "token reply unmarshal failed")
	}

	tkn, ok := m["key"]
	if !ok {
		return resp, body, token, errors.Wrap(err, "can't find token")
	}

	token = tkn
	return resp, body, token, nil
}

type ExecScriptInput struct {
	StackId    string
	Target     string
	Script     string
	ScriptName string
	InstUser   string
	Flag       string
}

type ScriptRes struct {
	Out string `json:"stdout"`
	Err string `json:"stderr"`
}

func (s *sesha3) ExecScript(in *ExecScriptInput) (*client.Response, []byte, ScriptRes, error) {
	var sresp ScriptRes

	if in == nil {
		return nil, nil, sresp, errors.New("input cannot be nil")
	}

	if in.Target == "" {
		return nil, nil, sresp, errors.New("target cannot be empty")
	}

	if in.Script == "" {
		return nil, nil, sresp, errors.New("script cannot be empty")
	}

	if in.ScriptName == "" {
		return nil, nil, sresp, errors.New("script cannot be empty")
	}

	if s.session.Config.ApiVersion >= 3 {
		if in.Flag == "" {
			return nil, nil, sresp, errors.New("flag cannot be empty")
		}
	}

	if in.InstUser == "" {
		return nil, nil, sresp, errors.New("instance username cannot be empty")
	}

	// get pem url from stack id
	almsvc := alm.New(s.session)
	inpem := alm.GetPemInput{
		StackId: in.StackId,
	}

	if s.session.Config.ApiVersion >= 3 {
		inpem.Flag = in.Flag
	}

	resp, body, _, err := almsvc.GetPem(&inpem)
	if err != nil {
		return resp, body, sresp, errors.Wrap(err, "get pem failed")
	}

	type rsaurl struct {
		Status string `json:"status"`
		Data   string `json:"data"`
	}

	var ru rsaurl
	err = json.Unmarshal(body, &ru)
	if err != nil {
		return resp, body, sresp, errors.Wrap(err, "url body unmarshal failed")
	}

	pemurl := strings.Replace(ru.Data, "\\", "", -1)

	// get sesha3 token
	_, _, token, err := s.GetToken()
	if err != nil {
		return resp, body, sresp, errors.Wrap(err, "get token failed")
	}

	type payload_t struct {
		Pem        string `json:"pem"`
		StackId    string `json:"stackid"`
		Target     string `json:"target"`
		Script     string `json:"script"`
		ScriptName string `json:"script_name"`
		User       string `json:"user"`
	}

	payload := payload_t{
		Pem:        pemurl,
		StackId:    in.StackId,
		Target:     in.Target,
		Script:     in.Script,
		ScriptName: in.ScriptName,
		User:       in.InstUser,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return resp, body, sresp, errors.Wrap(err, "payload marshal failed")
	}

	ep := s.session.Sesha3Endpoint() + "/exec"
	req, err := http.NewRequest(http.MethodGet, ep, bytes.NewBuffer(b))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	resp, body, err = s.client.Do(req)
	if err != nil {
		return resp, body, sresp, errors.Wrap(err, "client do failed")
	}

	err = json.Unmarshal(body, &sresp)
	if err != nil {
		return resp, body, sresp, errors.Wrap(err, "reply unmarshal failed")
	}

	return resp, body, sresp, nil
}

type GetSessionUrlInput struct {
	StackId  string
	IpAddr   string
	Flag     string
	InstUser string
	Timeout  int64
}

func (s *sesha3) GetSessionUrl(in *GetSessionUrlInput) (*client.Response, []byte, string, error) {
	var u string

	if in == nil {
		return nil, nil, u, errors.New("input cannot be nil")
	}

	if in.StackId == "" {
		return nil, nil, u, errors.New("stack id cannot be empty")
	}

	if in.IpAddr == "" {
		return nil, nil, u, errors.New("ip address cannot be empty")
	}

	if s.session.Config.ApiVersion >= 3 {
		if in.Flag == "" {
			return nil, nil, u, errors.New("flag cannot be empty")
		}
	}

	if in.InstUser == "" {
		return nil, nil, u, errors.New("instance username cannot be empty")
	}

	if in.Timeout == 0 {
		in.Timeout = 60
	}

	// get pem url from stack id
	almsvc := alm.New(s.session)
	inpem := alm.GetPemInput{
		StackId: in.StackId,
	}

	if s.session.Config.ApiVersion >= 3 {
		inpem.Flag = in.Flag
	}

	resp, body, _, err := almsvc.GetPem(&inpem)
	if err != nil {
		return resp, body, u, errors.Wrap(err, "get pem failed")
	}

	type rsaurl struct {
		Status string `json:"status"`
		Data   string `json:"data"`
	}

	var ru rsaurl
	err = json.Unmarshal(body, &ru)
	if err != nil {
		return resp, body, u, errors.Wrap(err, "url body unmarshal failed")
	}

	pemurl := strings.Replace(ru.Data, "\\", "", -1)

	// get sesha3 token
	_, _, token, err := s.GetToken()
	if err != nil {
		return resp, body, u, errors.Wrap(err, "get token failed")
	}

	type payload_t struct {
		Pem     string `json:"pem"`
		StackId string `json:"stackid"`
		Ip      string `json:"ip"`
		User    string `json:"user"`
		Timeout string `json:"timeout"`
	}

	payload := payload_t{
		Pem:     pemurl,
		StackId: in.StackId,
		Ip:      in.IpAddr,
		User:    in.InstUser,
		Timeout: fmt.Sprintf("%v", in.Timeout),
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return resp, body, u, errors.Wrap(err, "payload marshal failed")
	}

	ep := s.session.Sesha3Endpoint() + "/ttyurl"
	req, err := http.NewRequest(http.MethodGet, ep, bytes.NewBuffer(b))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	resp, body, err = s.client.Do(req)
	if err != nil {
		return resp, body, u, errors.Wrap(err, "client do failed")
	}

	type ttyurl_t struct {
		Url string `json:"tty_url"`
	}

	var tu ttyurl_t
	err = json.Unmarshal(body, &tu)
	if err != nil {
		return resp, body, u, errors.Wrap(err, "reply unmarshal failed")
	}

	if tu.Url != "" {
		u = tu.Url
	}

	return resp, body, u, nil
}

func New(s *session.Session) *sesha3 {
	if s == nil {
		return nil
	}

	var c client.HttpClient
	if s.Config.HttpClientConfig != nil {
		c = client.NewSimpleHttpClient(s.Config.HttpClientConfig)
	} else {
		c = client.NewSimpleHttpClient()
	}

	return &sesha3{session: s, client: c}
}
