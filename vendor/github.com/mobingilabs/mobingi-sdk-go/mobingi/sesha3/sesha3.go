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
	// Target is the target instance for the script execution.
	// Format: stack-id|user@ip:flag
	Targets []string

	// OutputCallback is the client callback function that is called for each target's output.
	// Parameter description:
	//   int - index in the provided `Targets`.
	//   *client.Response - http response object.
	//   []byte - http response body for easy access.
	//   *TargetHeader - object version of the indexed `Target`.
	//   error - operation error
	OutputCallback func(int, *client.Response, []byte, *TargetHeader, error)

	// Script is the script contents that will be executed to all targets.
	Script []byte
}

type TargetHeader struct {
	StackId string `json:"stack_id"`
	Ip      string `json:"ip"`
	VmUser  string `json:"vm_user"`
	Flag    string `json:"flag"`
	PemUrl  string `json:"pem_url"`
}

// ExecScriptPayload is the payload we send to sesha3 server. It will be formed from
// ExecScriptInput struct.
type ExecScriptPayload struct {
	// Target is a description of the script's run target.
	Target TargetHeader `json:"target"`

	// Script is the script contents that will be executed to all targets.
	Script []byte `json:"script"`
}

func (s *sesha3) ExecScript(in *ExecScriptInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

	if len(in.Targets) < 0 {
		return nil, nil, errors.New("targets cannot be empty")
	}

	if in.Script == nil {
		return nil, nil, errors.New("script cannot be empty")
	}

	almsvc := alm.New(s.session)
	var thdrs []TargetHeader

	// our pem url list
	for _, target := range in.Targets {
		// get stack id
		sidx := strings.Split(target, "|")
		if len(sidx) != 2 {
			return nil, nil, errors.New("invalid target format: " + target)
		}

		// get flag
		xflag := strings.Split(sidx[1], ":")
		if len(xflag) != 2 {
			return nil, nil, errors.New("invalid target format: " + sidx[1])
		}

		// get user + ip
		uip := strings.Split(xflag[0], "@")
		if len(uip) != 2 {
			return nil, nil, errors.New("invalid target format: " + xflag[0])
		}

		// get pem url via stack id
		inpem := alm.GetPemInput{
			StackId: sidx[0],
		}

		if s.session.Config.ApiVersion >= 3 {
			inpem.Flag = xflag[1]
		}

		resp, body, _, err := almsvc.GetPem(&inpem)
		if err != nil {
			return resp, body, errors.Wrap(err, "get pem failed")
		}

		type rsaurl struct {
			Status string `json:"status"`
			Data   string `json:"data"`
		}

		var ru rsaurl
		err = json.Unmarshal(body, &ru)
		if err != nil {
			return resp, body, errors.Wrap(err, "url body unmarshal failed")
		}

		pemurl := strings.Replace(ru.Data, "\\", "", -1)

		// create individual target
		targethdr := TargetHeader{
			StackId: sidx[0],
			Ip:      uip[1],
			VmUser:  uip[0],
			Flag:    xflag[1],
			PemUrl:  pemurl,
		}

		thdrs = append(thdrs, targethdr)
	}

	// get sesha3 token
	resp, body, token, err := s.GetToken()
	if err != nil {
		return resp, body, errors.Wrap(err, "get token failed")
	}

	s.session.AccessToken = token
	// TODO: change endpoint to POST
	ep := s.session.Sesha3Endpoint() + "/exec"

	for i, thdr := range thdrs {
		payload := ExecScriptPayload{
			Target: thdr,
			Script: in.Script,
		}

		b, err := json.Marshal(payload)
		if err != nil {
			if in.OutputCallback != nil {
				in.OutputCallback(i, resp, body, &thdr, errors.Wrap(err, "payload marshal failed"))
			}
		}

		req := s.session.SimpleAuthRequest(http.MethodGet, ep, bytes.NewBuffer(b))
		resp, body, err = s.client.Do(req)
		if err != nil {
			if in.OutputCallback != nil {
				in.OutputCallback(i, resp, body, &thdr, errors.Wrap(err, "client do failed"))
			}
		} else {
			if in.OutputCallback != nil {
				in.OutputCallback(i, resp, body, &thdr, nil)
			}
		}
	}

	return nil, nil, nil
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

	s.session.AccessToken = token
	ep := s.session.Sesha3Endpoint() + "/ttyurl"
	req := s.session.SimpleAuthRequest(http.MethodGet, ep, bytes.NewBuffer(b))
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
